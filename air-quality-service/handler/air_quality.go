package handler

import (
	"air-quality-service/model"
	"air-quality-service/service"
	"air-quality-service/util"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AirQualityHandler struct {
	db                *gorm.DB
	airQualityService *service.AirQualityService
}

func NewAirQualityHandler(db *gorm.DB, airQualityService *service.AirQualityService) *AirQualityHandler {
	return &AirQualityHandler{
		db:                db,
		airQualityService: airQualityService,
	}
}

type FetchAirQualityRequest struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (ah *AirQualityHandler) FetchAirQualityData(c echo.Context) error {
	// get Body
	var reqBody FetchAirQualityRequest
	err := c.Bind(&reqBody)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "Invalid Request", err.Error())
	}

	//validate if lat and lon is exist
	if reqBody.Latitude == 0 || reqBody.Longitude == 0 {
		return util.NewAppError(http.StatusBadRequest, "Invalid Request", "Latitude and Longitude is required")
	}

	//search if lat and lon combination exist, if not create new location

	//fetch air quality data
	airQuality, err := ah.airQualityService.FetchAirQuality(strconv.FormatFloat(reqBody.Latitude, 'f', -1, 64), strconv.FormatFloat(reqBody.Longitude, 'f', -1, 64))
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "Failed to fetch air quality data", err.Error())
	}

	//store data to db
	var airQualities []model.AirQuality
	createdTime := time.Now()
	for _, data := range airQuality.List {
		fetchedTime := time.Unix(data.Dt, 0).UTC()
		airQuality := model.AirQuality{
			LocationID: 1, //from searching result
			AQI:        data.Main.Aqi,
			CO:         data.Components.Co,
			NO:         data.Components.No,
			NO2:        data.Components.No2,
			O3:         data.Components.O3,
			SO2:        data.Components.So2,
			PM25:       data.Components.Pm25,
			PM10:       data.Components.Pm10,
			NH3:        data.Components.Nh3,
			CreatedAt:  &createdTime,
			FetchTime:  &fetchedTime,
		}
		airQualities = append(airQualities, airQuality)
	}

	if err := ah.db.Create(&airQualities).Error; err != nil {
		return util.NewAppError(http.StatusInternalServerError, "Failed to store air quality data", err.Error())
	}

	return c.JSON(http.StatusOK, util.Response{
		Message: "Air quality data fetched and stored to DB successfully",
		Data:    airQualities,
	})
}

func (ah *AirQualityHandler) GetAirQualities(c echo.Context) error {
	var airQualities []model.AirQuality
	err := ah.db.Find(&airQualities).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "Failed to fetch air quality data", err.Error())
	}

	return c.JSON(http.StatusOK, util.Response{
		Message: "Air quality data fetched successfully",
		Data:    airQualities,
	})
}
