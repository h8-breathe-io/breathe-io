package service

import (
	"email-notif-service/entity"
	"email-notif-service/pb"
	"log"
	"os"
	"time"

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
	//TODO
	// return dummy user for now
	return &entity.User{
		ID:          1,
		Username:    "Razif",
		Email:       "razif.dev@gmail.com",
		PhoneNumber: "12345",
		Tier:        "free",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
