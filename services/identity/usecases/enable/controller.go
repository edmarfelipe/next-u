package enable

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/http/response"
	"github.com/gofiber/fiber/v2"
)

func NewController(ct *infra.Container) *Controller {
	return &Controller{
		usecase: NewUsecase(ct.Logger, ct.UserDB),
	}
}

type Controller struct {
	usecase Usecase
}

func (ctrl Controller) Handler(c *fiber.Ctx) error {
	in := Input{
		ID: c.Params("id"),
	}

	err := ctrl.usecase.Execute(c.UserContext(), in)
	if err != nil {
		return response.SendError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
