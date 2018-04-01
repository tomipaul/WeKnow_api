package main

import (
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()

	// Load database config from .env
	dbConfig := map[string]string{
		"User":     os.Getenv("DB_USERNAME"),
		"Password": os.Getenv("DB_PASSWORD"),
		"Database": os.Getenv("DATABASE"),
	}

	// Create an instance of the application
	app := CreateApp(dbConfig)

	// Run app on port 3000
	app.run(":3000")
}
