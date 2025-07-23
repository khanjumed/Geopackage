package controllers

import (
	"net/http"
	"parking-service/internal/config"
	"parking-service/internal/models"
	"parking-service/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddParkingSpot(c *gin.Context) {
	var input models.ParkingSpot
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	coords, err := services.GeocodeAddress(input.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = config.DB.Exec(`INSERT INTO parking_spots (name, address, description, lat, lng) VALUES (?, ?, ?, ?, ?)`,
		input.Name, input.Address, input.Description, coords.Lat, coords.Lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Spot added"})
}

func GetNearbySpots(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	const radiusKm = 20.0
	var spots []models.ParkingSpot
	query := `SELECT id, name, lat, lng FROM parking_spots WHERE (6371 * acos(cos(radians(?)) * cos(radians(lat)) * cos(radians(lng) - radians(?)) + sin(radians(?)) * sin(radians(lat)))) < ? ORDER BY id ASC LIMIT 20`
	err := config.DB.Select(&spots, query, lat, lng, lat, radiusKm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": spots})
}

func UpdateParkingSpot(c *gin.Context) {
	id := c.Param("id")
	var input models.ParkingSpot
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	coords, err := services.GeocodeAddress(input.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geocoding failed"})
		return
	}
	_, err = config.DB.Exec(`UPDATE parking_spots SET name = ?, address = ?, description = ?, lat = ?, lng = ? WHERE id = ?`,
		input.Name, input.Address, input.Description, coords.Lat, coords.Lng, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Spot updated"})
}

func DeleteParkingSpot(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec(`DELETE FROM parking_spots WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Spot deleted"})
}
