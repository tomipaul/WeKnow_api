package controller

import (
	. "WeKnow_api/pgModel"
	utils "WeKnow_api/utilities"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"

	"github.com/gorilla/context"
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
		if err, errStatus := utils.ValidateSignInRequest(*user); errStatus != true {
			var foundUser User
			if err := db.Model(&foundUser).Where("Email = ?", user.Email).Select(); err != nil {
				if err.Error() == "pg: no rows in result set"{
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid signin parameters")
					return
				}
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else {
				if foundUser.CompareHashAndPassword(user.Password) == true {
					token, _ := foundUser.GenerateToken()
					utils.RespondWithSuccess(w, http.StatusOK, token, "token")
				} else {
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid signin parameters")
					return
				}
			}
		} else {
			utils.RespondWithJsonError(w, http.StatusUnauthorized, err)
		}
	}
	return
}

// ConnectUser create a connection between two users
func ConnectUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload struct{ UserId int64 }
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "A valid userId is required")
	} else {
		decodedClaims := context.Get(r, "decoded")
		UserId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
		initiatorId, recipientId := int64(UserId), payload.UserId
		if initiatorId == recipientId {
			utils.RespondWithError(
				w, http.StatusBadRequest, "You cannot connect to yourself",
			)
			return
		}
		connection := Connection{
			InitiatorId: initiatorId,
			RecipientId: recipientId,
		}
		values := []interface{}{
			connection.InitiatorId,
			connection.RecipientId,
			time.Now(),
			time.Now(),
		}

		q := `WITH connection as
		(INSERT INTO connections(initiator_id, recipient_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id)
		INSERT INTO user_connections(user_id, connection_id) VALUES`

		i, userIds := 5, []int64{initiatorId, recipientId}
		for _, userId := range userIds {
			q += fmt.Sprintf("($%d, (select connection.id from connection)),", i)
			values = append(values, userId)
			i++
		}
		q = strings.TrimSuffix(q, ",")

		stmt, err := db.Prepare(q)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		} else if _, err := stmt.Exec(values...); err != nil {
			if err.(pg.Error).Field('C') == "23505" {
				utils.RespondWithError(
					w, http.StatusConflict,
					"You are already connected with this user",
				)
			} else if err.(pg.Error).Field('C') == "23503" {
				utils.RespondWithError(
					w, http.StatusBadRequest,
					"User does not exist",
				)
			} else {
				utils.RespondWithError(
					w, http.StatusInternalServerError,
					"Something went wrong!",
				)
			}
			return
		} else {
			message := "connection successful"
			key := "message"
			utils.RespondWithSuccess(w, http.StatusOK, message, key)
		}
	}
	return
}
