package routes

import (
	"log/slog"
	"net/http"

	"github.com/armymp/event-booking-api/models"
	"github.com/armymp/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

// Users Request Handlers

func signup(context *gin.Context) {
	var req models.SignupRequest

	if !bindJSON(context, &req) {
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := user.Save(); err != nil {
		slog.Error("Failed to save user", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create user.",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created!",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func login(context *gin.Context) {
	var req models.LoginRequest

	if !bindJSON(context, &req) {
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := user.ValidateCredentials(); err != nil {
		slog.Warn("Login failed",
			"email", req.Email,
			"ip", context.ClientIP(),
		)

		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password.",
		})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not generate authentication token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"token":   token,
	})
}
