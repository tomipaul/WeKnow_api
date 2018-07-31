package main_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	. "WeKnow_api/libs/supertest"

	"github.com/subosito/gotenv"
)

func TestMain(m *testing.M) {
	fmt.Println("Loading environment variables...")
	gotenv.Load()

	fmt.Println("Setting up application...")
	setUpApplication()

	fmt.Println("Running tests...")
	code := m.Run()

	query := `DROP TABLE IF EXISTS users, messages, connections,
	comments, resources, collections, tags, resource_tags, 
	collection_tags, user_connections, resource_collections CASCADE`
	if _, err := app.Db.Exec(query); err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(code)
}

func TestUserProfile(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	_, userToken := addTestUser(t, testUser)

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
}

func TestUserPassword(t *testing.T) {
	initializeDatabase(t)
	testServer := httptest.NewServer(app.Router)
	defer closeDatabase(t)
	defer testServer.Close()

	testUser := dummyData["testUser"].(map[string]interface{})
	_, userToken := addTestUser(t, testUser)

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
}
