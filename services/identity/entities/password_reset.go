package entities

import "time"

type PasswordReset struct {
	ID       string    `bson:"id"`
	Email    string    `bason:"email"`
	CreateAt time.Time `bson:"createAt"`
	Done     bool      `bson:"done"`
}
