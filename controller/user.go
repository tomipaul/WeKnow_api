package controller

import (
	. "WeKnow_api/pgModel"
	utils "WeKnow_api/utilities"
	"fmt"

	"encoding/json"
	"net/http"
)

var db = Connect()

func UserSignUpEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
	} else {
		if err, errStatus := utils.ValidateSignUpRequest(*user); errStatus != true {
			if err := db.Insert(user); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			} else if token, err := user.GenerateToken(); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else {
				utils.RespondWithSuccess(w, http.StatusOK, token, "token")
			}
		} else {
			utils.RespondWithJsonError(w, 404, err)
		}
	}
	return
}

// UserSignInEndPoint user login
func UserSignInEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
	} else {
		if err, errStatus := utils.ValidateSignInRequest(*user); errStatus != true {
			user1 := User{Email: user.Email}
			if err := db.Model(&user1).Select(); err != nil {
				fmt.Print("I got here")
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else if token, err := user.GenerateToken(); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else {
				utils.RespondWithSuccess(w, http.StatusOK, token, "token")
			}
		} else {
			utils.RespondWithJsonError(w, 404, err)
		}
	}
	return
}
