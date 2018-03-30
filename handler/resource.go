package handler

import (
	. "WeKnow_api/model"
	utils "WeKnow_api/utilities"
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/gorilla/context"
)

// PostResource post a new resource
func (h *Handler) PostResource(w http.ResponseWriter, r *http.Request) {
	var resource struct {
		Resource
		Tags []string
	}
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		utils.RespondWithError(
			w, http.StatusBadRequest,
			"Invalid resource field(s) in request payload",
		)
	} else {
		decodedClaims := context.Get(r, "decoded")
		userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
		resource.UserId = int64(userId)
		err := utils.ValidateNewResource(&resource.Resource)
		if err != nil {
			utils.RespondWithJsonError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}
		if err := h.Db.Insert(&resource.Resource); err != nil {
			if err.(pg.Error).Field('C') == "23505" {
				utils.RespondWithError(
					w, http.StatusConflict,
					"A resource exists with provided link",
				)
			} else {
				utils.RespondWithError(
					w,
					http.StatusInternalServerError,
					"Something went wrong",
				)
			}
		} else {
			if tags, Ok := context.GetOk(r, "tags"); Ok {
				var resourceTags []interface{}
				for _, tag := range tags.([]interface{}) {
					resourceTags = append(resourceTags, &ResourceTag{
						TagId:      tag.(*Tag).Id,
						ResourceId: resource.Id,
					})
				}
				if err := h.Db.Insert(resourceTags...); err != nil {
					utils.RespondWithError(
						w, http.StatusInternalServerError,
						"Oops! we couldn't attach tags to the resource",
					)
					return
				}
			}
			payload := map[string]interface{}{
				"resource": resource.Resource,
				"tags":     resource.Tags,
				"message":  "Resource created",
			}
			utils.RespondWithJson(w, http.StatusCreated, payload)
		}
	}
	return
}
