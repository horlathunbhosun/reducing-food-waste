package main

import (
	"github.com/gin-gonic/gin"
	"github.com/horlathunbhosun/reducing-food-waste/database"
	"github.com/horlathunbhosun/reducing-food-waste/routes"
	"log"
)

func main() {
	database.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	err := server.Run(":9090")
	if err != nil {
		log.Fatal(err)
	}
}
