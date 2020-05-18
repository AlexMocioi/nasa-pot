package main

import (
	"log"
	"nasa-pot/src/api"
	"nasa-pot/src/logging"
	"nasa-pot/src/service"
	"os"
)

func init() {
	log.Println("Init main...")

	// Initialize the web server.
	err := api.Init()
	if err != nil {
		log.Fatal("failed to initialize web server", err)
		er := logging.Slack("failed to initialize web server", err.Error())
		if er != nil {
			log.Fatal("failed to post to Slack", er)
		}
		os.Exit(1)
	}
}

func main() {
	log.Println("Starting in parallel cache populator...")
	go service.CachePopulator()

	// Block forevermore.
	select {}
}
