package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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