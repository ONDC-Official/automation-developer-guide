package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExchangeCode struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code      string             `json:"code" bson:"code"`
	JWTToken  string             `json:"jwt_token" bson:"jwt_token"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
