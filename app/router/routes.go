package router

import (
	"github.com/gofiber/fiber/v2"
	"trust-verse-backend/app/module/user"
)

type Router struct {
	App        fiber.Router
	UserRouter *user.UserRouter
}

func NewRouter(fiber *fiber.App, userRouter *user.UserRouter) *Router {
	return &Router{
		App:        fiber,
		UserRouter: userRouter,
	}
}

func (r *Router) Register() {
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Hello World 👋")
	})
	r.UserRouter.RegisterUserRoutes()
}
