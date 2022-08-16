package reset

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/gofiber/fiber/v2"
)

func NewController(ct *infra.Container) Controller {
	return Controller{
		usecase: NewUsecase(ct.Logger, ct.Config, ct.UserDB, ct.MailService),
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
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}