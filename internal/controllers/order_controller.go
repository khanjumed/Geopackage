package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khanjumed/geopackage/internal/config"
	"github.com/khanjumed/geopackage/internal/models"
	"github.com/khanjumed/geopackage/internal/ws"
)

// POST /api/order
// body: { "orderId":"ORD123", "pickup":{"lat":..,"lng":..}, "drop":{"lat":..,"lng":..}, "status":"created" }
func UpsertOrder(c *gin.Context) {
	var in models.Order
	if err := c.ShouldBindJSON(&in); err != nil || in.OrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	_, err := config.DB.Exec(`
INSERT INTO orders (order_id, pickup_lat, pickup_lng, drop_lat, drop_lng, status)
VALUES (?, ?, ?, ?, ?, COALESCE(?, 'created'))
ON DUPLICATE KEY UPDATE
  pickup_lat = VALUES(pickup_lat),
  pickup_lng = VALUES(pickup_lng),
  drop_lat   = VALUES(drop_lat),
  drop_lng   = VALUES(drop_lng),
  status     = VALUES(status)
`, in.OrderID, in.Pickup.Lat, in.Pickup.Lng, in.Drop.Lat, in.Drop.Lng, in.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db upsert failed"})
		return
	}

	// push pickup/drop to everyone currently connected on this order
	ws.BroadcastOrderInfo(in.OrderID, in.Pickup, in.Drop)

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": in})
}

// GET /api/order/:id
func GetOrder(c *gin.Context) {
	oid := c.Param("id")
	var row struct {
		OrderID   string  `db:"order_id"`
		PickupLat float64 `db:"pickup_lat"`
		PickupLng float64 `db:"pickup_lng"`
		DropLat   float64 `db:"drop_lat"`
		DropLng   float64 `db:"drop_lng"`
		Status    string  `db:"status"`
	}
	err := config.DB.Get(&row, `SELECT order_id, pickup_lat, pickup_lng, drop_lat, drop_lng, status FROM orders WHERE order_id = ?`, oid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		}
		return
	}
	out := models.Order{
		OrderID: row.OrderID,
		Pickup:  models.LatLng{Lat: row.PickupLat, Lng: row.PickupLng},
		Drop:    models.LatLng{Lat: row.DropLat, Lng: row.DropLng},
		Status:  row.Status,
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": out})
}
