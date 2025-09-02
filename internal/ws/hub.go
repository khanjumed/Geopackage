package ws

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/khanjumed/geopackage/internal/config"
	"github.com/khanjumed/geopackage/internal/models"
)

type wsMsg struct {
	Type    string   `json:"type"` // "partner_location"
	Lat     float64  `json:"lat,omitempty"`
	Lng     float64  `json:"lng,omitempty"`
	Heading *float64 `json:"heading,omitempty"` // optional
	Speed   *float64 `json:"speed,omitempty"`   // optional
}

type client struct {
	conn      *websocket.Conn
	orderID   string
	role      string // "customer" | "partner"
	partnerID string // for partner
	lastSent  time.Time
	lastLat   float64
	lastLng   float64
}

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	roomsMu sync.RWMutex
	rooms   = map[string]map[*client]bool{} // orderID -> clients
)

// -------------------- helpers (rooms) --------------------

func register(c *client) {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	if rooms[c.orderID] == nil {
		rooms[c.orderID] = map[*client]bool{}
	}
	rooms[c.orderID][c] = true
}
func unregister(c *client) {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	if rs, ok := rooms[c.orderID]; ok {
		delete(rs, c)
		if len(rs) == 0 {
			delete(rooms, c.orderID)
		}
	}
	_ = c.conn.Close()
}
func broadcast(orderID string, payload any) {
	roomsMu.RLock()
	defer roomsMu.RUnlock()
	for cl := range rooms[orderID] {
		_ = cl.conn.WriteJSON(payload)
	}
}

// Public helper to push latest pickup/drop to a room
func BroadcastOrderInfo(orderID string, pickup models.LatLng, drop models.LatLng) {
	broadcast(orderID, gin.H{
		"type":   "order_info",
		"pickup": gin.H{"lat": pickup.Lat, "lng": pickup.Lng},
		"drop":   gin.H{"lat": drop.Lat, "lng": drop.Lng},
	})
}

// -------------------- helpers (Redis) --------------------

func redisKey(orderID string) string {
	return "partner_location:" + orderID
}

func storePartnerLocationInRedis(orderID string, msg wsMsg) {
	// HSET partner_location:<orderId> lat <v> lng <v> heading <v> speed <v>
	fields := map[string]interface{}{
		"lat": msg.Lat,
		"lng": msg.Lng,
	}
	// Only set optional fields if present
	if msg.Heading != nil {
		fields["heading"] = *msg.Heading
	}
	if msg.Speed != nil {
		fields["speed"] = *msg.Speed
	}
	if err := config.RDB.HSet(config.Ctx, redisKey(orderID), fields).Err(); err != nil {
		log.Println("redis HSet error:", err)
	}
	// TTL (optional): keep last known for some time; comment out if not desired.
	_ = config.RDB.Expire(config.Ctx, redisKey(orderID), 6*time.Hour).Err()
}

func getPartnerLocationFromRedis(orderID string) (lat float64, lng float64, heading *float64, speed *float64, ok bool) {
	data, err := config.RDB.HGetAll(config.Ctx, redisKey(orderID)).Result()
	if err != nil || len(data) == 0 {
		return 0, 0, nil, nil, false
	}
	if v, ok2 := data["lat"]; ok2 {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			lat = f
		}
	}
	if v, ok2 := data["lng"]; ok2 {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			lng = f
		}
	}
	if v, ok2 := data["heading"]; ok2 {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			heading = &f
		}
	}
	if v, ok2 := data["speed"]; ok2 {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			speed = &f
		}
	}
	if lat == 0 && lng == 0 {
		return 0, 0, nil, nil, false
	}
	return lat, lng, heading, speed, true
}

// -------------------- WebSocket handler --------------------

// GET /ws?orderId=ORD123&role=customer|partner[&partnerId=P1]
func HandleWS(c *gin.Context) {
	orderID := c.Query("orderId")
	role := c.Query("role")
	partnerID := c.Query("partnerId")

	if orderID == "" || (role != "customer" && role != "partner") || (role == "partner" && partnerID == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orderId and valid role (partnerId required for partner)"})
		return
	}

	// Upgrade to WS
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws upgrade:", err)
		return
	}
	cl := &client{conn: conn, orderID: orderID, role: role, partnerID: partnerID}
	register(cl)
	defer unregister(cl)

	// On connect: push current pickup/drop (from DB) if exists
	var row struct{ PickupLat, PickupLng, DropLat, DropLng float64 }
	_ = config.DB.Get(&row, `SELECT pickup_lat, pickup_lng, drop_lat, drop_lng FROM orders WHERE order_id=?`, orderID)
	if row.PickupLat != 0 || row.PickupLng != 0 || row.DropLat != 0 || row.DropLng != 0 {
		_ = cl.conn.WriteJSON(gin.H{
			"type":   "order_info",
			"pickup": gin.H{"lat": row.PickupLat, "lng": row.PickupLng},
			"drop":   gin.H{"lat": row.DropLat, "lng": row.DropLng},
		})
	}

	// Try Redis for last known partner location first
	if lat, lng, heading, speed, ok := getPartnerLocationFromRedis(orderID); ok {
		_ = cl.conn.WriteJSON(gin.H{
			"type":    "partner_location",
			"lat":     lat,
			"lng":     lng,
			"heading": heading,
			"speed":   speed,
		})
	} else {
		// Fallback to DB if Redis empty
		var last struct{ Lat, Lng float64 }
		_ = config.DB.Get(&last, `SELECT lat, lng FROM partner_locations WHERE order_id=? ORDER BY updated_at DESC LIMIT 1`, orderID)
		if last.Lat != 0 || last.Lng != 0 {
			_ = cl.conn.WriteJSON(gin.H{"type": "partner_location", "lat": last.Lat, "lng": last.Lng})
		}
	}

	// Read loop: partner sends live location; server caches & broadcasts
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var msg wsMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		// Only partner publishes live location
		if msg.Type == "partner_location" && cl.role == "partner" {
			if msg.Lat == 0 && msg.Lng == 0 {
				continue
			}
			if !shouldAcceptUpdate(cl, msg.Lat, msg.Lng) {
				continue
			}
			cl.lastLat, cl.lastLng = msg.Lat, msg.Lng
			cl.lastSent = time.Now()

			// 1) Fast path: cache latest in Redis (primary for realtime)
			storePartnerLocationInRedis(cl.orderID, msg)

			// 2) Broadcast to all viewers of this order
			broadcast(cl.orderID, gin.H{
				"type":    "partner_location",
				"lat":     msg.Lat,
				"lng":     msg.Lng,
				"heading": msg.Heading,
				"speed":   msg.Speed,
			})

			// 3) Optional persistence (kept from your code). You can throttle this more if desired.
			_, _ = config.DB.Exec(`
INSERT INTO partner_locations (order_id, partner_id, lat, lng, heading, speed)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE lat=VALUES(lat), lng=VALUES(lng),
  heading=VALUES(heading), speed=VALUES(speed), updated_at=CURRENT_TIMESTAMP
`, cl.orderID, cl.partnerID, msg.Lat, msg.Lng, msg.Heading, msg.Speed)
		}
	}
}

// throttle: accept if 2s passed OR moved >10m
func shouldAcceptUpdate(cl *client, lat, lng float64) bool {
	if time.Since(cl.lastSent) > 2*time.Second {
		return true
	}
	if cl.lastLat == 0 && cl.lastLng == 0 {
		return true
	}
	return distanceMeters(cl.lastLat, cl.lastLng, lat, lng) > 10.0
}

// Haversine (meters)
func distanceMeters(aLat, aLng, bLat, bLng float64) float64 {
	const R = 6371000.0
	toRad := func(d float64) float64 { return d * math.Pi / 180 }
	dLat := toRad(bLat - aLat)
	dLng := toRad(bLng - aLng)
	lat1 := toRad(aLat)
	lat2 := toRad(bLat)
	sinDLat := math.Sin(dLat / 2)
	sinDLng := math.Sin(dLng / 2)
	a := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLng*sinDLng
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
