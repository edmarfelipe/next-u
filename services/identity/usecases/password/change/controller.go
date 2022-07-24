package change

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/gofiber/fiber/v2"
)

func NewController(ct *infra.Container) *Controller {
	return &Controller{
		usecase: NewUsecase(ct.UserDB, ct.PasswordResetDB, ct.Validator, ct.PasswordHash),
	}
}

type Controller struct {
	usecase Usecase
}

func (ctrl Controller) Handler(c *fiber.Ctx) error {
	var input Input
	err := c.BodyParser(&input)
	if err != nil {
		return err
	}

	err = ctrl.usecase.Execute(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
