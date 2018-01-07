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

// mapEventMarkets creates a new map where every market has an array of events
// result will be { "market": ["eventID1", "eventID2"]}
func mapEventMarkets(events []Event) map[int][]string {
	match := map[int][]string{}
	for _, event := range events {
		for _, market := range event.Markets {
			_, ok := match[market]
			if ok {
				match[market] = append(match[market], event.ID)
				continue
			}
			match[market] = []string{event.ID}
		}
	}
	return match
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
	var mutex = &sync.Mutex{}
	for _, id := range events {
		go func(eventID int) {
			defer wg.Done()
			logInfo(fmt.Sprintf("Requesting %d event\n", eventID))
			event, err := api.GetEventByID(eventID)
			if err != nil {
				logError(fmt.Sprintf("Event error with id %d, might not existed", eventID), err)
				return
			}
			mutex.Lock()
			eventCollection = append(eventCollection, *event)
			mutex.Unlock()
		}(id)
	}
	wg.Wait()

	// Request all unique markets
	marketMap := mapEventMarkets(eventCollection)
	markets := make(chan *Market, len(marketMap))
	defer close(markets)
	for market := range marketMap {
		go func(id int) {
			logInfo(fmt.Sprintf("Requesting %d market\n", id))
			market, err := api.GetMarketByID(id)
			markets <- market
			if err != nil {
				logError(fmt.Sprintf("Market error with id %d", id), err)
				return
			}
		}(market)
	}
}
