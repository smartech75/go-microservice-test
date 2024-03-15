package middlewares

import (
	"encoding/json"
	"net/http"

	"io"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/color"
	"github.com/smartech75/go-microservice-test/task-management-service/handlers"
	"github.com/smartech75/go-microservice-test/task-management-service/models"
)

func IsAdmin(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userid string

		if userid = c.Request.Header.Get("userid"); userid != "" {
		} else {
			handlers.AuthorizationResponse("Admin Permission Required", c.Writer)
			return
		}

		var user models.User
		user, err := fetchUser(userid)
		if err != nil {
			handlers.AuthorizationResponse("Admin Permission Required", c.Writer)
			return
		}

		if user.Type == "admin" {
			next(c)
			return
		}
		handlers.AuthorizationResponse("Admin Permission Required", c.Writer)
	}
}

func IsValidUser(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userid string
		if userid = c.Request.Header.Get("userid"); userid != "" {
		} else {
			handlers.AuthorizationResponse("Invalid User", c.Writer)
			return
		}

		_, err := fetchUser(userid)
		if err != nil {
			handlers.AuthorizationResponse("Invalid User", c.Writer)
			return
		}

		next(c)
	}
}

func fetchUser(userid string) (models.User, error) {
	var user models.User
	userservice_url := handlers.DotEnvVariable("USERSERVICE_URL")

	resp, err := http.Get(userservice_url + "/api/v1/userservice/user/" + userid)
	if err != nil {
		return models.User{}, err
	}

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return models.User{}, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return models.User{}, err
	}

	data := response["data"].(map[string]interface{})

	user = models.User{
		ID:       data["_id"].(string),
		Username: data["username"].(string),
		Email:    data["email"].(string),
		Type:     data["type"].(string),
	}

	color.Red(user)
	return user, nil
}
