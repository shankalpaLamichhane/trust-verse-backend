package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDto struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Username      string             `bson:"username,omitempty"`
	Email         string             `bson:"email,omitempty"`
	UserProfile   UserProfile        `bson:"userProfile,omitempty"`
	Role          string             `bson:"role,omitempty"`
	StreetAddress string             `bson:"streetAddress, omitempty"`
	City          string             `bson:"city, omitempty"`
	State         string             `bson:"state, omitempty"`
	ZipCode       string             `bson:"zipcode,omitempty"`
	About         string             `bson:"about, omitempty"`
}

type UserProfile struct {
	ProfileImage string `bson:"profileImage,omitempty"`
}
