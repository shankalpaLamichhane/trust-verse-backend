package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
	"trust-verse-backend/app/models"
	"trust-verse-backend/app/utils"
)

func GetUser(ctx *fiber.Ctx) error {
	data := map[string]string{"foo": "foo", "bar": "bar"}
	return ctx.JSON(data)
}

func AuthenticateUser(ctx *fiber.Ctx) utils.Response {
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

/*
CreateUser | @Desc: Create new user |
@Method: POST |
@Route: "api/v1/users" |
@Auth: Public
*/
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)

	insertionResult, err := models.UserCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	// get the just inserted record in order to return it as response
	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := models.UserCollection.FindOne(c.Context(), filter)

	// decode the Mongo record into Employee
	createdUser := &models.User{}
	createdRecord.Decode(createdUser)

	return c.JSON(fiber.Map{"success": true, "data": createdUser})
}
