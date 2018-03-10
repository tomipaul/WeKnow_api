package pgModel

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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
		&UserConnection{},
	} {
		if err := db.CreateTable(
			model,
			&orm.CreateTableOptions{IfNotExists: true},
		); err != nil {
			return err
		}
	}
	content, err := ioutil.ReadFile("pgModel/sql.txt")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}
	return nil
}
