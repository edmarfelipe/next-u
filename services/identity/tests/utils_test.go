package tests_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	httpServer "github.com/edmarfelipe/next-u/services/identity/infra/http"
)

func NewHTTPTester(t *testing.T, ct *infra.Container) *httpexpect.Expect {
	app := httpServer.NewServer(ct)

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(app.Handler()),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
}
