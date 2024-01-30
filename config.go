package main

// Tool represents a tool to check versions for
type Tool struct {
	Name                  string `yaml:"name"`
	LatestVersionEndpoint string `yaml:"latestVersionEndpoint"`
	RemoteVersionEndpoint string `yaml:"remoteVersionEndpoint"`
	LatestVersionJSONKey  string `yaml:"latestVersionJSONKey"`
	RemoteVersionJSONKey  string `yaml:"remoteVersionJSONKey"`
}

// Config represents the YAML configuration format
type Config struct {
	Tools []Tool `yaml:"tools"`
}
