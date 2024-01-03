package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/horlathunbhosun/wastewarrior-api/pkg/response"
	"net/http"
)

func RegisterRoutes(server *gin.Engine) {

	var responseBody response.JsonResponse
	server.NoRoute(func(c *gin.Context) {
		responseBody.Error = true
		responseBody.Message = "Route not found"
		responseBody.Status = false
		c.JSON(http.StatusNotFound, responseBody)
	})

	server.NoMethod(func(c *gin.Context) {
		responseBody.Error = true
		responseBody.Message = "Method not allowed"
		responseBody.Status = false
		c.JSON(http.StatusMethodNotAllowed, responseBody)
	})

	v1 := server.Group("/v1")
	v1.GET("/", func(c *gin.Context) {
		responseBody.Message = "Welcome to Waste Warrior API"
		responseBody.Status = true
		c.JSON(http.StatusOK, responseBody)
	})
	v1.GET("/products")
}
