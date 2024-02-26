package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	fmt.Fprintf(w, "ok")
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	r.Close = true

	tpl := template.Must(template.ParseFiles("assets/home.html"))
	err := tpl.Execute(w, config.Endpoints)

	if err != nil {
		log.WithField("error", err).Error("Error executing homepage template")
	}
}
