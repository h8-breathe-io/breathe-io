package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"user-service/model"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type UserService interface {
	IsValidToken(token string) (*model.User, error)
	ValidateAndGetUser(c context.Context) (*model.User, error)
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	db *gorm.DB
}

func (u *userService) IsValidToken(token string) (*model.User, error) {
	//validate requests
	if token == "" {
		return nil, errors.New("token is required")
	}

	//validate token
	key := os.Getenv("JWT_SECRET")
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token '%s',", key)
	}

	if !t.Valid {
		return nil, fmt.Errorf("invalid token '%v',", t)
	}

	var user model.User
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token, invalid claims")
	}

	err = u.db.Where("email = ?", claims["email"]).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
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

func (u *userService) ValidateAndGetUser(c context.Context) (*model.User, error) {
	// extract token
	token, err := extractAuthToken(c)
	if err != nil {
		return nil, err
	}

	if token == os.Getenv("SERVICE_TOKEN") {
		log.Print("Service account found")
		// TODO return dummy user?
		return &model.User{
			Username: "ServiceAccount",
		}, nil
	}

	// call user Service
	user, err := u.IsValidToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}
