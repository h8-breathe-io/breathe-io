package model

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	LocationName string       `json:"location_name"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	AirQualities []AirQuality `json:"air_qualities" gorm:"foreignKey:LocationID"` // Has-many relation to AirQuality
}
