package utilities

import (
	. "WeKnow_api/model"
	"errors"
	"regexp"
	"strings"
)

const EXP_EMAIL = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var re = regexp.MustCompile(EXP_EMAIL)

// ValidateSignUpRequest validate inputs submitted for sign-up
func ValidateSignUpRequest(user *User) error {

	user.Username = strings.TrimSpace(user.Username)
	user.PhoneNumber = strings.TrimSpace(user.PhoneNumber)
	user.Password = strings.TrimSpace(user.Password)

	var err error
	switch {
	case user.Username == "":
		err = errors.New("Username is required")
	case user.Password == "":
		err = errors.New("Password is required")
	case re.MatchString(user.Email) != true:
		err = errors.New("Please enter a valid email")
	case len(user.PhoneNumber) < 11 || len(user.PhoneNumber) > 11:
		err = errors.New("Valid phone number is required")
	}
	return err
}

// ValidateSignInRequest validate inputs submitted for sign-in
func ValidateSignInRequest(user *User) error {

	user.Password = strings.TrimSpace(user.Password)

	var err error
	switch {
	case re.MatchString(user.Email) != true:
		err = errors.New("Please enter a valid email")
	case user.Password == "":
		err = errors.New("Password is required")
	}
	return err
}

// ValidateNewCollection validate inputs submitted to create new collection
func ValidateNewCollection(coll *Collection) error {

	coll.Name = strings.TrimSpace(coll.Name)
	var err error
	if coll.Name == "" {
		err = errors.New("Collection name is required")
	}
	return err
}

func ValidateProfileFields(user map[string]interface{}) error {

	var err error
	for key, value := range user {
		switch key {
		case "username":
			if value.(string) == "" {
				err = errors.New("Username cannot be empty")
			}
		case "email":
			if re.MatchString(value.(string)) != true {
				err = errors.New("Enter a valid email")
			}

		case "phoneNumber":
			if value.(string) == "" {
				err = errors.New("Phone number cannot be empty")
			} else if len(value.(string)) < 11 || len(value.(string)) > 11 {
				err = errors.New("Enter a valid phone number")
			}
		}
	}
	return err
}
