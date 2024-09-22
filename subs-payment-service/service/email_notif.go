package service

import (
	"context"
	"log"
	"os"
	"subs-payment-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailNotifService interface {
	NotifyPaymentSucccess(paymentID int)
}

func NewEmailNotifClient() pb.EmailNotifServiceClient {
	addr := os.Getenv("EMAIL_NOTIF_URL")
	log.Printf("email-notif service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewEmailNotifServiceClient(conn)

	return client
}

func NewEmailNotifService() EmailNotifService {
	return &emailNotifService{
		emailNotifClient: NewEmailNotifClient(),
	}
}

type emailNotifService struct {
	emailNotifClient pb.EmailNotifServiceClient
}

func (es *emailNotifService) NotifyPaymentSucccess(paymentID int) {
	go func() {
		res, err := es.emailNotifClient.NotifyPaymentComplete(context.TODO(), &pb.NotifyPaymentCompleteReq{PaymentId: int64(paymentID)})
		if err != nil {
			log.Printf("notify payment complete email failed: %s ", err.Error())
			return
		}
		log.Printf("notify payment complete email complete: %v ", res)
	}()
}
