package main

import (
	"reflect"
	"testing"
	"time"
)

func TestGetCache(t *testing.T) {
	// Test case 1: Cache does not exist for the URL
	result := GetCache("nonexistenturl")
	if result != nil {
		t.Errorf("Expected empty string for nonexistent URL, got %s", result)
	}

	// Test case 2: Cache exists but is expired
	url := "expiredurl"
	testData := map[string]interface{}{"test": "data"}
	SetCache(url, testData)
	cache[url].CachedAt = time.Now().Add(-time.Hour) // Set CachedAt to a time in the past
	result = GetCache(url)
	if result != nil {
		t.Errorf("Expected empty string for expired cache, got %s", result)
	}

	// Test case 3: Cache exists and is not expired
	url = "validurl"
	SetCache(url, testData)
	result = GetCache(url)
	if result == nil {
		t.Error("Expected non-empty string for valid cache, got empty string")
	}
}

func TestSetCache(t *testing.T) {
	url := "testurl"
	testData := map[string]interface{}{"test": "data"}

	// Verify cache is initially empty
	if GetCache(url) != nil {
		t.Error("Expected cache to be empty before setting")
	}

	// Set cache and verify it is not empty
	SetCache(url, testData)
	result := GetCache(url)
	if !reflect.DeepEqual(result, testData) {
		t.Errorf("Expected cache data %s, got %s", testData, result)
	}
}
