package helper 

import (
	"regexp"
	"strings"
	"net/http"
	"encoding/json"
	. "WeKnow_api/pgModel"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

const EXP_EMAIL = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
var fallback interface{}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJson(w, code, map[string]string{"error": msg})
}

func RespondWithJsonError(w http.ResponseWriter, code int, payload interface{}) {
	RespondWithJson(w, code, map[string]interface{}{"error": payload})
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithSuccess(w http.ResponseWriter, code int, msg string, key string) {
	RespondWithJson(w, code, map[string]string{key: msg})
}

func GenerateToken(w http.ResponseWriter, r *http.Request, user User){
    _ = json.NewDecoder(r.Body).Decode(&user)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
		"password": user.Password,
		"email": user.Email,
		"phoneNumber": user.PhoneNumber,
    })
    tokenString, error := token.SignedString([]byte("secret"))
    if error != nil {
        fmt.Println(error)
    }
    RespondWithSuccess(w,201,tokenString,"userToken")
}

func createErrorMessage(key string, value string) interface{}{
	return map[string]string{key: value}
}

func ValidateSignUpRequest(user User) (interface{},bool) {

	re := regexp.MustCompile(EXP_EMAIL)

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.PhoneNumber = strings.TrimSpace(user.PhoneNumber)

	if user.Username == "" {
		return createErrorMessage("username","Username is required"), true
	}
	if user.Password == "" {
		return createErrorMessage("password","Password is required"), true
	}
	if user.Email == "" {
		return createErrorMessage("email","Email is required"), true
	} else if re.MatchString(user.Email) != true {
        return createErrorMessage("email","Please enter a valid email"), true
	}
	if len(user.PhoneNumber) < 11 || len(user.PhoneNumber) > 11 {
		return createErrorMessage("phoneNumber","Valid phone number is required"), true
	}

	return fallback,false
}

func ValidateSignInRequest(user User) (interface{},bool){
	re := regexp.MustCompile(EXP_EMAIL)

    if user.Email == "" {
		return createErrorMessage("email","Email is required"), true
	} else if re.MatchString(user.Email) != true {
        return createErrorMessage("email","Please enter a valid email"), true
	}
	if user.Password == "" {
		return createErrorMessage("password","Password is required"), true
	}

	return fallback,false
}
