package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"user-service/config"
	"user-service/handler"
	pb "user-service/pb/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := config.CreateDBInstance()

	//make connection to location grpc server
	locationAddr := os.Getenv("LOCATION_GRPC_ADDR")
	locationConn, err := grpc.NewClient(locationAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to location service: %v", err)
	}
	defer locationConn.Close()
	locationServiceClient := pb.NewLocationServiceClient(locationConn)

	// instantiate dependencies
	userHandler := handler.NewUserHandler(db)
	businessFacilityHandler := handler.NewBusinessFacilitiesHandler(db, locationServiceClient)

	grpcServer := grpc.NewServer()

	pb.RegisterUserServer(grpcServer, userHandler)
	pb.RegisterBusinessFacilitiesServer(grpcServer, businessFacilityHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
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
