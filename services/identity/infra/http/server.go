package http

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/tracer"
	"github.com/edmarfelipe/next-u/services/identity/usecases/disable"
	"github.com/edmarfelipe/next-u/services/identity/usecases/enable"
	"github.com/edmarfelipe/next-u/services/identity/usecases/find"
	"github.com/edmarfelipe/next-u/services/identity/usecases/password/change"
	"github.com/edmarfelipe/next-u/services/identity/usecases/password/recover"
	"github.com/edmarfelipe/next-u/services/identity/usecases/signin"
	"github.com/edmarfelipe/next-u/services/identity/usecases/signup"

	"github.com/gofiber/fiber/v2"
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

func (srv *Server) registerRouters() {
	base := srv.fiber.Group("/identity/v1")

	base.Post("/sign-in", srv.adapter(signin.NewController(srv.ct)))
	base.Post("/sign-up", srv.adapter(signup.NewController(srv.ct)))
	base.Post("/recover-password", srv.adapter(recover.NewController(srv.ct)))

	base.Use(jwt.New(jwt.Config{
		SigningKey: []byte(srv.ct.Config.Server.JWTToken),
	}))

	base.Get("/", srv.adapter(find.NewController(srv.ct)))
	base.Post("/change-password", srv.adapter(change.NewController(srv.ct)))
	base.Patch("/enable/:username", srv.adapter(enable.NewController(srv.ct)))
	base.Patch("/disable/:username", srv.adapter(disable.NewController(srv.ct)))
}

func (srv *Server) adapter(ctrl Requester) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.StartSpan(c.Context(), "HTTP", c.Path())
		defer span.End()
		c.SetUserContext(ctx)
		return ctrl.Handler(c)
	}
}

// Handler returns the server handler.
func (srv *Server) Handler() fasthttp.RequestHandler {
	return srv.fiber.Handler()
}

func (srv *Server) Listen() {
	srv.fiber.Listen(srv.ct.Config.ServerAddress())
}
