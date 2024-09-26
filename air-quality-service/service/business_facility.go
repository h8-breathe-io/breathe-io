package service

import (
	"air-quality-service/entity"
	pb "air-quality-service/pb/generated"
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewBusinessFacilityClient() pb.BusinessFacilitiesClient {
	addr := os.Getenv("BUSINESS_FACILITY_SERVICE_URL")
	log.Printf("business facility service url: %s", addr)
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewBusinessFacilitiesClient(conn)

	return client
}

type BusinessFacilityService interface {
	GetBusinessFacilityByID(ctx context.Context, id int) (*entity.BusinessFacility, error)
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
func (b *businessFacilityService) GetBusinessFacilityByID(ctx context.Context, id int) (*entity.BusinessFacility, error) {

	res, err := b.businessFacilityClient.GetBusinessFacility(ctx, &pb.GetBFRequest{
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
