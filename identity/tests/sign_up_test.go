package tests_test

import (
	"net/http"
	"testing"

	"github.com/edmarfelipe/next-u/identity/infra"
	"github.com/edmarfelipe/next-u/identity/usecases/signup"
)

func TestSignUpRouter(t *testing.T) {
	server := NewHTTPTester(t, &infra.Container{})

	baseURL := "/identity/v1"

	t.Run("Should not create a user with a invalid boyd", func(t *testing.T) {
		body := signup.Input{
			Name:     "",
			Email:    "hello@google.com",
			Password: "1234",
		}

		server.POST(baseURL + "/sign-up").
			WithJSON(body).
			Expect().
			Status(http.StatusBadRequest).
			Body().
			Equal("Name is required")
	})

	t.Run("Should create a user with a valid body", func(t *testing.T) {
		body := signup.Input{
			Name:     "Ms. Hello",
			Email:    "hello@google.com",
			Password: "1234",
		}

		server.POST(baseURL + "/sign-up").
			WithJSON(body).
			Expect().
			Status(http.StatusCreated).
			Body().
			Equal("Created")
	})
}
