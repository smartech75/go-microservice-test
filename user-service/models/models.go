package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username" validate:"required,alphanum"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Type     string             `json:"type" bson:"type" validate:"required,oneof=admin normal"`
}
