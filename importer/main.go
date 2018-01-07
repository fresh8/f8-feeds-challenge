package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Roverr/f8-feeds-challenge/importer/config"
)

// Very minimal error logging
func logError(customMsg string, err error) {
	fmt.Printf("ERROR || %s || %s\n", customMsg, err.Error())
}
func logInfo(customMsg string) {
	fmt.Printf("INFO || %s", customMsg)
}

func main() {
	// Initialize configuration and API
	settings, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	api := NewAPI(settings.GetFeedURL(), settings.GetStoreURL(), Parser{})

	// Get all events
	events, err := api.GetEvents()
	if err != nil {
		log.Fatal(err)
	}
	if len(events) == 0 {
		fmt.Println("Events array was parsed as empty, exiting with 0")
		return
	}

	// Start requesting events
	var wg sync.WaitGroup
	wg.Add(len(events))
	var eventCollection []Event
	for _, id := range events {
		go func(eventID int) {
			defer wg.Done()
			logInfo(fmt.Sprintf("Requesting %d event\n", eventID))
			event, err := api.GetEventByID(eventID)
			if err != nil {
				logError(fmt.Sprintf("Event error with id %d, might not existed", eventID), err)
				return
			}
			eventCollection = append(eventCollection, *event)
		}(id)
	}
	wg.Wait()
}
