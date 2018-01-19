package mongoModel


import (
	"log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DAO struct {
	Server   string
	Database string
}

var db *mgo.Database
var (
	UserModel User
)

const (
	COLLECTION = "Users"
)

func (m *DAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *DAO) FindAll() ([]User, error) {
	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func (m *DAO) FindById(id string) (User, error) {
	var user User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (m *DAO) Insert(user User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

func (m *DAO) Delete(user User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

func (m *DAO) Update(user User) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}