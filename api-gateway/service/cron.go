package service

import (
	pb "api-gateway/pb"
	"api-gateway/util"
	"fmt"

	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CronServices struct {
	AQService       pb.AirQualityServiceClient
	LocationService pb.LocationServiceClient
	BFService       pb.BusinessFacilitiesClient
}

func NewCronServices(AQService pb.AirQualityServiceClient, LocationService pb.LocationServiceClient,
	BFService pb.BusinessFacilitiesClient) *CronServices {
	return &CronServices{
		AQService:       AQService,
		LocationService: LocationService,
		BFService:       BFService,
	}
}

func (cr *CronServices) RenewAQData() {
	log.Print("Record Location AQ")
	ctx := util.CreateServiceContext()
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

func (cr *CronServices) RenewBFAQData() {
	log.Print("Record Business Facility AQ")
	ctx := util.CreateServiceContext()
	//Get All BFs
	res, err := cr.BFService.GetBusinessFacilities(ctx, &pb.GetBFRequests{})
	if err != nil {
		log.Errorf("Cron Error when getting all business: %v\n", err)
	}

	//Get New AQ Data for each location
	for _, bf := range res.BusinessFacilities {
		fmt.Println("Getting air quality data for business with ID ", bf.Id)
		_, err := cr.AQService.SaveAirQualityForBusiness(ctx, &pb.SaveAirQualityForBusinessReq{BusinessId: int64(bf.Id)})
		if err != nil {
			log.Errorf("Cron Error when getting new AQ data for business id %s: %v\n", bf.Id, err)
		}
	}
}
