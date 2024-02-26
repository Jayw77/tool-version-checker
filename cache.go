package main

import (
	"time"
)

type Cache struct {
	CachedAt time.Time
	Data     map[string]interface{}
}

var cache = make(map[string]*Cache)           // key is url
var cacheTtl time.Duration = time.Minute * 10 // 10 minutes cache expiry time

func (cache *Cache) Expired() bool {
	return time.Now().After(cache.CachedAt.Add(cacheTtl))
}

func GetCache(url string) map[string]interface{} {
	if cache[url] != nil {
		// cache exists
		if !cache[url].Expired() {
			// cache not expired
			log.WithField("url", url).Debug("Cache used")
			return cache[url].Data
		}
	}

	return nil
}

func SetCache(url string, data map[string]interface{}) {
	log.WithField("url", url).Debug("Cache set")
	cache[url] = &Cache{
		CachedAt: time.Now(),
		Data:     data,
	}
}
