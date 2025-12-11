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

// Parse and validate event ID from URL
func parseEventID(context *gin.Context) (int64, bool) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("Failed to parse event ID from URL parameter",
			"http_method", context.Request.Method,
			"path", context.Request.URL.Path,
			"id_string", idStr,
			"error", err,
		)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID. Must be a number.",
		})
		return 0, false
	}
	return id, true
}

// Return event or 404
func getEventOr404(context *gin.Context, id int64) (*models.Event, bool) {
	event, err := models.GetEventByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, gin.H{
				"message": "Event not found.",
			})
			return nil, false
		}

		slog.Error("Database error while retrieving event",
			"event_id", id,
			"error", err,
			"http_method", context.Request.Method,
		)

		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error.",
		})
		return nil, false
	}
	return event, true
}