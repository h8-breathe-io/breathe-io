package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"subs-payment-service/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type User struct {
	ID          int
	Username    string
	Email       string
	PhoneNumber string
	Tier        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
	GetUserByID(id int) (*User, error)
	IsValidToken(token string) (*User, error)
	ValidateAndGetUser(c context.Context) (*User, error)
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
func (u *userService) GetUserByID(id int) (*User, error) {

	res, err := u.userClient.GetUser(context.TODO(), &pb.GetUserRequest{
		Id: uint64(id),
	})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:          id,
		Username:    res.Username,
		Email:       res.Email,
		PhoneNumber: res.Phonenumber,
		Tier:        res.Tier,
	}, nil
}

func (u *userService) IsValidToken(token string) (*User, error) {
	res, err := u.userClient.IsValidToken(context.TODO(), &pb.IsValidTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:          int(res.User.Id),
		Username:    res.User.Username,
		Email:       res.User.Email,
		PhoneNumber: res.User.Phonenumber,
		Tier:        res.User.Tier,
	}, nil
}

func extractAuthToken(ctx context.Context) (string, error) {
	// Retrieve metadata from the context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("no metadata found in context")
		return "", fmt.Errorf("no metadata found in context")
	}

	// Check if 'auth_token' exists in the metadata
	tokens := md["auth_token"]
	if len(tokens) == 0 {
		log.Print("auth_token not found in metadata")
		return "", fmt.Errorf("auth_token not found in metadata")
	}

	log.Printf("auth_token found '%s'", tokens[0])
	// Return the first token (in case there are multiple)
	return tokens[0], nil
}

func (u *userService) ValidateAndGetUser(c context.Context) (*User, error) {
	// extract token
	token, err := extractAuthToken(c)
	if err != nil {
		return nil, err
	}

	// call user Service
	user, err := u.IsValidToken(token)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Tier:        user.Tier,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
