package controllers

import (
	"github.com/gin-gonic/gin"
	middlewares "github.com/smartech75/go-microservice-test/user-service/handlers"
)

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}
