package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "WeKnow_api/libs/supertest"

	"testing"
)

func TestAddComment(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testResource := dummyData["testResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	resource := addTestResource(t, testResource)

	commentURI := "/api/v1/comment"
	t.Run("cannot add comment with invalid field types", func(t *testing.T) {
		comment := `{
			"text": [],
			"resourceId": ""
		}`
		Request(testServer.URL, t).
			Post(commentURI).
			Set("authorization", userToken).
			Send(comment).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid field(s) in request payload"}`).
			End()
	})

	t.Run("cannot add comment with empty request body", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(commentURI).
			Set("authorization", userToken).
			Send(`{}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Empty Request Payload"}`).
			End()
	})

	t.Run("cannot add comment with no fields in payload ", func(t *testing.T) {
		payload := `{}`

		client := testServer.Client()
		uri := testServer.URL + commentURI
		reader := strings.NewReader(payload)
		request, _ := http.NewRequest("POST", uri, reader)
		request.Header.Set("authorization", userToken)

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err.Error())
		}

		var responseMap map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
			t.Fatal(err.Error())
		}

		errorMessage, OK := responseMap["error"].(string)
		if !OK {
			t.Fatalf("Expected key error in response")
		} else {
			expectedErrorMessage := "No fields in request payload"
			if errorMessage != expectedErrorMessage {
				t.Fatalf("Expected error %q; Got error %q",
					expectedErrorMessage, errorMessage)
			}
		}

		expectedStatusCode := 400
		obtainedStatusCode := response.StatusCode
		if expectedStatusCode != obtainedStatusCode {
			t.Fatalf("Expected status %v; Got status %v",
				expectedStatusCode, obtainedStatusCode)
		}

		expectedContentType := "application/json"
		obtainedContentType := response.Header.Get("Content-Type")
		if expectedContentType != obtainedContentType {
			t.Fatalf("Expected content-type %q; Got content-type %q",
				expectedContentType, obtainedContentType)
		}
	})

	t.Run("cannot add comment with empty text", func(t *testing.T) {
		comment := map[string]interface{}{
			"text":       "",
			"resourceId": resource.Id,
		}
		Request(testServer.URL, t).
			Post(commentURI).
			Set("authorization", userToken).
			Send(comment).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"comment Text is required, it cannot be empty"}`).
			End()
	})

	t.Run("cannot add comment with invalid resourceId", func(t *testing.T) {
		comment := map[string]interface{}{
			"text":       "This is a comment",
			"resourceId": 0,
		}
		Request(testServer.URL, t).
			Post(commentURI).
			Set("authorization", userToken).
			Send(comment).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resourceId is required, comment must be associated with a valid resource"}`).
			End()
	})

	t.Run("cannot add comment for nonexistent resource", func(t *testing.T) {
		comment := map[string]interface{}{
			"text":       "This is a comment",
			"resourceId": 15,
		}
		Request(testServer.URL, t).
			Post(commentURI).
			Set("authorization", userToken).
			Send(comment).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Resource with id 15 does not exist"}`).
			End()
	})

	t.Run("can add comment for an existing resource", func(t *testing.T) {
		comment := map[string]interface{}{
			"text":       "This is a comment",
			"resourceId": resource.Id,
		}

		client := testServer.Client()
		uri := testServer.URL + commentURI
		reader := bytes.NewBuffer([]byte{})
		json.NewEncoder(reader).Encode(comment)
		request, _ := http.NewRequest("POST", uri, reader)
		request.Header.Set("authorization", userToken)

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err.Error())
		}

		var responseMap map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
			t.Fatal(err.Error())
		}

		responseComment := responseMap["comment"].(map[string]interface{})

		expectedText := comment["text"].(string)
		obtainedText := responseComment["Text"]
		if obtainedText != expectedText {
			t.Fatalf("Expected comment text %q; Got %q", expectedText, obtainedText)
		}

		expectedResourceId := resource.Id
		obtainedResourceId := int64(responseComment["ResourceId"].(float64))
		if obtainedResourceId != expectedResourceId {
			t.Fatalf("Expected resource id %d; Got resource id %d", expectedResourceId, obtainedResourceId)
		}

		expectedUserId := user.Id
		obtainedUserId := int64(responseComment["UserId"].(float64))
		if obtainedUserId != expectedUserId {
			t.Fatalf("Expected user id %d; Got user id %d", expectedUserId, obtainedUserId)
		}

		expectedLikes := 0
		obtainedLikes := int(responseComment["Likes"].(float64))
		if obtainedLikes != expectedLikes {
			t.Fatalf("Expected likes count %d; Got likes count %d", expectedLikes, obtainedLikes)
		}

		expectedStatusCode := 200
		obtainedStatusCode := response.StatusCode
		if expectedStatusCode != obtainedStatusCode {
			t.Fatalf("Expected status %v; Got status %v",
				expectedStatusCode, obtainedStatusCode)
		}

		expectedContentType := "application/json"
		obtainedContentType := response.Header.Get("Content-Type")
		if expectedContentType != obtainedContentType {
			t.Fatalf("Expected content-type %q; Got content-type %q",
				expectedContentType, obtainedContentType)
		}

		obtainedResponseMessage := responseMap["message"].(string)
		expectedMessage := "Comment added to resource"
		if expectedMessage != obtainedResponseMessage {
			t.Fatalf("Expected response message %q; Got message %q",
				expectedMessage, obtainedResponseMessage)
		}
	})
}
