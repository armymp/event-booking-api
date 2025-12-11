package main

import (
	"fmt"

	"github.com/armymp/event-booking-api/config"
	"github.com/armymp/event-booking-api/db"
	"github.com/armymp/event-booking-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	
	db.InitDB()

	port := config.AppConfig.Server.Port
	if port == 0 {
		port = 8080
	}

	fmt.Printf("Server starting in %s mode on port %d\n", config.AppConfig.Server.Env, port)

	server := gin.Default()
	routes.RegisterRoutes(server)
	
	server.Run(fmt.Sprintf(":%d", port))
}

