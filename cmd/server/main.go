package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/khanjumed/geopackage/internal/config"
	"github.com/khanjumed/geopackage/internal/routes"
)

func main() {
	_ = godotenv.Load()

	// Init MySQL
	if err := config.InitDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	// Init Redis
	if err := config.InitRedis(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer config.Close()

	r := gin.Default()

	// Load templates
	r.LoadHTMLGlob("templates/*.html")
	log.Println("Templates loaded successfully.")

	// Home (customer map)
	r.GET("/", func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoogleMapsApiKey": apiKey,
		})
	})

	// Partner/customer tracking page
	r.GET("/track", func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "track_partner.html", gin.H{
			"google_maps_api_key": apiKey,
		})
	})

	// API routes
	routes.Register(r)

	// Fallback (serve index)
	r.NoRoute(func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoogleMapsApiKey": apiKey,
		})
	})

	// Start server
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
