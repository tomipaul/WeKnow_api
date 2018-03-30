package middleware

import (
	. "WeKnow_api/model"
	utils "WeKnow_api/utilities"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-pg/pg"
	"github.com/gorilla/context"
)

// CreateAndSelectAddedTags create/select added tags for resources or collections
func (mw *Middleware) CreateAndSelectAddedTags(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tagTitles struct {
			Tags []string
		}
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if err := json.Unmarshal(bodyBytes, &tagTitles); err != nil {
			utils.RespondWithError(
				w,
				http.StatusBadRequest,
				"Tags should be an array of tag titles",
			)
			return
		}
		if tagTitles.Tags == nil || len(tagTitles.Tags) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		if err := utils.ValidateNewTags(tagTitles.Tags); err != nil {
			utils.RespondWithJsonError(
				w, http.StatusBadRequest,
				err.Error(),
			)
			return
		}
		var allTags []interface{}
		for tagIndex, tagTitle := range tagTitles.Tags {
			title := strings.TrimSpace(strings.Title(tagTitle))
			tagTitles.Tags[tagIndex] = title
			allTags = append(allTags, &Tag{Title: title})
		}
		res, insertTagsErr := mw.Db.Model(allTags...).
			OnConflict("DO NOTHING").
			Insert()
		var selectTagsErr error
		if len(allTags) != res.RowsAffected() {
			selectTagsErr = mw.Db.Model(allTags...).
				Where("title in (?)", pg.In(tagTitles.Tags)).
				Select()
		}
		if insertTagsErr != nil || selectTagsErr != nil {
			utils.RespondWithError(
				w,
				http.StatusInternalServerError,
				"Oops! we couldn't create or select resource tags",
			)
		} else {
			context.Set(r, "tags", allTags)
			next.ServeHTTP(w, r)
		}
	})
}
