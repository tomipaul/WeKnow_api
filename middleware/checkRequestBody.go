package middleware

import (
	utils "WeKnow_api/utilities"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// CheckRequestBody check if request body is empty for post and put requests
func (mw *Middleware) CheckRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "PUT" {
			next.ServeHTTP(w, r)
			return
		}
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if len(bodyBytes) == 0 {
			utils.RespondWithError(
				w,
				http.StatusBadRequest,
				"Empty Request Payload",
			)
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			utils.RespondWithError(
				w,
				http.StatusBadRequest,
				"Invalid request payload",
			)
			return
		}
		if len(payload) == 0 || payload == nil {
			utils.RespondWithError(w, http.StatusBadRequest,
				"No fields in request payload")
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
