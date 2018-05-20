package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var valueMap interface{}

func TestCollection(t *testing.T) {

	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["naruto"].(map[string]interface{})
	_, userToken := addTestUser(t, testUser)

	testCollection1 := dummyData["collection1"].(map[string]interface{})
	testCollection2 := dummyData["collection2"].(map[string]interface{})
	addTestCollection(t, testCollection1)
	addTestCollection(t, testCollection2)

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
