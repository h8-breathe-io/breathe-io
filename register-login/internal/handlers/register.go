package handlers

import (
	"context"
	"fmt"
	"register-login/internal/helpers"
	"register-login/internal/models"
	"register-login/proto/pb"

	"gorm.io/gorm"
)

type AuthServiceServer struct {
	pb.UnimplementedRegisterLoginServer
	DB *gorm.DB
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Hash the password
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user model
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		Phonenumber:  req.Phonenumber,
		PasswordHash: hashedPassword,
		Tier:         "free",
	}

	// Save to DB
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &pb.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}
