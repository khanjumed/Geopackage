package models

type ParkingSpot struct {
    ID          int     `json:"id" db:"id"`
    Name        string  `json:"name" db:"name"`
    Address     string  `json:"address" db:"address"`
    Description string  `json:"description" db:"description"`
    Lat         float64 `json:"lat" db:"lat"`
    Lng         float64 `json:"lng" db:"lng"`
}