package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
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

func GetUserDetail(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	user, err := GetUserByUserName(ctx, username)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	currentUserName, err := middlewares.ExtractUsernameFromToken(ctx.Get("Authorization"))
	follows := false
	for _, followedBy := range user.FollowedBy {
		if followedBy == currentUserName {
			follows = true
			break
		}
	}
	userData := make(map[string]interface{})
	userData["ID"] = user.ID
	userData["Name"] = user.Name
	userData["Username"] = user.Username
	userData["Phone"] = user.Phone
	userData["Email"] = user.Email
	userData["UserProfile"] = user.UserProfile
	userData["Follows"] = follows
	// Return the map as JSON response
	return ctx.JSON(fiber.Map{"success": true, "data": userData})
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

	existingfilter := bson.D{{Key: "username", Value: user.Username}}

	existingUser := models.UserCollection.FindOne(c.Context(), existingfilter)
	if existingUser.Err() == mongo.ErrNoDocuments {
		// User does not exist
		// Continue with registration logic
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
	} else if existingUser.Err() != nil {
		// Handle error
		return c.Status(500).JSON(fiber.Map{"success": false, "data": "Internal Server Error"})
	} else {
		// User already exists
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "The user already exists"})
	}

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

	err = models.UserCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	return c.JSON(fiber.Map{"success": true, "data": "User updated successfully."})
}

func FollowUser(c *fiber.Ctx) error {
	followRequest := new(dto.FollowRequest)

	fmt.Println("11")
	if err := c.BodyParser(followRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		fmt.Println("NEW ERRRRRRRRRR")
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	filer := bson.D{{Key: "username", Value: followRequest.Username}}
	fmt.Println("USER NAME IS ", username)

	userRecord := models.UserCollection.FindOne(c.Context(), filer)

	if userRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": userRecord.Err()})
	}

	// Find the employee and update its data
	var update = bson.D{}
	if username != followRequest.Username {
		update = bson.D{{"$addToSet", bson.D{{"followedBy", username}}}}
	}
	query := bson.D{{Key: "username", Value: followRequest.Username}}

	err = models.UserCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	return c.JSON(fiber.Map{"success": true, "data": "User updated successfully."})
}

func GetUsersForSupport(c *fiber.Ctx) error {
	var usersWithSentiment []*models.UserWithSentiment
	// Extract offset from query parameters, default to 0 if not provided
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(10000)

	cursor, err := models.UserCollection.Find(context.Background(), bson.D{}, findOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err != nil {
		// Log the error for debugging
		log.Println("Error querying MongoDB:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	defer cursor.Close(ctx)
	// Apply offset and limit

	for cursor.Next(ctx) {
		var b models.User // Create a new instance of the User struct
		err := cursor.Decode(&b)

		if err != nil {
			log.Fatal(err)
		}
		modelCursor, err := models.TrustVerseModelCollection.Find(c.Context(), bson.M{"userId": b.ID})

		var totalNegScore float64
		var count int
		var averageNegScore float64

		for modelCursor.Next(ctx) {
			var t models.TrustVerseModel // Create a new instance of the Post struct
			err := modelCursor.Decode(&t)
			if err != nil {
				log.Fatal(err)
			}
			negSentimentScore := t.RobertaNeg
			totalNegScore += negSentimentScore
			count++
		}

		if count > 0 {
			averageNegScore = totalNegScore / float64(count)
			fmt.Println("Average negative sentiment score:", averageNegScore)
		} else {
			fmt.Println("No negative sentiment scores found")
		}
		userWithSentiment := &models.UserWithSentiment{
			User:                     &b,
			AverageNegSentimentScore: averageNegScore,
		}
		if averageNegScore > 0.7 {
			usersWithSentiment = append(usersWithSentiment, userWithSentiment)
		}
	}

	if len(usersWithSentiment) == 0 {
		usersWithSentiment = []*models.UserWithSentiment{}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	// Log the number of documents retrieved (for debugging)
	return c.JSON(fiber.Map{"success": true, "data": usersWithSentiment})
}
