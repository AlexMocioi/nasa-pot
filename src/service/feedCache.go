package service

import (
	"encoding/json"
	"github.com/allegro/bigcache"
	"log"
	"time"
)

var (
	cache     *bigcache.BigCache
	syncCache chan map[string]interface{}
)

func init() {
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(60 * time.Minute))
	syncCache = make(chan map[string]interface{}, 10)
}

// Listens on a chan for new events to cache data
func CachePopulator() {
	for {
		select {
		case cacheObj := <-syncCache:
			log.Println("Got new data on chan for caching...")
			date := cacheObj["date"].(string)
			objectsBytes, _ := json.Marshal(cacheObj["nearEarthObjects"])
			cache.Set(date, objectsBytes)
		}
	}
}
