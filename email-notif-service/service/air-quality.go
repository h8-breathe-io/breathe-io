package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"email-notif-service/entity"
	"email-notif-service/pb"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AirQualityService interface {
	GetAirQualityByID(id int) (*entity.AirQuality, error)
}

func NewAirQualityClient() pb.AirQualityServiceClient {
	addr := os.Getenv("AIR_QUALITY_SERVICE_URL")
	log.Printf("air-quality service url: %s", addr)
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

	client := pb.NewAirQualityServiceClient(conn)

	return client
}

func NewAirQualityService() AirQualityService {
	return &airQualityService{
		client: NewAirQualityClient(),
	}
}

type airQualityService struct {
	client pb.AirQualityServiceClient
}

// GetPaymentByID implements SubsPaymentService.
func (as *airQualityService) GetAirQualityByID(id int) (*entity.AirQuality, error) {
	res, err := as.client.GetAirQualityByID(context.TODO(), &pb.GetAirQualityByIDReq{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	return &entity.AirQuality{
		Id:         res.AirQuality.Id,
		LocationId: res.AirQuality.LocationId,
		Aqi:        res.AirQuality.Aqi,
		Co:         res.AirQuality.Co,
		No:         res.AirQuality.No,
		No2:        res.AirQuality.No2,
		O3:         res.AirQuality.O3,
		So2:        res.AirQuality.So2,
		Pm25:       res.AirQuality.Pm25,
		Pm10:       res.AirQuality.Pm10,
		Nh3:        res.AirQuality.Nh3,
		FetchTime:  res.AirQuality.FetchTime,
	}, nil
}
