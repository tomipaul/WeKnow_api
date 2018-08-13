package model

import (
	"io/ioutil"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

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
		&ResourceCollection{},
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
		&ResourceCollection{},
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
