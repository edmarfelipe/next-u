package enable

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type usecase struct {
	userRepository db.UserDB
	validator      infra.Validatorer
}

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

func NewUsecase(userRepository db.UserDB, validator infra.Validatorer) Usecase {
	return &usecase{
		userRepository: userRepository,
		validator:      validator,
	}
}

type Input struct {
	Username string `validate:"required"`
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	err := usc.validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userRepository.FindByUsername(ctx, in.Username)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("could not found user")
	}

	user.Active = true

	err = usc.userRepository.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
