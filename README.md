# WeKnow
[![CircleCI](https://circleci.com/gh/tomipaul/WeKnow_api.svg?style=svg)](https://circleci.com/gh/tomipaul/WeKnow_api)

Introduction
------------
* WeKnow is a social education platform where people can share resources and interact to learn

Main Features
-------------
* The platform has the following main features

    - Users can share textual, audio and video resources
    - Users can follow one other to get feed on shared content
    - Users can recommend and comment on a shared resource
    - Users can create a collection of shared resources
    - Users can access top rated and top viewed resources

Technologies
------------
- The application was developed with [Golang](https://golang.org/)
- [gorilla/mux](https://github.com/gorilla/mux) was used for routing
- The [Postgres](http://postgresql.com) database was used with [go-pg](https://github.com/go-pg/pg) as the ORM

Installation
------------
1. Ensure you have [Golang](https://golang.org/doc/install) and [Postgres](https://www.postgresql.org/download/) installed
2. Clone the project and place it in the $GOPATH/src
3. Change your directory `cd WeKnow_api`
4. Create a .env file in the root of the directory following the format in the provided .env.example file.
5. Create a database and include database credentials in the .env
6. Then run `go run main.go app.go` on the terminal to start application
7. You can now use WeKnow_api by visiting http://localhost:port (where port is the PORT environment variable in your .env file; defaults to 3000 if not set).

Database Migrations
-------------------
To run all available migrations, run command `go run migrations/*.go` which is equivalent to `go run migrations/*.go up`
Generally migration commands take the form `go run migrations/*.go <command>` where command can be one of:

- up - runs all available migrations.
- down - reverts last migration.
- reset - reverts all migrations.
- version - prints current db version.
- set_version [version] - sets db version without running migrations.
- create [migration name] - creates a new migration file with consecutive version number and migration name

To get more information on how these migrations work, you can read up the README.md at https://github.com/go-pg/migrations

Tests
-----
*  Tests have been written to ensure the API endpoints accept the approriate input and give the right output
*  Run the test with `go test`

