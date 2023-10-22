package routes

import (
	"github.com/gofiber/fiber/v2"
	"trust-verse-backend/app/controllers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")
	users := v1.Group("users")
	users.Get("/", controllers.GetUser)
	auth := v1.Group("auths")
	auth.Post("/login", controllers.Login)

}
