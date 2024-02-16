package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	if resp.StatusCode == http.StatusBadRequest {
		log.WithFields(logrus.Fields{"endpoint": endpoint, "statusCode": resp.StatusCode}).Error("Bad request")
		return "", fmt.Errorf("bad request: %v", resp.Status)
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
