package routes

import (
	"log/slog"
	"net/http"

	"github.com/armymp/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

// Request Handlers

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		slog.Error("Failed to retrieve events from database",
			"http_method", context.Request.Method,
			"request_path", context.Request.URL.Path,
			"error_details", err,
		)

		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not retrieve events. Try again later.",
		})
		return
	}

	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	eventID, ok := parseEventID(context)
	if !ok {
		return
	}

	event, ok := getEventOr404(context, eventID)
	if !ok {
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

func updateEvent(context *gin.Context) {
	eventID, ok := parseEventID(context)
	if !ok {
		return
	}

	_, ok = getEventOr404(context, eventID)
	if !ok {
		return
	}

	var updatedEvent models.Event

	err := context.ShouldBindJSON(&updatedEvent)
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

	updatedEvent.ID = eventID
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update event.",
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
	})
}

func deleteEvent(context *gin.Context) {
	eventID, ok := parseEventID(context)
	if !ok { return }

	event, ok := getEventOr404(context, eventID)
	if !ok { return }

	if err := event.Delete(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete event.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully.",
	})
}
