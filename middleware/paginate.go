package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Paginate pagination for bulk data retrieval
func (mw *Middleware) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			next.ServeHTTP(w, r)
			return
		}
		queryValues := r.URL.Query()
		limit := queryValues.Get("limit")
		page := queryValues.Get("page")
		if len(mux.Vars(r)) == 0 {
			if limit == "" {
				limit = os.Getenv("DEFAULT_LIMIT")
			}
			if page == "" {
				page = os.Getenv("DEFAULT_PAGE")
			}
			queryValues.Set("limit", limit)
			queryValues.Set("page", page)
			r.URL.RawQuery = queryValues.Encode()
		}
		next.ServeHTTP(w, r)
		return
	})
}
