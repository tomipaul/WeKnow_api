package controller

import (
	. "WeKnow_api/pgModel"
	utils "WeKnow_api/utilities"

	"encoding/json"
	"net/http"
	"github.com/go-pg/pg"
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
				if err.(pg.Error).Field('C') == "23505" {
					utils.RespondWithError(w, http.StatusConflict, "User already exists")
				}
			} else if token, err := user.GenerateToken(); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
			} else {
				utils.RespondWithSuccess(w, http.StatusOK, token, "token")
			}
		} else {
			utils.RespondWithJsonError(w, http.StatusBadRequest, err)
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
		if err,errStatus := utils.ValidateSignInRequest(*user); errStatus != true {
			var foundUser User
			if err := db.Model(&foundUser).Where("Email = ?", user.Email).Select(); err != nil {
				if err.Error() == "pg: no rows in result set"{
					utils.RespondWithJsonError(w, 401, utils.CreateErrorMessage("message","Invalid signin parameters"))
					return
				}
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else {
				if foundUser.CompareHashAndPassword(user.Password) == true {
					token, _ := foundUser.GenerateToken()
					utils.RespondWithSuccess(w, http.StatusOK, token, "token")
				} else {
					utils.RespondWithJsonError(w, 401, utils.CreateErrorMessage("message","Invalid signin parameters"))
					return
				}
			}
		} else {
			utils.RespondWithJsonError(w, 400, err)
		}
	}
	return
}