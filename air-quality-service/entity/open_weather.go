package entity

type AirQuality struct {
	Coordinate struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	List []struct {
		Main struct {
			Aqi int `json:"aqi"`
		} `json:"main"`
		Components struct {
			Co   float64 `json:"co"`
			No   float64 `json:"no"`
			No2  float64 `json:"no2"`
			O3   float64 `json:"o3"`
			So2  float64 `json:"so2"`
			Pm25 float64 `json:"pm2_5"`
			Pm10 float64 `json:"pm10"`
			Nh3  float64 `json:"nh3"`
		} `json:"components"`
		Dt int64 `json:"dt"`
	} `json:"list"`
}
