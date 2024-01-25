package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.New()
var currentToolData []ToolData // Global variable to store the latest data

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	yamlFile, err := os.ReadFile("config.yaml")
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

	startWebServer()
}
