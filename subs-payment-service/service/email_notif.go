package service

import (
	"context"
	"log"
	"os"
	pb "subs-payment-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailNotifService interface {
	NotifyPaymentSucccess(paymentID int)
	NotifyRegister(userID int)
	NotifyAirQuality(userID int, airQualityID int)
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

// NotifyAirQuality implements EmailNotifService.
func (es *emailNotifService) NotifyAirQuality(userID int, airQualityID int) {
	go func() {
		es.emailNotifClient.NotifyAirQuality(context.TODO(), &pb.NotifyAirQualityReq{
			UserId:       int64(userID),
			AirQualityId: int64(airQualityID),
		})
	}()
}

// NotifyRegister implements EmailNotifService.
func (es *emailNotifService) NotifyRegister(userID int) {
	go func() {
		es.emailNotifClient.NotifyRegister(context.TODO(), &pb.NotifyRegisterReq{
			UserId: int64(userID),
		})
	}()
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
