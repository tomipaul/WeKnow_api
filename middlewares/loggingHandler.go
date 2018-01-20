package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggingHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
