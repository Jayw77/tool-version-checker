package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ToolData contains the display information for each tool
type ToolData struct {
	Name          string
	LatestVersion string
	RemoteVersion string
	UpToDate      bool
	Comment       string
}

// HomePageData represents the data for the home page
type HomePageData struct {
	Tools []ToolData
}

var (
	latestVersionsCache = make(map[string]string)
	mu                  sync.Mutex // Global mutex for synchronizing access to shared resources
	currentToolData     []ToolData // Global variable to store the latest data
)

func fetchToolDataPeriodically(config Config) {
	log.Info("Starting periodic data fetch...")

	fetchInterval := 10 * time.Minute

	// allow override of time to check for updates
	if config.FetchInterval > 0 {
		fetchInterval = time.Duration(config.FetchInterval) * time.Minute
	}

	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	// Initial fetch
	currentToolData = fetchToolData(config)

	for {
		select {
		case <-ticker.C:
			log.Info("Fetching tool data...")
			currentToolData = fetchToolData(config)
			log.Info("Tool data fetched successfully.")
			log.Infof("Will check again in %v...", fetchInterval)
		}
	}
}

func fetchLatestVersions(config Config) {
	var wg sync.WaitGroup
	for _, endpoint := range config.LatestVersionEndpoints {
		wg.Add(1)
		go func(e LatestVersionEndpoints) {
			defer wg.Done()
			version, err := fetchVersion(e.Url, e.JSONKey)
			if err != nil {
				log.WithFields(logrus.Fields{"endpoint": e.Url, "error": err}).Error("Error fetching latest version")
				version = "Error fetching version"
			}
			mu.Lock()
			latestVersionsCache[e.Name] = version
			mu.Unlock()
		}(endpoint)
	}
	wg.Wait()
}

func fetchToolData(config Config) []ToolData {
	log.Info("Fetching tool data for all tools...")
	fetchLatestVersions(config) // Ensure the latest versions are fetched and cached

	var toolData []ToolData
	var wg sync.WaitGroup

	for _, tool := range config.MyEndpoints {
		wg.Add(1)
		go func(t MyEndpoints) {
			defer wg.Done()

			mu.Lock()
			latestVersion := latestVersionsCache[t.Name] // Retrieve from cache
			mu.Unlock()

			var remoteVersion string
			var err error
			if t.Url != "" { // Check if there's a specific endpoint to fetch the tool's current version
				remoteVersion, err = fetchVersion(t.Url, t.JSONKey)
				if err != nil {
					log.WithFields(logrus.Fields{"tool": t.Name, "error": err}).Error("Error fetching remote version")
					remoteVersion = "Error fetching version"
				}
			}

			upToDate := latestVersion == remoteVersion

			mu.Lock()
			toolData = append(toolData, ToolData{
				Name:          t.Name,
				LatestVersion: latestVersion,
				RemoteVersion: remoteVersion,
				UpToDate:      upToDate,
			})
			mu.Unlock()
			log.WithField("tool", t.Name).Info("Processed data for tool")
		}(tool)
	}

	wg.Wait()
	log.Info("All tool data fetched.")
	return toolData
}
