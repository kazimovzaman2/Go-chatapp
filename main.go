package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kazimovzaman2/Go-jwt-gorm/config"
	"github.com/kazimovzaman2/Go-jwt-gorm/database"
	_ "github.com/kazimovzaman2/Go-jwt-gorm/docs"
	"github.com/kazimovzaman2/Go-jwt-gorm/router"
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	database.ConnectDB(&config)
}

// @title App API
// @version 1.0
// @description This is a simple app API

// @contact.name Zaman Kazimov
// @contact.email kazimovzaman2@gmail.com

// @license.name GPlv3
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @BasePath /api
// @host localhost:8000
func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}))

	app.Static("/media/avatars", "./media/avatars")

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}
