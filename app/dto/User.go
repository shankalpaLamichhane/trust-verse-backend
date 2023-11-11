package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDto struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Email       string             `bson:"email,omitempty"`
	UserProfile UserProfile        `bson:"userProfile,omitempty"`
}

type UserProfile struct {
	ProfileImage string `bson:"profileImage,omitempty"`
}
