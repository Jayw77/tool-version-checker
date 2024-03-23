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

	d := struct {
		Endpoints              []*Endpoint
		KubernetesImages       []KubernetesImage
		KubernetesClusterCount int
	}{
		Endpoints:              config.Endpoints,
		KubernetesImages:       KubernetesImages,
		KubernetesClusterCount: len(config.Kubernetes.Clusters),
	}

	tpl := template.Must(template.ParseFiles("assets/home.html"))
	err := tpl.Execute(w, d)

	if err != nil {
		log.WithField("error", err).Error("Error executing homepage template")
	}
}
