package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"trust-verse-backend/app/database"
)

// UserCollection | @desc: the user ccollection on the database
var PostCollection *mongo.Collection

// User | @desc: user model struct
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	User      *User              `bson:"user,omitempty"`
	Content   string             `bson:"content,omitempty"`
	Image     string             `bson:"image,omitempty"`
	Comments  []Comment          `bson:"comments,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
	Type      string             `bson:"type,omitempty"`
}

type PostWithModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	User      *User              `bson:"user,omitempty"`
	Content   string             `bson:"content,omitempty"`
	Image     string             `bson:"image,omitempty"`
	Comments  []Comment          `bson:"comments,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
	Type      string             `bson:"type,omitempty"`
	Model     TrustVerseModel    `bson:"model,omitempty"`
}

type Comment struct {
	UserName  string    `bson:"username,omitempty"`
	Text      string    `bson:"text,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
}

// CreatePostSchema | @desc: adds schema validation and indexes to collection
func CreatePostSchema() {
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"content"},
		"properties": bson.M{
			"content": bson.M{
				"bsonType":    "string",
				"description": "must be a string and is required",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	database.DB.CreateCollection(database.Ctx, "posts", options.CreateCollection().SetValidator(validator))

	PostCollection = database.DB.Collection("posts")

	// Make indexes for user collection
	_, _ = PostCollection.Indexes().CreateOne(database.Ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "image", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
