package service

import (
	pb "api-gateway/pb"
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CronServices struct {
	AQService       pb.AirQualityServiceClient
	LocationService pb.LocationServiceClient
}

func NewCronServices(AQService pb.AirQualityServiceClient, LocationService pb.LocationServiceClient) *CronServices {
	return &CronServices{
		AQService:       AQService,
		LocationService: LocationService,
	}
}

func (cr *CronServices) RenewAQData() {
	//Get All Locations
	res, err := cr.LocationService.GetLocations(context.TODO(), &emptypb.Empty{})
	if err != nil {
		log.Errorf("Cron Error when getting all locations: %v\n", err)
	}

	//Get New AQ Data for each location
	for _, location := range res.Locations {
		fmt.Println("Getting air quality data for location with ID ", location.LocationId)
		_, err := cr.AQService.SaveAirQualities(context.TODO(), &pb.SaveAirQualitiesRequest{Latitude: location.Latitude, Longitude: location.Longitude})
		if err != nil {
			log.Errorf("Cron Error when getting new AQ data for location %s: %v\n", location.LocationName, err)
		}
	}
}
