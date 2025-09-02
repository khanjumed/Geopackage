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
