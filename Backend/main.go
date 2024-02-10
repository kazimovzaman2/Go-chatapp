package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kazimovzaman2/Go-chatapp/config"
	"github.com/kazimovzaman2/Go-chatapp/database"
	"github.com/kazimovzaman2/Go-chatapp/router"
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	database.ConnectDB(&config)
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Chat App",
	})

	app.Static("/media/avatars", "./media/avatars")

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}
