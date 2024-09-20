package service

import (
	"air-quality-service/entity"
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
)

type AirQualityService struct {
	Client  *resty.Client
	APIKey  string
	BaseUrl string
}

func NewAirQualityService() *AirQualityService {
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	baseUrl := os.Getenv("OPEN_WEATHER_BASE_URL")

	client := resty.New()

	return &AirQualityService{
		Client:  client,
		APIKey:  apiKey,
		BaseUrl: baseUrl,
	}
}

func (a *AirQualityService) FetchAirQuality(lat, lon string) (*entity.AirQuality, error) {
	resp, err := a.Client.R().
		SetQueryParam("lat", lat).
		SetQueryParam("lon", lon).
		SetQueryParam("appid", a.APIKey).
		Get(a.BaseUrl)

	if err != nil {
		return nil, err
	}

	var airQuality entity.AirQuality
	err = json.Unmarshal(resp.Body(), &airQuality)
	if err != nil {
		return nil, err
	}

	return &airQuality, nil
}
