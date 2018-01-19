package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HomeHandler handle GET request to the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Welcome to the WeKnow API")
}

func main() {
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":3000", r))
}
