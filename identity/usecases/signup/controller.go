package signup

import (
	"github.com/edmarfelipe/next-u/identity/infra"
	"github.com/edmarfelipe/next-u/identity/infra/http/response"
	"github.com/gofiber/fiber/v2"
)

func NewController(ct *infra.Container) *Controller {
	return &Controller{
		usecase: NewUsecase(ct.Logger, ct.UserDB, ct.PasswordHash),
	}
}

type Controller struct {
	usecase Usecase
}

func (ctrl Controller) Handler(c *fiber.Ctx) error {
	var in Input
	err := c.BodyParser(&in)
	if err != nil {
		return err
	}

	err = ctrl.usecase.Execute(c.UserContext(), in)
	if err != nil {
		return response.SendError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(in)
}
