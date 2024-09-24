package service

import (
	"context"
	"email-notif-service/entity"
	"email-notif-service/pb"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AirQualityService interface {
	GetAirQualityByID(id int) (*entity.AirQuality, error)
}

func NewAirQualityClient() pb.AirQualityServiceClient {
	addr := os.Getenv("AIR_QUALITY_SERVICE_URL")
	log.Printf("air-quality service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	var fetchTime time.Time
	if res.AirQuality.FetchTime != nil {
		fetchTime = res.AirQuality.FetchTime.AsTime()
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
		FetchTime:  &fetchTime,
	}, nil
}
