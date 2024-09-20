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

//	@title			H8 P2 Final Project App
//	@version		1.0
//	@description	Hacktiv8 Phase 2 Final Project

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securitydefinitions.basic	BasicAuth
//	@tokenUrl					https://localhost:8080/users/login
//	@scope.read					Grants read access
//	@scope.write				Grants write access

func main() {
	db := config.CreateDBInstance()

	// instantiate dependencies
	emailNotifService := service.NewEmailNotifService()
	invoiceService := service.NewInvoiceService()
	// payments, for call backs by xendit
	paymentServer := server.NewPaymentServer(db, emailNotifService, invoiceService)
	// swagger docs
	// e.GET("/swagger/*", echoSwagger.WrapHandler)

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
