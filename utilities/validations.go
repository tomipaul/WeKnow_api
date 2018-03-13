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
		return createErrorMessage("username", "Username is required"), true
	}
	if user.Password == "" {
		return createErrorMessage("password", "Password is required"), true
	}
	if user.Email == "" {
		return createErrorMessage("email", "Email is required"), true
	} else if re.MatchString(user.Email) != true {
		return createErrorMessage("email", "Please enter a valid email"), true
	}
	if len(user.PhoneNumber) < 11 || len(user.PhoneNumber) > 11 {
		return createErrorMessage("phoneNumber", "Valid phone number is required"), true
	}

	return fallback, false
}

func ValidateSignInRequest(user User) (interface{}, bool) {
	re := regexp.MustCompile(EXP_EMAIL)

	if user.Email == "" {
		return createErrorMessage("email", "Email is required"), true
	} else if re.MatchString(user.Email) != true {
		return createErrorMessage("email", "Please enter a valid email"), true
	}
	if user.Password == "" {
		return createErrorMessage("password", "Password is required"), true
	}

	return fallback, false
}

func ValidateNewCollection(coll Collection) (interface{}, bool) {
	
	coll.Name = strings.TrimSpace(coll.Name)

	if coll.Name == "" {
		return createErrorMessage("name", "Collection name is required"), true
	}

	return fallback, false
}