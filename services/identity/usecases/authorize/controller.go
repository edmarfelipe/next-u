package authorize

import (
	"time"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/errors"
	"github.com/edmarfelipe/next-u/services/identity/infra/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func NewController(ct *infra.Container) *Controller {
	return &Controller{
		config:  ct.Config,
		usecase: NewUsecase(ct.Logger, ct.UserDB, ct.PasswordHash),
	}
}

type Controller struct {
	config  *infra.Config
	usecase Usecase
}

var (
	userNotAuthorized = errors.NewNotAuthorizedError("user not authorized")
)

func (ctrl Controller) Handler(c *fiber.Ctx) error {
	var in Input
	err := c.BodyParser(&in)
	if err != nil {
		return err
	}

	output, err := ctrl.usecase.Execute(c.UserContext(), in)
	if err != nil {
		return response.SendError(c, err)
	}

	if output == nil {
		return response.SendError(c, userNotAuthorized)
	}

	claims := jwt.MapClaims{
		"name":  output.Name,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(ctrl.config.Server.JWTToken))
	if err != nil {
		return response.SendError(c, userNotAuthorized)
	}

	return c.JSON(fiber.Map{"token": t})
}
