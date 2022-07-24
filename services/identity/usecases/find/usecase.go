package find

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Output struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

type Usecase interface {
	Execute(ctx context.Context, input Input) ([]Output, error)
}

type usecase struct {
	userRepository db.UserDB
	validator      infra.Validatorer
}

func NewUsecase(userRepository db.UserDB, validator infra.Validatorer) Usecase {
	return &usecase{
		userRepository: userRepository,
		validator:      validator,
	}
}

type Input struct{}

func (usc usecase) Execute(ctx context.Context, in Input) ([]Output, error) {
	err := usc.validator.IsValid(in)
	if err != nil {
		return nil, err
	}

	users, err := usc.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	outUsers := make([]Output, 0)
	for _, user := range *users {
		outUsers = append(outUsers, Output{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Active:   user.Active,
		})
	}

	return outUsers, nil
}
