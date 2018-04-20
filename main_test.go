package main_test

import (
	main "WeKnow_api"
	. "WeKnow_api/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/haoxins/supertest"
	"github.com/subosito/gotenv"
)

var app main.App

func addTestUser(t *testing.T) (User, string) {
	user := User{
		Username:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "08123425634",
		Password:    "test",
	}

	err := app.Db.Insert(&user)

	if err != nil {
		t.Fatal(err.Error())
	}

	token, _ := user.GenerateToken()
	userToken := "Bearer " + token
	return user, userToken
}

func addAnotherTestUser(t *testing.T) (User, string) {
	user := User{
		Username:    "anotherTest",
		Email:       "anotherTest@gmail.com",
		PhoneNumber: "08134567901",
		Password:    "anotherTest",
	}
	err := app.Db.Insert(&user)

	if err != nil {
		t.Fatal(err.Error())
	}

	token, _ := user.GenerateToken()
	userToken := "Bearer " + token
	return user, userToken
}

func setUpApplication() {
	gotenv.Load()
	dbConfig := map[string]string{
		"User":     os.Getenv("TEST_DB_USERNAME"),
		"Password": os.Getenv("TEST_DB_PASSWORD"),
		"Database": os.Getenv("TEST_DATABASE"),
	}
	app = main.CreateApp(dbConfig)
}

func dropAllDatabaseTables(t *testing.T) {
	query := `DROP TABLE IF EXISTS users, messages, connections,
	comments, resources, collections, tags,
	resource_tags, collection_tags, user_connections`
	if _, err := app.Db.Exec(query); err != nil {
		t.Fatal(err.Error())
	}
}

func TestMain(m *testing.M) {
	setUpApplication()

	code := m.Run()

	query := `DROP TABLE IF EXISTS users, messages, connections,
	comments, resources, collections, tags,
	resource_tags, collection_tags, user_connections`
	if _, err := app.Db.Exec(query); err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(code)
}

func TestUserProfile(t *testing.T) {
	setUpApplication()
	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

	_, userToken := addTestUser(t)

	t.Run("cannot update with no username", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/profile").
			Set("authorization", userToken).
			Send(`{"username": ""}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Username cannot be empty"}`).
			End()
	})

	t.Run("cannot update with no phone number", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/profile").
			Set("authorization", userToken).
			Send(`{"phoneNumber": ""}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Phone number cannot be empty"}`).
			End()
	})

	t.Run("cannot update with empty email", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/profile").
			Set("authorization", userToken).
			Send(`{"email": "" }`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error": "Enter a valid email"}`).
			End()
	})

	t.Run("cannot update without valid email", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/profile").
			Set("authorization", userToken).
			Send(`{"email": "testemail" }`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error": "Enter a valid email"}`).
			End()
	})

	t.Run("updates with valid username", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/profile").
			Set("authorization", userToken).
			Send(`{"username": "tester"}`).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(`{"message":"ProfileUpdatedsuccessfully","updatedProfile":{"username":"tester"}}`).
			End()
	})

	t.Run("updates with valid phone number", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/profile").
			Set("authorization", userToken).
			Send(`{"phoneNumber": "09023450022" }`).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(`{"message":"ProfileUpdatedsuccessfully","updatedProfile":{"phoneNumber":"09023450022"}}`).
			End()
	})

	dropAllDatabaseTables(t)
}

func TestUserPassword(t *testing.T) {
	setUpApplication()
	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

	_, userToken := addTestUser(t)

	t.Run("cannot be reset without valid password input", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/password/reset").
			Set("Authorization", userToken).
			Send(`{"password": ""}`).
			Expect(400).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Password is required"}`).
			End()

	})

	t.Run("can be reset", func(t *testing.T) {
		Request(testServer.URL, t).
			Put("/api/v1/user/password/reset").
			Set("Authorization", userToken).
			Send(`{"password": "new password"}`).
			Expect(200).
			Expect("Content-Type", "application/json").
			Expect(`{"message":"Password updated successfully"}`).
			End()

	})

	dropAllDatabaseTables(t)
}

func TestPostResource(t *testing.T) {
	setUpApplication()
	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

	_, userToken := addTestUser(t)

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

	dropAllDatabaseTables(t)
}

func TestUpdateResource(t *testing.T) {
	setUpApplication()
	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

	user, userToken := addTestUser(t)

	resource := Resource{
		Title:   "A new resource",
		Type:    "textual",
		Link:    "https://localhost.textual/material/6.pdf",
		Privacy: "public",
		UserId:  user.Id,
	}
	Tags := []string{"Python", "Fortran", "Lisp"}

	var tags []interface{}
	for _, title := range Tags {
		tags = append(tags, &Tag{Title: title})
	}
	if err := app.Db.Insert(&resource); err != nil {
		t.Log(err.Error())
	}
	if err := app.Db.Insert(tags...); err != nil {
		t.Log(err.Error())
	}
	var resourceTags []interface{}
	for _, tag := range tags {
		resourceTags = append(resourceTags, &ResourceTag{
			TagId:      tag.(*Tag).Id,
			ResourceId: resource.Id,
		})
	}
	if err := app.Db.Insert(resourceTags...); err != nil {
		t.Log(err.Error())
	}
	resourceURI := fmt.Sprintf("/api/v1/resource/%v", resource.Id)

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
			Expect(403).
			Expect("Content-Type", "application/json").
			Expect(`{"error":"Either this resource does not exist or you cannot access it"}`).
			End()
	})

	t.Run("can only update own resource", func(t *testing.T) {
		_, userToken := addAnotherTestUser(t)
		payload := `{
			"title": "An updated resource",
			"type": "audio"
		}`
		Request(testServer.URL, t).
			Put(resourceURI).
			Set("authorization", userToken).
			Send(payload).
			Expect(403).
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

	dropAllDatabaseTables(t)
}
