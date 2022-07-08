package http

import (
	"github.com/edmarfelipe/next-u/services/identity/usercases/disable"
	"github.com/edmarfelipe/next-u/services/identity/usercases/enable"
	"github.com/edmarfelipe/next-u/services/identity/usercases/find"
	"github.com/edmarfelipe/next-u/services/identity/usercases/password/change"
	"github.com/edmarfelipe/next-u/services/identity/usercases/signin"
	"github.com/edmarfelipe/next-u/services/identity/usercases/signup"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var JWT_KEY = []byte("secret")

type Requester interface {
	Handler(c *fiber.Ctx) error
}

func Adapter(ctrl Requester) func(c *fiber.Ctx) error {
	return ctrl.Handler
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()

	base := app.Group("/users/v1")

	base.Post("/sign-in", Adapter(signin.NewController()))
	base.Post("/sign-up", Adapter(signup.NewController()))
	// base.Post("/recover-password", Adapter(recover.NewController()))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: JWT_KEY,
	}))

	base.Get("/", Adapter(find.NewController()))
	base.Post("/change-password", Adapter(change.NewController()))
	base.Patch("/enable/:id", Adapter(enable.NewController()))
	base.Patch("/disable/:id", Adapter(disable.NewController()))

	return app
}
