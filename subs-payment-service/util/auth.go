package util

import (
	"context"
	"fmt"
	"log"
	"subs-payment-service/entity"
	"subs-payment-service/service"

	"google.golang.org/grpc/metadata"
)

// Function to extract 'auth_token' from gRPC metadata
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

func ValidateAndGetUser(c context.Context, userService service.UserService) (*entity.User, error) {
	// extract token
	token, err := extractAuthToken(c)
	if err != nil {
		return nil, err
	}

	// call user Service
	user, err := userService.IsValidToken(token)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Tier:        user.Tier,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
