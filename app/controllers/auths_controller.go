package controllers

import "github.com/gofiber/fiber/v2"

func Login(ctx *fiber.Ctx) error {
	response := AuthenticateUser()
	return ctx.JSON(response)
}
