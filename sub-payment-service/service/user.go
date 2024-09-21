package service

import (
	"sub-payment-service/entity"
	"time"
)

type UserService interface {
	GetUserByID(id int) (*entity.User, error)
}

func NewUserService() UserService {
	return &userService{}
}

type userService struct {
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(id int) (*entity.User, error) {
	//TODO
	// return dummy user for now
	return &entity.User{
		ID:          1,
		Username:    "dummyuser",
		Email:       "dummyuser@mail.com",
		PhoneNumber: "12345",
		Tier:        "free",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
