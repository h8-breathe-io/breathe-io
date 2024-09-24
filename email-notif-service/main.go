package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "email-notif-service/pb"
	"email-notif-service/server"
	"email-notif-service/service"

	// _ "api-gateway/docs"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()
	// instantiate services
	paymentService := service.NewSubsPaymentService()
	userService := service.NewUserService()
	emailService := service.NewEmailService(os.Getenv("MAIL_API_URL"))

	// email notif grpc server handler
	emailNotifServer := server.NewEmailNotifServer(
		paymentService,
		userService,
		emailService,
		service.NewAirQualityService(),
	)

	opts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		// grpc.UnaryInterceptor(jwtIntercept.EnsureValidToken),
	}
	port := os.Getenv("LISTEN_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEmailNotifServiceServer(grpcServer, emailNotifServer)
	log.Printf("starting gRPC server on %s", port)
	log.Fatal(grpcServer.Serve(lis))
}
