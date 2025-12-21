package routes

import (
	"log/slog"
	"net/http"

	"github.com/armymp/event-booking-api/models"
	"github.com/armymp/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

// Events Request Handlers

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
	// extract token (makes function protected so that its only executed if valid token exists)
	token := context.Request.Header.Get("Authorization")
	err := utils.VerifyToken(token)
	if err != nil || token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized",
		})
		return
	}

	var event models.Event
	if !bindJSON(context, &event) {
		return
	}

	// TODO: Remove hardcoded ID and UserID
	event.ID = 1
	event.UserID = 1

	if err := event.Save(); err != nil {
		slog.Error("Failed to save event", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create event.",
		})
		return
	}

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
	if !bindJSON(context, &updatedEvent) {
		return
	}

	updatedEvent.ID = eventID

	if err := updatedEvent.Update(); err != nil {
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
	if !ok {
		return
	}

	event, ok := getEventOr404(context, eventID)
	if !ok {
		return
	}

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
