package utilities

import (
	. "WeKnow_api/pgModel"
	"regexp"
	"strings"
)

const EXP_EMAIL = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var fallback interface{}

func ValidateSignUpRequest(user User) (interface{}, bool) {

	re := regexp.MustCompile(EXP_EMAIL)

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.PhoneNumber = strings.TrimSpace(user.PhoneNumber)

	if user.Username == "" {
		return CreateErrorMessage("username", "Username is required"), true
	}
	if user.Password == "" {
		return CreateErrorMessage("password", "Password is required"), true
	}
	if user.Email == "" {
		return CreateErrorMessage("email", "Email is required"), true
	} else if re.MatchString(user.Email) != true {
		return CreateErrorMessage("email", "Please enter a valid email"), true
	}
	if len(user.PhoneNumber) < 11 || len(user.PhoneNumber) > 11 {
		return CreateErrorMessage("phoneNumber", "Valid phone number is required"), true
	}

	return fallback, false
}

func ValidateSignInRequest(user User) (interface{}, bool) {
	re := regexp.MustCompile(EXP_EMAIL)

	if user.Email == "" {
		return CreateErrorMessage("email", "Email is required"), true
	} else if re.MatchString(user.Email) != true {
		return CreateErrorMessage("email", "Please enter a valid email"), true
	}
	if user.Password == "" {
		return CreateErrorMessage("password", "Password is required"), true
	}

	return fallback, false
}
