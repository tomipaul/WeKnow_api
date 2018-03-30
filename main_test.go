package main_test

import (
	main "WeKnow_api"
	. "WeKnow_api/model"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/haoxins/supertest"
	"github.com/subosito/gotenv"
)

var app main.App

func TestMain(m *testing.M) {
	gotenv.Load()
	dbConfig := map[string]string{
		"User":     os.Getenv("TEST_DB_USERNAME"),
		"Password": os.Getenv("TEST_DB_PASSWORD"),
		"Database": os.Getenv("TEST_DATABASE"),
	}
	app = main.CreateApp(dbConfig)

	query := `DROP TABLE IF EXISTS users, messages, connections,
	comments, resources, collections, tags,
	resource_tags, collection_tags, user_connections`

	code := m.Run()

	if _, err := app.Db.Exec(query); err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(code)
}

func TestUserProfile(t *testing.T) {

	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

	user := User{
		Username:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "08123425634",
		Password:    "test",
	}

	err := app.Db.Insert(&user)

	if err != nil {
		t.Log(err.Error())
	}

	token, _ := user.GenerateToken()
	userToken := "Bearer " + token

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

	if err := app.Db.Delete(&user); err != nil {
		t.Log(err.Error())
	}
}

func TestUserPassword(t *testing.T) {

	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

	user := User{
		Username:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "08123425634",
		Password:    "test",
	}

	err := app.Db.Insert(&user)

	if err != nil {
		t.Log(err.Error())
	}

	token, _ := user.GenerateToken()
	userToken := "Bearer " + token

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

	if err := app.Db.Delete(&user); err != nil {
		t.Log(err.Error())
	}
}

func TestPostResource(t *testing.T) {

	testServer := httptest.NewServer(app.Router)
	defer testServer.Close()

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
