package main

import "os"

// Tool struct used inside the tools list for each tool
type Tool struct {
	Name                  string  `yaml:"name"`
	LatestVersionEndpoint string  `yaml:"latestVersionEndpoint"`
	MyVersionEndpoint     string  `yaml:"myVersionEndpoint"`
	LatestVersionJSONKey  string  `yaml:"latestVersionJSONKey"`
	MyVersionJSONKey      string  `yaml:"myVersionJSONKey"`
	CurrentVersion        *string `yaml:"currentVersion"` // using string pointer allows us to differentiate between null & ""
	Comment               string  `yaml:"comment"`
}

// Config struct for the top level yaml file
type Config struct {
	Tools         []Tool `yaml:"tools"`
	FetchInterval int    `yaml:"fetchInterval"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
