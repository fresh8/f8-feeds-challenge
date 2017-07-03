package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// EventRootHandler returns a hardcoded list of ID's that should exist
// NOTE: They don't (intentionally)
func EventRootHandler(w http.ResponseWriter, r *http.Request) {
	eventIds := "[1, 2, 3, 4, 5]"
	w.Write([]byte(eventIds))
}

// EventHandler will take an ID and return the event, this method does no data
// validation at all and will pass back the data as added to the file. This
// enables passing non-ideal data
func EventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	data, err := ioutil.ReadFile(fmt.Sprintf("feed/files/events/%s.json", id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// MarketHandler will take an ID and return the market, much the same as the
// EventHandler does
func MarketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	data, err := ioutil.ReadFile(fmt.Sprintf("feed/files/markets/%s.json", id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
