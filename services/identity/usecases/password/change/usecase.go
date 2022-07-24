package change

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

var (
	errUserNotFound     = errors.New("could not found user")
	errPasswordNotMatch = errors.New("password does not match")
	errTokenExpired     = errors.New("token is expired")
)

type Input struct {
	Username    string `json:"username" `
	Token       string `json:"token"`
	OldPassword string `json:"oldPassword" validate:"required_with=username"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

type usecase struct {
	userDB          db.UserDB
	passwordResetDB db.PasswordResetDB
	validator       infra.Validatorer
	passwordHash    passwordhash.PasswordHash
}

func NewUsecase(userDB db.UserDB, passwordResetDB db.PasswordResetDB, validator infra.Validatorer, passwordHash passwordhash.PasswordHash) Usecase {
	return &usecase{
		userDB:          userDB,
		passwordResetDB: passwordResetDB,
		validator:       validator,
		passwordHash:    passwordHash,
	}
}

func (usc *usecase) updatePasswordByUsername(ctx context.Context, in Input) error {
	user, err := usc.userDB.FindByUsername(ctx, in.Username)
	if err != nil {
		return err
	}

	if user == nil {
		return errUserNotFound
	}

	if !usc.passwordHash.CheckHash(user.Password, in.OldPassword) {
		return errPasswordNotMatch
	}

	err = usc.updatePassword(ctx, user, in.NewPassword)
	if err != nil {
		return errUserNotFound
	}

	return nil
}

func (usc *usecase) updatePasswordByToken(ctx context.Context, in Input) error {
	reset, err := usc.passwordResetDB.FindByTokenNotDone(ctx, in.Token)
	if err != nil {
		return err
	}

	if reset == nil {
		return errTokenExpired
	}

	maxExpiration := 60 * time.Minute
	tokenExpiration := time.Now().UTC().Sub(reset.CreateAt)

	fmt.Printf("%d \n", tokenExpiration-maxExpiration)
	fmt.Println(tokenExpiration - maxExpiration)

	if tokenExpiration-maxExpiration > 0 {
		return errTokenExpired
	}

	user, err := usc.userDB.FindOne(ctx, reset.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return errUserNotFound
	}

	err = usc.updatePassword(ctx, user, in.NewPassword)
	if err != nil {
		return err
	}

	reset.Done = true

	err = usc.passwordResetDB.Update(ctx, *reset)
	if err != nil {
		return err
	}

	return nil
}

func (usc *usecase) updatePassword(ctx context.Context, user *entity.User, newPassword string) error {
	hashedPassword, err := usc.passwordHash.Hash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = usc.userDB.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	err := usc.validator.IsValid(in)
	if err != nil {
		return err
	}

	if len(in.Token) == 0 {
		return usc.updatePasswordByUsername(ctx, in)
	}

	return usc.updatePasswordByToken(ctx, in)
}
