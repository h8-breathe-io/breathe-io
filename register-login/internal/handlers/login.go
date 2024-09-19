package handlers

import (
	"context"
	"fmt"
	"register-login/internal/helpers"
	"register-login/internal/models"
	"register-login/proto/pb"
)

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	// Find user by email
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check password
	if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("incorrect password")
	}

	// No token generation here, just return success message
	return &pb.LoginResponse{
		Message: "Login successful",
	}, nil
}
