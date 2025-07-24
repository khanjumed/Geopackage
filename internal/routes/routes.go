package routes

import (
	"github.com/khanjumed/Geopackage/internal/config"
	"github.com/khanjumed/Geopackage/internal/controllers"
	"github.com/khanjumed/Geopackage/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	RegisterAuth(r)

	r.GET("/api/parking-spots/nearby", controllers.GetNearbySpots)

	admin := r.Group("/api/admin")

	// Apply auth if enabled
	if config.EnableAuth {
		admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())
	}

	admin.POST("/add-spot", controllers.AddParkingSpot)
	admin.PUT("/update-spot/:id", controllers.UpdateParkingSpot)
	admin.DELETE("/delete-spot/:id", controllers.DeleteParkingSpot)
}
