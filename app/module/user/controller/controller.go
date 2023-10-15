package controller

import "trust-verse-backend/app/module/user/service"

type Controller struct {
	User *UserController
}

func NewController(service *service.UserService) *Controller {
	return &Controller{
		User: &UserController{userService: service},
	}
}
