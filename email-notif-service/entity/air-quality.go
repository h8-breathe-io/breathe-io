package entity

type AirQuality struct {
	Id         uint64  `json:"id,omitempty"`
	LocationId int64   `json:"location_id,omitempty"`
	Aqi        int64   `json:"aqi,omitempty"`
	Co         float64 `json:"co,omitempty"`
	No         float64 `json:"no,omitempty"`
	No2        float64 `json:"no2,omitempty"`
	O3         float64 `json:"o3,omitempty"`
	So2        float64 `json:"so2,omitempty"`
	Pm25       float64 `json:"pm25,omitempty"`
	Pm10       float64 `json:"pm10,omitempty"`
	Nh3        float64 `json:"nh3,omitempty"`
	FetchTime  string  `json:"fetch_time,omitempty"`
}
