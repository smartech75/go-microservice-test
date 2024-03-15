package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       string `json:"_id" bson:"_id"`
	Username string `json:"username" bson:"username" validate:"required,alphanum"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Type     string `json:"type" bson:"type" validate:"required,oneof=admin normal"`
}

type Task struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	UserID      string             `json:"userid" bson:"userid" validate:"required"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Priority    int                `json:"priority" bson:"priority" default:"0"`
	DueDate     time.Time          `json:"duedate" bson:"duedate" `
	Completed   bool               `json:"completed" bson:"completed" default:"false"`
}

type SearchParams struct {
	UserID    string `json:"userid" bson:"userid"`
	Title     string `json:"title" bson:"title"`
	Priority  string `json:"priority" bson:"priority" validate:"oneof=asce desc"`
	DueDate   string `json:"duedate" bson:"duedate" `
	Completed string `json:"completed" bson:"completed" `
}
