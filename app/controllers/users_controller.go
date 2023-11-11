package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"trust-verse-backend/app/dto"
	"trust-verse-backend/app/middlewares"
	"trust-verse-backend/app/models"
	"trust-verse-backend/app/utils"
)

func FetchUsr(c *fiber.Ctx) error {
	user := GetUser(c)
	return c.JSON(user)
}

/*
GetUser | @Desc: Get user by id |
@Method: GET |
@Route: "api/v1/users/:id" |
@Auth: Public
*/
// todo : move this to service and refine the whole architecture of the app.
func GetUser(c *fiber.Ctx) utils.Response {

	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		return utils.Response{
			Code:     fiber.StatusBadRequest,
			Messages: utils.Messages{err.Error()},
			Data:     nil,
		}
	}
	filer := bson.D{{Key: "username", Value: username}}
	userRecord := models.UserCollection.FindOne(c.Context(), filer)
	if userRecord.Err() != nil {
		return utils.Response{
			Code:     fiber.StatusBadRequest,
			Messages: utils.Messages{userRecord.Err()},
			Data:     nil,
		}
	}
	user := &dto.UserDto{}
	userRecord.Decode(user)

	return utils.Response{
		Code:     fiber.StatusOK,
		Messages: utils.Messages{"User fetch successful"},
		Data:     user,
	}
}

func GetUserByUserName(ctx *fiber.Ctx, username string) (*models.User, error) {
	filer := bson.D{{Key: "username", Value: username}}
	userRecord := models.UserCollection.FindOne(ctx.Context(), filer)
	if userRecord.Err() != nil {
		return &models.User{}, userRecord.Err()
	}
	user := &models.User{}
	userRecord.Decode(user)
	return user, nil
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
	// todo : refactor and extract this.
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

/*
UpdateUser | @Desc: Update existing user |
@Method: PUT |
@Route: "api/v1/users/:id" |
@Auth: Private
*/
func UpdateUserProfile(c *fiber.Ctx) error {
	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		fmt.Println("NEW ERRRRRRRRRR")
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	filer := bson.D{{Key: "username", Value: username}}
	fmt.Println("USER NAME IS ", username)
	fmt.Println("-----")

	userRecord := models.UserCollection.FindOne(c.Context(), filer)

	if userRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": userRecord.Err()})
	}

	user := new(models.User)
	fmt.Println("11")
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	// Find the employee and update its data
	query := bson.D{{Key: "username", Value: username}}

	update := bson.D{{"$set", bson.D{{"userProfile", user.UserProfile}}}}

	fmt.Println("22")
	err = models.UserCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	return c.JSON(fiber.Map{"success": true, "data": "User updated successfully."})
}
