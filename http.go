package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
)

var funcMap = template.FuncMap{
	"sort": func(a interface{}) any {
		switch b := a.(type) {
		case []*Endpoint:
			endpoints := b
			sort.Slice(endpoints, func(i, j int) bool {
				return endpoints[i].Name < endpoints[j].Name
			})
			return endpoints
		case []KubernetesImage:
			images := b
			sort.Slice(images, func(i, j int) bool {
				return images[i].Name < images[j].Name
			})
			return images
		}
		return nil
	},
}

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

	tpl := template.Must(template.New("t").Funcs(funcMap).ParseFiles("assets/home.html"))
	err := tpl.ExecuteTemplate(w, "home.html", d)

	if err != nil {
		log.WithField("error", err).Error("Error executing homepage template")
	}
}
