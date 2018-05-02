package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.NewRoute().
		Name("EventRootHandler").
		Methods("GET").
		Path("/football/events").
		HandlerFunc(EventRootHandler)

	r.NewRoute().
		Name("EventHandler").
		Methods("GET").
		Path("/football/events/{id:[0-9]+}").
		HandlerFunc(EventHandler)

	r.NewRoute().
		Name("MarketHandler").
		Methods("GET").
		Path("/football/markets/{id:[0-9]+}").
		HandlerFunc(MarketHandler)

	log.Println("server running on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
