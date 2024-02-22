package main

import "os"

// Tool struct used inside the tools list for each tool
type LatestVersionEndpoints struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	JSONKey string `yaml:"JSONKey"`
}

type MyEndpoints struct {
	Name                      string  `yaml:"name"`
	Url                       string  `yaml:"url"`
	JSONKey                   string  `yaml:"JSONKey"`
	LatestVersionEndpointName *string `yaml:"latestVersionEndpointName"` // using string pointer allows us to differentiate between null & ""
	Comment                   string  `yaml:"comment"`
}

// Config struct for the top level yaml file
type Config struct {
	LatestVersionEndpoints []LatestVersionEndpoints `yaml:"latestVersionEndpoints"`
	MyEndpoints            []MyEndpoints            `yaml:"myEndpoints"`
	FetchInterval          int                      `yaml:"fetchInterval"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
