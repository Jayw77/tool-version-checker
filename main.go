package main

import (
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/yousuf64/shift"
)

var log = logrus.New()
var config Config

func main() {
	// set json formatter
	log.SetFormatter(&logrus.JSONFormatter{})

	// load config
	loadConfig()

	// begin checking versions
	go fetchAll()

	// tmpl := template.Must(template.ParseFiles("home.html"))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	err := tmpl.Execute(w, HomePageData{Tools: currentToolData})
	// 	if err != nil {
	// 		log.WithField("error", err).Error("Error executing template")
	// 	}
	// })

	router := shift.New()
	router.GET("/", shift.HTTPHandlerFunc(handlerHome))
	router.GET("/health/live", shift.HTTPHandlerFunc(handlerHealth))
	router.GET("/health/ready", shift.HTTPHandlerFunc(handlerHealth))
	router.GET("/styles.css", shift.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "assets/styles.css") }))

	var wg sync.WaitGroup
	wg.Add(1)
	go http.ListenAndServe(":8080", router.Serve())
	log.Info("Server started on port 8080...")
	wg.Wait()
}
