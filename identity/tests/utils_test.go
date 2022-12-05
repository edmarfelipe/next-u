package tests_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"github.com/edmarfelipe/next-u/identity/infra"
	"github.com/edmarfelipe/next-u/identity/infra/http/server"
)

func NewHTTPTester(t *testing.T, ct *infra.Container) *httpexpect.Expect {
	app := server.New(ct)

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(app.Handler()),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
}
