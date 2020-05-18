package api

import (
	"encoding/json"
	"fmt"
	"log"
	"nasa-pot/src/config"
	"nasa-pot/src/service"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome Asteroids Fan! Use /grabLatest/:startDate to find real-time data\n")
}

func GrabLatestsEndpoint(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	startDate, err := time.Parse("2006-01-02", ps.ByName("startDate"))
	if err != nil {
		log.Fatal(err)
	}
	endDate := startDate

	writeJsonContent(w, service.GrabLatest(startDate, endDate))
}

func GrabLatestsEndpointAsync(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	startDate, err := time.Parse("2006-01-02", ps.ByName("startDate"))
	if err != nil {
		log.Fatal(err)
	}
	endDate := startDate

	writeJsonContent(w, service.GrabLatestAsync(startDate, endDate))
}

func writeJsonContent(w http.ResponseWriter, jsonStruct interface{}) {
	js, err := json.Marshal(jsonStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Initialize web server.
func Init() error {

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/grabLatest/:startDate/sync", GrabLatestsEndpoint)
	router.GET("/grabLatest/:startDate/async", GrabLatestsEndpointAsync)

	endpoint := config.Configuration.Api.Endpoint
	log.Println("Starting web demo on endpoint ", endpoint, "...")
	go func() {
		log.Fatal(http.ListenAndServe(endpoint, router))
	}()
	return nil
}
