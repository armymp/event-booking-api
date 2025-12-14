package routes

import (
	"log/slog"
	"net/http"

	"github.com/armymp/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

// Users Request Handlers

func signup(context *gin.Context) {
	var user models.User

	if !bindJSON(context, &user) {
		return
	}

	if err := user.Save(); err != nil {
		slog.Error("Failed to save user", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H {
			"message": "Could not create user.",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H {
		"message": "User created!",
		"user": user,
	})
}