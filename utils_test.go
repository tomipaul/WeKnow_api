package main_test

import (
	main "WeKnow_api"
	. "WeKnow_api/model"
	"WeKnow_api/utilities"
	"fmt"
	"os"
	"testing"
)

var app main.App

var dummyData = map[string]interface{}{
	"testUser": map[string]interface{}{
		"username":    "test",
		"email":       "test@gmail.com",
		"phoneNumber": "08123425634",
		"password":    "test",
	},
	"anotherTestUser": map[string]interface{}{
		"username":    "anotherUser",
		"email":       "anotherUser@gmail.com",
		"phoneNumber": "08134567901",
		"password":    "anotherUser",
	},
	"thirdTestUser": map[string]interface{}{
		"username":    "thirdUser",
		"email":       "thirdUser@gmail.com",
		"phoneNumber": "08100002348",
		"password":    "thirdUser",
	},
	"testResource": map[string]interface{}{
		"title":   "A new resource",
		"type":    "textual",
		"link":    "https://localhost.textual/material/6.pdf",
		"privacy": "public",
		"tags":    []string{"Python", "Fortran", "Lisp"},
	},
	"privateResource": map[string]interface{}{
		"title":   "A private resource",
		"type":    "textual",
		"link":    "https://localhost.textual/material/7.pdf",
		"privacy": "private",
		"tags":    []string{"Golang", "Rust"},
	},
	"followersResource": map[string]interface{}{
		"title":   "A resource for followers only",
		"type":    "audio",
		"link":    "https://localhost.textual/material/8.pdf",
		"privacy": "followers",
		"tags":    []string{"JavaScript", "Scala"},
	},
	"collection1": map[string]interface{}{
		"name": "first collection",
	},
	"collection2": map[string]interface{}{
		"name": "second collection",
	},
	"naruto": map[string]interface{}{
		"username":    "naruto",
		"email":       "naruto@gmail.com",
		"phoneNumber": "08123425634",
		"password":    "uzumaki",
	},
	"testComment1": map[string]interface{}{
		"text": "This is the first comment",
	},
	"testComment2": map[string]interface{}{
		"text": "This is the second comment",
	},
	"testComment3": map[string]interface{}{
		"text": "This is the third comment",
	},
}

func setUpApplication() {
	var dbConfig = map[string]string{
		"User":     os.Getenv("TEST_DB_USERNAME"),
		"Password": os.Getenv("TEST_DB_PASSWORD"),
		"Database": os.Getenv("TEST_DATABASE"),
	}
	app = main.CreateApp(dbConfig)
}

func closeDatabase(t *testing.T) {
	DropSchema(app.Db)
	if err := app.Db.Close(); err != nil {
		t.Fatal(err.Error())
	}
}

func initializeDatabase(t *testing.T) {
	var dbConfig = map[string]string{
		"User":     os.Getenv("TEST_DB_USERNAME"),
		"Password": os.Getenv("TEST_DB_PASSWORD"),
		"Database": os.Getenv("TEST_DATABASE"),
	}
	app.Db = utilities.Connect(dbConfig)
	CreateSchema(app.Db)
}

func addTestUser(t *testing.T, testData map[string]interface{}) (User, string) {
	user := User{
		Username:    testData["username"].(string),
		Email:       testData["email"].(string),
		PhoneNumber: testData["phoneNumber"].(string),
		Password:    testData["password"].(string),
	}

	err := app.Db.Insert(&user)

	if err != nil {
		t.Fatal(err.Error())
	}

	token, _ := user.GenerateToken()
	userToken := "Bearer " + token
	return user, userToken
}

func addTestResource(t *testing.T, testData map[string]interface{}) Resource {
	userId := testData["userId"].(int64)

	resource := Resource{
		Title:   testData["title"].(string),
		Type:    testData["type"].(string),
		Link:    testData["link"].(string),
		Privacy: testData["privacy"].(string),
		UserId:  userId,
	}

	if err := app.Db.Insert(&resource); err != nil {
		t.Fatal(err.Error())
	}

	Tags, OK := testData["tags"].([]string)

	if OK {
		var tags []interface{}
		for _, title := range Tags {
			tags = append(tags, &Tag{Title: title})
		}

		if err := app.Db.Insert(tags...); err != nil {
			t.Fatal(err.Error())
		}
		var resourceTags []interface{}
		for _, tag := range tags {
			resourceTags = append(resourceTags, &ResourceTag{
				TagId:      tag.(*Tag).Id,
				ResourceId: resource.Id,
			})
			resource.Tags = append(resource.Tags, tag.(*Tag))
		}
		if err := app.Db.Insert(resourceTags...); err != nil {
			t.Fatal(err.Error())
		}
	}

	return resource
}

func addTestComment(t *testing.T, testData map[string]interface{}) Comment {
	resourceId := testData["resourceId"].(int64)
	userId := testData["userId"].(int64)

	comment := Comment{
		Text:       testData["text"].(string),
		UserId:     userId,
		ResourceId: resourceId,
	}
	if err := app.Db.Insert(&comment); err != nil {
		t.Fatal(err.Error())
	}
	return comment
}

func addTestCollection(t *testing.T, testData map[string]interface{}) Collection {
	testCollection := Collection{
		Name:   testData["name"].(string),
		UserId: testData["userId"].(int64),
	}

	if err := app.Db.Insert(&testCollection); err != nil {
		t.Fatal(err.Error())
	}

	return testCollection

}

func addTestConnection(t *testing.T, testData map[string]interface{}) {
	testConnection := Connection{
		InitiatorId: testData["initiatorId"].(int64),
		RecipientId: testData["recipientId"].(int64),
	}

	if err := app.Db.Insert(&testConnection); err != nil {
		t.Fatal(err.Error())
	}
}

func customizeEnvVariables(t *testing.T, variables map[string]string) {
	for key, val := range variables {
		if err := os.Setenv(key, val); err != nil {
			t.Fatal(fmt.Sprintf(
				"Failure setting environment variable %v",
				key,
			))
		}
	}
}
