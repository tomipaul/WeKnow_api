package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "WeKnow_api/libs/supertest"
)

var valueMap interface{}

func TestCollection(t *testing.T) {

	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["naruto"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testCollection1 := dummyData["collection1"].(map[string]interface{})
	testCollection2 := dummyData["collection2"].(map[string]interface{})
	addTestCollection(t, testCollection1, user.Id)
	addTestCollection(t, testCollection2, user.Id)

	t.Run("can get all user's collections", func(t *testing.T) {

		client := testServer.Client()
		uri := testServer.URL + "/api/v1/collection"
		request, _ := http.NewRequest("GET", uri, nil)
		request.Header.Set("authorization", userToken)

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err.Error())
		}

		var responseMap map[string]interface{}

		if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
			t.Fatal(err.Error())
		}

		keys := []string{}
		expectedKey := "collections"
		expectedFirstCollectionName := "first collection"
		expectedSecondCollectionName := "second collection"

		for key, value := range responseMap {

			keys = append(keys, key)

			valueMap = value

		}

		obtainedFirstCollectionName := responseMap["collections"].([]interface{})[0].(map[string]interface{})["Name"]
		obtainedSecondCollectionName := responseMap["collections"].([]interface{})[1].(map[string]interface{})["Name"]

		if obtainedFirstCollectionName != expectedFirstCollectionName {
			t.Fatalf("Expected name %q; Got name %q", expectedFirstCollectionName, obtainedFirstCollectionName)
		}

		if obtainedSecondCollectionName != expectedSecondCollectionName {
			t.Fatalf("Expected name %q; Got name %q", expectedSecondCollectionName, obtainedSecondCollectionName)
		}

		if keys[0] != expectedKey {
			t.Fatalf("Expected key %q; Got key %q", expectedKey, keys[0])
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
	})
}

func TestUpdateCollection(t *testing.T) {

	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()
	defer closeDatabase(t)

	testUser := dummyData["testUser"].(map[string]interface{})

	anotherTestUser := dummyData["anotherTestUser"].(map[string]interface{})

	user, userToken := addTestUser(t, testUser)
	anotherUser, anotherUserToken := addTestUser(t, anotherTestUser)

	testCollection := dummyData["testCollection"].(map[string]interface{})
	anotherTestCollection := dummyData["anotherTestCollection"].(map[string]interface{})

	addTestCollection(t, testCollection, user.Id)
	addTestCollection(t, anotherTestCollection, anotherUser.Id)

	t.Run("cannot be updated when it doesn't belong to the user", func(t *testing.T) {

		Request(testServer.URL, t).
			Put("api/v1/collection/2").
			Set("authorization", userToken).
			Send(`{"name": "new name"}`).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Collection not found"}`).
			End()
	})

	t.Run("can be updated by user who created it", func(t *testing.T) {

		payload := `{"name": "another collection"}`

		client := testServer.Client()
		uri := testServer.URL + "/api/v1/collection/2"
		reader := strings.NewReader(payload)
		request, _ := http.NewRequest("PUT", uri, reader)
		request.Header.Set("authorization", anotherUserToken)

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err.Error())
		}

		var responseMap map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
			t.Fatal(err.Error())
		}

		expectedMessage := "Collection Updated Successfully"
		obtainedMessage := responseMap["message"]
		if obtainedMessage != expectedMessage {
			t.Fatalf("Expected message %q; Got message %q", expectedMessage, obtainedMessage)
		}

		expectedName := "another collection"
		obtainedName := responseMap["updatedCollection"].(map[string]interface{})["Name"]
		if obtainedName != expectedName {
			t.Fatalf("Expected name %q; Got name %q", expectedName, obtainedName)
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
	})

	t.Run("cannot update collection that does not exist", func(t *testing.T) {

		Request(testServer.URL, t).
			Put("api/v1/collection/5").
			Set("authorization", userToken).
			Send(`{"name": "new name"}`).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Collection not found"}`).
			End()
	})

	t.Run("cannot update collection without valid request payload", func(t *testing.T) {

		Request(testServer.URL, t).
			Put("api/v1/collection/5").
			Set("authorization", userToken).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Empty Request Payload"}`).
			End()
	})

	t.Run("name can be updated ", func(t *testing.T) {

		payload := `{"name": "new collection"}`

		client := testServer.Client()
		uri := testServer.URL + "/api/v1/collection/1"
		reader := strings.NewReader(payload)
		request, _ := http.NewRequest("PUT", uri, reader)
		request.Header.Set("authorization", userToken)

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err.Error())
		}

		var responseMap map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
			t.Fatal(err.Error())
		}

		expectedMessage := "Collection Updated Successfully"
		obtainedMessage := responseMap["message"]
		if obtainedMessage != expectedMessage {
			t.Fatalf("Expected message %q; Got message %q", expectedMessage, obtainedMessage)
		}

		expectedName := "new collection"
		obtainedName := responseMap["updatedCollection"].(map[string]interface{})["Name"]
		if obtainedName != expectedName {
			t.Fatalf("Expected name %q; Got name %q", expectedName, obtainedName)
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
	})

	t.Run("cannot be updated with invalid request params", func(t *testing.T) {

		Request(testServer.URL, t).
			Put("api/v1/collection/invalid").
			Set("authorization", userToken).
			Send(`{"name": "new name"}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{ "error": "Please enter valid collection ID"}`).
			End()
	})

	t.Run("cannot be updated with empty name", func(t *testing.T) {

		Request(testServer.URL, t).
			Put("api/v1/collection/1").
			Set("authorization", userToken).
			Send(`{"name": ""}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{ "error": "Please enter valid collection name"}`).
			End()
	})

}
