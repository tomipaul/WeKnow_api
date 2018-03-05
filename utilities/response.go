package utilities

import (
	"encoding/json"
	"net/http"
)

// RespondWithError  send error message in JSON response
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJson(w, code, map[string]string{"error": msg})
}

// RespondWithJsonError send error object in JSON response
func RespondWithJsonError(w http.ResponseWriter, code int, payload interface{}) {
	RespondWithJson(w, code, map[string]interface{}{"error": payload})
}

// RespondWithJson  set headers and send JSON reponse
func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithSuccess send success message in JSON response
func RespondWithSuccess(w http.ResponseWriter, code int, msg string, key string) {
	RespondWithJson(w, code, map[string]string{key: msg})
}
