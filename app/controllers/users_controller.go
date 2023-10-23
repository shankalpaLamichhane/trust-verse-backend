package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
	"trust-verse-backend/app/dto"
	"trust-verse-backend/app/utils"
)

func GetUser(ctx *fiber.Ctx) error {
	data := map[string]string{"foo": "foo", "bar": "bar"}
	return ctx.JSON(data)
}

func AuthenticateUser() utils.Response {
	claims := dto.Claims{
		Issuer:         strconv.Itoa(1),
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	fmt.Print("token is ", token)
	signedString, err := token.SignedString([]byte("l6xTbagpRDzvB4z5p9j4Nz8VWOGF6rG3JowU8yKsD3Y="))
	fmt.Println("signed string is ", signedString)
	if err != nil {
		fmt.Print("the err is ", err)
		data := map[string]string{"access_token": "foo_token"}
		return utils.Response{
			Code:     fiber.StatusOK,
			Messages: utils.Messages{"Could not generate token"},
			Data:     data,
		}
	}
	return utils.Response{
		Code:     fiber.StatusOK,
		Messages: utils.Messages{"Authentication successful"},
		Data:     signedString,
	}
}
