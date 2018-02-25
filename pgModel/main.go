package pgModel

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/subosito/gotenv"
)

// Connect Load env variables and connect to database
func Connect() *pg.DB {
	gotenv.Load()

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DATABASE"),
	})

	err := createSchema(db)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		&User{},
		&Message{},
		&Connection{},
		&Comment{},
		&Resource{},
		&Collection{},
		&Tag{},
		&ResourceTag{},
		&CollectionTag{},
	} {
		if err := db.CreateTable(model, nil); err != nil {
			return nil
		}
	}
	return nil
}
