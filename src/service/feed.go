package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"nasa-pot/src/config"
	"nasa-pot/src/logging"
	"nasa-pot/src/service/asteroid"
	"net/http"
	"net/url"
	"time"
)

var (
	endpoint   string
	apiKey     string
	httpClient http.Client
)

func init() {
	endpoint = config.Configuration.Services.FeedEndpoint
	log.Println("Feed Endpoint: ", endpoint)

	apiKey = config.Configuration.Services.ApiKey
	log.Println("API KEY: ", apiKey)

	httpClient = http.Client{
		Timeout: time.Second * 5,
	}
}

// Gets data from Nasa Service, pushes it on a chan to be cached, extracts a summary and returns to client
func GrabLatest(start, end time.Time) []asteroid.Specs {

	endpointWithQuery := formatUrl(endpoint, start, end, apiKey)
	body := doHttpRequest(endpointWithQuery)

	return parseFeedApiResponseAndGetSummary(body, start, end)
}

func parseFeedApiResponseAndGetSummary(body []byte, start, end time.Time) []asteroid.Specs {
	jsonResp := asteroid.ApiResponse{}
	jsonErr := json.Unmarshal(body, &jsonResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	responseData := make([]asteroid.Specs, 0)
	if jsonResp.ElementCount == 0 {
		return responseData
	}

	dateCursor := start
	for dateCursor.Before(end) || dateCursor.Equal(end) {
		events := jsonResp.NearEarthObjects[dateCursor.Format("2006-01-02")]

		sendToCachePopulator(dateCursor, events)

		for _, event := range events {
			responseData = append(responseData, asteroid.NewSpecs(event.(map[string]interface{})))

			if event.(map[string]interface{})["is_potentially_hazardous_asteroid"].(bool) {
				logging.Slack("Attention humans, asteroid getting hazardous closely to earth: ", event)
			}
		}
		dateCursor = dateCursor.AddDate(0, 0, 1)
	}

	return responseData
}

func sendToCachePopulator(cursor time.Time, events []interface{}) {
	// send to cache populator to be accessible async next time
	cacheObj := make(map[string]interface{})
	cacheObj["date"] = cursor.Format("2006-01-02")
	cacheObj["nearEarthObjects"] = events
	syncCache <- cacheObj
	log.Println("Wrote new data on chan to be cached")
}

func doHttpRequest(endpointWithQuery string) []byte {
	req, err := http.NewRequest(http.MethodGet, endpointWithQuery, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return body
}

func formatUrl(e string, start time.Time, end time.Time, key string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal("failed to parse Feed Endpoint URL", err)
	}
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("api_key", apiKey)
	q.Add("start_date", start.Format("2006-01-02"))
	q.Add("end_date", end.Format("2006-01-02"))
	u.RawQuery = q.Encode()
	return u.String()
}
