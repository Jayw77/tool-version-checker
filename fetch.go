package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func fetchVersion(endpoint string, jsonKey string) (string, error) {
	log.WithField("endpoint", endpoint).Info("Fetching version")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "error": err}).Error("Error fetching from endpoint")
		return "", err
	}
	defer resp.Body.Close()

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "error": err}).Error("Error decoding JSON")
		return "", err
	}

	// Check if the status code is http.StatusBadRequest
	if resp.StatusCode == http.StatusBadRequest {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "statusCode": resp.StatusCode}).Error("Bad request")
		return "", fmt.Errorf("bad request: %v", resp.Status)
	}

	if val, ok := jsonData[jsonKey]; ok {
		version, ok := val.(string)
		if !ok {
			log.WithFields(logrus.Fields{"endpoint": endpoint, "key": jsonKey}).Error("Unexpected type for key")
			return "", fmt.Errorf("unexpected type for key %s", jsonKey)
		}
		return version, nil
	} else {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "key": jsonKey}).Error("Key not found in JSON")
		return "", fmt.Errorf("key %s not found", jsonKey)
	}
}
