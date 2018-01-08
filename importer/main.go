package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

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
			logInfo(fmt.Sprintf("Requesting event number %d\n", eventID))
			event, err := api.GetEventByID(eventID)
			if err != nil {
				logError(fmt.Sprintf("Event error with id %d", eventID), err)
				return
			}
			if event == nil {
				logInfo(fmt.Sprintf("Event number %d not exists\n", eventID))
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
			logInfo(fmt.Sprintf("Requesting market number %d\n", id))
			market, err := api.GetMarketByID(id)
			markets <- market
			if err != nil {
				logError(fmt.Sprintf("Market error with id %d", id), err)
				return
			}
			if market == nil {
				logInfo(fmt.Sprintf("Market number %d was not found\n", id))
			}
		}(market)
	}

	// Process markets and create populated events
	populatedMap := map[string]PopulatedEvent{}
	for i := 0; i < cap(markets); i++ {
		market := <-markets
		if market == nil {
			continue
		}

		// Get ID in integer form for map
		id, err := market.GetIntegerID()
		if err != nil {
			logError(fmt.Sprintf("Market ID error at %s", market.ID), err)
			return
		}
		eventIDs, ok := marketMap[id]
		if !ok {
			logError(fmt.Sprintf("Market number %s missed from market map", market.ID), nil)
		}
		for _, id := range eventIDs {
			// If event already occoured, append current market to it
			elem, ok := populatedMap[id]
			if ok {
				elem.Markets = append(elem.Markets, *market)
				populatedMap[id] = elem
				continue
			}
			event := SearchEvent(eventCollection, id)
			if event == nil {
				continue
			}
			elem = PopulatedEvent{
				Event:              *event,
				TimeRepresentation: event.Time.Format(time.RFC3339),
				Markets:            []Market{*market},
			}
			populatedMap[id] = elem
		}
	}

	// Check for unvalid events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, event := range events {
			_, ok := populatedMap[strconv.Itoa(event)]
			if !ok {
				logInfo(fmt.Sprintf("Event number %d was not valid\n", event))
			}
		}
	}()

	// Post correctly formatted data to the store
	wg.Add(len(populatedMap))
	for _, elem := range populatedMap {
		go func(event PopulatedEvent) {
			defer wg.Done()
			err := api.PostToStore(event)
			if err != nil {
				logError(fmt.Sprintf("Storing of event number %s failed", event.ID), err)
				return
			}
			logInfo(fmt.Sprintf("Event number %s stored\n", event.ID))
		}(elem)
	}
	wg.Wait()
}
