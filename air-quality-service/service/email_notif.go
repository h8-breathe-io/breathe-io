package service

import (
	pb "air-quality-service/pb/generated"
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type EmailNotifService interface {
	NotifyPaymentSucccess(paymentID int)
	NotifyRegister(userID int)
	NotifyAirQuality(userID int, airQualityID int)
	NotifyAirQualityBusiness(businessID int, airQualityID int)
}

func NewEmailNotifClient() pb.EmailNotifServiceClient {
	addr := os.Getenv("EMAIL_NOTIF_URL")
	log.Printf("email-notif service url: %s", addr)
	// Set up a connection to the server.
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.NewClient(addr, opts...)
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

// NotifyAirQualityBusiness implements EmailNotifService.
func (es *emailNotifService) NotifyAirQualityBusiness(businessID int, airQualityID int) {
	go func() {
		res, err := es.emailNotifClient.NotifyAirQualityBusiness(context.TODO(), &pb.NotifyAirQualityBusinessReq{
			BusinessId:   int64(businessID),
			AirQualityId: int64(airQualityID),
		})
		if err != nil {
			log.Printf("notify aq business email failed: %s ", err.Error())
			return
		}
		log.Printf("notify aq business emailcomplete: %v ", res)
	}()
}

// NotifyAirQuality implements EmailNotifService.
func (es *emailNotifService) NotifyAirQuality(userID int, airQualityID int) {
	go func() {
		res, err := es.emailNotifClient.NotifyAirQuality(context.TODO(), &pb.NotifyAirQualityReq{
			UserId:       int64(userID),
			AirQualityId: int64(airQualityID),
		})
		if err != nil {
			log.Printf("notify air quality email failed: %s ", err.Error())
			return
		}
		log.Printf("notify air quality email complete: %v ", res)
	}()
}

// NotifyRegister implements EmailNotifService.
func (es *emailNotifService) NotifyRegister(userID int) {
	go func() {
		res, err := es.emailNotifClient.NotifyRegister(context.TODO(), &pb.NotifyRegisterReq{
			UserId: int64(userID),
		})
		if err != nil {
			log.Printf("notify register email failed: %s ", err.Error())
			return
		}
		log.Printf("notify register emailcomplete: %v ", res)
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
