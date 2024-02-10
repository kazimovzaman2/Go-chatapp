package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kazimovzaman2/Go-chatapp/config"
	"github.com/kazimovzaman2/Go-chatapp/handler"
	"github.com/kazimovzaman2/Go-chatapp/middleware"
)

func SetupRoutes(app *fiber.App) {
	config, _ := config.LoadConfig(".")
	protected := middleware.NewAuthMiddleware(config.JwtAccessSecret)

	api := app.Group("/api", logger.New())
	api.Get("/hello/", handler.Hello)

	auth := api.Group("/jwt")
	auth.Post("/create/", handler.Login)
	auth.Post("/refresh/", handler.RefreshToken)

	users := api.Group("/users")
	users.Get("/", handler.GetAllUsers)
	users.Post("/", handler.CreateUser)
	users.Get("/me/", protected, handler.GetMe)
	users.Delete("/me/", protected, handler.DeleteMe)
	users.Patch("/me/", protected, handler.UpdateMe)
	users.Get("/:id/", handler.GetUser)
}
