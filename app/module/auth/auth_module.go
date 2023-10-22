package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"trust-verse-backend/app/module/auth/controller"
	"trust-verse-backend/app/module/auth/service"
)

type AuthRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

var NewAuthModule = fx.Options(
	fx.Provide(service.NewAuthService),
	fx.Provide(controller.NewController),
	fx.Provide(NewAuthRouter),
)

func NewAuthRouter(fiber *fiber.App, controller *controller.Controller) *AuthRouter {
	return &AuthRouter{
		App:        fiber,
		Controller: controller,
	}
}

func (r *AuthRouter) RegisterAuthRoutes() {
	authController := r.Controller.Auth
	r.App.Route("/auth", func(router fiber.Router) {
		router.Post("/login", authController.Login)
		router.Post("/register", authController.Register)
	})
}
