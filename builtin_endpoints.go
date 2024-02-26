package main

var CurrentVersionEndpoints = map[string]EndpointConfig{
	"grafana": {
		Endpoint: "/api/health",
		JsonKey:  "version",
	},
	"prometheus": {
		Endpoint: "/api/v1/status/buildinfo",
		JsonKey:  "data.version",
	},
	"alertmanager": {
		Endpoint: "/api/v2/status",
		JsonKey:  "versionInfo.version",
	},
	"loki": {
		Endpoint: "/loki/api/v1/status/buildinfo",
		JsonKey:  "version",
	},
	"traefik": {
		Endpoint: "/api/version",
		JsonKey:  "version", // needs checking
	},
}

var LatestVersionEndpoints = map[string]EndpointConfig{
	"grafana": {
		Endpoint: "https://api.github.com/repos/grafana/grafana/releases/latest",
		JsonKey:  "tag_name",
	},
	"prometheus": {
		Endpoint: "https://api.github.com/repos/prometheus/prometheus/releases/latest",
		JsonKey:  "tag_name",
	},
	"alertmanager": {
		Endpoint: "https://api.github.com/repos/prometheus/alertmanager/releases/latest",
		JsonKey:  "tag_name",
	},
	"loki": {
		Endpoint: "https://api.github.com/repos/grafana/loki/releases/latest",
		JsonKey:  "tag_name",
	},
	"traefik": {
		Endpoint: "https://api.github.com/repos/traefik/traefik/releases/latest",
		JsonKey:  "tag_name",
	},
}
