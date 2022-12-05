package response

import (
	"context"

	"github.com/edmarfelipe/next-u/identity/infra/errors"
	"github.com/gofiber/fiber/v2"
)

func getRequestID(ctx context.Context) string {
	value := ctx.Value("requestid")
	if id, ok := value.(string); ok {
		return id
	}
	return ""
}

type response struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

func SendError(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	message := "Internal Server Error"

	switch e := err.(type) {
	case errors.InvalidInputError:
		message = e.Error()
		statusCode = 400
	case errors.BusinessRuleError:
		message = e.Error()
		statusCode = 422
	case errors.NotAuthorizedError:
		message = e.Error()
		statusCode = 401
	case errors.InsufficientPermissionError:
		message = e.Error()
		statusCode = 403
	case errors.InternalError:
	default:
		statusCode = 500
	}

	resp := response{
		RequestID: getRequestID(c.UserContext()),
		Message:   message,
	}

	return c.Status(statusCode).JSON(resp)
}
