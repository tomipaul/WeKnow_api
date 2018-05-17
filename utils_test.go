package main_test

import (
	main "WeKnow_api"
	. "WeKnow_api/model"
	"os"
	"testing"
)

var dummyData = map[string]interface{}{
	"testUser": map[string]interface{}{
		"username":    "test",
		"email":       "test@gmail.com",
		"phoneNumber": "08123425634",
		"password":    "test",
	},
	"anotherTestUser": map[string]interface{}{
		"username":    "anotherTest",
		"email":       "anotherTest@gmail.com",
		"phoneNumber": "08134567901",
		"password":    "anotherTest",
	},
	"testResource": map[string]interface{}{
		"title":   "A new resource",
		"type":    "textual",
		"link":    "https://localhost.textual/material/6.pdf",
		"privacy": "public",
		"tags":    []string{"Python", "Fortran", "Lisp"},
	},
	"collection1": map[string]interface{}{
		"name":   "first collection",
		"userId": 1,
	},
	"collection2": map[string]interface{}{
		"name":   "second collection",
		"userId": 1,
	},
	"naruto": map[string]interface{}{
		"username":    "naruto",
		"email":       "naruto@gmail.com",
		"phoneNumber": "08123425634",
		"password":    "uzumaki",
	},
	"testCollection": map[string]interface{}{
		"name": "new collection",
	},
	"anotherTestCollection": map[string]interface{}{
		"name": "another collection",
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
	query := `DROP TABLE IF EXISTS users, messages, connections,
	comments, resources, collections, tags,
	resource_tags, collection_tags, user_connections`
	if _, err := app.Db.Exec(query); err != nil {
		t.Fatal(err.Error())
	}
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
	app.Db = Connect(dbConfig)
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
		Text:       testData["commentText"].(string),
		UserId:     userId,
		ResourceId: resourceId,
	}
	if err := app.Db.Insert(&comment); err != nil {
		t.Fatal(err.Error())
	}
	return comment
}

func addTestCollection(t *testing.T, testData map[string]interface{}, userId int64) {

	testCollection := Collection{
		Name:   testData["name"].(string),
		UserId: userId,
	}

	if err := app.Db.Insert(&testCollection); err != nil {
		t.Fatal(err.Error())
	}

}
