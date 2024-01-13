package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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

	startWebServer()
}
