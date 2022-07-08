package tests_test

import (
	"net/http"
	"testing"

	httpServer "github.com/edmarfelipe/next-u/services/identity/infra/http"

	"github.com/gavv/httpexpect"
)

func fiberHTTPTester(t *testing.T) *httpexpect.Expect {
	app := httpServer.SetupHTTPServer()

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(app.Handler()),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func TestSignUpRouter(t *testing.T) {
	server := fiberHTTPTester(t)

	t.Run("Should not create a user with an empty name", func(t *testing.T) {
		baseURL := "/users/v1"
		body := map[string]string{
			"email": "",
			"pass":  "",
			"name":  "",
			"user":  "",
		}

		server.POST(baseURL + "/sign-up").
			WithJSON(body).
			Expect().
			Status(http.StatusBadRequest).
			Body().
			Equal("Bad Request")
	})

	t.Run("Should create a user with a valid body", func(t *testing.T) {
		server := fiberHTTPTester(t)

		baseURL := "/users/v1"
		body := map[string]string{
			"name":  "Ms. Hello",
			"user":  "hello",
			"email": "hello@google.com",
			"pass":  "1234",
		}

		server.POST(baseURL + "/sign-up").
			WithJSON(body).
			Expect().
			Status(http.StatusCreated).
			Body().
			Equal("Created")
	})
}
