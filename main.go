package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"trust-verse-backend/app/routes"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
	}))
	routes.Setup(app)
	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
