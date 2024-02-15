package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})
	var yamlFile []byte
	var err error

	if fileExists("config/config.yaml") {
		yamlFile, err = os.ReadFile("config/config.yaml")
		if err != nil {
			log.WithField("error", err).Error("Error reading YAML file from config directory")
			return
		}
		fmt.Println("Using config/config.yaml")
	} else if fileExists("default_config.yaml") {
		yamlFile, err = os.ReadFile("default_config.yaml")
		if err != nil {
			log.WithField("error", err).Error("Error reading YAML file from default config")
			return
		}
		fmt.Println("Using default_config.yaml")
	} else {
		log.Error("No configuration file found")
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
