package model

import (
	"fmt"
	"io/ioutil"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Connect connect to database and set up tables
func Connect(config map[string]string) *pg.DB {
	var dbConfig = &pg.Options{
		User:     config["User"],
		Password: config["Password"],
		Database: config["Database"],
	}
	db := pg.Connect(dbConfig)

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
	content, err := ioutil.ReadFile("model/sql.txt")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}
	return nil
}
