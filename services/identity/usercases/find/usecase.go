package find

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/entities"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecaser interface {
	Execute(ctx context.Context, input Input) (*[]entities.User, error)
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

type Input struct{}

func (usc Usecase) Execute(ctx context.Context, input Input) (*[]entities.User, error) {
	err := usc.validator.IsValid(input)
	if err != nil {
		return nil, err
	}

	users, err := usc.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
