package entity

import "time"

type BusinessFacility struct {
	ID            int
	UserID        uint64
	CompanyType   string
	TotalEmission float64
	LocationID    uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
