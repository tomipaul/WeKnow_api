package middlewares

import (
	utils "WeKnow_api/utilities"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func ValidateEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if bearerToken[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Bearer should be added before token")
				return
			}
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok ||
						!token.Claims.(jwt.MapClaims).VerifyIssuer(os.Getenv("ISSUER"), true) {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if error != nil {
					utils.RespondWithError(w, http.StatusUnauthorized, error.Error())
				} else if token.Valid {
					context.Set(r, "decoded", token.Claims)
					next.ServeHTTP(w, r)
				}
			} else {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
			}
		} else {
			utils.RespondWithError(w, http.StatusUnauthorized, "An authorization header is required")
		}
		return
	})
}
