package main

import (
	. "WeKnow_api/model"
	"fmt"
	"io/ioutil"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"
)

func createSchema(db migrations.DB) error {
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
		if _, err := orm.CreateTable(
			db,
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

func dropSchema(db migrations.DB) error {
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
		if _, err := orm.DropTable(
			db,
			model,
			&orm.DropTableOptions{IfExists: true, Cascade: true},
		); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating tables...")
		err := createSchema(db)
		return err

	}, func(db migrations.DB) error {
		fmt.Println("dropping tables...")
		err := dropSchema(db)
		return err
	})
}
