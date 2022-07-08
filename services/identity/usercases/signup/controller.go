package signup

import (
	"github.com/edmarfelipe/next-u/libs/common"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
	"github.com/gofiber/fiber/v2"
)

func NewController() Controller {
	return Controller{
		usecase: NewUsecase(
			db.NewUserRepository(),
			infra.NewValidator(),
			common.NewPasswordHash("7be66b96-4cab-436f-850a-e89bd22968d1"),
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
