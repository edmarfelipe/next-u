package validator

import (
	"fmt"

	"github.com/edmarfelipe/next-u/identity/infra/errors"
	"github.com/go-playground/validator/v10"
)

// TODO: refact
func IsValid(model interface{}) error {
	val := validator.New()
	err := val.Struct(model)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		fmt.Println(err.Field())
		fmt.Println(err.Tag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())

		return errors.NewInvalidInputError(err.Field() + " is " + err.Tag())
	}

	return nil
}
