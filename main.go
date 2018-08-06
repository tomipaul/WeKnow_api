package main

import (
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	// Load env vars from .env
	gotenv.Load()

	// Get database credentials from env vars
	var dbConfig map[string]string
	if val, ok := os.LookupEnv("DATABASE_URL"); ok {
		dbConfig = map[string]string{
			"DATABASE_URL": val,
		}
	} else {
		dbConfig = map[string]string{
			"User":     os.Getenv("DB_USERNAME"),
			"Password": os.Getenv("DB_PASSWORD"),
			"Database": os.Getenv("DATABASE"),
		}
	}

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
