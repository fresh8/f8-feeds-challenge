package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type event struct {
	ID      json.Number   `json:"id"`
	Name    string        `json:"name"`
	Time    string        `json:"time"`
	Markets []interface{} `json:"markets"`
}

type market struct {
	ID      json.Number `json:"id"`
	Type    string      `json:"type"`
	Options []option    `json:"options"`
}

type option struct {
	ID   json.Number `json:"id"`
	Name string      `json:"name"`
	Odds string      `json:"odds,omitempty"`
	Num  int         `json:"num"`
	Den  int         `json:"den"`
}

func main() {
	storeAddress := os.Getenv("STORE_ADDR")

	if storeAddress == "" {
		storeAddress = "localhost:8001"
		log.Println("STORE_ADDR environemnt variable not found. Using localhost:8001 as the store address.")
	}

	// Http client with 10 second timeout
	var client = &http.Client{Timeout: 10 * time.Second}

	importEvents(client, storeAddress)
}

func importEvents(client *http.Client, storeAddress string) {
	events := getEventsList(client)

	// Go and get each of the events from the feed and store it if valid
	for k := range events {
		eventID := events[k]
		log.Printf("Importing event %d", eventID)
		r, err := client.Get(fmt.Sprintf("http://localhost:8000/football/events/%d", eventID))
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		var retreivedEvent event
		err = json.NewDecoder(r.Body).Decode(&retreivedEvent)
		if err != nil {
			log.Printf("Can't decode Event data. Skipping event %d.", eventID)
			continue
		}

		populateEventMarkets(client, &retreivedEvent)

		storeEvent(client, storeAddress, &retreivedEvent)
	}

	log.Printf("Event import finished.")
}

func getEventsList(client *http.Client) []int {
	r, err := client.Get("http://localhost:8000/football/events")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	events := make([]int, 0)
	err = json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		panic(err)
	}
	return events
}

func populateEventMarkets(client *http.Client, retreivedEvent *event) {
	for k := range retreivedEvent.Markets {
		marketRaw := retreivedEvent.Markets[k]
		var marketID = int(marketRaw.(float64))
		log.Printf("Importing market %d", marketID)
		r, err := client.Get(fmt.Sprintf("http://localhost:8000/football/markets/%d", marketID))
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		var retreivedMarket market
		err = json.NewDecoder(r.Body).Decode(&retreivedMarket)
		if err != nil {
			log.Printf("Can't decode Market data. Skipping market %d.", marketID)
			// As we failed to find and decode this market, lets remove the id
			retreivedEvent.Markets = append(retreivedEvent.Markets[:k], retreivedEvent.Markets[k+1:]...)
			continue
		}

		formatMarketOptionsOdds(&retreivedMarket)

		retreivedEvent.Markets[k] = retreivedMarket
	}
}

func formatMarketOptionsOdds(retreivedMarket *market) {
	for k := range retreivedMarket.Options {
		option := retreivedMarket.Options[k]
		oddsSlice := strings.Split(option.Odds, "/")
		retreivedMarket.Options[k].Num, _ = strconv.Atoi(oddsSlice[0])
		retreivedMarket.Options[k].Den, _ = strconv.Atoi(oddsSlice[1])
		// As we don't want to send the odds field, let's set it to it's nil value
		retreivedMarket.Options[k].Odds = ""
	}
}

func storeEvent(client *http.Client, storeAddress string, retreivedEvent *event) {
	log.Printf("%+v", retreivedEvent)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(retreivedEvent)
	r, err := client.Post("http://"+storeAddress+"/event", "application/json", b)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
}
