package controller

import (
	"github.com/gofiber/fiber/v2"
	"trust-verse-backend/app/module/user/service"
	"trust-verse-backend/app/utils"
)

type UserController struct {
	userService *service.UserService
}

func (con *UserController) GetUser(c *fiber.Ctx) error {
	user, err := con.userService.GetUser()
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, utils.Response{Messages: utils.Messages{"This will fetch user list!"}, Data: user})
}
