package middleware

import (
	"github.com/go-pg/pg"
)

//Middleware type Middleware
type Middleware struct {
	Db *pg.DB
}
