package handler

import (
	. "WeKnow_api/model"
	utils "WeKnow_api/utilities"
	"encoding/json"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/gorilla/context"

	"net/http"
)

// AddComment add a comment to a resource
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		utils.RespondWithError(
			w, http.StatusBadRequest,
			"Invalid field(s) in request payload",
		)
		return
	}
	decodedClaims := context.Get(r, "decoded")
	userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
	comment.UserId, comment.Likes = int64(userId), 0
	err = utils.ValidateNewComment(&comment)
	if err != nil {
		utils.RespondWithJsonError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	if err := h.Db.Insert(&comment); err != nil {
		if err.(pg.Error).Field('C') == "23503" {
			errorMsg := fmt.Sprintf(
				"Resource with id %d does not exist",
				comment.ResourceId,
			)
			utils.RespondWithError(
				w, http.StatusNotFound,
				errorMsg,
			)
		} else {
			utils.RespondWithError(
				w, http.StatusInternalServerError,
				"Something went wrong",
			)
		}
	} else {
		payload := map[string]interface{}{
			"comment": comment,
			"message": "Comment added to resource",
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	}
	return
}
