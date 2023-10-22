package controller

import (
	"github.com/gofiber/fiber/v2"
	"trust-verse-backend/app/module/auth/dto"
	"trust-verse-backend/app/module/auth/service"
	"trust-verse-backend/app/utils"
)

type AuthController struct {
	authService *service.AuthService
}

type IAuthController interface {
	Login(c *fiber.Ctx)
	Register(c *fiber.Ctx)
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (con *AuthController) Login(c *fiber.Ctx) error {
	req := new(dto.LoginRequest)
	token, err := con.authService.Login(*req)
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, utils.Response{Messages: utils.Messages{"Successfully logged in!"}, Data: token})
}

func (con *AuthController) Register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	token, err := con.authService.Register(*req)
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, utils.Response{Messages: utils.Messages{"Successfully logged in!"}, Data: token})
}
