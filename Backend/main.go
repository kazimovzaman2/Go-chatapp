package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kazimovzaman2/Go-chatapp/database"
	"github.com/kazimovzaman2/Go-chatapp/router"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Chat App",
	})

	database.ConnectDB()
	app.Static("/media/avatars", "./media/avatars")

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}
