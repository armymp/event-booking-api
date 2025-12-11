package routes

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/armymp/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		slog.Error("Failed to retrieve events from database",
			"http_method", context.Request.Method,
			"request_path", context.Request.URL.Path,
			"error_details", err.Error(),
		)

		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not retrieve events. Try again later.",
		})
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

		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id. Ensure it is a number.",
		})
		return
	}

	event, err := models.GetEventByID(eventID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("Event ID not found in database",
				"event_id", eventID,
				"request_path", context.Request.URL.Path)

			context.JSON(http.StatusNotFound, gin.H{
				"message": "Could not find event with provided event id",
			})
			return
		}

		slog.Error("Database error while retrieving event ID",
			"event_id", eventID,
			"error_details", err.Error(),
			"http_method", context.Request.Method)

		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "An internal server error occurred.",
		})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		slog.Warn("Failed to bind requestJSON to event struct",
			"http_method", context.Request.Method,
			"request_path", context.Request.URL.Path,
			"error_details", err.Error(),
		)

		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event data. Ensure all required fields are included.",
		})
		return
	}

	// TODO: Remove hardcoded ID and UserID
	event.ID = 1
	event.UserID = 1000

	err = event.Save()
	if err != nil {
		slog.Error("Database error while saving event",
			"event_name", event.Name,
			"event_location", event.Location,
			"error_details", err.Error(),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create event. Try again later.",
		})
		return
	}

	slog.Info("Event created successfully",
		"event_id", event.ID,
		"event_name", event.Name,
		"user_id", event.UserID,
	)

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created!",
		"event":   event,
	})
}
