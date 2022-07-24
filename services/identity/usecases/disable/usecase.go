package disable

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecase interface {
	Execute(ctx context.Context, in Input) error
}

type usecase struct {
	userDB    db.UserDB
	validator infra.Validatorer
}

func NewUsecase(userDB db.UserDB, validator infra.Validatorer) Usecase {
	return &usecase{
		userDB:    userDB,
		validator: validator,
	}
}

type Input struct {
	Username string `validate:"required"`
}

func (usc usecase) Execute(ctx context.Context, in Input) error {
	err := usc.validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindByUsername(ctx, in.Username)
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
