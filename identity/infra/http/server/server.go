package server

import (
	"context"
	"fmt"

	"github.com/edmarfelipe/next-u/identity/infra"
	"github.com/edmarfelipe/next-u/identity/usecases/authorize"
	"github.com/edmarfelipe/next-u/identity/usecases/changerole"
	"github.com/edmarfelipe/next-u/identity/usecases/disable"
	"github.com/edmarfelipe/next-u/identity/usecases/enable"
	"github.com/edmarfelipe/next-u/identity/usecases/find"
	"github.com/edmarfelipe/next-u/identity/usecases/password/change"
	"github.com/edmarfelipe/next-u/identity/usecases/password/changewithtoken"
	"github.com/edmarfelipe/next-u/identity/usecases/password/reset"
	"github.com/edmarfelipe/next-u/identity/usecases/signup"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/valyala/fasthttp"
)

type Requester interface {
	Handler(c *fiber.Ctx) error
}

type server struct {
	fiber *fiber.App
	ct    *infra.Container
}

func New(ct *infra.Container) *server {
	server := &server{
		ct:    ct,
		fiber: fiber.New(),
	}

	server.registerRouters()

	return server
}

func (s *server) registerRouters() {
	base := s.fiber.Group("/identity/v1")

	base.Use(requestid.New())
	base.Use(otelfiber.Middleware(s.ct.Config.Server.Host))

	base.Post("/authorize", s.adapter(authorize.NewController(s.ct)))
	base.Post("/signup", s.adapter(signup.NewController(s.ct)))
	base.Post("/password/reset", s.adapter(reset.NewController(s.ct)))
	base.Use(authMiddleware(s.ct.Config.Server.JWTToken))
	base.Get("/", s.adapter(find.NewController(s.ct)))
	base.Post("/password/change", s.adapter(change.NewController(s.ct)))
	base.Post("/password/change/:token", s.adapter(changewithtoken.NewController(s.ct)))
	base.Patch("/enable/:id", s.adapter(enable.NewController(s.ct)))
	base.Patch("/disable/:id", s.adapter(disable.NewController(s.ct)))
	base.Post("/change-role/:id", s.adapter(changerole.NewController(s.ct)))
}

func (s *server) adapter(ctrl Requester) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s.ct.Logger.Info(c.UserContext(), fmt.Sprintf("%s %s", c.Method(), c.Path()), "http.path", c.Path(), "http.method", c.Method(), "http.statusCode", c.Response().StatusCode())
		return ctrl.Handler(c)
	}
}

// Handler returns the server handler.
func (s *server) Handler() fasthttp.RequestHandler {
	return s.fiber.Handler()
}

func (s *server) Listen() {
	err := s.fiber.Listen(s.ct.Config.ServerAddress())
	if err != nil {
		s.ct.Logger.Error(context.Background(), "err", err)
	}
}
