package model

type User struct {
	Id      string `json:"_id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Picture string `json:"picture" bson:"picture"`
	Email   string `json:"-" bson:"email"`
	Date    string `json:"-" bson:"date"`
}
