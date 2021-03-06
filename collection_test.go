package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	. "WeKnow_api/libs/supertest"
)

var valueMap interface{}

func TestGetAllCollection(t *testing.T) {

	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	type ExpectedCollection struct {
		Id     int64
		Name   string
		UserId int64
	}
	type ExpectedResponse struct {
		Collections []ExpectedCollection
		TotalCount  int
	}

	testUser := dummyData["naruto"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testCollDummy1 := dummyData["collection1"].(map[string]interface{})
	testCollDummy2 := dummyData["collection2"].(map[string]interface{})
	testCollDummy1["userId"] = user.Id
	testCollDummy2["userId"] = user.Id
	testColl1 := addTestCollection(t, testCollDummy1)
	testColl2 := addTestCollection(t, testCollDummy2)
	totalTestCollections := 2

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
		expectedKey := `["collections","totalCount"]`
		expectedFirstCollectionName := "first collection"
		expectedSecondCollectionName := "second collection"

		for key, value := range responseMap {

			keys = append(keys, key)

			valueMap = value

		}
		// ensure that response keys are sorted in slice
		sort.Strings(keys)

		obtainedFirstCollectionName := responseMap["collections"].([]interface{})[0].(map[string]interface{})["Name"]
		obtainedSecondCollectionName := responseMap["collections"].([]interface{})[1].(map[string]interface{})["Name"]

		if obtainedFirstCollectionName != expectedFirstCollectionName {
			t.Fatalf("Expected name %q; Got name %q", expectedFirstCollectionName, obtainedFirstCollectionName)
		}

		if obtainedSecondCollectionName != expectedSecondCollectionName {
			t.Fatalf("Expected name %q; Got name %q", expectedSecondCollectionName, obtainedSecondCollectionName)
		}

		keyByte, _ := json.Marshal(keys)
		if string(keyByte) != expectedKey {
			t.Fatalf("Expected keys %q; Got keys %q", expectedKey, string(keyByte))
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

	t.Run("can paginate user's collections", func(t *testing.T) {
		expectedCollections := []ExpectedCollection{
			ExpectedCollection{
				testColl2.Id,
				testColl2.Name,
				testColl2.UserId,
			},
		}
		expectedResponse := ExpectedResponse{
			expectedCollections,
			totalTestCollections,
		}

		Request(testServer.URL, t).
			Get(fmt.Sprintf(
				"api/v1/collection?limit=%v&page=%v", 1, 2)).
			Set("authorization", userToken).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(expectedResponse).
			End()
	})

	t.Run("can paginate user's collections", func(t *testing.T) {
		expectedCollections := []ExpectedCollection{
			ExpectedCollection{
				testColl1.Id,
				testColl1.Name,
				testColl1.UserId,
			},
			ExpectedCollection{
				testColl2.Id,
				testColl2.Name,
				testColl2.UserId,
			},
		}
		expectedResponse := ExpectedResponse{
			expectedCollections,
			totalTestCollections,
		}

		Request(testServer.URL, t).
			Get(fmt.Sprintf(
				"api/v1/collection?limit=%v&page=%v", 2, 1)).
			Set("authorization", userToken).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(expectedResponse).
			End()
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

	testCollection := dummyData["collection1"].(map[string]interface{})
	anotherTestCollection := dummyData["collection2"].(map[string]interface{})
	testCollection["userId"] = user.Id
	anotherTestCollection["userId"] = anotherUser.Id
	addTestCollection(t, testCollection)
	addTestCollection(t, anotherTestCollection)

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

func TestResource(t *testing.T) {

	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testResource := dummyData["testResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	resource := addTestResource(t, testResource)

	testCollection := dummyData["collection1"].(map[string]interface{})
	testCollection["userId"] = user.Id
	collection := addTestCollection(t, testCollection)

	testCollection2 := dummyData["collection2"].(map[string]interface{})
	testCollection2["userId"] = user.Id
	collection2 := addTestCollection(t, testCollection2)

	collectionURI := fmt.Sprintf("/api/v1/collection/add/%v", collection.Id)

	collectionURI2 := fmt.Sprintf("/api/v1/collection/add/%v", collection2.Id)

	invalidCollectionURI := fmt.Sprintf("/api/v1/collection/add/%v", 100)

	payload := `{
		"ResourceId": %v
	}`

	payload = fmt.Sprintf(payload, resource.Id)

	invalidPayload := `{
		"ResourceId": 0
	}`

	invalidPayload1 := `{
		"ResourceId": 1000
	}`

	t.Run("can be added to a collection", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(collectionURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(`{"message":"resource added to collection"}`).
			End()
	})

	t.Run("cannot be added more than once", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(collectionURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(409).
			Expect("Content-Type", "application/json").
			Expect(`{"error": "Resource already added to collection"}`).
			End()
	})

	t.Run("can be added to a another or different collection", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(collectionURI2).
			Set("authorization", userToken).
			Send(payload).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(`{"message":"resource added to collection"}`).
			End()
	})

	t.Run("cannot be added to unknown collection", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(invalidCollectionURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error": "Resource or Collection does not exist"}`).
			End()
	})

	t.Run("cannot be added to a collection because it doesn't exist", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(collectionURI).
			Set("authorization", userToken).
			Send(invalidPayload1).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error": "Resource or Collection does not exist"}`).
			End()
	})

	t.Run("cannot be added because it is invalid", func(t *testing.T) {
		Request(testServer.URL, t).
			Post(collectionURI).
			Set("authorization", userToken).
			Send(invalidPayload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error": "A valid ResourceId is required"}`).
			End()
	})

}
