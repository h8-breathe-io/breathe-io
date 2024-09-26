package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"reporting/internal/database"
	"reporting/internal/handlers"
	"reporting/internal/service"
	"reporting/proto/pb"

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

	// Get the database connection from the global variable
	db := database.DB

	// Set up gRPC server
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Initialize the ReportService with the database connection
	reportService := &handlers.ReportService{DB: db, UserService: service.NewUserService()}

	// Register the ReportService with the gRPC server
	pb.RegisterReportServiceServer(grpcServer, reportService)

	// Start the gRPC server
	log.Printf("gRPC server started on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
