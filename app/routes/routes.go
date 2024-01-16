package routes

import (
	"github.com/gofiber/fiber/v2"
	"trust-verse-backend/app/controllers"
	"trust-verse-backend/app/middlewares"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")
	users := v1.Group("users")
	users.Get("/", middlewares.AuthRequired(), controllers.FetchUsr)
	users.Post("/", controllers.CreateUser)
	users.Post("/profile", middlewares.AuthRequired(), controllers.UpdateUserProfile)
	auth := v1.Group("auths")
	auth.Post("/login", controllers.Login)

	posts := v1.Group("posts")
	posts.Post("/", middlewares.AuthRequired(), controllers.CreatePost)
	posts.Get("/", middlewares.AuthRequired(), controllers.GetPosts)

}
