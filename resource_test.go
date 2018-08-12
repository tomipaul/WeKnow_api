package main_test

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "WeKnow_api/libs/supertest"
	. "WeKnow_api/model"
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

func TestGetResource(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	type ExpectedResource struct {
		Id      int64
		UserId  int64
		Title   string
		Link    string
		Privacy string
		Type    string
		Tags    []*Tag
	}
	type ExpectedResponse struct {
		Resource ExpectedResource
	}

	testUser := dummyData["testUser"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	anotherTestUser := dummyData["anotherTestUser"].(map[string]interface{})
	anotherUser, anotherUserToken := addTestUser(t, anotherTestUser)

	thirdTestUser := dummyData["thirdTestUser"].(map[string]interface{})
	_, thirdUserToken := addTestUser(t, thirdTestUser)

	testConnectionData := map[string]interface{}{
		"initiatorId": anotherUser.Id,
		"recipientId": user.Id,
	}
	addTestConnection(t, testConnectionData)

	testResource := dummyData["testResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	publicResource := addTestResource(t, testResource)

	testResource = dummyData["privateResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	privateResource := addTestResource(t, testResource)

	testResource = dummyData["followersResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	followersResource := addTestResource(t, testResource)

	resourceURI := "/api/v1/resource/"

	t.Run("cannot get resource with id 0", func(t *testing.T) {
		Request(testServer.URL, t).
			Get(resourceURI+"0").
			Set("authorization", userToken).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid resource Id in request"}`).
			End()
	})

	t.Run("cannot get nonexistent resource", func(t *testing.T) {
		Request(testServer.URL, t).
			Get(resourceURI+"238").
			Set("authorization", userToken).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
			End()
	})

	t.Run("can get a public resource", func(t *testing.T) {
		expectedResponse := ExpectedResponse{
			ExpectedResource{
				publicResource.Id,
				publicResource.UserId,
				publicResource.Title,
				publicResource.Link,
				publicResource.Privacy,
				publicResource.Type,
				publicResource.Tags,
			},
		}

		t.Run("if owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, publicResource.Id)).
				Set("authorization", userToken).
				Expect(200).
				Expect("Content-Type", "application/json").
				Expect(expectedResponse).
				End()
		})

		t.Run("if not owner but a follower of owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, publicResource.Id)).
				Set("authorization", anotherUserToken).
				Expect(200).
				Expect("Content-Type", "application/json").
				Expect(expectedResponse).
				End()
		})

		t.Run("if not owner and not a follower of owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, publicResource.Id)).
				Set("authorization", thirdUserToken).
				Expect(200).
				Expect("Content-Type", "application/json").
				Expect(expectedResponse).
				End()
		})
	})

	t.Run("can get a private resource if owner", func(t *testing.T) {
		expectedResponse := ExpectedResponse{
			ExpectedResource{
				privateResource.Id,
				privateResource.UserId,
				privateResource.Title,
				privateResource.Link,
				privateResource.Privacy,
				privateResource.Type,
				privateResource.Tags,
			},
		}

		Request(testServer.URL, t).
			Get(fmt.Sprintf("%v%v", resourceURI, privateResource.Id)).
			Set("authorization", userToken).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(expectedResponse).
			End()
	})

	t.Run("cannot get private resource if not owner", func(t *testing.T) {
		t.Run("but a follower of owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, privateResource.Id)).
				Set("authorization", anotherUserToken).
				Expect(404).
				Expect("Content-Type", "application/json").
				Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
				End()
		})

		t.Run("and not a follower of owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, privateResource.Id)).
				Set("authorization", thirdUserToken).
				Expect(404).
				Expect("Content-Type", "application/json").
				Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
				End()
		})
	})

	t.Run("can get a resource with privacy follower", func(t *testing.T) {
		expectedResponse := ExpectedResponse{
			ExpectedResource{
				followersResource.Id,
				followersResource.UserId,
				followersResource.Title,
				followersResource.Link,
				followersResource.Privacy,
				followersResource.Type,
				followersResource.Tags,
			},
		}

		t.Run("if owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, followersResource.Id)).
				Set("authorization", userToken).
				Expect(200).
				Expect("Content-Type", "application/json").
				Expect(expectedResponse).
				End()
		})

		t.Run("if a follower of owner", func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, followersResource.Id)).
				Set("authorization", anotherUserToken).
				Expect(200).
				Expect("Content-Type", "application/json").
				Expect(expectedResponse).
				End()
		})
	})

	t.Run("cannot get a resource with privacy follower if not a follower of the owner",
		func(t *testing.T) {
			Request(testServer.URL, t).
				Get(fmt.Sprintf("%v%v", resourceURI, followersResource.Id)).
				Set("authorization", thirdUserToken).
				Expect(404).
				Expect("Content-Type", "application/json").
				Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
				End()
		},
	)
}

func TestUpdateResource(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testResource := dummyData["testResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	resource := addTestResource(t, testResource)

	resourceURI := fmt.Sprintf("/api/v1/resource/%v", resource.Id)

	t.Run("cannot update resource with id 0", func(t *testing.T) {
		payload := `{
			"title": "An updated resource",
			"type": "audio"
		}`
		Request(testServer.URL, t).
			Put("/api/v1/resource/0").
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid resource Id in request"}`).
			End()
	})

	t.Run("cannot update resource with invalid JSON payload", func(t *testing.T) {
		payload := `{
			"title": "invalid payload"
			"type": "invalid payload"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid request payload"}`).
			End()
	})

	t.Run("cannot update resource with no fields in payload ", func(t *testing.T) {
		payload := `{}`

		client := testServer.Client()
		uri := testServer.URL + resourceURI
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

	t.Run("cannot update resource with invalid fields in payload ", func(t *testing.T) {
		payload := `{
			"tle": "invalid payload field"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid keys in request payload"}`).
			End()
	})

	t.Run("cannot update resource with empty request body", func(t *testing.T) {
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(`{}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Empty Request Payload"}`).
			End()
	})

	t.Run("cannot update resource with empty title", func(t *testing.T) {
		payload := `{
			"title": "",
			"privacy": "public"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"A valid title is required"}`).
			End()
	})

	t.Run("cannot update resource with empty type", func(t *testing.T) {
		payload := `{
			"type": "",
			"privacy": "public"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"A valid type is required"}`).
			End()
	})

	t.Run("cannot update resource with invalid type", func(t *testing.T) {
		payload := `{
			"type": "ebook",
			"privacy": "public"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resource Type must be one of 'video', 'audio' or 'textual'"}`).
			End()
	})

	t.Run("cannot update resource with empty link", func(t *testing.T) {
		payload := `{
			"link": "",
			"privacy": "public"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"A valid link is required"}`).
			End()
	})

	t.Run("cannot update resource with empty privacy", func(t *testing.T) {
		payload := `{
			"privacy": ""
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"A valid privacy is required"}`).
			End()
	})

	t.Run("cannot update resource with empty collectionId", func(t *testing.T) {
		payload := `{
			"collectionId": ""
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"A valid collection Id is required"}`).
			End()
	})

	t.Run("cannot update resource with empty tag title strings", func(t *testing.T) {
		payload := `{
			"title": "A new resource",
			"type": "textual",
			"link": "https://localhost.textual/material/6.pdf",
			"privacy": "public",
			"tags": ["", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Tag titles must be non-empty strings"}`).
			End()
	})

	t.Run("cannot update resource with invalid tags field", func(t *testing.T) {
		payload := `{
			"title": "A new resource",
			"type": "textual",
			"link": "https://localhost.textual/material/6.pdf",
			"privacy": "public",
			"tags": "lisp"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Tags should be an array of tag titles"}`).
			End()
	})

	t.Run("cannot update resource with invalid removedTags field", func(t *testing.T) {
		payload := `{
			"removedTags": "pascal"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"RemovedTags should be an array of tag titles"}`).
			End()
	})

	t.Run("cannot update nonexistent resource", func(t *testing.T) {
		payload := `{
			"title": "An updated resource",
			"type": "audio"
		}`
		Request(testServer.URL, t).
			Put("/api/v1/resource/238").
			Set("authorization", userToken).
			Send(payload).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
			End()
	})

	t.Run("can only update own resource", func(t *testing.T) {
		anotherTestUser := dummyData["anotherTestUser"].(map[string]interface{})
		_, userToken := addTestUser(t, anotherTestUser)

		payload := `{
			"title": "An updated resource",
			"type": "audio"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(404).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
			End()
	})

	t.Run("update resource fields and return valid response", func(t *testing.T) {
		payload := `{
			"title": "An updated resource",
			"type": "audio"
		}`
		updatedResource := resource
		updatedResource.Title = "An updated resource"
		updatedResource.Type = "audio"

		client := testServer.Client()
		uri := testServer.URL + resourceURI
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

		expectedTitle := "An updated resource"
		obtainedTitle := responseMap["updatedResource"].(map[string]interface{})["Title"]
		if obtainedTitle != expectedTitle {
			t.Fatalf("Expected title %q; Got title %q", expectedTitle, obtainedTitle)
		}

		expectedType := "audio"
		obtainedType := responseMap["updatedResource"].(map[string]interface{})["Type"]
		if obtainedType != expectedType {
			t.Fatalf("Expected type %q; Got type %q", expectedType, obtainedType)
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

	t.Run("update resource tags and return valid response", func(t *testing.T) {
		payload := `{
			"tags": ["dart", "Go"],
			"removedTags": ["fortran"]
		}`

		expectedResponse := map[string]interface{}{
			"updatedResource": Resource{
				Id:     resource.Id,
				UserId: resource.UserId,
			},
			"addedTags":   []string{"Dart", "Go"},
			"removedTags": []string{"Fortran"},
			"message":     "resource updated successfully",
		}

		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(expectedResponse).
			End()
	})
}

func TestPostResource(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	_, userToken := addTestUser(t, testUser)

	t.Run("cannot create resource with invalid field types", func(t *testing.T) {
		resource := `{
			"Title": 2,
			"Type": [],
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid resource field(s) in request payload"}`).
			End()
	})

	t.Run("cannot create resource with empty title", func(t *testing.T) {
		resource := `{
			"Title": "",
			"Type": "textual",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resource Title is required"}`).
			End()
	})

	t.Run("cannot create resource with empty type", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resource Type is required"}`).
			End()
	})

	t.Run("cannot create resource with invalid type", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "ebook",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resource Type must be one of 'video', 'audio' or 'textual'"}`).
			End()
	})

	t.Run("cannot create resource with empty link", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "textual",
			"Link": "",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resource Link is required"}`).
			End()
	})

	t.Run("cannot create resource with empty privacy", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "textual",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "",
			"Tags": ["python", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"resource Privacy is required"}`).
			End()
	})

	t.Run("cannot create resource with empty tag strings", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "textual",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["", "fortran", "lisp"]
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Tag titles must be non-empty strings"}`).
			End()
	})

	t.Run("cannot create resource with invalid tags", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "textual",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": "lisp"
		}`
		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Tags should be an array of tag titles"}`).
			End()
	})

	t.Run("create resource and return valid response", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "textual",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`

		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(201).
			Expect("Content-Type", "application/json").
			End()
	})

	t.Run("cannot add two resources with same link", func(t *testing.T) {
		resource := `{
			"Title": "A new resource",
			"Type": "textual",
			"Link": "https://localhost.textual/material/6.pdf",
			"Privacy": "public",
			"Tags": ["python", "fortran", "lisp"]
		}`

		Request(testServer.URL, t).
			Post("/api/v1/resource").
			Set("authorization", userToken).
			Send(resource).
			Expect(409).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"A resource exists with provided link"}`).
			End()
	})
}

func TestDeleteResource(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	user, userToken := addTestUser(t, testUser)

	testResource := dummyData["testResource"].(map[string]interface{})
	testResource["userId"] = user.Id
	resource := addTestResource(t, testResource)

	resourceURI := fmt.Sprintf("/api/v1/resource/%v", resource.Id)

	t.Run("cannot delete nonexistent resource", func(t *testing.T) {
		Request(testServer.URL, t).
			Delete("/api/v1/resource/238").
			Set("authorization", userToken).
			Expect(403).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
			End()
	})

	t.Run("cannot delete resource with id 0", func(t *testing.T) {
		Request(testServer.URL, t).
			Delete("/api/v1/resource/0").
			Set("authorization", userToken).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Invalid resource Id in request"}`).
			End()
	})

	t.Run("can only delete own resource", func(t *testing.T) {
		anotherTestUser := dummyData["anotherTestUser"].(map[string]interface{})
		_, userToken := addTestUser(t, anotherTestUser)

		Request(testServer.URL, t).
			Delete(resourceURI).
			Set("authorization", userToken).
			Expect(403).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
			End()
	})

	t.Run("can delete resource", func(t *testing.T) {
		expectedResponse := map[string]interface{}{
			"deletedResource": resource.Id,
			"message":         "Resource deleted successfully",
		}
		Request(testServer.URL, t).
			Delete(resourceURI).
			Set("authorization", userToken).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(expectedResponse).
			End()
	})
}
