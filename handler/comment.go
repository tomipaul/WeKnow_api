package handler

import (
	. "WeKnow_api/model"
	utils "WeKnow_api/utilities"
	"encoding/json"
	"fmt"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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

// GetComments get comments filtered by resource
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	if err := utils.ValidateQueryParams(queryValues); err != nil {
		utils.RespondWithError(
			w, http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	var values []interface{}
	var condition string
	resourceId := queryValues.Get("resourceId")
	if resourceId != "" {
		id, _ := strconv.ParseInt(resourceId, 10, 64)
		if err := utils.ValidateResourceId(id); err != nil {
			utils.RespondWithError(
				w, http.StatusBadRequest,
				err.Error(),
			)
			return
		}
		condition += "resource_id = ?0"
		values = append(values, id)
	}
	if condition == "" {
		utils.RespondWithError(
			w, http.StatusBadRequest,
			"No expected query parameters in request",
		)
		return
	}
	var comments []Comment
	count, err := h.Db.Model(&comments).
		Where(condition, values...).
		Apply(orm.Pagination(r.URL.Query())).
		SelectAndCount()
	if err != nil {
		utils.RespondWithError(
			w, http.StatusInternalServerError,
			"Something went wrong",
		)
	} else {
		payload := map[string]interface{}{
			"totalCount": count,
			"comments":   comments,
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	}
}
