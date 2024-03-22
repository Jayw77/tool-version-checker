package main

var CurrentVersionEndpoints = map[string]EndpointConfig{
	"alertmanager": {
		Endpoint: "/api/v2/status",
		JsonKey:  "versionInfo.version",
	},
	"grafana": {
		Endpoint: "/api/health",
		JsonKey:  "version",
	},
	"loki": {
		Endpoint: "/loki/api/v1/status/buildinfo",
		JsonKey:  "version",
	},
	"prometheus": {
		Endpoint: "/version",
		JsonKey:  "version",
	},
	"sonarqube": {
		Endpoint: "/api/system/status",
		JsonKey:  "version",
	},
	"traefik": {
		Endpoint: "/api/version",
		JsonKey:  "version", // needs checking
	},
}

// we automatically recognise some images and map them to LatestVersionEndpoints
var ContainerImageLatestVersionEndpointNames = map[string]string{
	"quay.io/prometheus/alertmanager":                        "alertmanager",
	"quay.io/jetstack/cert-manager-cainjector":               "cert-manager",
	"quay.io/jetstack/cert-manager-webhook":                  "cert-manager",
	"quay.io/jetstack/cert-manager-controller":               "cert-manager",
	"registry.k8s.io/descheduler/descheduler":                "descheduler",
	"docker.io/grafana/grafana":                              "grafana",
	"registry.k8s.io/kube-state-metrics/kube-state-metrics":  "kube-state-metrics",
	"ghcr.io/kube-vip/kube-vip":                              "kube-vip",
	"docker.io/grafana/loki":                                 "loki",
	"docker.io/grafana/loki-canary":                          "loki",
	"docker.io/grafana/promtail":                             "loki",
	"quay.io/prometheus/node-exporter":                       "node-exporter",
	"quay.io/prometheus/prometheus":                          "prometheus",
	"quay.io/prometheus-operator/prometheus-config-reloader": "prometheus-operator",
	"quay.io/prometheus-operator/prometheus-operator":        "prometheus-operator",
	"tailscale/tailscale":                                    "tailscale",
	"tailscale/k8s-operator":                                 "tailscale",
	"jayw77/version-checker":                                 "tool-version-checker",
	"linuxserver/transmission":                               "transmission",
	"koenkk/zigbee2mqtt":                                     "zigbee2mqtt",
}

var LatestVersionEndpoints = map[string]EndpointConfig{
	"alertmanager": {
		Endpoint: "https://api.github.com/repos/prometheus/alertmanager/releases/latest",
		JsonKey:  "tag_name",
	},
	"cert-manager": {
		Endpoint: "https://api.github.com/repos/cert-manager/cert-manager/releases/latest",
		JsonKey:  "tag_name",
	},
	"descheduler": {
		Endpoint: "https://api.github.com/repos/kubernetes-sigs/descheduler/releases/latest",
		JsonKey:  "tag_name",
	},
	"grafana": {
		Endpoint: "https://api.github.com/repos/grafana/grafana/releases/latest",
		JsonKey:  "tag_name",
	},
	"kubernetes": {
		Endpoint: "https://api.github.com/repos/kubernetes/kubernetes/releases/latest",
		JsonKey:  "tag_name",
	},
	"kube-state-metrics": {
		Endpoint: "https://api.github.com/repos/kubernetes/kube-state-metrics/releases/latest",
		JsonKey:  "tag_name",
	},
	"kube-vip": {
		Endpoint: "https://api.github.com/repos/kube-vip/kube-vip/releases/latest",
		JsonKey:  "tag_name",
	},
	"loki": {
		Endpoint: "https://api.github.com/repos/grafana/loki/releases/latest",
		JsonKey:  "tag_name",
	},
	"node-exporter": {
		Endpoint: "https://api.github.com/repos/prometheus/node_exporter/releases/latest",
		JsonKey:  "tag_name",
	},
	"prometheus": {
		Endpoint: "https://api.github.com/repos/prometheus/prometheus/releases/latest",
		JsonKey:  "tag_name",
	},
	"prometheus-operator": {
		Endpoint: "https://api.github.com/repos/prometheus-operator/prometheus-operator/releases/latest",
		JsonKey:  "tag_name",
	},
	"sonarqube": {
		Endpoint: "https://api.github.com/repos/sonarsource/sonarqube/releases/latest",
		JsonKey:  "tag_name",
	},
	"tailscale": {
		Endpoint: "https://api.github.com/repos/tailscale/tailscale/releases/latest",
		JsonKey:  "tag_name",
	},
	"tool-version-checker": {
		Endpoint: "https://api.github.com/repos/Jayw77/tool-version-checker/releases/latest",
		JsonKey:  "tag_name",
	},
	"traefik": {
		Endpoint: "https://api.github.com/repos/traefik/traefik/releases/latest",
		JsonKey:  "tag_name",
	},
	"transmission": {
		Endpoint: "https://api.github.com/repos/transmission/transmission/releases/latest",
		JsonKey:  "tag_name",
	},
	"zigbee2mqtt": {
		Endpoint: "https://api.github.com/repos/Koenkk/zigbee2mqtt/releases/latest",
		JsonKey:  "tag_name",
	},
}
