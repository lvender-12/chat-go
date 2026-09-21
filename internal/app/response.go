package app

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JSON(c fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	} else if err != nil {
		message = err.Error()
	}

	return JSON(c, code, message, nil)
}
