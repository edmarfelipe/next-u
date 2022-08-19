package http

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/usecases/authorize"
	"github.com/edmarfelipe/next-u/services/identity/usecases/disable"
	"github.com/edmarfelipe/next-u/services/identity/usecases/enable"
	"github.com/edmarfelipe/next-u/services/identity/usecases/find"
	"github.com/edmarfelipe/next-u/services/identity/usecases/password/change"
	"github.com/edmarfelipe/next-u/services/identity/usecases/password/changewithtoken"
	"github.com/edmarfelipe/next-u/services/identity/usecases/password/reset"
	"github.com/edmarfelipe/next-u/services/identity/usecases/signup"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwt "github.com/gofiber/jwt/v3"
	"github.com/valyala/fasthttp"
)

type Requester interface {
	Handler(c *fiber.Ctx) error
}

type Server struct {
	fiber *fiber.App
	ct    *infra.Container
}

func NewServer(ct *infra.Container) *Server {
	server := &Server{
		ct:    ct,
		fiber: fiber.New(),
	}

	server.registerRouters()

	return server
}

func (s *Server) registerRouters() {
	base := s.fiber.Group("/identity/v1")

	base.Use(requestid.New())
	base.Use(otelfiber.Middleware(s.ct.Config.Server.Host))

	base.Post("/authorize", s.adapter(authorize.NewController(s.ct)))
	base.Post("/signup", s.adapter(signup.NewController(s.ct)))
	base.Post("/password/reset", s.adapter(reset.NewController(s.ct)))
	base.Use(jwt.New(jwt.Config{
		SigningKey: []byte(s.ct.Config.Server.JWTToken),
	}))

	base.Get("/", s.adapter(find.NewController(s.ct)))
	base.Post("/password/change", s.adapter(change.NewController(s.ct)))
	base.Post("/password/change/:token", s.adapter(changewithtoken.NewController(s.ct)))
	base.Patch("/enable/:id", s.adapter(enable.NewController(s.ct)))
	base.Patch("/disable/:id", s.adapter(disable.NewController(s.ct)))
}

// Handler returns the server handler.
func (s *Server) Handler() fasthttp.RequestHandler {
	return s.fiber.Handler()
}

func (s *Server) Listen() {
	err := s.fiber.Listen(s.ct.Config.ServerAddress())
	if err != nil {
		s.ct.Logger.Error(context.Background(), "err", err)
	}
}
