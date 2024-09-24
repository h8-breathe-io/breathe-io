package service

import (
	"context"
	"log"
	"os"
	"subs-payment-service/entity"
	"subs-payment-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient() pb.UserClient {
	addr := os.Getenv("USER_SERVICE_URL")
	log.Printf("user service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	res, err := u.userClient.GetUser(context.TODO(), &pb.GetUserRequest{
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

func (u *userService) IsValidToken(token string) (*entity.User, error) {
	res, err := u.userClient.IsValidToken(context.TODO(), &pb.IsValidTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:          int(res.User.Id),
		Username:    res.User.Username,
		Email:       res.User.Email,
		PhoneNumber: res.User.Phonenumber,
		Tier:        res.User.Tier,
	}, nil
}
