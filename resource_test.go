package main_test

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	. "WeKnow_api/libs/supertest"
)

func TestRecommendResource(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testResource := dummyData["testResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	resource := addTestResource(t, testResource)

	resourceURI := "/api/v1/resource/recommend/"

	t.Run("cannot recommend resource with id 0", func(t *testing.T) {
		Request(testServer.URL, t).
			Get(resourceURI+"0").
			Set("authorization", userToken).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid resource Id in request"}`).
			End()
	})

	t.Run("cannot recommend nonexistent resource", func(t *testing.T) {
		Request(testServer.URL, t).
			Get(resourceURI+"238").
			Set("authorization", userToken).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Resource does not exist"}`).
			End()
	})

	t.Run("can recommend a resource", func(t *testing.T) {
		Request(testServer.URL, t).
			Get(fmt.Sprintf("%v%v", resourceURI, resource.Id)).
			Set("authorization", userToken).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(`{"message":"Recommend resource successful",
			"recommendationCount":1}`).
			End()
	})

	t.Run("can handle simultaneous requests to recommend same resource", func(t *testing.T) {
		anotherTestUser := dummyData["anotherTestUser"].(map[string]interface{})
		_, anotherUserToken := addTestUser(t, anotherTestUser)
		thirdTestUser := dummyData["thirdTestUser"].(map[string]interface{})
		_, thirdUserToken := addTestUser(t, thirdTestUser)

		var firstRecommendationCount, secondRecommendationCount float64

		t.Run("", func(t *testing.T) {
			t.Run("first request", func(t *testing.T) {
				t.Parallel()

				client := testServer.Client()
				uri := fmt.Sprintf("%v%v%v", testServer.URL, resourceURI, resource.Id)
				request, _ := http.NewRequest("GET", uri, nil)
				request.Header.Set("authorization", thirdUserToken)

				response, err := client.Do(request)
				if err != nil {
					t.Fatal(err.Error())
				}
				var responseMap map[string]interface{}
				if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
					t.Fatal(err.Error())
				}
				firstRecommendationCount = responseMap["recommendationCount"].(float64)
				obtainedResponseMessage := responseMap["message"].(string)
				expectedMessage := "Recommend resource successful"
				if expectedMessage != obtainedResponseMessage {
					t.Fatalf("Expected response message %q; Got message %q",
						expectedMessage, obtainedResponseMessage)
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

			t.Run("second request", func(t *testing.T) {
				t.Parallel()

				client := testServer.Client()
				uri := fmt.Sprintf("%v%v%v", testServer.URL, resourceURI, resource.Id)
				request, _ := http.NewRequest("GET", uri, nil)
				request.Header.Set("authorization", anotherUserToken)

				response, err := client.Do(request)
				if err != nil {
					t.Fatal(err.Error())
				}
				var responseMap map[string]interface{}
				if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
					t.Fatal(err.Error())
				}
				secondRecommendationCount = responseMap["recommendationCount"].(float64)
				obtainedResponseMessage := responseMap["message"].(string)
				expectedMessage := "Recommend resource successful"
				if expectedMessage != obtainedResponseMessage {
					t.Fatalf("Expected response message %q; Got message %q",
						expectedMessage, obtainedResponseMessage)
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
		})
		countDiff := firstRecommendationCount - secondRecommendationCount
		if diff := math.Abs(countDiff); diff != 1 {
			t.Fatal("Recommendation counts in response should differ by 1")
		}
	})

	t.Run("cannot recommend a resource twice", func(t *testing.T) {
		Request(testServer.URL, t).
			Get(fmt.Sprintf("%v%v", resourceURI, resource.Id)).
			Set("authorization", userToken).
			Expect(409).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"You have recommended this resource"}`).
			End()
	})
}
