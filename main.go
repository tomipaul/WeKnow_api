package main

import (
	ctrl "WeKnow_api/controller"
	mwr "WeKnow_api/middlewares"
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
	pr := r.NewRoute().Subrouter()

	// Routes consist of a path and a handler function.
	// Log all requests to the application
	r.Use(mwr.LoggingHandler)

	// Handle GET request to the main endpoint
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// Handle authentication requests
	authSubRouter := r.PathPrefix("/api/v1/user").Subrouter()
	authSubRouter.
		HandleFunc("/signup", ctrl.UserSignUpEndPoint).
		Methods("POST")
	authSubRouter.
		HandleFunc("/signin", ctrl.UserSignInEndPoint).
		Methods("POST")

	// Protect data endpoints
	pr.Use(mwr.ValidateEndpoint)

	// Handle connection requests
	connectionSubRouter := pr.PathPrefix("/api/v1/connection").Subrouter()
	connectionSubRouter.
		HandleFunc("", ctrl.ConnectUser).
		Methods("POST")
	connectionSubRouter.
		HandleFunc("/favorites", ctrl.GetAllFavorites).
		Methods("GET")
	connectionSubRouter.
		HandleFunc("/followers", ctrl.GetAllFollowers).
		Methods("GET")

	// Handle collection requests
	collectionSubRouter := pr.PathPrefix("/api/v1/collection").Subrouter()
	collectionSubRouter.
		HandleFunc("", ctrl.CreateCollectionEndPoint).
		Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":3000", r))
}
