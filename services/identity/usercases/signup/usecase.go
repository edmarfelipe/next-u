package signup

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/libs/common"
	"github.com/edmarfelipe/next-u/services/identity/entities"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecaser interface {
	Execute(ctx context.Context, input Input) error
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
	Name  string `json:"name" validate:"required"`
	User  string `json:"user" validate:"required"`
	Pass  string `json:"pass" validate:"required"`
	Email string `json:"email" validate:"email"`
}

func (usc Usecase) Execute(ctx context.Context, input Input) error {
	err := usc.validator.IsValid(input)
	if err != nil {
		return err
	}

	existingUser, err := usc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("email already in use")
	}

	existingUser, err = usc.userRepository.FindByUser(ctx, input.User)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("user already in use")
	}

	hashedPassword, err := usc.passwordHasher.Hash(input.Pass)
	if err != nil {
		return err
	}

	model := entities.NewUser(
		input.Name,
		input.User,
		hashedPassword,
		input.Email,
	)

	err = usc.userRepository.Create(ctx, *model)
	if err != nil {
		return err
	}

	return nil
}
