package main

import (
	"flag"
	"fmt"
	"os"

	"WeKnow_api/utilities"

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

	flag.Usage = usage
	flag.Parse()

	config := utilities.GetDatabaseCredentials()
	db := utilities.Connect(config)

	var oldVersion, newVersion int64
	err := db.RunInTransaction(func(tx *pg.Tx) (err error) {
		oldVersion, newVersion, err = migrations.Run(db, flag.Args()...)
		return err
	})
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
