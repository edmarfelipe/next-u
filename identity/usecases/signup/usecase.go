package signup

import (
	"context"

	"github.com/edmarfelipe/next-u/identity/entity"
	"github.com/edmarfelipe/next-u/identity/infra/db"
	"github.com/edmarfelipe/next-u/identity/infra/errors"
	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/validator"
)

var (
	errEmailAlreadyInUse = errors.NewBusinessRuleError("email already in use")
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
	usc.logger.Info(ctx, "creating user with email : "+in.Email)

	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	existingUser, err := usc.userDB.FindByEmail(ctx, in.Email, false)
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
