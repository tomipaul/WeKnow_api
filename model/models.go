package model

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/orm"
	"golang.org/x/crypto/bcrypt"
)

type BaseModel struct {
	CreatedAt *time.Time `json:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty"`
}

func (bm *BaseModel) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if bm.CreatedAt == nil {
		bm.CreatedAt = &now
	}
	if bm.UpdatedAt == nil {
		bm.UpdatedAt = &now
	}
	return nil
}

func (bm *BaseModel) BeforeUpdate(db orm.DB) error {
	now := time.Now()
	bm.UpdatedAt = &now
	return nil
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

type User struct {
	Id          int64         `json:",omitempty"`
	Username    string        `sql:",unique,notnull" json:",omitempty"`
	Email       string        `sql:",unique,notnull" json:",omitempty"`
	Password    string        `json:",omitempty"`
	PhoneNumber string        `json:",omitempty"`
	Connections []*Connection `pg:",many2many:user_connections" json:",omitempty"`
	Comments    []*Comment    `json:",omitempty"`
	Collections []*Collection `json:",omitempty"`
	Resources   []*Resource   `json:",omitempty"`
	BaseModel
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %s>", u.Id, u.Username, u.Email)
}

func (u *User) BeforeInsert(db orm.DB) error {
	if err := u.BaseModel.BeforeInsert(db); err != nil {
		return err
	}
	hashed, error := u.HashPassword()
	if error != nil {
		return error
	}
	u.Password = hashed
	return nil
}

func (u *User) BeforeUpdate(db orm.DB) error {
	if err := u.BaseModel.BeforeUpdate(db); err != nil {
		return err
	}
	hashed, error := u.HashPassword()
	if error != nil {
		return error
	}
	u.Password = hashed
	return nil
}

// HashPassword hash password before storage to database
func (u *User) HashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	return string(bytes), err
}

// CompareHashAndPassword compare stored hash with plain text password
func (u *User) CompareHashAndPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GenerateToken generate authorization token
func (u User) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":      u.Id,
		"username":    u.Username,
		"email":       u.Email,
		"phoneNumber": u.PhoneNumber,
		"iss":         os.Getenv("ISSUER"),
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})
	HMACSecret := os.Getenv("JWT_SECRET")
	tokenString, error := token.SignedString([]byte(HMACSecret))
	if error != nil {
		return "", error
	}
	return tokenString, nil
}

type Connection struct {
	Id          int64  `json:",omitempty"`
	InitiatorId int64  `sql:"unique:connected_users" json:",omitempty"`
	RecipientId int64  `sql:"unique:connected_users" json:",omitempty"`
	Recipient   *User  `json:",omitempty"`
	Initiator   *User  `json:",omitempty"`
	Users       []User `pg:",many2many:user_connections" json:",omitempty"`
	BaseModel
}

type Message struct {
	Id         int64
	Content    string
	Connection *Connection
	BaseModel
}

func (m Message) String() string {
	return fmt.Sprintf("Message<%d %s>", m.Id, m.Content)
}

type Comment struct {
	Id         int64
	UserId     int64     `sql:",notnull"`
	ResourceId int64     `sql:",notnull"`
	Text       string    `sql:",notnull"`
	Likes      int64     `sql:",notnull"`
	Resource   *Resource `json:",omitempty"`
	BaseModel
}

func (c Comment) String() string {
	return fmt.Sprintf("Comment<%d %s %d>", c.Id, c.Text, c.UserId)
}

type Resource struct {
	Id              int64
	UserId          int64      `sql:",notnull" json:",omitempty"`
	Title           string     `sql:",notnull" json:",omitempty"`
	Link            string     `sql:",unique,notnull" json:",omitempty"`
	Privacy         string     `sql:",notnull" json:",omitempty"`
	Type            string     `sql:",notnull" json:",omitempty"`
	Views           int64      `json:",omitempty"`
	Recommendations int64      `json:",omitempty"`
	User            *User      `json:",omitempty"`
	Comments        []*Comment `json:",omitempty"`
	Tags            []*Tag     `pg:",many2many:resource_tags" json:",omitempty"`
	BaseModel
}

func (r Resource) String() string {
	return fmt.Sprintf("Resource<%d %s %s>", r.Id, r.Title, r.Link)
}

type Collection struct {
	Id        int64
	Name      string `sql:",unique,notnull"`
	UserId    int64
	Resources []*Resource
	Tags      []Tag `pg:",many2many:collection_tags"`
	BaseModel
}

func (c Collection) String() string {
	return fmt.Sprintf("Collection<%d %s>", c.Id, c.Name)
}

type Tag struct {
	Id    int64
	Title string `sql:",unique,notnull"`
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag<%d %s>", t.Id, t.Title)
}

type ResourceTag struct {
	TagId      int64 `sql:",pk"`
	ResourceId int64 `sql:",pk"`
}

type CollectionTag struct {
	TagId        int64 `sql:",pk"`
	CollectionId int64 `sql:",pk"`
}

type UserConnection struct {
	UserId       int64 `sql:",pk"`
	ConnectionId int64 `sql:",pk"`
}

type Recommendation struct {
	ResourceId int64 `sql:",pk"`
	UserId     int64 `sql:",pk"`
}

type ResourceCollection struct {
	ResourceId   int64 `sql:",pk"`
	CollectionId int64 `sql:",pk"`
}
