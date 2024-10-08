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

type SubsPaymentService interface {
	GetPaymentByID(id int) (*entity.Payment, error)
}

func NewSubsPaymentClient() pb.SubPaymentClient {
	addr := os.Getenv("SUBS_PAYMENT_SERVICE_URL")
	log.Printf("subs-payment service url: %s", addr)
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

	client := pb.NewSubPaymentClient(conn)

	return client
}

func NewSubsPaymentService() SubsPaymentService {
	return &subsPaymentService{
		client: NewSubsPaymentClient(),
	}
}

type subsPaymentService struct {
	client pb.SubPaymentClient
}

// GetPaymentByID implements SubsPaymentService.
func (s *subsPaymentService) GetPaymentByID(id int) (*entity.Payment, error) {
	res, err := s.client.GetPaymentByID(context.TODO(), &pb.GetPaymentByIDReq{PaymentId: int64(id)})
	if err != nil {
		return nil, err
	}
	asTime := res.Payment.TransactionDate.AsTime()
	return &entity.Payment{
		Id:              res.Payment.Id,
		UserId:          res.Payment.UserId,
		PaymentGateway:  res.Payment.PaymentGateway,
		Amount:          res.Payment.Amount,
		Currency:        res.Payment.Currency,
		TransactionDate: &asTime,
		Status:          res.Payment.Status,
		Url:             res.Payment.Url,
	}, nil
}
