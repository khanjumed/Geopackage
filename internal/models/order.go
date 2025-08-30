package models

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Order struct {
	OrderID string `json:"orderId" db:"order_id"`
	Pickup  LatLng `json:"pickup"`
	Drop    LatLng `json:"drop"`
	Status  string `json:"status" db:"status"`
}
