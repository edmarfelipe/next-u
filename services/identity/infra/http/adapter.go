package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) adapter(ctrl Requester) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s.ct.Logger.Info(c.UserContext(), fmt.Sprintf("%s %s", c.Method(), c.Path()), "http.path", c.Path(), "http.method", c.Method(), "http.statusCode", c.Response().StatusCode())
		return ctrl.Handler(c)
	}
}
