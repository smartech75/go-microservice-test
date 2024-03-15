package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/color"

	middlewares "github.com/smartech75/go-microservice-test/user-service/handlers"
	"github.com/smartech75/go-microservice-test/user-service/routes"
)

func main() {

	port := middlewares.DotEnvVariable("PORT")

	fmt.Println("Port number is: " + port)
	color.Cyan("🌏 Server running on localhost:" + port)

	router := gin.Default()

	routes.Routes(router)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
