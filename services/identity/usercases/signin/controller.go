package signin

import (
	"time"

	"github.com/edmarfelipe/next-u/libs/common"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
	input := Input{
		Username: c.FormValue("user"),
		Password: c.FormValue("pass"),
	}

	output, err := ctrl.usecase.Execute(c.Context(), input)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if output == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"name":  output.Name,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
