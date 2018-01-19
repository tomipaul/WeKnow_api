package mongoModel 


import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Username    string        `bson:"username" json:"username"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password" json:"password"`
}
