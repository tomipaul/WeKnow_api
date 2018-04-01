package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg"
)

// Handler type Handler
type Handler struct {
	Db *pg.DB
}

// HomeHandler handle GET request to the root endpoint
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Welcome to the WeKnow API")
}
