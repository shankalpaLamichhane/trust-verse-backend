package middlewares

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

const jwtSecret = "asecret"

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(jwtSecret),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		var errorList []*fiber.Error
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "Missing or Malformed Authentication Token",
			},
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})

	}

	fmt.Print("THE ERR IS ", err.Error())
	var errorList []*fiber.Error
	errorList = append(
		errorList,
		&fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid or Expired Authentication Token",
		},
	)
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
}

// ExtractUsernameFromToken extracts the username from the JWT token
func ExtractUsernameFromToken(tokenString string) (string, error) {
	//reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, "Bearer ")
	if len(splitToken) < 1 {
		return "", errors.New("UnAuthorized")
	}
	reqToken := splitToken[1]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		fmt.Println("ERROR IS ", err)
	}
	username := fmt.Sprintf("%v", claims["username"])
	return username, nil
}
