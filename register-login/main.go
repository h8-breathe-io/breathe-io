package main

import (
	"fmt"
	"log"
	"net"
	"register-login/internal/database"
	"register-login/internal/handlers"
	"register-login/proto/pb"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database using GORM
	database.ConnectDB()

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the RegisterLoginServer with gRPC
	authService := &handlers.AuthServiceServer{
		DB: database.DB,
	}

	pb.RegisterRegisterLoginServer(grpcServer, authService)

	fmt.Println("gRPC server is running on port 50051")

	// Start gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
