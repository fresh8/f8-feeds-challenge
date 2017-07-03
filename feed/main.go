package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/football/events", EventRootHandler)
	r.HandleFunc("/football/events/{id:[0-9]+}", EventHandler)
	r.HandleFunc("/football/markets/{id:[0-9]+}", MarketHandler)

	log.Println("server running on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
