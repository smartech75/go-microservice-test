package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smartech75/go-microservice-test/user-service/controllers"
	"github.com/smartech75/go-microservice-test/user-service/middlewares"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/userservice/test", controllers.Test)
	router.GET("/api/v1/userservice/user/:id", controllers.GetUserByID)

	router.POST("/api/v1/userservice/adduser", middlewares.IsAuthorized(controllers.AddNewUser))
	router.DELETE("/api/v1/userservice/deleteuser/:id", middlewares.IsAuthorized(controllers.DeleteUser))
}
