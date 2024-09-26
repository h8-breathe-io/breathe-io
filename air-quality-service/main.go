package main

import (
	"air-quality-service/config"
	"air-quality-service/handler"
	pb "air-quality-service/pb/generated"
	"air-quality-service/service"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	db := config.CreateDBInstance()

	// instantiate dependencies
	airQualityService := service.NewAirQualityService()
	userService := service.NewUserService()
	businessFacilityService := service.NewBusinessFacilityService()
	airQualityHandler := handler.NewAirQualityHandler(db, airQualityService)
	locationHandler := handler.NewLocationHandler(db, userService, businessFacilityService)
	grpcServer := grpc.NewServer()

	pb.RegisterAirQualityServiceServer(grpcServer, airQualityHandler)
	pb.RegisterLocationServiceServer(grpcServer, locationHandler)

	//declare cron services
	c := cron.New()
	//running cron job exactly every start of the
	c.AddFunc("0 0 * * * *", func() {
		//Get All Locations
		res, err := locationHandler.GetLocations(context.TODO(), &emptypb.Empty{})
		if err != nil {
			log.Printf("Cron Error when getting all locations: %v\n", err)
		}

		//Get New AQ Data for each location
		for _, location := range res.Locations {
			fmt.Println("Getting air quality data for location with ID ", location.LocationId)
			_, err := airQualityHandler.SaveAirQualities(context.TODO(), &pb.SaveAirQualitiesRequest{Latitude: location.Latitude, Longitude: location.Longitude})
			if err != nil {
				log.Printf("Cron Error when getting new AQ data for location %s: %v\n", location.LocationName, err)
			}
		}
	})
	c.Start()

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	// start server
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %s", listen.Addr().String())
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
