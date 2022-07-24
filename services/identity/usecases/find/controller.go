package find

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
	in := Input{}
	output, err := ctrl.usecase.Execute(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(output)
}
