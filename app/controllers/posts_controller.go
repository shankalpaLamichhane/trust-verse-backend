package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
	"trust-verse-backend/app/models"
)

func CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	_, err := models.PostCollection.InsertOne(c.Context(), post)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": post})

}

func GetPosts(c *fiber.Ctx) error {
	var results []*models.Post
	cursor, err := models.PostCollection.Find(c.Context(), bson.D{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err != nil {
		// Log the error for debugging
		log.Println("Error querying MongoDB:", err)

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	defer cursor.Close(ctx)
	// bson.M
	var b *models.Post
	for cursor.Next(ctx) {
		err := cursor.Decode(&b)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, b)

	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Log the number of documents retrieved (for debugging)
	return c.JSON(fiber.Map{"success": true, "data": results})
}
