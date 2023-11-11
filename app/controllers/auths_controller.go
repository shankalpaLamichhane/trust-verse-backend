package controllers

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"trust-verse-backend/app/dto"
	"trust-verse-backend/app/utils"
)

func Login(ctx *fiber.Ctx) error {
	response := AuthenticateUser(ctx)
	return ctx.JSON(response)
}

func AuthenticateUser(c *fiber.Ctx) utils.Response {
	const jwtSecret = "asecret"

	var request dto.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Response{
			Code:     fiber.StatusBadRequest,
			Messages: utils.Messages{"Could not parse request"},
			Data:     nil,
		}
	}

	user, err := GetUserByUserName(c, request.Username)
	spew.Dump("user ", user)

	if err != nil {
		return utils.Response{
			Code:     fiber.StatusBadRequest,
			Messages: utils.Messages{"Could not generate token"},
			Data:     "data",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		return utils.Response{
			Code:     fiber.StatusBadRequest,
			Messages: utils.Messages{"Invalid username or password"},
			Data:     "data",
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["username"] = request.Username
	encryptedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return utils.Response{
			Code:     fiber.StatusInternalServerError,
			Messages: utils.Messages{"Could not generate token"},
			Data:     "data",
		}
	}
	return utils.Response{
		Code:     fiber.StatusOK,
		Messages: utils.Messages{"Authentication successful"},
		Data:     encryptedToken,
	}
}
