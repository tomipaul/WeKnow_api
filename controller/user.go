package controller

import (
	"fmt"
	. "WeKnow_api/helper"
	. "WeKnow_api/pgModel"
	"encoding/json"
	"net/http"
)

var db = Connect()

func UserSignUpEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else {
		
		if err,errStatus := ValidateSignUpRequest(*user); errStatus != true {
			if err := db.Insert(user); err != nil {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			GenerateToken(w, r, *user)
		} else {
			RespondWithJsonError(w,404,err)
		}

	}
}

func UserSignInEndPoint(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else {
		
		if err,errStatus := ValidateSignInRequest(*user); errStatus != true {
			user1 := User{Email: user.Email}
			if err := db.Model(&user1).Select(); err != nil {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			GenerateToken(w, r, user1)
		} else {
			RespondWithJsonError(w,404,err)
		}

	}
}
