package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kazimovzaman2/Go-chatapp/model"
)

func Hello(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "Hello, World!",
		Data:    nil,
	})
}
