package entity

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	LocationName string  `json:"location_name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}
