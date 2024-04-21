package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"trust-verse-backend/app/database"
)

var TrustVerseModelCollection *mongo.Collection

type TrustVerseModel struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PostId        primitive.ObjectID `json:"postId" bson:"postId"`
	UserId        primitive.ObjectID `json:"userId" bson:"userId""`
	RobertaNeg    float64            `json:"roberta_neg" bson:"roberta_neg"`
	RobertaNeu    float64            `json:"roberta_neu" bson:"roberta_neu"`
	RobertaPos    float64            `json:"roberta_pos" bson:"roberta_pos"`
	IsFake        string             `json:"is_fake" bson:"is_fake"`
	SuicidalScore float64            `json:"suicidal_score",bson:"suicidal_score"`
	SuicidalLabel string             `json:"suicidal_label",bson:"suicidal_score"`
	CreationTime  time.Time          `json:"creation_time" bson:"creation_time"`
}

// CreateTrustVerseModel | @desc: adds schema validation and indexes to collection
func CreateTrustVerseModelSchema() {

	database.DB.CreateCollection(database.Ctx, "trustVerseModel", options.CreateCollection())

	TrustVerseModelCollection = database.DB.Collection("trustVerseModel")

	// Make indexes for user collection
	_, _ = PostCollection.Indexes().CreateOne(database.Ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "IsFake", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
