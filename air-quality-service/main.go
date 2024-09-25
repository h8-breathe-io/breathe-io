package main

import (
	"air-quality-service/config"
	"air-quality-service/handler"
	pb "air-quality-service/pb/generated"
	"air-quality-service/service"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
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
