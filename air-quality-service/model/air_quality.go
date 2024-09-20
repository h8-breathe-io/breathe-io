package model

import "time"

type AirQuality struct {
	ID         uint       `json:"id"`
	LocationID int        `json:"location_id"`
	AQI        int        `json:"aqi"`
	CO         float64    `json:"co"`
	NO         float64    `json:"no"`
	NO2        float64    `json:"no2"`
	O3         float64    `json:"o3"`
	SO2        float64    `json:"so2"`
	PM25       float64    `json:"pm25"`
	PM10       float64    `json:"pm10"`
	NH3        float64    `json:"nh3"`
	FetchTime  *time.Time `json:"fetch_time"`
	CreatedAt  *time.Time `json:"created_at"`
}
