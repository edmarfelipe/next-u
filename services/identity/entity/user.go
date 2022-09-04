package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleAdmin   = "admin"
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `bson:"name"`
	Password       string             `bson:"password"`
	Email          string             `bson:"email"`
	Active         bool               `bson:"active"`
	Role           string             `bson:"role"`
	PasswordResets []PasswordReset    `bson:"passwordResets"`
}

type PasswordReset struct {
	Token    string    `bson:"token"`
	CreateAt time.Time `bson:"createAt"`
	Done     bool      `bson:"done"`
}

func MakeUser(name string, email string, password string) User {
	return User{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     RoleStudent,
		Active:   true,
	}
}

func (r *PasswordReset) IsExpired() bool {
	maxExpiration := 60 * time.Minute
	tokenExpiration := time.Now().UTC().Sub(r.CreateAt)
	return tokenExpiration-maxExpiration > 0
}

func (u *User) CreatePasswordToken() *PasswordReset {
	newReset := PasswordReset{
		Token:    uuid.New().String(),
		CreateAt: time.Now().UTC(),
	}
	u.PasswordResets = append(u.PasswordResets, newReset)
	return &newReset
}

func (u *User) FindNotDoneToken() *PasswordReset {
	for _, reset := range u.PasswordResets {
		if !reset.Done {
			return &reset
		}
	}
	return nil
}

func (u *User) MarkTokenHasDone(token string) {
	resets := make([]PasswordReset, 0)

	for _, reset := range u.PasswordResets {
		if reset.Token == token {
			reset.Done = true
		}

		resets = append(resets, reset)
	}

	u.PasswordResets = resets
}
