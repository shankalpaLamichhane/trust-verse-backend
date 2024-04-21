package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"trust-verse-backend/app/database"
)

// UserCollection | @desc: the user ccollection on the database
var UserCollection *mongo.Collection

// User | @desc: user model struct
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Username    string             `bson:"username,omitempty"`
	Phone       string             `bson:"phone,omitempty"`
	Email       string             `bson:"email,omitempty"`
	Password    string             `bson:"password,omitempty"`
	UserProfile UserProfile        `bson:"userProfile"`
	FollowedBy  []string           `bson:"followedBy, omitempty"`
	Role        string             `bson:"role, omitempty"`
}

type UserWithSentiment struct {
	*User
	AverageNegSentimentScore float64 `json:"averageNegSentimentScore"`
}

type UserProfile struct {
	ProfileImage string `bson:"profileImage,omitempty"`
}

type LoginRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

// CreateUserSchema | @desc: adds schema validation and indexes to collection
func CreateUserSchema() {
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"name", "email"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "must be a string and is required",
			},
			"email": bson.M{
				"bsonType":    "string",
				"description": "must be a string and is required",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	database.DB.CreateCollection(database.Ctx, "users", options.CreateCollection().SetValidator(validator))

	UserCollection = database.DB.Collection("users")

	// Make indexes for user collection
	_, _ = UserCollection.Indexes().CreateOne(database.Ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

}
