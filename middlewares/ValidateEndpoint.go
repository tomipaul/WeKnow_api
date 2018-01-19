package middlewares

import (
	"fmt"
    . "WeKnow_api/helper"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"os"
	"github.com/subosito/gotenv"
)

func ValidateEndpoint(next http.HandlerFunc) http.HandlerFunc {
	gotenv.Load()
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authorizationHeader := r.Header.Get("authorization")
        if authorizationHeader != "" {
            bearerToken := strings.Split(authorizationHeader, " ")
            if len(bearerToken) == 2 {
                token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("There was an error")
                    }
                    return []byte(os.Getenv("JWT_SECRET")), nil
                })
                if error != nil {
                    RespondWithError(w,401,error.Error())
                    return
                }
                if token.Valid {
                    context.Set(r, "decoded", token.Claims)
                    next(w, r)
                } else {
                    RespondWithError(w,401,"Invalid authorization token")
                }
            }
        } else {
            RespondWithError(w,401,"An authorization header is required")
        }
    })
}