package enable

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/validator"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
	"github.com/edmarfelipe/next-u/services/identity/infra/errors"
)

var (
	errCouldNotFoundUser = errors.NewBusinessRuleError("could not found user")
)

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

type usecase struct {
	logger logger.Logger
	userDB db.UserDB
}

func NewUsecase(logger logger.Logger, userDB db.UserDB) Usecase {
	return &usecase{
		logger: logger,
		userDB: userDB,
	}
}

type Input struct {
	ID string `validate:"required"`
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	usc.logger.Info(ctx, "Enabling user "+in.ID)

	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindOne(ctx, in.ID)
	if err != nil {
		return err
	}

	if user == nil {
		return errCouldNotFoundUser
	}

	user.Active = true

	err = usc.userDB.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
