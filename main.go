package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"trust-verse-backend/app/database"
	"trust-verse-backend/app/models"
	"trust-verse-backend/app/routes"
)

//; git-token=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

func main() {
	godotenv.Load(".env")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
	}))
	routes.Setup(app)
	database.Connect()
	defer database.Cancel()
	defer database.Client.Disconnect(database.Ctx)

	models.CreateUserSchema()
	models.CreatePostSchema()
	models.CreateTrustVerseModelSchema()

	err := app.Listen(":8080")

	if err != nil {
		return
	}
}
