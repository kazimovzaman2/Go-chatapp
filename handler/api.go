package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kazimovzaman2/Go-jwt-gorm/model"
)

// Hello godoc
// @Summary Hello, World!
// @Description Get Hello, World!
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Router /hello/ [get]
func Hello(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "Hello, World!",
		Data:    nil,
	})
}
