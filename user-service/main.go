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
)

func main() {
	db := config.CreateDBInstance()

	// instantiate dependencies
	userHandler := handler.NewUserHandler(db)
	businessFacilityHandler := handler.NewBusinessFacilitiesHandler(db)

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
