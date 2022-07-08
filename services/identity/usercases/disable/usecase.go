package disable

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecaser interface {
	Execute(ctx context.Context, input Input) error
}

type Usecase struct {
	userRepository db.UserRepositorier
	validator      infra.Validatorer
}

func NewUsecase(userRepository db.UserRepositorier, validator infra.Validatorer) Usecaser {
	return &Usecase{
		userRepository: userRepository,
		validator:      validator,
	}
}

type Input struct {
	ID string `json:"id" validate:"required"`
}

func (usc Usecase) Execute(ctx context.Context, input Input) error {
	err := usc.validator.IsValid(input)
	if err != nil {
		return err
	}

	user, err := usc.userRepository.FindOne(ctx, input.ID)
	if err != nil {
		return err
	}

	user.Active = false

	err = usc.userRepository.Update(ctx, input.ID, *user)
	if err != nil {
		return err
	}

	return nil
}
