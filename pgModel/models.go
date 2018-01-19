package pgModel

type JwtToken struct {
    Token string `json:"token"`
}

type Exception struct {
    Message string `json:"message"`
}

type User struct {
	Id     int64 `sql:",pk,unique"`
	Username   string
	Email string `sql:",unique,notnull"`
	Password string
	PhoneNumber string
	Messages []*Message
	Connections []*Connection
	Comments []*Comment
	Collections []*Collection
	Resources []*Resource
}

type Connection struct {
	Id int64
	FirstUserId int64
	SecondUserId int64
}

type Message struct {
	Id int64
	SenderId int64
	ReceiverId int64
	Content string
}

type Resource struct {
	Id int64
	UserId int64
	CollectionId int64
	Title string
	Link string
	Privacy string
	Type string
	Views int64
	Recommendations int64
	Comments []*Comment
}

type Comment struct {
	Id int64
	UserId int64
	ResourceId int64
	Text string
	Likes int64
}

type Collection struct {
	Id int64
	Name string `sql:",unique,notnull"`
	UserId int64
	Resources []*Resource
}

type Tag struct {
	Id int64
	Title string `sql:",unique,notnull"`
}

type ResourceTag struct {
	TagId int64
	ResourceId int64
}

type CollectionTag struct {
	TagId int64
	CollectionId int64
}

