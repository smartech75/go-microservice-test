package controllers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartech75/go-microservice-test/user-service/db"
	"github.com/smartech75/go-microservice-test/user-service/handlers"
	"github.com/smartech75/go-microservice-test/user-service/models"
	"github.com/smartech75/go-microservice-test/user-service/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var client = db.Dbconnect()

var Test = gin.HandlerFunc(func(c *gin.Context) {
	handlers.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
})

var GetUserByID = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handlers.ErrorResponse("Invalid User ID", c.Writer)
		return
	}

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		handlers.ErrorResponse("User does not exist", c.Writer)
		return
	}

	handlers.SuccessItemRespond(user, "User", c.Writer)
})

func AddNewUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil { //Bad Request
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	if ok, errors := validators.ValidateInputs(user); !ok {
		handlers.ValidationResponse(errors, c.Writer)
		return
	}

	collection := client.Database("userservice").Collection("users")

	user.ID = primitive.NewObjectID()

	var ex_user models.User

	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&ex_user)
	if err == nil {
		handlers.ErrorResponse("User already exist with the username", c.Writer)
		return
	}

	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&ex_user)
	if err == nil {
		handlers.ErrorResponse("User already exist with the email", c.Writer)
		return
	}

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	handlers.SuccessMessageResponse(`Inserted new user at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
}

var DeleteUser = gin.HandlerFunc(func(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		handlers.ErrorResponse("User does not exist", c.Writer)
		return
	}
	_, err = collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	handlers.SuccessMessageResponse("Deleted", c.Writer)
})

func IsAdmin(userid string) (bool, error) {

	id, _ := primitive.ObjectIDFromHex(userid)

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		return false, err
	}

	return user.Type == "admin", nil
}
