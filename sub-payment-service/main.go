package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"sub-payment-service/config"
	pb "sub-payment-service/pb"
	"sub-payment-service/server"
	"sub-payment-service/service"

	// _ "h8-p2-finalproj-app/docs"

	"google.golang.org/grpc"
)

func main() {
	db := config.CreateDBInstance()

	// instantiate services
	emailNotifService := service.NewEmailNotifService()
	invoiceService := service.NewInvoiceService()
	userService := service.NewUserService()

	// subs-payments grpc server handler
	paymentServer := server.NewPaymentServer(
		db,
		emailNotifService,
		invoiceService,
		userService,
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
	pb.RegisterSubPaymentServer(grpcServer, paymentServer)
	log.Printf("starting gRPC server on %s", port)
	log.Fatal(grpcServer.Serve(lis))
}
