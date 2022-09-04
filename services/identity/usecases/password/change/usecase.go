package change

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/validator"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
	"github.com/edmarfelipe/next-u/services/identity/infra/errors"
)

var (
	errUserNotFound     = errors.NewBusinessRuleError("could not found user")
	errPasswordNotMatch = errors.NewBusinessRuleError("password does not match")
)

type Input struct {
	Email       string `json:"email" validate:"required"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

type usecase struct {
	logger       logger.Logger
	userDB       db.UserDB
	passwordHash passwordhash.PasswordHash
}

func NewUsecase(logger logger.Logger, userDB db.UserDB, passwordHash passwordhash.PasswordHash) Usecase {
	return &usecase{
		logger:       logger,
		userDB:       userDB,
		passwordHash: passwordHash,
	}
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	usc.logger.Info(ctx, "Updating password from user "+in.Email)
	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindByEmail(ctx, in.Email, false)
	if err != nil {
		return err
	}

	if user == nil {
		usc.logger.Error(ctx, "failed to change password", "err", errUserNotFound)
		return errUserNotFound
	}

	if !usc.passwordHash.CheckHash(ctx, user.Password, in.OldPassword) {
		return errPasswordNotMatch
	}

	user.Password, err = usc.passwordHash.Hash(ctx, in.NewPassword)
	if err != nil {
		return err
	}

	err = usc.userDB.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
