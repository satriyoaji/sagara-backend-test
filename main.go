package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/satriyoaji/sagara-backend-test/database"
	"github.com/satriyoaji/sagara-backend-test/helper"
	"github.com/satriyoaji/sagara-backend-test/routes"
)

func main() {
	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":"+helper.GoDotEnvVariable("PORT"))

}