package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kazimovzaman2/Go-chatapp/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Get("/hello/", handler.Hello)

	users := api.Group("/users")
	users.Get("/", handler.GetAllUsers)
	users.Get("/:id/", handler.GetUser)
	users.Post("/", handler.CreateUser)
}
