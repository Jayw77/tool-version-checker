package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"net/http"
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

// HomePageData represents the data for the home page
type HomePageData struct {
	Tools []ToolData
}

var log = logrus.New()
var currentToolData []ToolData // Global variable to store the latest data

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.WithField("error", err).Error("Error reading YAML file")
		return
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.WithField("error", err).Error("Error unmarshalling YAML")
		return
	}

	go fetchToolDataPeriodically(config)

	tmpl := template.Must(template.ParseFiles("home.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, HomePageData{Tools: currentToolData})
		if err != nil {
			log.WithField("error", err).Error("Error executing template")
		}
	})

	log.Info("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.WithField("error", err).Error("Failed to start server")
	}
}

func fetchVersion(endpoint string, jsonKey string) (string, error) {
	log.WithField("endpoint", endpoint).Info("Fetching version")

	resp, err := http.Get(endpoint)
	if err != nil {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "error": err}).Error("Error fetching from endpoint")
		return "", err
	}
	defer resp.Body.Close()

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "error": err}).Error("Error decoding JSON")
		return "", err
	}

	if val, ok := jsonData[jsonKey]; ok {
		version, ok := val.(string)
		if !ok {
			log.WithFields(logrus.Fields{"endpoint": endpoint, "key": jsonKey}).Error("Unexpected type for key")
			return "", fmt.Errorf("unexpected type for key %s", jsonKey)
		}
		return version, nil
	} else {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "key": jsonKey}).Error("Key not found in JSON")
		return "", fmt.Errorf("key %s not found", jsonKey)
	}
}

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
