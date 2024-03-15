package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/smartech75/go-microservice-test/user-service/controllers"
	"github.com/smartech75/go-microservice-test/user-service/handlers"
)

func IsAuthorized(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userid string
		if userid = c.Request.Header.Get("userid"); userid != "" {
		} else {
			handlers.AuthorizationResponse("Admin Permission Required", c.Writer)
			return
		}

		is_admin, _ := controllers.IsAdmin(userid)

		if is_admin {
			next(c)
			return
		}

		handlers.AuthorizationResponse("Admin Permission Required", c.Writer)
	}
}
