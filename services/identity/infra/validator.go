package infra

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validatorer interface {
	IsValid(interface{}) error
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validatorer {
	return &Validator{
		validate: validator.New(),
	}
}

func (val *Validator) IsValid(model interface{}) error {
	err := val.validate.Struct(model)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		fmt.Println(err.Field())
		fmt.Println(err.Tag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())

		return errors.New(err.Field() + " is " + err.Tag())
	}

	return nil
}
