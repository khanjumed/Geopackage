// package geopackage

// import (
// 	"log"

// 	"github.com/gin-gonic/gin"
// 	"github.com/khanjumed/geopackage/internal/config" // Import the config package
// 	"github.com/khanjumed/geopackage/internal/routes"
// )

// // RegisterAllRoutes initializes the DB and registers all the routes
// func RegisterAllRoutes(r *gin.Engine) {
// 	// Initialize the DB connectionc
// 	err := config.InitDB()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize DB: %v", err) // Stop the app if DB init fails
// 	}

// 	// Register all the routes
// 	routes.Register(r)
// }

package geopackage

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/khanjumed/geopackage/internal/routes"
)

func RegisterAllRoutes(r *gin.Engine) {
	// ✅ Load HTML templates
	r.LoadHTMLFiles("templates/index.html")

	// ✅ Expose the HTML publicly
	r.GET("/public/index.html", func(c *gin.Context) {
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoogleMapsApiKey": apiKey,
		})
	})

	// ✅ Register all API routes
	routes.Register(r)
}
