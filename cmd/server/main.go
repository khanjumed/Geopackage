package main

import (
	"log"
	"net/http"
	"os"
	"parking-service/internal/config"
	"parking-service/internal/routes"

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
	r.LoadHTMLFiles("templates/index.html")

	// Serve index.html with injected API key
	r.GET("/", func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoogleMapsApiKey": apiKey,
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

	r.Run(":8080")
}
