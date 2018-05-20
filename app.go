package main

import (
	"WeKnow_api/handler"
	"WeKnow_api/middleware"
	"WeKnow_api/model"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// CreateApp create a new instance of the app
func CreateApp(config map[string]string) App {
	router := mux.NewRouter()
	db := model.Connect(config)
	app := App{
		router,
		db,
	}
	app.declareRoutes()
	return app
}

// App type application
type App struct {
	Router *mux.Router
	Db     *pg.DB
}

// run start application
func (app App) run(address string) {
	// listen at address and handle with handler app.Router
	err := http.ListenAndServe(address, app.Router)
	if err != nil {
		panic(err)
	}
}

// declareRoutes declare application endpoints
func (app App) declareRoutes() {
	hr := &handler.Handler{Db: app.Db}
	mwr := &middleware.Middleware{Db: app.Db}

	r := app.Router
	pr := r.NewRoute().Subrouter()

	// Routes consist of a path and a handler function.
	// Middleware Log all requests to the application
	r.Use(mwr.LogRequest)

	// Middleware Check if request body is empty for POST and PUT requests
	r.Use(mwr.CheckRequestBody)

	// Handle GET request to the main endpoint
	r.HandleFunc("/", hr.HomeHandler).Methods("GET")

	// Handle authentication requests
	authSubRouter := r.PathPrefix("/api/v1/auth").Subrouter()
	authSubRouter.
		HandleFunc("/signup", hr.UserSignUpEndPoint).
		Methods("POST")
	authSubRouter.
		HandleFunc("/signin", hr.UserSignInEndPoint).
		Methods("POST")

	// Middleware Protect data endpoints
	pr.Use(mwr.AuthorizeRequest)

	// Handle connection requests
	connectionSubRouter := pr.PathPrefix("/api/v1/connection").Subrouter()
	connectionSubRouter.
		HandleFunc("", hr.ConnectUser).
		Methods("POST")
	connectionSubRouter.
		HandleFunc("/favorites", hr.GetAllFavorites).
		Methods("GET")
	connectionSubRouter.
		HandleFunc("/followers", hr.GetAllFollowers).
		Methods("GET")

	// Handle collection requests
	collectionSubRouter := pr.PathPrefix("/api/v1/collection").Subrouter()
	collectionSubRouter.
		HandleFunc("", hr.CreateCollectionEndPoint).
		Methods("POST")
	collectionSubRouter.
		HandleFunc("", hr.GetAllCollections).
		Methods("GET")

	userSubRouter := pr.PathPrefix("/api/v1/user").Subrouter()
	userSubRouter.
		HandleFunc("/profile", hr.UpdateProfile).Methods("PUT")
	userSubRouter.
		HandleFunc("/password/reset", hr.ResetPassword).Methods("PUT")

	// Handle resource requests
	resourceSubRouter := pr.PathPrefix("/api/v1/resource").Subrouter()
	resourceSubRouter.
		HandleFunc("/{resourceId:[0-9]+}", hr.DeleteResource).
		Methods("DELETE")

	resourceTagsSubRouter := resourceSubRouter.NewRoute().Subrouter()
	// Middleware For added tags; select if exists else create and select
	resourceTagsSubRouter.Use(mwr.CreateAndSelectAddedTags)
	resourceTagsSubRouter.
		HandleFunc("", hr.PostResource).
		Methods("POST")
	// Middleware select removed tags
	resourceTagsSubRouter.Use(mwr.SelectRemovedTags)
	resourceTagsSubRouter.
		HandleFunc("/{resourceId:[0-9]+}", hr.UpdateResource).
		Methods("PUT")
}
