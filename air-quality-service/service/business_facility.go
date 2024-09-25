package service

import (
	"air-quality-service/entity"
	pb "air-quality-service/pb/generated"
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewBusinessFacilityClient() pb.BusinessFacilitiesClient {
	addr := os.Getenv("BUSINESS_FACILITY_SERVICE_URL")
	log.Printf("business facility service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewBusinessFacilitiesClient(conn)

	return client
}

type BusinessFacilityService interface {
	GetBusinessFacilityByID(id int) (*entity.BusinessFacility, error)
}

func NewBusinessFacilityService() BusinessFacilityService {
	return &businessFacilityService{
		businessFacilityClient: NewBusinessFacilityClient(),
	}
}

type businessFacilityService struct {
	businessFacilityClient pb.BusinessFacilitiesClient
}

// GetBusinessFacilityByID implements BusinessFacilityService.
func (b *businessFacilityService) GetBusinessFacilityByID(id int) (*entity.BusinessFacility, error) {

	res, err := b.businessFacilityClient.GetBusinessFacility(context.TODO(), &pb.GetBFRequest{
		Id: uint64(id),
	})
	if err != nil {
		return nil, err
	}

	return &entity.BusinessFacility{
		ID:            id,
		UserID:        res.UserId,
		CompanyType:   res.CompanyType,
		TotalEmission: res.TotalEmission,
		LocationID:    res.LocationId,
	}, nil
}
