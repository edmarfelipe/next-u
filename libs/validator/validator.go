package validator

import (
	"errors"
	"fmt"

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

		return errors.New(err.Field() + " is " + err.Tag())
	}

	return nil
}
