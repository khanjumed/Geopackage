package main

import (
	"log"
	"net/http"
	"os"

	"github.com/khanjumed/geopackage/internal/config"
	"github.com/khanjumed/geopackage/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := config.InitDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	r := gin.Default()

	// Load templates from templates/ folder
	r.LoadHTMLGlob("templates/*.html")

	log.Println("Templates loaded successfully.")

	// Serve index.html with injected API key
	r.GET("/", func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoogleMapsApiKey": apiKey,
		})
	})
	// Serve the customer tracking page (track_partner.html)
	r.GET("/track", func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "track_partner.html", gin.H{
			"google_maps_api_key": apiKey, // Pass API key for Google Maps
		})
	})
	// Serve API routes
	routes.Register(r)

	// Fallback for non-matching routes
	r.NoRoute(func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoogleMapsApiKey": apiKey,
		})
	})

	r.Run(":8082")
}
