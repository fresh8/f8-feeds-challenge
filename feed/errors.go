package main

import (
	"log"
	"net/http"
)

// ServeNotFound writes not found error.
func ServeRecordNotFound(w http.ResponseWriter, err error) {
	log.Print(err)
	http.Error(w, "Record not found", http.StatusNotFound)
}

// ServeBadRequest writes bad request.
func ServeBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

// ServeStatusOkay writes Status okay.
func ServeStatusOkay(w http.ResponseWriter, err error) {
	log.Print(err)
	http.Error(w, "Okay", http.StatusOK)
}
