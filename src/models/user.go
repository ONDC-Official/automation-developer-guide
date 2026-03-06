package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login     string             `json:"login" bson:"login"`
	Name      string             `json:"name" bson:"name"`
	AvatarURL string             `json:"avatar_url" bson:"avatar_url"`
	Email     string             `json:"email" bson:"email"`
}
