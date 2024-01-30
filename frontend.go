package main

import (
	"html/template"
	"net/http"
)

// HomePageData represents the data for the home page
type HomePageData struct {
	Tools []ToolData
}

// StartServer Initialize and start the HTTP server
func startWebServer() {
	tmpl := template.Must(template.ParseFiles("home.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, HomePageData{Tools: currentToolData})
		if err != nil {
			log.WithField("error", err).Error("Error executing template")
		}
	})

	log.Info("Server is starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.WithField("error", err).Error("Failed to start server")
	}
}
