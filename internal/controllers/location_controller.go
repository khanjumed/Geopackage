// package controllers

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"github.com/khanjumed/geopackage/internal/models"
// )

// // In-memory storage for partner locations (replace with a database if needed)
// var deliveryPartnerLocations = make(map[string]models.Location)

// // WebSocket clients (connected customers)
// var clients = make(map[*websocket.Conn]bool)

// // Handle WebSocket connection for customers
// func WebSocketHandler(c *gin.Context) {
// 	upgrader := websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket Upgrade Error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	clients[conn] = true
// 	defer delete(clients, conn)

// 	// Keep the WebSocket connection open and handle incoming messages
// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			break
// 		}
// 	}
// }

// // Update delivery partner location (POST)
// func UpdateDeliveryPartnerLocation(c *gin.Context) {
// 	var location models.Location
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
// 		return
// 	}

// 	// Assuming partner_id is passed as a URL parameter
// 	partnerID := c.Param("partner_id")

// 	// Store the updated location in memory (or a database)
// 	deliveryPartnerLocations[partnerID] = location

// 	// Broadcast the updated location to all connected customers via WebSocket
// 	SendPartnerLocation(partnerID, location)

// 	c.JSON(http.StatusOK, gin.H{"message": "Location updated"})
// }

// // Broadcast partner's location to all WebSocket clients (customers)
// func SendPartnerLocation(partnerID string, location models.Location) {
// 	for client := range clients {
// 		err := client.WriteJSON(map[string]interface{}{
// 			"partner_id": partnerID,
// 			"lat":        location.Lat,
// 			"lng":        location.Lng,
// 		})
// 		if err != nil {
// 			log.Println("Error sending location to client:", err)
// 			client.Close()
// 			delete(clients, client)
// 		}
// 	}
// }

// package controllers

// import (
// 	"log"
// 	"net/http"
// 	"sync" // Importing sync package for Mutex

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"github.com/khanjumed/geopackage/internal/models"
// )

// // In-memory storage for partner locations (replace with a database if needed)
// var deliveryPartnerLocations = make(map[string]models.Location)

// // WebSocket clients (connected customers)
// var clients = make(map[*websocket.Conn]bool)

// // Mutex to synchronize WebSocket writes
// var mu sync.Mutex

// // Handle WebSocket connection for customers
// func WebSocketHandler(c *gin.Context) {
// 	upgrader := websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket Upgrade Error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	clients[conn] = true
// 	defer delete(clients, conn)

// 	// Keep the WebSocket connection open and handle incoming messages
// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			break
// 		}
// 	}
// }

// // Update delivery partner location (POST)
// func UpdateDeliveryPartnerLocation(c *gin.Context) {
// 	var location models.Location
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
// 		return
// 	}

// 	// Assuming partner_id is passed as a URL parameter
// 	partnerID := c.Param("partner_id")

// 	// Store the updated location in memory (or a database)
// 	deliveryPartnerLocations[partnerID] = location

// 	// Broadcast the updated location to all connected customers via WebSocket
// 	SendPartnerLocation(partnerID, location)

// 	c.JSON(http.StatusOK, gin.H{"message": "Location updated"})
// }

// // Broadcast partner's location to all WebSocket clients (customers) with Mutex synchronization
// func SendPartnerLocation(partnerID string, location models.Location) {
// 	// Lock the mutex to ensure only one WebSocket write happens at a time
// 	mu.Lock()
// 	defer mu.Unlock()

// 	for client := range clients {
// 		err := client.WriteJSON(map[string]interface{}{
// 			"partner_id": partnerID,
// 			"lat":        location.Lat,
// 			"lng":        location.Lng,
// 		})
// 		if err != nil {
// 			log.Println("Error sending location to client:", err)
// 			client.Close()
// 			delete(clients, client)
// 		}
// 	}
// }

// package controllers

// import (
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"github.com/khanjumed/geopackage/internal/models"
// )

// // In-memory partner latest locations (optional)
// var deliveryPartnerLocations = make(map[string]models.Location)

// // Customers viewing map (legacy). Keep if you still broadcast to all customers.
// var customers = make(map[*websocket.Conn]bool)

// // NEW: delivery partners grouped by partner_id
// // (allow multiple sockets per partner: mobile + web, etc.)
// var partnerConns = make(map[string]map[*websocket.Conn]bool)

// var mu sync.Mutex

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// // WebSocket endpoint for both roles:
// //
// //	/websocket?role=partner&partner_id=123
// //	/websocket?role=customer
// func WebSocketHandler(c *gin.Context) {
// 	role := c.Query("role")
// 	partnerID := c.Query("partner_id")

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket Upgrade Error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	mu.Lock()
// 	if role == "partner" && partnerID != "" {
// 		if partnerConns[partnerID] == nil {
// 			partnerConns[partnerID] = make(map[*websocket.Conn]bool)
// 		}
// 		partnerConns[partnerID][conn] = true
// 	} else {
// 		// default: treat as customer
// 		customers[conn] = true
// 	}
// 	mu.Unlock()

// 	// Simple read loop to keep connection alive
// 	for {
// 		if _, _, err := conn.ReadMessage(); err != nil {
// 			break
// 		}
// 	}

// 	// Cleanup on disconnect
// 	mu.Lock()
// 	if role == "partner" && partnerID != "" {
// 		if set, ok := partnerConns[partnerID]; ok {
// 			delete(set, conn)
// 			if len(set) == 0 {
// 				delete(partnerConns, partnerID)
// 			}
// 		}
// 	} else {
// 		delete(customers, conn)
// 	}
// 	mu.Unlock()
// }

// // POST /update-partner-location/:partner_id
// // (partner app pings this; we may broadcast to customers)
// func UpdateDeliveryPartnerLocation(c *gin.Context) {
// 	var location models.Location
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
// 		return
// 	}
// 	partnerID := c.Param("partner_id")

// 	// store last known
// 	mu.Lock()
// 	deliveryPartnerLocations[partnerID] = location
// 	mu.Unlock()

// 	// Optional: broadcast to ALL customers (legacy behavior)
// 	broadcastPartnerLocationToCustomers(partnerID, location)

// 	c.JSON(http.StatusOK, gin.H{"message": "Location updated"})
// }

// func broadcastPartnerLocationToCustomers(partnerID string, location models.Location) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	for cli := range customers {
// 		if err := cli.WriteJSON(map[string]interface{}{
// 			"type":       "partner.location",
// 			"partner_id": partnerID,
// 			"lat":        location.Lat,
// 			"lng":        location.Lng,
// 		}); err != nil {
// 			log.Println("Error sending to customer:", err)
// 			cli.Close()
// 			delete(customers, cli)
// 		}
// 	}
// }

// // NEW: POST /send-location-to-delivery-partner/:partner_id
// // Body: { "customer_location": {...}, "parking_location": {...} }
// func SendCustomerAndParkingLocationToDeliveryPartner(c *gin.Context) {
// 	var payload models.AssignmentPayload
// 	if err := c.ShouldBindJSON(&payload); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
// 		return
// 	}
// 	partnerID := c.Param("partner_id")

// 	msg := models.WSAssignmentMessage{
// 		Type:     "assignment",
// 		Customer: payload.Customer,
// 		Parking:  payload.Parking,
// 	}

// 	// Targeted send to the specific partner
// 	if !sendToPartner(partnerID, msg) {
// 		// No WS connections for this partner (e.g., offline). Decide your behavior:
// 		c.JSON(http.StatusAccepted, gin.H{"message": "Partner offline; assignment queued or logged"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Locations sent to delivery partner"})
// }

// func sendToPartner(partnerID string, v interface{}) bool {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	set, ok := partnerConns[partnerID]
// 	if !ok || len(set) == 0 {
// 		return false
// 	}
// 	for cli := range set {
// 		if err := cli.WriteJSON(v); err != nil {
// 			log.Println("Error sending to partner:", err)
// 			cli.Close()
// 			delete(set, cli)
// 		}
// 	}
// 	return true
// }

// package controllers

// import (
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"github.com/khanjumed/geopackage/internal/models"
// )

// // In-memory partner latest locations (optional)
// var deliveryPartnerLocations = make(map[string]models.Location)

// // Customers viewing map (legacy). Keep if you still broadcast to all customers.
// var customers = make(map[*websocket.Conn]bool)

// // NEW: delivery partners grouped by partner_id
// // (allow multiple sockets per partner: mobile + web, etc.)
// var partnerConns = make(map[string]map[*websocket.Conn]bool)

// var mu sync.Mutex

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// // WebSocket endpoint for both roles:
// //
// //	/websocket?role=partner&partner_id=123
// //	/websocket?role=customer
// func WebSocketHandler(c *gin.Context) {
// 	role := c.Query("role")
// 	partnerID := c.Query("partner_id")

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket Upgrade Error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	mu.Lock()
// 	if role == "partner" && partnerID != "" {
// 		if partnerConns[partnerID] == nil {
// 			partnerConns[partnerID] = make(map[*websocket.Conn]bool)
// 		}
// 		partnerConns[partnerID][conn] = true
// 	} else {
// 		// default: treat as customer
// 		customers[conn] = true
// 	}
// 	mu.Unlock()

// 	// Simple read loop to keep connection alive
// 	for {
// 		if _, _, err := conn.ReadMessage(); err != nil {
// 			break
// 		}
// 	}

// 	// Cleanup on disconnect
// 	mu.Lock()
// 	if role == "partner" && partnerID != "" {
// 		if set, ok := partnerConns[partnerID]; ok {
// 			delete(set, conn)
// 			if len(set) == 0 {
// 				delete(partnerConns, partnerID)
// 			}
// 		}
// 	} else {
// 		delete(customers, conn)
// 	}
// 	mu.Unlock()
// }

// // POST /update-partner-location/:partner_id
// // (partner app pings this; we may broadcast to customers)
// func UpdateDeliveryPartnerLocation(c *gin.Context) {
// 	var location models.Location
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
// 		return
// 	}
// 	partnerID := c.Param("partner_id")

// 	// store last known
// 	mu.Lock()
// 	deliveryPartnerLocations[partnerID] = location
// 	mu.Unlock()

// 	// Optional: broadcast to ALL customers (legacy behavior)
// 	broadcastPartnerLocationToCustomers(partnerID, location)

// 	c.JSON(http.StatusOK, gin.H{"message": "Location updated"})
// }

// func broadcastPartnerLocationToCustomers(partnerID string, location models.Location) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	for cli := range customers {
// 		if err := cli.WriteJSON(map[string]interface{}{
// 			"type":       "partner.location",
// 			"partner_id": partnerID,
// 			"lat":        location.Lat,
// 			"lng":        location.Lng,
// 		}); err != nil {
// 			log.Println("Error sending to customer:", err)
// 			cli.Close()
// 			delete(customers, cli)
// 		}
// 	}
// }

// // NEW: POST /send-location-to-delivery-partner/:partner_id
// // Body: { "customer_location": {...}, "parking_location": {...} }
// func SendCustomerAndParkingLocationToDeliveryPartner(c *gin.Context) {
// 	var payload models.AssignmentPayload
// 	if err := c.ShouldBindJSON(&payload); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
// 		return
// 	}
// 	partnerID := c.Param("partner_id")

// 	// Create the message with customer and parking details
// 	msg := models.WSAssignmentMessage{
// 		Type:     "assignment",
// 		Customer: payload.Customer,
// 		Parking:  payload.Parking,
// 	}

// 	// Targeted send to the specific partner
// 	if !sendToPartner(partnerID, msg) {
// 		// No WebSocket connections for this partner (e.g., offline)
// 		c.JSON(http.StatusAccepted, gin.H{"message": "Partner offline; assignment queued or logged"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Locations sent to delivery partner"})
// }

// func sendToPartner(partnerID string, v interface{}) bool {
// 	mu.Lock()
// 	defer mu.Unlock()

//		// Get the delivery partner's WebSocket connection
//		set, ok := partnerConns[partnerID]
//		if !ok || len(set) == 0 {
//			return false
//		}
//		for cli := range set {
//			if err := cli.WriteJSON(v); err != nil {
//				log.Println("Error sending to partner:", err)
//				cli.Close()
//				delete(set, cli)
//			}
//		}
//		return true
//	}
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khanjumed/geopackage/internal/config"
	"github.com/khanjumed/geopackage/internal/models"
)

// GET /api/order/:id/partner-location
func GetLastPartnerLocation(c *gin.Context) {
	oid := c.Param("id")
	var row models.PartnerLocation
	err := config.DB.Get(&row, `SELECT order_id, partner_id, lat, lng, heading, speed, updated_at
	                            FROM partner_locations
	                            WHERE order_id=? ORDER BY updated_at DESC LIMIT 1`, oid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no location yet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": row})
}
