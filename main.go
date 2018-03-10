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

	// Routes consist of a path and a handler function.
	r.Use(mwr.LoggingHandler)
	r.HandleFunc("/", HomeHandler).Methods("GET")

	authSubRouter := r.PathPrefix("/api/v1/user").Subrouter()
	authSubRouter.
		HandleFunc("/signup", ctrl.UserSignUpEndPoint).
		Methods("POST")
	authSubRouter.
		HandleFunc("/signin", ctrl.UserSignInEndPoint).
		Methods("POST")

	connectionSubRouter := r.NewRoute().Subrouter()
	connectionSubRouter.Use(mwr.ValidateEndpoint)
	connectionSubRouter.
		HandleFunc("/api/v1/connection", ctrl.ConnectUser).
		Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":3000", r))
}
