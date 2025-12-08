package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
	server.GET("/events/:id", getEvent)
	server.POST("/events", createEvent)

	server.Run(fmt.Sprintf(":%d", port))
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve events. Try again later."})
		return
	}

	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	eventIDStr := context.Param("id")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		slog.Error("Failed to parse event ID from URL parameter",
					"http_method", context.Request.Method,
					"request_path", context.Request.URL.Path,
					"event_id_string", eventIDStr,
					"error_details", err.Error(), 
		)

		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id. Ensure it is a number."})
		return
	}

	event, err := models.GetEventByID(eventID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("Event ID not found in database",
						"event_id", eventID,
						"request_path", context.Request.URL.Path)
			
			context.JSON(http.StatusNotFound, gin.H{"message": "Could not find event with provided event id"})
			return
		}

		slog.Error("Database error while retrieving event ID",
					"event_id", eventID,
					"error_details", err.Error(),
					"http_method", context.Request.Method)

		context.JSON(http.StatusInternalServerError, gin.H{"message": "An internal server error occurred."})
		return
	}

	context.JSON(http.StatusOK, event)
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

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}
