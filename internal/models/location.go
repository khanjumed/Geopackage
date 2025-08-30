// package models

// type Location struct {
// 	Lat float64 `json:"lat"`
// 	Lng float64 `json:"lng"`
// }

// package models

// type Location struct {
// 	Lat float64 `json:"lat"`
// 	Lng float64 `json:"lng"`
// }

// type AssignmentPayload struct {
// 	Customer Location `json:"customer_location"`
// 	Parking  Location `json:"parking_location"`
// }

// type WSAssignmentMessage struct {
// 	Type     string   `json:"type"`
// 	Customer Location `json:"customer_location"`
// 	Parking  Location `json:"parking_location"`
// }

package models

import "time"

type PartnerLocation struct {
	OrderID   string    `db:"order_id"   json:"orderId"`
	PartnerID string    `db:"partner_id" json:"partnerId"`
	Lat       float64   `db:"lat"        json:"lat"`
	Lng       float64   `db:"lng"        json:"lng"`
	Heading   *float64  `db:"heading"    json:"heading,omitempty"`
	Speed     *float64  `db:"speed"      json:"speed,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
