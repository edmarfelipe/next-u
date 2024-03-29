package changewithtoken

import (
	"context"

	"github.com/edmarfelipe/next-u/identity/infra/db"
	"github.com/edmarfelipe/next-u/identity/infra/errors"
	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/validator"
)

var (
	errUserNotFoundToken = errors.NewBusinessRuleError("could not found user with this token")
	errTokenNotFound     = errors.NewBusinessRuleError("token not found")
	errTokenExpired      = errors.NewBusinessRuleError("token is expired")
)

type Input struct {
	Token       string `json:"-" `
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
	usc.logger.Info(ctx, "Updating password with token: "+in.Token[:5])
	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindByTokenNotDone(ctx, in.Token)
	if err != nil {
		return err
	}

	if user == nil {
		usc.logger.Error(ctx, "failed to change password", "err", errUserNotFoundToken)
		return errUserNotFoundToken
	}

	reset := user.FindNotDoneToken()
	if reset == nil {
		usc.logger.Error(ctx, "failed to change password", "err", errTokenNotFound)
		return errTokenNotFound
	}

	if reset.IsExpired() {
		usc.logger.Error(ctx, "failed to change password", "err", errTokenExpired)
		return errTokenExpired
	}

	user.MarkTokenHasDone(reset.Token)

	hashedPassword, err := usc.passwordHash.Hash(ctx, in.NewPassword)
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
