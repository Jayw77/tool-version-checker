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
}

// HomePageData represents the data for the home page
type HomePageData struct {
	Tools []ToolData
}

var currentToolData []ToolData // Global variable to store the latest data

func fetchToolDataPeriodically(config Config) {
	log.Info("Starting periodic data fetch...")
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Initial fetch
	currentToolData = fetchToolData(config)

	for {
		select {
		case <-ticker.C:
			log.Info("Fetching tool data...")
			currentToolData = fetchToolData(config)
			log.Info("Tool data fetched successfully.")
		}
	}
}

func fetchToolData(config Config) []ToolData {
	log.Info("Fetching tool data for all tools...")
	var toolData []ToolData
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect toolData slice

	for _, tool := range config.Tools {
		wg.Add(1)
		go func(t Tool) {
			defer wg.Done()
			latestVersion, err := fetchVersion(t.LatestVersionEndpoint, t.LatestVersionJSONKey)
			if err != nil {
				log.WithFields(logrus.Fields{"tool": t.Name, "error": err}).Error("Error fetching latest version")
				latestVersion = "Error fetching version"
			}

			remoteVersion, err := fetchVersion(t.RemoteVersionEndpoint, t.RemoteVersionJSONKey)
			if err != nil {
				log.WithFields(logrus.Fields{"tool": t.Name, "error": err}).Error("Error fetching remote version")
				remoteVersion = "Error fetching version"
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
