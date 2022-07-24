package signin

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type usecase struct {
	userRepository db.UserDB
	validator      infra.Validatorer
	passwordHasher passwordhash.PasswordHash
}

func NewUsecase(userRepository db.UserDB, validator infra.Validatorer, passwordHasher passwordhash.PasswordHash) Usecase {
	return &usecase{
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

func (usc *usecase) Execute(ctx context.Context, in Input) (*Output, error) {
	err := usc.validator.IsValid(in)
	if err != nil {
		return nil, err
	}

	result, err := usc.userRepository.FindByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("could not found user")
	}

	if !usc.passwordHasher.CheckHash(result.Password, in.Password) {
		return nil, errors.New("invalid password")
	}

	out := Output{
		Name:     result.Name,
		Username: result.Username,
	}

	return &out, nil
}
