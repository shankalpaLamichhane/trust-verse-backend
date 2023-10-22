package controller

import "trust-verse-backend/app/module/auth/service"

type Controller struct {
	Auth *AuthController
}

func NewController(authService *service.AuthService) *Controller {
	return &Controller{
		Auth: &AuthController{authService: authService},
	}
}
