package handler

import (
	. "WeKnow_api/model"
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

// UserSignUpEndPoint user sign up
func (h *Handler) UserSignUpEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
	} else {
		if err := utils.ValidateSignUpRequest(*user); err == nil {
			if err := h.Db.Insert(user); err != nil {
				if err.(pg.Error).Field('C') == "23505" {
					utils.RespondWithError(w, http.StatusConflict, "User already exists")
				}
			} else if token, err := user.GenerateToken(); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
			} else {
				payload := map[string]interface{}{
					"token":   token,
					"message": "Authentication successful",
				}
				utils.RespondWithJson(w, http.StatusOK, payload)
			}
		} else {
			utils.RespondWithJsonError(w, http.StatusBadRequest, err.Error())
		}
	}
	return
}

// UserSignInEndPoint user login
func (h *Handler) UserSignInEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
	} else {
		if err := utils.ValidateSignInRequest(*user); err == nil {
			var foundUser User
			if err := h.Db.Model(&foundUser).Where("Email = ?", user.Email).Select(); err != nil {
				if err.Error() == "pg: no rows in result set" {
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid signin parameters")
					return
				}
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else {
				if foundUser.CompareHashAndPassword(user.Password) == true {
					token, _ := foundUser.GenerateToken()
					payload := map[string]interface{}{
						"token":   token,
						"message": "Authentication successful",
					}
					utils.RespondWithJson(w, http.StatusOK, payload)
				} else {
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid signin parameters")
					return
				}
			}
		} else {
			utils.RespondWithJsonError(w, http.StatusUnauthorized, err.Error())
		}
	}
	return
}

// ConnectUser create a connection between two users
func (h *Handler) ConnectUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload struct{ UserId int64 }
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest,
			"A valid userId is required",
		)
	} else {
		decodedClaims := context.Get(r, "decoded")
		UserId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
		initiatorId, recipientId := int64(UserId), payload.UserId
		if initiatorId == recipientId {
			utils.RespondWithError(
				w, http.StatusBadRequest,
				"You cannot connect to yourself",
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

		stmt, err := h.Db.Prepare(q)
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

// GetAllFavorites get all the favorites of a user
func (h *Handler) GetAllFavorites(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decodedClaims := context.Get(r, "decoded")
	userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
	var connection []Connection
	err := h.Db.Model(&connection).
		Column(
			"connection.id",
			"initiator_id",
			"recipient_id",
			"connection.created_at",
			"Recipient.id",
			"Recipient.email",
			"Recipient.username",
		).
		Where("initiator_id = ?", int(userId)).
		Select()
	if err != nil {
		utils.RespondWithError(
			w, http.StatusInternalServerError,
			"Something went wong",
		)
	} else {
		payload := map[string][]Connection{
			"connections": connection,
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	}
}

// GetAllFollowers get all the followers of a user
func (h *Handler) GetAllFollowers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decodedClaims := context.Get(r, "decoded")
	userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
	var connection []Connection
	err := h.Db.Model(&connection).
		Column(
			"connection.id",
			"initiator_id",
			"recipient_id",
			"connection.created_at",
			"Initiator.id",
			"Initiator.email",
			"Initiator.username",
		).
		Where("recipient_id = ?", int(userId)).
		Select()
	if err != nil {
		utils.RespondWithError(
			w, http.StatusInternalServerError,
			"Something went wong",
		)
	} else {
		payload := map[string][]Connection{
			"connections": connection,
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	}
}

// UpdateProfile helps user update their profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decodedClaims := context.Get(r, "decoded")
	userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
	var user map[string]interface{}
	foundUser := &User{Id: int64(userId)}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := utils.ValidateProfileFields(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	updatedFields := []string{}
	for key, value := range user {
		switch key {
		case "username":
			foundUser.Username = value.(string)
			updatedFields = append(updatedFields, "username")
		case "phoneNumber":
			foundUser.PhoneNumber = value.(string)
			updatedFields = append(updatedFields, "phone_number")
		case "email":
			foundUser.Email = value.(string)
			updatedFields = append(updatedFields, "email")
		}
	}
	res, err := h.Db.Model(foundUser).Column(updatedFields...).Update()

	if err == nil {
		if res.RowsAffected() == 0 {
			utils.RespondWithJsonError(w, http.StatusNotFound, "User not found")
			return
		}
		payload := map[string]interface{}{
			"updatedProfile": user,
			"message":        "Profile Updated successfully",
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	} else {
		utils.RespondWithError(
			w, http.StatusInternalServerError,
			"Something went wong",
		)
	}
	return
}
