// package main

// import (
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// type ParkingSpot struct {
// 	ID   int     `json:"id"`
// 	Name string  `json:"name"`
// 	Lat  float64 `json:"lat"`
// 	Lng  float64 `json:"lng"`
// }

// func main() {
// 	_ = godotenv.Load()

// 	if err := InitDB(); err != nil {
// 		log.Fatalf("DB connection failed: %v", err)
// 	}

// 	router := gin.Default()

// 	// Serve static HTML on root `/`
// 	router.LoadHTMLFiles("static/index.html")
// 	router.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.html", nil)
// 	})

// 	// Serve additional static files if needed (e.g., CSS, JS)
// 	router.Static("/static", "./static")

// 	// API endpoints
// 	router.GET("/api/parking-spots/nearby", getNearbyParkingSpots)
// 	router.POST("/api/admin/add-spot", addParkingSpot)

// 	// Run server on port 8080
// 	router.Run(":8080")
// }

// func getNearbyParkingSpots(c *gin.Context) {
// 	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
// 	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
// 	const radiusKm = 20.0
// 	var spots []ParkingSpot

// 	query := `
//     SELECT id, name, lat, lng
//     FROM parking_spots
//     WHERE
//         (6371 * acos(
//             cos(radians(?)) * cos(radians(lat)) *
//             cos(radians(lng) - radians(?)) +
//             sin(radians(?)) * sin(radians(lat))
//         )) < ?
//     ORDER BY id ASC
//     LIMIT 10;
//     `

// 	// Log query and parameters
// 	log.Printf("Running SQL Query: %s\nWith Params: lat=%.6f, lng=%.6f, lat again=%.6f, radius=%.2f\n", query, lat, lng, lat, radiusKm)

// 	err := db.Select(&spots, query, lat, lng, lat, radiusKm)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": spots})
// }

// func addParkingSpot(c *gin.Context) {
// 	var spot ParkingSpot
// 	if err := c.BindJSON(&spot); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	_, err := db.Exec("INSERT INTO parking_spots (name, lat, lng) VALUES (?, ?, ?)",
// 		spot.Name, spot.Lat, spot.Lng)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Parking spot added"})
// }

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type ParkingSpot struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

func main() {
	_ = godotenv.Load()

	if err := InitDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	router := gin.Default()

	router.LoadHTMLFiles("static/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Static("/static", "./static")

	router.GET("/api/parking-spots/nearby", getNearbyParkingSpots)
	router.POST("/api/admin/add-spot", addParkingSpot)

	router.Run(":8080")
}

func getNearbyParkingSpots(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	const radiusKm = 20.0
	var spots []ParkingSpot

	query := `
		SELECT id, name, lat, lng
		FROM parking_spots
		WHERE 
			(6371 * acos(
				cos(radians(?)) * cos(radians(lat)) *
				cos(radians(lng) - radians(?)) +
				sin(radians(?)) * sin(radians(lat))
			)) < ?
		ORDER BY id ASC
		LIMIT 10;
	`

	// Log query params
	log.Printf("Params: lat=%.6f, lng=%.6f, lat2=%.6f, radius=%.2f\n", lat, lng, lat, radiusKm)

	// Run query and log full error
	err := db.Select(&spots, query, lat, lng, lat, radiusKm)
	if err != nil {
		log.Printf("DB query failed: %v\n", err) // ADD THIS LINE
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": spots})
}

func addParkingSpot(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	coords, err := geocodeAddress(input.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geocoding failed: " + err.Error()})
		return
	}

	_, err = db.Exec(`
		INSERT INTO parking_spots (name, address, description, lat, lng)
		VALUES (?, ?, ?, ?, ?)`,
		input.Name, input.Address, input.Description, coords.Lat, coords.Lng)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking spot added with geolocation"})
}

type GeocodeResult struct {
	Lat float64
	Lng float64
}

func geocodeAddress(address string) (GeocodeResult, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	reqURL := fmt.Sprintf("%s?address=%s&key=%s", baseURL, url.QueryEscape(address), apiKey)

	resp, err := http.Get(reqURL)
	if err != nil {
		return GeocodeResult{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeocodeResult{}, err
	}

	if result.Status != "OK" || len(result.Results) == 0 {
		return GeocodeResult{}, fmt.Errorf("Google Maps API returned: %s", result.Status)
	}

	loc := result.Results[0].Geometry.Location
	return GeocodeResult{Lat: loc.Lat, Lng: loc.Lng}, nil
}
