package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordReset struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   string             `bson:"userId"`
	Token    string             `bson:"token"`
	CreateAt time.Time          `bson:"createAt"`
	Done     bool               `bson:"done"`
}

func MakePasswordReset(userID primitive.ObjectID) PasswordReset {
	return PasswordReset{
		ID:       primitive.NewObjectID(),
		UserID:   userID.Hex(),
		Token:    uuid.New().String(),
		CreateAt: time.Now().UTC(),
	}
}
