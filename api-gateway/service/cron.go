package service

import (
	pb "api-gateway/pb"
	"context"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/metadata"
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

func (cr *CronServices) createServiceContext() context.Context {
	token := os.Getenv("AQ_SERVICE_TOKEN")
	md := metadata.Pairs("auth_token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}

func (cr *CronServices) RenewAQData() {

	ctx := cr.createServiceContext()
	//Get All Locations
	res, err := cr.LocationService.GetLocations(ctx, &emptypb.Empty{})
	if err != nil {
		log.Errorf("Cron Error when getting all locations: %v\n", err)
	}

	//Get New AQ Data for each location
	for _, location := range res.Locations {
		fmt.Println("Getting air quality data for location with ID ", location.LocationId)
		_, err := cr.AQService.SaveAirQualities(ctx, &pb.SaveAirQualitiesRequest{Latitude: location.Latitude, Longitude: location.Longitude})
		if err != nil {
			log.Errorf("Cron Error when getting new AQ data for location %s: %v\n", location.LocationName, err)
		}
	}
}
