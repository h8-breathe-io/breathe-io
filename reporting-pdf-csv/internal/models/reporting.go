package models

import "time"

// User model
type User struct {
	ID           int    `gorm:"primary_key"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PhoneNumber  string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Tier         string `gorm:"not null;check(tier IN ('free', 'business'))"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Location model
type Location struct {
	ID           int `gorm:"primary_key"`
	LocationName string
	Latitude     float64 `gorm:"not null"`
	Longitude    float64 `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BusinessFacility model
type BusinessFacility struct {
	ID            int `gorm:"primary_key"`
	UserID        int `gorm:"not null"`
	CompanyType   string
	TotalEmission float64 `gorm:"not null"`
	LocationID    int     `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// AirQuality model
type AirQuality struct {
	ID         int `gorm:"primary_key"`
	LocationID int `gorm:"not null"`
	AQI        int `gorm:"not null"`
	CO         float64
	NO         float64
	NO2        float64
	O3         float64
	SO2        float64
	PM25       float64
	PM10       float64
	NH3        float64
	FetchTime  time.Time
	CreatedAt  time.Time
}


// Tentukan nama tabel secara eksplisit
func (AirQuality) TableName() string {
    return "air_quality" // Ini memastikan bahwa GORM menggunakan tabel "air_quality"
}