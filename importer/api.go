package main

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
	return nil, nil
}

// GetEvents is for making a request for event IDs
func (a API) GetEvents() ([]int, error) {
	return nil, nil
}

// GetEventByID returns an Event for a given ID
func (a API) GetEventByID(id int) (*Event, error) {
	return nil, nil
}

// GetMarketByID returns a market for a given ID
func (a API) GetMarketByID(id int) (*Market, error) {
	return nil, nil
}

// PostToStore is posting data to the store
func (a API) PostToStore(event PopulatedEvent) error {
	return nil
}
