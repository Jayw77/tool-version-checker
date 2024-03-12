package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func fetchAll() {
	var currentVersionEndpoint, latestVersionEndpoint EndpointConfig

	for {
		for _, e := range config.Endpoints {
			// get endpoint configs
			if e.Type == "custom" {
				currentVersionEndpoint = e.Custom.MyVersion
				latestVersionEndpoint = e.Custom.LatestVersion
			} else {
				currentVersionEndpoint = CurrentVersionEndpoints[e.Type]
				latestVersionEndpoint = LatestVersionEndpoints[e.Type]
			}
			currentVersionEndpoint.Endpoint = e.Url + currentVersionEndpoint.Endpoint

			// call endpoints to get versions
			e.Version.Current, _ = fetchVersion(currentVersionEndpoint.Endpoint, currentVersionEndpoint.JsonKey)
			e.Version.Latest, _ = fetchVersion(latestVersionEndpoint.Endpoint, latestVersionEndpoint.JsonKey)
			e.Version.UpToDate = e.Version.Current == e.Version.Latest

			log.WithFields(logrus.Fields{"endpoint": e.Url, "type": e.Type, "currentVersion": e.Version.Current, "latestVersion": e.Version.Latest, "UpToDate": e.Version.UpToDate}).Info("Version information collected")
		}
		time.Sleep(config.FetchInterval)
	}
}

func fetchVersion(endpoint string, jsonKey string) (string, error) {
	log.WithField("endpoint", endpoint).Info("Fetching version")

	var jsonData map[string]interface{}

	if GetCache(endpoint) != nil {
		// use cached data
		jsonData = GetCache(endpoint)
	} else {
		// make call to endpoint
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(endpoint)
		if err != nil {
			log.WithFields(logrus.Fields{"endpoint": endpoint, "error": err}).Error("Error fetching from endpoint")
			return "", err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&jsonData)
		if err != nil {
			log.WithFields(logrus.Fields{"endpoint": endpoint, "error": err}).Error("Error decoding JSON")
			return "", err
		}

		// set the cache
		SetCache(endpoint, jsonData)

		if resp.StatusCode == http.StatusBadRequest {
			log.WithFields(logrus.Fields{"endpoint": endpoint, "statusCode": resp.StatusCode}).Error("Bad request")
			return "", fmt.Errorf("bad request: %v", resp.Status)
		}
	}

	// Split the jsonKey by "." to support nested keys
	keys := strings.Split(jsonKey, ".")

	var current interface{} = jsonData
	for _, key := range keys {
		if currMap, ok := current.(map[string]interface{}); ok {
			current, ok = currMap[key]
			if !ok {
				log.WithFields(logrus.Fields{"endpoint": endpoint, "key": key}).Error("Key not found in JSON")
				return "", fmt.Errorf("key %s not found in JSON", key)
			}
		} else {
			log.WithFields(logrus.Fields{"endpoint": endpoint, "key": key}).Error("Unexpected JSON structure")
			return "", fmt.Errorf("unexpected JSON structure for key %s", key)
		}
	}

	version, ok := current.(string)
	if !ok {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "key": jsonKey}).Error("Unexpected type for key")
		return "", fmt.Errorf("unexpected type for key %s", jsonKey)
	}
	return version, nil
}
