package authorize

import (
	"context"

	"github.com/edmarfelipe/next-u/identity/infra/db"
	"github.com/edmarfelipe/next-u/identity/infra/errors"
	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/validator"
)

var (
	errCouldNotFoundUser = errors.NewBusinessRuleError("could not found user")
	errInvalidPassword   = errors.NewBusinessRuleError("invalid password")
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type usecase struct {
	logger         logger.Logger
	userDB         db.UserDB
	passwordHasher passwordhash.PasswordHash
}

func NewUsecase(logger logger.Logger, userRepository db.UserDB, passwordHasher passwordhash.PasswordHash) Usecase {
	return &usecase{
		logger:         logger,
		userDB:         userRepository,
		passwordHasher: passwordHasher,
	}
}

type Input struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Output struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (usc *usecase) Execute(ctx context.Context, in Input) (*Output, error) {
	usc.logger.Info(ctx, "signing with user: "+in.Email)
	err := validator.IsValid(in)
	if err != nil {
		return nil, err
	}

	result, err := usc.userDB.FindByEmail(ctx, in.Email, true)
	if err != nil {
		return nil, err
	}

	if result == nil {
		usc.logger.Error(ctx, "failed to signin", "err", errCouldNotFoundUser)
		return nil, errCouldNotFoundUser
	}

	if !usc.passwordHasher.CheckHash(ctx, result.Password, in.Password) {
		usc.logger.Error(ctx, "failed to signin", "err", errInvalidPassword)
		return nil, errInvalidPassword
	}

	out := Output{
		ID:    result.ID.String(),
		Name:  result.Name,
		Email: result.Email,
	}

	return &out, nil
}
