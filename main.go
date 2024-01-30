package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.New()

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
