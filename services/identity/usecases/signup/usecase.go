package signup

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/validator"
	"github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

var (
	errEmailAlreadyInUse = errors.New("email already in use")
	errUserAlreadyInUse  = errors.New("user already in use")
)

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

type usecase struct {
	logger         logger.Logger
	userDB         db.UserDB
	passwordHasher passwordhash.PasswordHash
}

func NewUsecase(logger logger.Logger, userDB db.UserDB, passwordHasher passwordhash.PasswordHash) Usecase {
	return &usecase{
		logger:         logger,
		userDB:         userDB,
		passwordHasher: passwordHasher,
	}
}

type Input struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email"`
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	usc.logger.Info(ctx, "signing up with user: "+in.Email)

	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	existingUser, err := usc.userDB.FindByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		usc.logger.Error(ctx, "failed to signup", "err", errEmailAlreadyInUse)
		return errEmailAlreadyInUse
	}

	hashedPassword, err := usc.passwordHasher.Hash(ctx, in.Password)
	if err != nil {
		return err
	}

	model := entity.MakeUser(
		in.Name,
		in.Email,
		hashedPassword,
	)

	err = usc.userDB.Create(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
