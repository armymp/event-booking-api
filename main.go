package main

import (
	"fmt"
	"net/http"

	"github.com/armymp/event-booking-api/config"
	"github.com/armymp/event-booking-api/db"
	"github.com/armymp/event-booking-api/models"
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


	server.GET("/events", getEvents)
	server.POST("/events", createEvent)
	server.Run(fmt.Sprintf(":%d", port))
}

func getEvents(context *gin.Context) {
	events := models.GetAllEvents()
	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		fmt.Println("BIND ERROR:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	event.ID = 1
	event.UserID = 1000

	event.Save()

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}
