package user

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"trust-verse-backend/app/module/user/controller"
	"trust-verse-backend/app/module/user/service"
)

type UserRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

var NewUserModule = fx.Options(
	fx.Provide(service.NewUserService),
	fx.Provide(controller.NewController),
	fx.Provide(NewUserRouter),
)

func NewUserRouter(fiber *fiber.App, controller *controller.Controller) *UserRouter {
	return &UserRouter{
		App:        fiber,
		Controller: controller,
	}
}
func (r *UserRouter) RegisterUserRoutes() {
	userController := r.Controller.User

	r.App.Route("/users", func(router fiber.Router) {
		router.Get("/", userController.GetUser)
	})
}
