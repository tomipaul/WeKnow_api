package model

import (
	"io/ioutil"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Connect connect to database and set up tables
func Connect(config map[string]string) *pg.DB {
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
	db := pg.Connect(dbConfig)
	return db
}

// CreateSchema create database tables
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		&User{},
		&Message{},
		&Connection{},
		&Resource{},
		&Comment{},
		&Collection{},
		&Tag{},
		&ResourceTag{},
		&CollectionTag{},
		&UserConnection{},
		&Recommendation{},
	} {
		if err := db.CreateTable(
			model,
			&orm.CreateTableOptions{IfNotExists: true, FKConstraints: true},
		); err != nil {
			return err
		}
	}
	content, err := ioutil.ReadFile("migrations/sql.txt")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}
	return nil
}

// DropSchema drop database tables
func DropSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		&User{},
		&Message{},
		&Connection{},
		&Resource{},
		&Comment{},
		&Collection{},
		&Tag{},
		&ResourceTag{},
		&CollectionTag{},
		&UserConnection{},
		&Recommendation{},
	} {
		if err := db.DropTable(
			model,
			&orm.DropTableOptions{IfExists: true, Cascade: true},
		); err != nil {
			return err
		}
	}
	return nil
}
