package model

import "gorm.io/gorm"

type BusinessFacility struct {
	gorm.Model
	UserID        uint64  `gorm:"not null" json:"user_id"`
	CompanyType   string  `gorm:"not null" json:"company_type"`
	TotalEmission float64 `gorm:"not null" json:"total_emission"`
	LocationID    uint64  `gorm:"not null" json:"location_id"`
}
