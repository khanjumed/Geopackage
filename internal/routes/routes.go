// package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/khanjumed/geopackage/internal/config"
// 	"github.com/khanjumed/geopackage/internal/controllers"
// 	"github.com/khanjumed/geopackage/internal/middleware"
// )

// func Register(r *gin.Engine) {
// 	RegisterAuth(r)

// 	// // Route for tracking delivery partner (Customer Page)
// 	// r.GET("/track", func(c *gin.Context) {
// 	// 	c.HTML(200, "track_partner.html", nil) // Serve the tracking page template
// 	// })

// 	// Existing routes
// 	r.GET("/api/parking-spots/nearby", controllers.GetNearbySpots)

// 	admin := r.Group("/api/admin")

// 	// Apply auth if enabled
// 	if config.EnableAuth {
// 		admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())
// 	}

// 	admin.POST("/add-spot", controllers.AddParkingSpot)
// 	admin.PUT("/update-spot/:id", controllers.UpdateParkingSpot)
// 	admin.DELETE("/delete-spot/:id", controllers.DeleteParkingSpot)

// 	// Delivery partner location update route
// 	r.POST("/update-partner-location/:partner_id", controllers.UpdateDeliveryPartnerLocation)

// 	// WebSocket route for customer tracking delivery partner's location
// 	r.GET("/websocket", controllers.WebSocketHandler)
// }

// package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/khanjumed/geopackage/internal/config"
// 	"github.com/khanjumed/geopackage/internal/controllers"
// 	"github.com/khanjumed/geopackage/internal/middleware"
// )

// func Register(r *gin.Engine) {
// 	RegisterAuth(r)

// 	// Existing routes
// 	r.GET("/api/parking-spots/nearby", controllers.GetNearbySpots)

// 	admin := r.Group("/api/admin")

// 	// Apply auth if enabled
// 	if config.EnableAuth {
// 		admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())
// 	}

// 	admin.POST("/add-spot", controllers.AddParkingSpot)
// 	admin.PUT("/update-spot/:id", controllers.UpdateParkingSpot)
// 	admin.DELETE("/delete-spot/:id", controllers.DeleteParkingSpot)

// 	// Delivery partner location update route
// 	r.POST("/update-partner-location/:partner_id", controllers.UpdateDeliveryPartnerLocation)

// 	// WebSocket route for customer tracking delivery partner's location
// 	r.GET("/websocket", controllers.WebSocketHandler)

// 	// New route for sending customer and parking location to delivery partner
// 	r.POST("/send-location-to-delivery-partner/:partner_id", controllers.SendCustomerAndParkingLocationToDeliveryPartner)
// }

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
