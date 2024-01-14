package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Mock for fetchVersion function
var mockFetchVersion = func(endpoint, jsonVersionKey string) (string, error) {
	// Mock behavior based on the endpoint or jsonVersionKey
	// For example, return a dummy version string and nil error
	return "v1.0.0", nil
}

func TestFetchToolData(t *testing.T) {
	// Setup
	originalFetchVersion := fetchVersion
	fetchVersion = mockFetchVersion
	defer func() { fetchVersion = originalFetchVersion }() // Reset after test

	config := Config{
		Tools: []Tool{
			{
				Name:                  "ExampleTool",
				LatestVersionEndpoint: "http://example.com/latest",
				RemoteVersionEndpoint: "http://example.com/remote",
				JSONVersionKey:        "version",
			},
			// Add more tools as needed
		},
	}

	// Execution
	toolData := fetchToolData(config)

	// Assertions
	assert.NotEmpty(t, toolData, "ToolData should not be empty")
	for _, data := range toolData {
		assert.Equal(t, "1.0.0", data.LatestVersion, "LatestVersion should be 1.0.0")
		assert.Equal(t, "1.0.0", data.RemoteVersion, "RemoteVersion should be 1.0.0")
		assert.True(t, data.UpToDate, "UpToDate should be true")
		// Add more assertions based on your requirements
	}
}
