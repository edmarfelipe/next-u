package find

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/validator"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Output struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

type Usecase interface {
	Execute(ctx context.Context, input Input) ([]Output, error)
}

type usecase struct {
	logger logger.Logger
	userDB db.UserDB
}

func NewUsecase(logger logger.Logger, userDB db.UserDB) Usecase {
	return &usecase{
		logger: logger,
		userDB: userDB,
	}
}

type Input struct{}

func (usc usecase) Execute(ctx context.Context, in Input) ([]Output, error) {
	usc.logger.Info(ctx, "Finding users ")
	err := validator.IsValid(in)
	if err != nil {
		return nil, err
	}

	users, err := usc.userDB.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	outUsers := make([]Output, 0)
	for _, user := range users {
		outUsers = append(outUsers, Output{
			ID:     user.ID.Hex(),
			Name:   user.Name,
			Email:  user.Email,
			Active: user.Active,
		})
	}

	return outUsers, nil
}
