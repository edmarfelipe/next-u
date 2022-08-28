package server

import (
	"github.com/edmarfelipe/next-u/services/identity/infra/errors"
	"github.com/edmarfelipe/next-u/services/identity/infra/http/response"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
)

func authMiddleware(jwtToken string) fiber.Handler {
	return jwt.New(jwt.Config{
		SigningKey:   []byte(jwtToken),
		ErrorHandler: authErrorHandle,
	})
}

func authErrorHandle(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return response.SendError(c, errors.NewInvalidInputError("Missing or malformed JWT"))
	}

	return response.SendError(c, errors.NewNotAuthorizedError("user not authorized"))
}
