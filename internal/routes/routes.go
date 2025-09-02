package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/khanjumed/geopackage/internal/config"
	"github.com/khanjumed/geopackage/internal/controllers"
	"github.com/khanjumed/geopackage/internal/middleware"
	"github.com/khanjumed/geopackage/internal/ws"
)

func Register(r *gin.Engine) {
	RegisterAuth(r)

	// Public APIs
	r.GET("/api/parking-spots/nearby", controllers.GetNearbySpots)

	// Orders + Live tracking
	r.POST("/api/order", controllers.UpsertOrder)
	r.GET("/api/order/:id", controllers.GetOrder)
	// Optional REST to get latest location without WS
	r.GET("/api/order/:id/partner-location", controllers.GetLastPartnerLocation)

	// WebSocket: /ws?orderId=ORD123&role=customer|partner[&partnerId=P1]
	r.GET("/ws", ws.HandleWS)

	// Admin APIs
	admin := r.Group("/api/admin")
	if config.EnableAuth {
		admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())
	}
	admin.POST("/add-spot", controllers.AddParkingSpot)
	admin.PUT("/update-spot/:id", controllers.UpdateParkingSpot)
	admin.DELETE("/delete-spot/:id", controllers.DeleteParkingSpot)
}
