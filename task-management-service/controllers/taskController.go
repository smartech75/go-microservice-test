package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartech75/go-microservice-test/task-management-service/db"
	"github.com/smartech75/go-microservice-test/task-management-service/handlers"
	"github.com/smartech75/go-microservice-test/task-management-service/models"
	"github.com/smartech75/go-microservice-test/task-management-service/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client = db.Dbconnect()

var Test = gin.HandlerFunc(func(c *gin.Context) {
	handlers.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
})

var GetTaskByID = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handlers.ErrorResponse("Invalid Task ID", c.Writer)
		return
	}

	var task models.Task
	collection := client.Database("taskservice").Collection("tasks")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&task)
	if err != nil {
		handlers.ErrorResponse("Task does not exist", c.Writer)
		return
	}

	handlers.SuccessItemRespond(task, "Task", c.Writer)
})

func AddNewTask(c *gin.Context) {
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil { //Bad Request
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	if ok, errors := validators.ValidateInputs(task); !ok {
		handlers.ValidationResponse(errors, c.Writer)
		return
	}

	collection := client.Database("taskservice").Collection("tasks")
	task.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	handlers.SuccessMessageResponse(`Inserted new task at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
}

var DeleteTask = gin.HandlerFunc(func(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var task models.Task
	collection := client.Database("taskservice").Collection("tasks")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&task)
	if err != nil {
		handlers.ErrorResponse("Task does not exist", c.Writer)
		return
	}
	_, err = collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	handlers.SuccessMessageResponse("Deleted", c.Writer)
})

func GetTaskByIDFunc(id string) (models.Task, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, err
	}

	var task models.Task
	collection := client.Database("taskservice").Collection("tasks")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

var EditTask = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var new_task models.Task

	err := c.BindJSON(&new_task)
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	_, err = GetTaskByIDFunc(c.Param("id"))

	if err != nil {
		handlers.ErrorResponse("Task does not exist", c.Writer)
		return
	}

	new_task.ID = id

	update := bson.D{{Key: "$set", Value: new_task}}
	collection := client.Database("taskservice").Collection("tasks")
	_, err = collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	handlers.SuccessMessageResponse("Updated", c.Writer)
})

func SortStringToVal(str string) int {
	switch str {
	case "asce":
		return 1
	case "desc":
		return -1
	default:
		return 0
	}
}

func SearchTask(c *gin.Context) {

	var tasks []*models.Task

	var params models.SearchParams
	err := c.BindJSON(&params)
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	filter := bson.D{
		{Key: "title", Value: bson.D{{"$regex", fmt.Sprintf(".*%s.*", params.Title)}}},
	}

	sort := bson.D{}

	if params.Priority != "" {
		sort = append(sort, bson.E{Key: "priority", Value: SortStringToVal(params.Priority)})
	}
	if params.DueDate != "" {
		sort = append(sort, bson.E{Key: "duedate", Value: SortStringToVal(params.DueDate)})
	}
	if params.Completed != "" {
		sort = append(sort, bson.E{Key: "completed", Value: SortStringToVal(params.Completed)})
	}

	collection := client.Database("taskservice").Collection("tasks")
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetSort(sort))
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var task *models.Task
		err := cursor.Decode(&task)

		if err != nil {
			continue
		}

		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	handlers.SuccessItemsRespond(tasks, "Task", c.Writer)
}

var MarkCompleted = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	task, err := GetTaskByIDFunc(c.Param("id"))
	if err != nil {
		handlers.ErrorResponse("Task does not exist", c.Writer)
		return
	}

	if userid := c.Request.Header.Get("userid"); userid != "" {
		if task.UserID != userid {
			handlers.ErrorResponse("Invalid User to Mark", c.Writer)
			return
		}
	}
	task.Completed = true

	update := bson.D{{Key: "$set", Value: task}}
	collection := client.Database("taskservice").Collection("tasks")
	_, err = collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		handlers.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	handlers.SuccessMessageResponse("Completed", c.Writer)
})
