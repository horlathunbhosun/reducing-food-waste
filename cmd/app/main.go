package main

import (
	"github.com/gin-gonic/gin"
	"github.com/horlathunbhosun/wastewarrior-api/internal/database"
	"github.com/horlathunbhosun/wastewarrior-api/internal/routes"
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
