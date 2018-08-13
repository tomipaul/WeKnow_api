package utilities

import (
	"os"

	"github.com/go-pg/pg"
)

// Connect connect to database
func Connect(config map[string]string) (db *pg.DB) {
	var dbConfig *pg.Options
	if dbURL, ok := config["DATABASE_URL"]; ok {
		var err error
		dbConfig, err = pg.ParseURL(dbURL)
		if err != nil {
			panic(err)
		}
	} else {
		dbConfig = &pg.Options{
			User:     config["User"],
			Password: config["Password"],
			Database: config["Database"],
		}
	}
	db = pg.Connect(dbConfig)
	return
}

// GetDatabaseCredentials Get database credentials from env vars
func GetDatabaseCredentials() (dbConfig map[string]string) {
	if val, ok := os.LookupEnv("USE_DATABASE_URL"); ok {
		dbConfig = map[string]string{
			"DATABASE_URL": os.Getenv(val),
		}
	} else {
		dbConfig = map[string]string{
			"User":     os.Getenv("DB_USERNAME"),
			"Password": os.Getenv("DB_PASSWORD"),
			"Database": os.Getenv("DATABASE"),
		}
	}
	return
}
