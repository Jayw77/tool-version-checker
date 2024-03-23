package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config struct for the top level yaml file
type Config struct {
	Endpoints     []*Endpoint   `yaml:"endpoints"`
	Kubernetes    *Kubernetes   `yaml:"kubernetes"`
	FetchInterval time.Duration `yaml:"fetchInterval"`
}

type Endpoint struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Url     string `yaml:"url"`
	Image   string `yaml:"image"`
	Custom  Custom `yaml:"custom"`
	Version Version
}

type Custom struct {
	MyVersion     EndpointConfig `yaml:"myVersion"`
	LatestVersion EndpointConfig `yaml:"latestVersion"`
}

type Kubernetes struct {
	Clusters     []*KubernetesCluster `yaml:"clusters"`
	CustomImages []*CustomImage       `yaml:"customImages"`
}

type KubernetesCluster struct {
	Name                string `yaml:"name"`
	Host                string
	ClientCertificate   string
	ClientKey           string
	ClientCACertificate string
	KubeConfig          string
	KubeConfigPath      string `yaml:"kubeConfig"`
}

type CustomImage struct {
	Name          string         `yaml:"name"`
	Image         string         `yaml:"image"`
	LatestVersion EndpointConfig `yaml:"latestVersion"`
}

type EndpointConfig struct {
	Endpoint string `yaml:"endpoint"`
	JsonKey  string `yaml:"jsonKey"`
}

type Version struct {
	Current  string
	Latest   string
	UpToDate bool
}

func loadConfig() {
	var yamlFile []byte
	var err error

	if fileExists("config/config.yaml") {
		yamlFile, err = os.ReadFile("config/config.yaml")
		if err != nil {
			log.WithField("error", err).Error("Error reading YAML file from config directory")
			return
		}
		log.Info("Using config/config.yaml")
	} else if fileExists("default_config.yaml") {
		yamlFile, err = os.ReadFile("default_config.yaml")
		if err != nil {
			log.WithField("error", err).Error("Error reading YAML file from default config")
			return
		}
		log.Info("Using default_config.yaml")
	} else {
		log.Error("No configuration file found")
		return
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.WithField("error", err).Error("Error unmarshalling YAML")
		return
	}

	log.Info("Fetch interval: " + config.FetchInterval.String())
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
