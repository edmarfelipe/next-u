package signin

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/libs/common"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecaser interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type Usecase struct {
	userRepository db.UserRepositorier
	validator      infra.Validatorer
	passwordHasher common.PasswordHasher
}

func NewUsecase(userRepository db.UserRepositorier, validator infra.Validatorer, passwordHasher common.PasswordHasher) Usecaser {
	return &Usecase{
		userRepository: userRepository,
		validator:      validator,
		passwordHasher: passwordHasher,
	}
}

type Input struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Output struct {
	Name     string
	Username string
}

func (usc Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	err := usc.validator.IsValid(input)
	if err != nil {
		return nil, err
	}

	result, err := usc.userRepository.FindByUser(ctx, input.Username)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("could not found user")
	}

	if !usc.passwordHasher.CheckHash(result.Password, input.Password) {
		return nil, errors.New("invalid password")
	}

	output := Output{
		Name:     result.Name,
		Username: result.Username,
	}

	return &output, nil
}
