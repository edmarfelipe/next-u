package recover

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
	"github.com/gofiber/fiber/v2"
)

func NewController(config infra.Config) Controller {
	return Controller{
		usecase: NewUsecase(
			db.NewUserRepository(),
			infra.NewMail(config),
			infra.NewValidator(),
		),
	}
}

type Controller struct {
	usecase Usecaser
}

func (ctrl Controller) Handler(c *fiber.Ctx) error {
	var input Input
	err := c.BodyParser(&input)
	if err != nil {
		return err
	}

	err = ctrl.usecase.Execute(c.Context(), input)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusCreated)
}
