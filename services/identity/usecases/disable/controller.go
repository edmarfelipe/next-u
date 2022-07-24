package disable

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/gofiber/fiber/v2"
)

func NewController(ct *infra.Container) *Controller {
	return &Controller{
		usecase: NewUsecase(ct.UserDB, ct.Validator),
	}
}

type Controller struct {
	usecase Usecase
}

func (ctrl Controller) Handler(c *fiber.Ctx) error {
	input := Input{
		Username: c.Params("username"),
	}

	err := ctrl.usecase.Execute(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}
