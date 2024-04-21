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
	users.Get("/:username", middlewares.AuthRequired(), controllers.GetUserDetail)
	users.Post("/", controllers.CreateUser)
	users.Post("/follow", controllers.FollowUser)
	users.Post("/profile", middlewares.AuthRequired(), controllers.UpdateUserProfile)
	users.Post("/support", middlewares.AuthRequired(), controllers.GetUsersForSupport)
	auth := v1.Group("auths")
	auth.Post("/login", controllers.Login)

	posts := v1.Group("posts")
	posts.Post("/", middlewares.AuthRequired(), controllers.CreatePost)
	posts.Post("/incident", middlewares.AuthRequired(), controllers.CreateIncidentReport)
	posts.Post("/comments", middlewares.AuthRequired(), controllers.CreatePostComment)
	posts.Get("/", middlewares.AuthRequired(), controllers.GetPosts)
	posts.Get("/:id", middlewares.AuthRequired(), controllers.GetPostDetail)
	posts.Get("/user/my-sentiment", middlewares.AuthRequired(), controllers.GetMySentiment)
	posts.Get("/user/my-feed", middlewares.AuthRequired(), controllers.MyFeed)

	storages := v1.Group("storages")
	storages.Post("/", middlewares.AuthRequired(), controllers.UploadFile)
	storages.Post("/profile", middlewares.AuthRequired(), controllers.UploadFileForProfile)

	ai := v1.Group("ai")
	ai.Post("/", middlewares.AuthRequired(), controllers.CreateTrustVerseAIResult)
}
