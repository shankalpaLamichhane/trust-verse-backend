package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"time"
	"trust-verse-backend/app/utils"
)

func GetUser(ctx *fiber.Ctx) error {
	data := map[string]string{"foo": "foo", "bar": "bar"}
	return ctx.JSON(data)
}

func AuthenticateUser() utils.Response {
	//claims := dto.Claims{
	//	Issuer:         strconv.Itoa(1),
	//	StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix()},
	//}

	const jwtSecret = "asecret"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7)
	s, err := token.SignedString([]byte(jwtSecret))
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
		Data:     s,
	}

}
