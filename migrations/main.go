package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/subosito/gotenv"
)

const usageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

func main() {
	gotenv.Load()
	var env string
	var db *pg.DB

	flag.StringVar(
		&env,
		"env",
		"dev",
		"Run migration for specified environment")

	flag.Usage = usage
	flag.Parse()

	if env == "test" {
		db = pg.Connect(&pg.Options{
			User:     os.Getenv("TEST_DB_USERNAME"),
			Password: os.Getenv("TEST_DB_PASSWORD"),
			Database: os.Getenv("TEST_DATABASE"),
		})
	} else {
		db = pg.Connect(&pg.Options{
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DATABASE"),
		})
	}

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Printf(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
