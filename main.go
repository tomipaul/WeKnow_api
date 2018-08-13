package main

import (
	"os"

	"WeKnow_api/utilities"

	"github.com/subosito/gotenv"
)

func main() {
	// Load env vars from .env
	gotenv.Load()

	// Get database credentials from env vars
	dbConfig := utilities.GetDatabaseCredentials()

	// Create an instance of the application
	app := CreateApp(dbConfig)

	// Run app on port; $PORT or fallback to 3000
	var address string
	if port, ok := os.LookupEnv("PORT"); ok {
		address = ":" + port
	} else {
		address = ":" + "3000"
	}
	app.run(address)

}
