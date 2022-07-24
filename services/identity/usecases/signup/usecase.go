package signup

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) error
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
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email"`
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	err := usc.validator.IsValid(in)
	if err != nil {
		return err
	}

	existingUser, err := usc.userRepository.FindByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("email already in use")
	}

	existingUser, err = usc.userRepository.FindByUsername(ctx, in.Username)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("user already in use")
	}

	hashedPassword, err := usc.passwordHasher.Hash(in.Password)
	if err != nil {
		return err
	}

	model := entity.MakeUser(
		in.Name,
		in.Username,
		hashedPassword,
		in.Email,
	)

	err = usc.userRepository.Create(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
