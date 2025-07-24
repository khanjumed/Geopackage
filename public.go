package Geopackage

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/khanjumed/geopackage/internal/config" // Import the config package
	"github.com/khanjumed/geopackage/internal/routes"
)

// RegisterAllRoutes initializes the DB and registers all the routes
func RegisterAllRoutes(r *gin.Engine) {
	// Initialize the DB connectionc
	err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err) // Stop the app if DB init fails
	}

	// Register all the routes
	routes.Register(r)
}
