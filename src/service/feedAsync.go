package service

import (
	"encoding/json"
	"log"
	"nasa-pot/src/service/asteroid"
	"time"
)

// Gets data from in memory big cache if available
func GrabLatestAsync(start, end time.Time) []asteroid.Specs {

	responseData := make([]asteroid.Specs, 0)
	dateCursor := start

	for dateCursor.Before(end) || dateCursor.Equal(end) {
		events := getEventsFromCacheFor(dateCursor)
		log.Println("Read from cache")
		for _, event := range events {
			responseData = append(responseData, asteroid.NewSpecs(event.(map[string]interface{})))
		}
		dateCursor = dateCursor.AddDate(0, 0, 1)
	}
	return responseData
}

func getEventsFromCacheFor(date time.Time) []interface{} {
	events := make([]interface{}, 0)
	cachedBytes, _ := cache.Get(date.Format("2006-01-02"))
	json.Unmarshal(cachedBytes, &events)
	return events
}
