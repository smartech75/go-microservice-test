package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smartech75/go-microservice-test/user-service/controllers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/userservice/test", controllers.Test)

	router.GET("/api/v1/userservice/user/:id", controllers.GetUserByID)

}
