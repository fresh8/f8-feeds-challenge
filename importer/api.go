package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Feed describes the endpoints of the feed server
type feed struct {
	root    string
	events  string
	markets string
}

// Store describes the endpoint of the store
type store struct {
	root  string
	event string
}

// API describes the API used by the application
type API struct {
	feed   feed
	store  store
	parser Parser
}

// NewAPI creates a new instance of the application's API
func NewAPI(feedHost, storeHost string, parser Parser) API {
	return API{
		feed: feed{
			root:    feedHost,
			events:  feedHost + `/football/events`,
			markets: feedHost + `/football/markets`,
		},
		store: store{
			root:  storeHost,
			event: storeHost + `/event`,
		},
		parser: parser,
	}
}

// GetRequest makes a simple get request for a given url and returns the body
func (a API) GetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

// GetEvents is for making a request for event IDs
func (a API) GetEvents() ([]int, error) {
	body, err := a.GetRequest(a.feed.events)
	if err != nil {
		return nil, err
	}
	return a.parser.Events(body)
}

// GetEventByID returns an Event for a given ID
func (a API) GetEventByID(id int) (*Event, error) {
	url := fmt.Sprintf("%s/%d", a.feed.events, id)
	body, err := a.GetRequest(url)
	if err != nil {
		return nil, err
	}
	return a.parser.Event(body)
}

// GetMarketByID returns a market for a given ID
func (a API) GetMarketByID(id int) (*Market, error) {
	return nil, nil
}

// PostToStore is posting data to the store
func (a API) PostToStore(event PopulatedEvent) error {
	return nil
}
