package disable

import (
	"context"

	"github.com/edmarfelipe/next-u/identity/infra/db"
	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/validator"
)

type Usecase interface {
	Execute(ctx context.Context, in Input) error
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

func (usc usecase) Execute(ctx context.Context, in Input) error {
	usc.logger.Info(ctx, "Disabling user "+in.ID)

	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindOne(ctx, in.ID)
	if err != nil {
		return err
	}

	user.Active = false

	err = usc.userDB.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
