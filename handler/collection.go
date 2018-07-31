package handler

import (
	. "WeKnow_api/model"
	utils "WeKnow_api/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// CreateCollectionEndPoint create a new collection
func (h *Handler) CreateCollectionEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	collection := &Collection{}

	if err := json.NewDecoder(r.Body).Decode(collection); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else {
		decodedClaims := context.Get(r, "decoded")
		userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)
		if err := utils.ValidateNewCollection(collection); err == nil {
			collection.UserId = int64(userId)
			if err := h.Db.Insert(collection); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			utils.RespondWithSuccess(w, http.StatusCreated, fmt.Sprintf("%s collection was created successfully", string(collection.Name)), "message")
		} else {
			utils.RespondWithJsonError(w, http.StatusBadRequest, err.Error())
		}

	}
}

func (h *Handler) GetAllCollections(w http.ResponseWriter, r *http.Request) {

	decodedClaims := context.Get(r, "decoded")
	userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)

	var collections []Collection

	err := h.Db.Model(&collections).
		Column(
			"collection.*",
		).
		Where("collection.user_id = ?", int(userId)).
		Select()
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Something went wrong",
		)
	} else {
		payload := map[string][]Collection{
			"collections": collections,
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	}

}

// UpdateCollectionEndpoint updates a new collection
func (h *Handler) UpdateCollectionEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var collection Collection

	params := mux.Vars(r)

	collectionID, err := strconv.ParseInt(params["collectionID"], 10, 64)

	decodedClaims := context.Get(r, "decoded")
	userId := decodedClaims.(jwt.MapClaims)["userId"].(float64)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Please enter valid collection ID")
		return
	}

	foundCollection := Collection{Id: collectionID}

	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if collection.Name != "" {
		foundCollection.Name = collection.Name
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Please enter valid collection name")
		return
	}

	res, err := h.Db.Model(&foundCollection).Where("id = ? and user_id = ?", collectionID, userId).Column("name").Update()

	if err == nil {
		if res.RowsAffected() == 0 {
			utils.RespondWithJsonError(w, http.StatusNotFound, "Collection not found")
			return
		}
		payload := map[string]interface{}{
			"updatedCollection": foundCollection,
			"message":           "Collection Updated Successfully",
		}
		utils.RespondWithJson(w, http.StatusOK, payload)
	} else {
		fmt.Println(err.Error())
		utils.RespondWithError(
			w, http.StatusInternalServerError,
			"Something went wong",
		)
	}
	return
}

// AddResourceToCollection - It allows users add a resource to a collection
func (h *Handler) AddResourceToCollection(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload struct{ ResourceId int64 }
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.ResourceId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest,
			"A valid ResourceId is required",
		)
		return
	}

	collectionId, _ := strconv.ParseInt(mux.Vars(r)["collectionId"], 10, 64)
	query := `INSERT INTO resource_collections(resource_id, collection_id) VALUES ($1, $2)`

	values := []interface{}{
		payload.ResourceId,
		collectionId,
	}

	stmt, err := h.Db.Prepare(query)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else if _, err := stmt.Exec(values...); err != nil {
		if err.(pg.Error).Field('C') == "23505" {
			utils.RespondWithError(
				w, http.StatusConflict,
				"Resource already added to collection",
			)
		} else if err.(pg.Error).Field('C') == "23503" {
			utils.RespondWithError(
				w, http.StatusNotFound,
				"Resource or Collection does not exist",
			)
		} else {
			utils.RespondWithError(
				w, http.StatusInternalServerError,
				"Something went wrong!",
			)
		}
		return
	} else {
		message := "resource added to collection"
		key := "message"
		utils.RespondWithSuccess(w, http.StatusOK, message, key)
	}

}
