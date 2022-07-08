package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
	Active   bool               `bson:"active"`
}

func NewUser(name string, username string, password string, email string) *User {
	return &User{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Username: username,
		Password: password,
		Email:    email,
		Active:   true,
	}
}
