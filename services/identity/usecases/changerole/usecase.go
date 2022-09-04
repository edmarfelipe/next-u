package changerole

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/validator"
	"github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
	"github.com/edmarfelipe/next-u/services/identity/infra/errors"
)

var (
	errCouldNotFoundUser         = errors.NewBusinessRuleError("could not found user")
	errRoleNotValid              = errors.NewInvalidInputError("role not valid")
	errUserDoesNotHavePermission = errors.NewInsufficientPermissionError("user does not have permission")
)

type Input struct {
	ID   string `validate:"required"`
	Role string `validate:"required"`
}

type Output struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type usecase struct {
	logger       logger.Logger
	userDB       db.UserDB
	passwordHash passwordhash.PasswordHash
}

func NewUsecase(logger logger.Logger, userDB db.UserDB, passwordHash passwordhash.PasswordHash) Usecase {
	return &usecase{
		logger:       logger,
		userDB:       userDB,
		passwordHash: passwordHash,
	}
}

func (usc *usecase) Execute(ctx context.Context, in Input) (*Output, error) {
	usc.logger.Info(ctx, "Changing role to "+in.Role+" for user "+in.ID)
	err := validator.IsValid(in)
	if err != nil {
		return nil, err
	}

	if in.Role != entity.RoleAdmin &&
		in.Role != entity.RoleStudent &&
		in.Role != entity.RoleTeacher {
		return nil, errRoleNotValid
	}

	user, err := usc.userDB.FindOne(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errCouldNotFoundUser
	}

	user.Role = in.Role

	err = usc.userDB.Update(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &Output{
		ID:     user.ID.Hex(),
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Active: user.Active,
	}, nil
}
