package routes

import (
	"gitlab.plainsurf.com/plainsurf/poc/jumed/poc-project/internal/config"
	"gitlab.plainsurf.com/plainsurf/poc/jumed/poc-project/internal/controllers"
	"gitlab.plainsurf.com/plainsurf/poc/jumed/poc-project/internal/middleware"

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
