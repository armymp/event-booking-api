package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// Event routes
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)
	server.POST("/events", createEvent)
	server.PUT("/events/:id", updateEvent)
	server.DELETE("/events/:id", deleteEvent)

	// User routes
	server.POST("/signup", signup)
}
