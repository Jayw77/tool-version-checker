package main

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Tool represents a tool to check versions for
type Tool struct {
	Name                  string `yaml:"name"`
	LatestVersionEndpoint string `yaml:"latestVersionEndpoint"`
	RemoteVersionEndpoint string `yaml:"remoteVersionEndpoint"`
	JSONVersionKey        string `yaml:"jsonVersionKey"`
}

// Config represents the YAML configuration format
type Config struct {
	Tools []Tool `yaml:"tools"`
}

// ToolData contains the display information for each tool
type ToolData struct {
	Name          string
	LatestVersion string
	RemoteVersion string
	UpToDate      bool
}

// Check if version is up to date by comparing versions from fetch-versions
func fetchToolData(config Config) []ToolData {
	log.Info("Fetching tool data for all tools...")
	var toolData []ToolData
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect toolData slice

	for _, tool := range config.Tools {
		wg.Add(2)
		go func(t Tool) {
			defer wg.Done()
			latestVersion, err := fetchVersion(t.LatestVersionEndpoint, t.JSONVersionKey)
			if err != nil {
				log.WithFields(logrus.Fields{"tool": t.Name, "error": err}).Error("Error fetching latest version")
				latestVersion = "Error fetching version"
			}

			remoteVersion, err := fetchVersion(t.RemoteVersionEndpoint, t.JSONVersionKey)
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

// Periodically refresh the data to check if versions have changes
func fetchToolDataPeriodically(config Config) {
	log.Info("Starting periodic data fetch...")
	ticker := time.NewTicker(11 * time.Minute)
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
