package handler

import (
	"air-quality-service/model"
	"air-quality-service/service"
	"air-quality-service/util"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	pb "air-quality-service/pb/generated"

	"gorm.io/gorm"
)

type AirQualityHandler struct {
	pb.UnimplementedAirQualityServiceServer
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

func (ah *AirQualityHandler) SaveAirQualities(ctx context.Context, req *pb.SaveAirQualitiesRequest) (*pb.SaveAirQualitiesResponse, error) {
	//validate lat and long from req
	// check if latitude is off limit
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, errors.New("latitude value is off limit")
	}

	// check if longitude is off limit
	if req.Longitude == -180 || req.Longitude > 180 {
		return nil, errors.New("longitude value is off limit")
	}

	// get Body
	var reqBody = FetchAirQualityRequest{
		Latitude:  float64(req.Latitude),
		Longitude: float64(req.Longitude),
	}

	//search if lat and lon combination exist, if not create new location
	var location model.Location
	err := ah.db.Where("latitude = ? AND longitude = ?", reqBody.Latitude, reqBody.Longitude).First(&location).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			location = model.Location{
				LocationName: "Location",
				Latitude:     reqBody.Latitude,
				Longitude:    reqBody.Longitude,
			}
			if err := ah.db.Create(&location).Error; err != nil {
				return nil, errors.New("failed to create location")
			}
		} else {
			return nil, errors.New("failed to search location")
		}
	}

	//fetch air quality data
	airQuality, err := ah.airQualityService.FetchAirQuality(strconv.FormatFloat(reqBody.Latitude, 'f', -1, 64), strconv.FormatFloat(reqBody.Longitude, 'f', -1, 64))
	if err != nil {
		return nil, errors.New("failed to fetch air quality data")
	}

	fmt.Println(airQuality)

	//store data to db
	var airQualities []model.AirQuality
	createdTime := time.Now()
	for _, data := range airQuality.List {
		fetchedTime := time.Unix(data.Dt, 0).UTC()
		airQuality := model.AirQuality{
			LocationID: int(location.ID), //from searching result
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
		return nil, errors.New("failed to save air qualities")
	}

	res := &pb.SaveAirQualitiesResponse{
		Success: true,
	}

	return res, nil
}

func (ah *AirQualityHandler) GetAirQualities(ctx context.Context, req *pb.GetAirQualitiesRequest) (*pb.GetAirQualitiesResponse, error) {
	// Validate location id
	if req.LocationId == 0 {
		return nil, errors.New("location id is required")
	}

	// Search if location exists
	var location model.Location
	err := ah.db.Where("id = ?", req.LocationId).First(&location).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("location not found")
		}
		return nil, errors.New("failed to search location")
	}

	// Initialize a query for air quality data
	query := ah.db.Where("location_id = ?", req.LocationId)

	// Handle start_date if provided
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start date format, expected yyyy-mm-dd")
		}

		// Check if the provided start_date is lower than the oldest available data
		var oldestAirQuality model.AirQuality
		err = ah.db.Where("location_id = ?", req.LocationId).Order("fetch_time ASC").First(&oldestAirQuality).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("failed to fetch oldest air quality data")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) && startDate.Before(*oldestAirQuality.FetchTime) {
			// Use the fetch_time of the oldest record if start_date is before the oldest available data
			startDate = *oldestAirQuality.FetchTime
		}

		query = query.Where("fetch_time >= ?", startDate)
	}

	// Handle end_date if provided
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format, expected yyyy-mm-dd")
		}

		// Check if the provided end_date is higher than the newest available data
		var newestAirQuality model.AirQuality
		err = ah.db.Where("location_id = ?", req.LocationId).Order("fetch_time DESC").First(&newestAirQuality).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("failed to fetch newest air quality data")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) && endDate.After(*newestAirQuality.FetchTime) {
			// Use the fetch_time of the newest record if end_date is after the newest available data
			endDate = *newestAirQuality.FetchTime
		}

		query = query.Where("fetch_time <= ?", endDate)
	}

	// Execute the query
	var airQualities []model.AirQuality
	err = query.Order("fetch_time DESC").Find(&airQualities).Error
	if err != nil {
		return nil, errors.New("failed to get air qualities")
	}

	// Build the response
	res := &pb.GetAirQualitiesResponse{
		AirQualities: make([]*pb.AirQuality, 0, len(airQualities)),
	}

	for _, airQuality := range airQualities {
		fetchedTime := airQuality.FetchTime.Format("2006-01-02 15:04:05")
		res.AirQualities = append(res.AirQualities, &pb.AirQuality{
			Id:         uint64(airQuality.ID),
			LocationId: int64(airQuality.LocationID),
			Aqi:        int64(airQuality.AQI),
			Co:         airQuality.CO,
			No:         airQuality.NO,
			No2:        airQuality.NO2,
			O3:         airQuality.O3,
			So2:        airQuality.SO2,
			Pm25:       airQuality.PM25,
			Pm10:       airQuality.PM10,
			Nh3:        airQuality.NH3,
			FetchTime:  fetchedTime,
		})
	}

	return res, nil
}

func (ah *AirQualityHandler) SaveHistoricalAirQualities(ctx context.Context, req *pb.SaveHistoricalAirQualitiesRequest) (*pb.SaveAirQualitiesResponse, error) {
	//validate lat and long from req
	// check if latitude is off limit
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, errors.New("latitude value is off limit")
	}

	// check if longitude is off limit
	if req.Longitude == -180 || req.Longitude > 180 {
		return nil, errors.New("longitude value is off limit")
	}

	//search if lat and lon combination exist, if not create new location
	var location model.Location
	err := ah.db.Where("latitude = ? AND longitude = ?", req.Latitude, req.Longitude).First(&location).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			location = model.Location{
				LocationName: "Location",
				Latitude:     req.Latitude,
				Longitude:    req.Longitude,
			}
			if err := ah.db.Create(&location).Error; err != nil {
				return nil, errors.New("failed to create location")
			}
		} else {
			return nil, errors.New("failed to search location")
		}
	}

	startDateUnix, err := util.ParseStrToUnixDate(req.StartDate)
	if err != nil {
		return nil, errors.New("failed to parse start date")
	}

	endDateUnix, err := util.ParseStrToUnixDate(req.EndDate)
	if err != nil {
		return nil, errors.New("failed to parse end date")
	}

	//fetch air quality historical data
	airQuality, err := ah.airQualityService.FetchAirQualityByRange(strconv.FormatFloat(req.Latitude, 'f', -1, 64), strconv.FormatFloat(req.Longitude, 'f', -1, 64), startDateUnix, endDateUnix)
	if err != nil {
		return nil, errors.New("failed to fetch air quality data")
	}

	//store data to db
	var airQualities []model.AirQuality
	createdTime := time.Now()
	for _, data := range airQuality.List {
		fetchedTime := time.Unix(data.Dt, 0).UTC()
		airQuality := model.AirQuality{
			LocationID: int(location.ID), //from searching result
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
		return nil, errors.New("failed to save air qualities")
	}

	res := &pb.SaveAirQualitiesResponse{
		Success: true,
	}

	return res, nil
}

func (ah *AirQualityHandler) GetAirQualityByID(c context.Context, req *pb.GetAirQualityByIDReq) (*pb.GetAirQualityByIDResp, error) {
	// validate location id
	if req.Id == 0 {
		return nil, errors.New("aiur quality id is required")
	}

	var airQuality model.AirQuality
	err := ah.db.Where("id = ?", req.Id).First(&airQuality).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get air qualitiy id %d", req.Id)
	}

	fetchedTime := airQuality.FetchTime.Format("2006-01-02 15:04:05")
	return &pb.GetAirQualityByIDResp{
		AirQuality: &pb.AirQuality{
			Id:         uint64(airQuality.ID),
			LocationId: int64(airQuality.LocationID),
			Aqi:        int64(airQuality.AQI),
			Co:         airQuality.CO,
			No:         airQuality.NO,
			No2:        airQuality.NO2,
			O3:         airQuality.O3,
			So2:        airQuality.SO2,
			Pm25:       airQuality.PM25,
			Pm10:       airQuality.PM10,
			Nh3:        airQuality.NH3,
			FetchTime:  fetchedTime,
		}}, nil
}
