package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smartech75/go-microservice-test/task-management-service/controllers"
	"github.com/smartech75/go-microservice-test/task-management-service/middlewares"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/taskservice/test", controllers.Test)
	router.GET("/api/v1/taskservice/task/:id", controllers.GetTaskByID)

	router.POST("/api/v1/taskservice/addtask", middlewares.IsAdmin(controllers.AddNewTask))
	router.DELETE("/api/v1/taskservice/deletetask/:id", middlewares.IsAdmin(controllers.DeleteTask))
	router.PUT("/api/v1/taskservice/edittask/:id", middlewares.IsAdmin(controllers.EditTask))

	router.GET("/api/v1/taskservice/search", middlewares.IsValidUser(controllers.SearchTask))
	router.PUT("/api/v1/taskservice/markcompleted/:id", middlewares.IsValidUser(controllers.MarkCompleted))

}
