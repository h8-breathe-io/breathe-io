package service

import (
	"crypto/tls"
	"crypto/x509"
	"email-notif-service/entity"
	"email-notif-service/pb"
	"email-notif-service/util"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewUserClient() pb.UserClient {
	addr := os.Getenv("USER_SERVICE_URL")
	log.Printf("user service url: %s", addr)
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

	client := pb.NewUserClient(conn)

	return client
}

type UserService interface {
	GetUserByID(id int) (*entity.User, error)
}

func NewUserService() UserService {
	return &userService{
		userClient: NewUserClient(),
	}
}

type userService struct {
	userClient pb.UserClient
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(id int) (*entity.User, error) {

	ctx := util.CreateServiceContext()
	res, err := u.userClient.GetUser(ctx, &pb.GetUserRequest{
		Id: uint64(id),
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:          id,
		Username:    res.Username,
		Email:       res.Email,
		PhoneNumber: res.Phonenumber,
		Tier:        res.Tier,
	}, nil
}
