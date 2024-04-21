package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"trust-verse-backend/app/dto"
	"trust-verse-backend/app/middlewares"
	"trust-verse-backend/app/models"
)

func CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		return err
	}

	filer := bson.D{{Key: "username", Value: username}}
	userRecord := models.UserCollection.FindOne(c.Context(), filer)
	user := &models.User{}
	userRecord.Decode(user)
	post.User = user
	post.CreatedAt = time.Now()
	_, err = models.PostCollection.InsertOne(c.Context(), post)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": post})
}

func CreateIncidentReport(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		return err
	}

	filer := bson.D{{Key: "username", Value: username}}
	userRecord := models.UserCollection.FindOne(c.Context(), filer)
	user := &models.User{}
	userRecord.Decode(user)
	post.User = user
	post.Type = "incident"
	post.CreatedAt = time.Now()
	result, err := models.PostCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	// Get the ID of the inserted document
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	payload := map[string]string{
		"text": post.Content,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Make a POST request to the sentiment analysis API
	resp, err := http.Post("http://127.0.0.1:5000/sentiments", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the response body into SentimentResponse struct
	var sentimentResponse dto.SentimentResponse
	err = json.Unmarshal(body, &sentimentResponse)
	if err != nil {
		panic(err)
	}
	// Print the parsed sentiment response for demonstration
	fmt.Printf("Sentiment Response: for input %+v\n", sentimentResponse)

	// Make a POST request to the sentiment analysis API
	resp, err = http.Post("http://127.0.0.1:5000/predict", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the response body into SentimentResponse struct
	var fakeNewsResponse dto.FakeNewsResponse
	err = json.Unmarshal(body, &fakeNewsResponse)
	if err != nil {
		panic(err)
	}
	// Print the parsed sentiment response for demonstration
	fmt.Printf("Fake news Response: %+v\n", fakeNewsResponse)

	resp, err = http.Post("http://127.0.0.1:8081/predict-health", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the response body into SentimentResponse struct
	var healthResponse dto.HealthResponse
	err = json.Unmarshal(body, &healthResponse)
	if err != nil {
		panic(err)
	}
	// Print the parsed sentiment response for demonstration
	fmt.Printf("Health Response: %+v\n", healthResponse)

	model := new(models.TrustVerseModel)
	model.UserId = user.ID
	model.PostId = insertedID
	model.RobertaNeu = sentimentResponse.RobertaNeu
	model.RobertaPos = sentimentResponse.RobertaPos
	model.RobertaNeg = sentimentResponse.RobertaNeg
	model.SuicidalScore = healthResponse.PredictionScore
	model.SuicidalLabel = healthResponse.Prediction
	model.IsFake = fakeNewsResponse.Prediction
	model.CreationTime = time.Now()
	_, err = models.TrustVerseModelCollection.InsertOne(c.Context(), model)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": post})
}

func CreatePostComment(c *fiber.Ctx) error {
	var request dto.Comment
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	filter := bson.D{{Key: "username", Value: username}}
	postID, err := primitive.ObjectIDFromHex(request.PostId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	postFilter := bson.D{{Key: "_id", Value: postID}}

	// Attempt to find user and post records
	userRecord := models.UserCollection.FindOne(c.Context(), filter)
	if userRecord.Err() != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "data": "Error finding user: " + userRecord.Err().Error()})
	}

	postRecord := models.PostCollection.FindOne(c.Context(), postFilter)
	if postRecord.Err() != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "data": "Error finding post: " + postRecord.Err().Error()})
	}

	// Decode user and post
	user := &models.User{}
	if err := userRecord.Decode(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "data": "Error decoding user: " + err.Error()})
	}

	post := &models.Post{}
	if err := postRecord.Decode(post); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "data": "Error decoding post: " + err.Error()})
	}

	// Add comment to post
	comment := &models.Comment{}
	comment.UserName = user.Username
	comment.Text = request.Text
	comment.CreatedAt = time.Now()
	post.Comments = append(post.Comments, *comment)

	// Create an update document to set the comments field
	update := bson.M{"$set": bson.M{"comments": post.Comments}}

	// Perform the update operation
	_, err = models.PostCollection.UpdateOne(c.Context(), bson.M{"_id": post.ID}, update)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "messages": "Updated successfully"})

}

func GetPosts(c *fiber.Ctx) error {
	var results []*models.PostWithModel
	filter := bson.M{}
	postType := c.Query("type")
	// If type is not provided, set it to an empty string
	if postType != "" {
		filter = bson.M{"type": postType}
	}
	// Define a filter to find posts by type

	// Find posts in the collection matching the filter
	findOptions := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}})
	cursor, err := models.PostCollection.Find(c.Context(), filter, findOptions)
	defer cursor.Close(c.Context())

	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal Server Error"})
	}

	for cursor.Next(c.Context()) {
		var post models.PostWithModel
		if err := cursor.Decode(&post); err != nil {
			log.Println("Error decoding post:", err)
			continue
		}
		if postType == "incident" {
			var model models.TrustVerseModel

			//err = models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": b.ID}).Decode(&models.TrustVerseModel{})

			err := models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": post.ID}).Decode(&model)
			if err == nil {
				post.Model = model
			}
		}
		results = append(results, &post)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal Server Error"})
	}

	return c.JSON(fiber.Map{"success": true, "data": results})
}

func GetMySentiment(c *fiber.Ctx) error {
	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		fmt.Println("NEW ERRRRRRRRRR")
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	filer := bson.D{{Key: "username", Value: username}}
	fmt.Println("USER NAME IS ", username)
	var user models.User
	err = models.UserCollection.FindOne(c.Context(), filer).Decode(&user)

	var results []*models.PostWithModel
	// If type is not provided, set it to an empty string

	filter := bson.M{}

	// Find posts in the collection matching the filter
	findOptions := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}})
	cursor, err := models.PostCollection.Find(c.Context(), filter, findOptions)
	defer cursor.Close(c.Context())

	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal Server Error"})
	}

	for cursor.Next(c.Context()) {
		var post models.PostWithModel
		if err := cursor.Decode(&post); err != nil {
			log.Println("Error decoding post:", err)
			continue
		}
		var model models.TrustVerseModel

		//err = models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": b.ID}).Decode(&models.TrustVerseModel{})

		err := models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": post.ID}).Decode(&model)
		if err == nil {
			post.Model = model
		}
		results = append(results, &post)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal Server Error"})
	}

	return c.JSON(fiber.Map{"success": true, "data": results})
}

func MyFeed(c *fiber.Ctx) error {
	username, err := middlewares.ExtractUsernameFromToken(c.Get("Authorization"))
	if err != nil {
		fmt.Println("NEW ERRRRRRRRRR")
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}
	filer := bson.D{{Key: "username", Value: username}}
	fmt.Println("USER NAME IS ", username)
	var user models.User
	err = models.UserCollection.FindOne(c.Context(), filer).Decode(&user)

	var sentimentScore float64
	var count int
	// Find posts in the collection matching the filter
	cursor, err := models.TrustVerseModelCollection.Find(c.Context(), bson.M{"userId": user.ID})
	defer cursor.Close(c.Context())
	for cursor.Next(c.Context()) {
		var sentiment models.TrustVerseModel
		if err := cursor.Decode(&sentiment); err != nil {
			log.Println("Error decoding post:", err)
			continue
		}
		sentimentScore += sentiment.RobertaNeg
		count++
	}

	avgSentimentScore := sentimentScore / (float64(count))

	fmt.Printf(" THE AVG SENTIMENT SCORE IS %#v and THE 80 PERCENTAGE OF SENTIMENT SCORE IS %#v", avgSentimentScore, avgSentimentScore*0.7)

	var results []*models.PostWithModel
	filter := bson.M{}
	postType := c.Query("type")
	// If type is not provided, set it to an empty string
	if postType != "" {
		filter = bson.M{"type": postType}
	}
	// Define a filter to find posts by type

	// Find posts in the collection matching the filter
	findOptions := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}})
	cursor, err = models.PostCollection.Find(c.Context(), filter, findOptions)
	defer cursor.Close(c.Context())

	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal Server Error"})
	}

	for cursor.Next(c.Context()) {
		var post models.PostWithModel
		if err := cursor.Decode(&post); err != nil {
			log.Println("Error decoding post:", err)
			continue
		}
		var model models.TrustVerseModel

		err := models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": post.ID}).Decode(&model)
		if err == nil {
			post.Model = model
		}
		if model.RobertaNeg <= avgSentimentScore*0.7 && avgSentimentScore > 0 {
			results = append(results, &post)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal Server Error"})
	}

	return c.JSON(fiber.Map{"success": true, "data": results})
}

func GetPostDetail(c *fiber.Ctx) error {
	postIDParam := c.Params("id")
	if postIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "postId parameter is required",
		})
	}

	postID, err := primitive.ObjectIDFromHex(postIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}
	post := models.Post{}
	// Get the post from MongoDB
	err = models.PostCollection.FindOne(c.Context(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	var trustVerseModel models.TrustVerseModel
	err = models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": postID}).Decode(&trustVerseModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"post": post,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Return the post and TrustVerseModel
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post":            post,
		"trustVerseModel": trustVerseModel,
	})
}

func CreateTrustVerseAIResult(c *fiber.Ctx) error {
	// Extract offset from query parameters, default to 0 if not provided
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(10000)

	cursor, err := models.PostCollection.Find(context.Background(), bson.D{}, findOptions)
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
		var b models.Post // Create a new instance of the Post struct
		err := cursor.Decode(&b)
		if err != nil {
			log.Fatal(err)
		}

		err = models.TrustVerseModelCollection.FindOne(c.Context(), bson.M{"postId": b.ID}).Decode(&models.TrustVerseModel{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Define the request payload
				payload := map[string]string{
					"text": b.Content,
				}

				// Convert the payload to JSON
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					panic(err)
				}

				// Make a POST request to the sentiment analysis API
				resp, err := http.Post("http://127.0.0.1:5000/sentiments", "application/json", bytes.NewBuffer(payloadBytes))
				if err != nil {
					panic(err)
				}

				// Read the response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				// Parse the response body into SentimentResponse struct
				var sentimentResponse dto.SentimentResponse
				err = json.Unmarshal(body, &sentimentResponse)
				if err != nil {
					panic(err)
				}
				fmt.Println("THE INPUT IS " + b.Content)
				// Print the parsed sentiment response for demonstration
				fmt.Printf("Sentiment Response: for input %+v\n", sentimentResponse)

				// Make a POST request to the sentiment analysis API
				resp, err = http.Post("http://127.0.0.1:5000/predict", "application/json", bytes.NewBuffer(payloadBytes))
				if err != nil {
					panic(err)
				}

				// Read the response body
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				// Parse the response body into SentimentResponse struct
				var fakeNewsResponse dto.FakeNewsResponse
				err = json.Unmarshal(body, &fakeNewsResponse)
				if err != nil {
					panic(err)
				}
				// Print the parsed sentiment response for demonstration
				fmt.Printf("Fake news Response: %+v\n", fakeNewsResponse)

				resp, err = http.Post("http://127.0.0.1:8081/predict-health", "application/json", bytes.NewBuffer(payloadBytes))
				if err != nil {
					panic(err)
				}

				// Read the response body
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				// Parse the response body into SentimentResponse struct
				var healthResponse dto.HealthResponse
				err = json.Unmarshal(body, &healthResponse)
				if err != nil {
					panic(err)
				}
				// Print the parsed sentiment response for demonstration
				fmt.Printf("Health Response: %+v\n", healthResponse)

				model := new(models.TrustVerseModel)
				model.UserId = b.User.ID
				model.PostId = b.ID
				model.RobertaNeu = sentimentResponse.RobertaNeu
				model.RobertaPos = sentimentResponse.RobertaPos
				model.RobertaNeg = sentimentResponse.RobertaNeg
				model.SuicidalScore = healthResponse.PredictionScore
				model.SuicidalLabel = healthResponse.Prediction
				model.IsFake = fakeNewsResponse.Prediction
				model.CreationTime = time.Now()
				_, err = models.TrustVerseModelCollection.InsertOne(c.Context(), model)
				if err != nil {
					return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
				}
			}
			if err := cursor.Err(); err != nil {
				log.Fatal(err)
			}
		} else {
			continue
		}
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "data": "Success"})
}
