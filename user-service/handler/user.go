package handler

import (
	"context"
	"errors"
	"os"
	"user-service/model"
	pb "user-service/pb/generated"
	"user-service/util"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserHandler struct {
	pb.UnimplementedUserServer
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	//validate requests
	if req.Username == "" {
		return nil, errors.New("username is required")
	}

	if req.Phonenumber == "" {
		return nil, errors.New("phonenumber is required")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	// check if user already exists
	var user model.User
	err := u.db.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		return nil, errors.New("user already exists")
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// create new user
	user = model.User{
		Username:     req.Username,
		Phonenumber:  req.Phonenumber,
		Email:        req.Email,
		PasswordHash: hashed,
		Tier:         "free",
	}

	err = u.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Tier:        user.Tier,
	}, nil
}

func (u *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//validate requests
	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	// check if user exists
	var user model.User
	err := u.db.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !util.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	key := os.Getenv("JWT_SECRET")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.Email,
		"tier":     user.Tier,
		"username": user.Username,
		"id":       user.ID,
	})

	s, err := t.SignedString([]byte(key))
	if err != nil {
		return nil, errors.New("failed to sign token")
	}

	return &pb.LoginResponse{
		Token: s,
	}, nil
}

func (u *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("id is required")
	}

	var user model.User
	err := u.db.Where("id = ?", req.Id).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &pb.UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Tier:        user.Tier,
	}, nil
}

func (u *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("id is required")
	}

	var user model.User
	err := u.db.Where("id = ?", req.Id).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Username == "" {
		req.Username = user.Username
	}

	if req.Phonenumber == "" {
		req.Phonenumber = user.Phonenumber
	}

	if req.Tier == "" {
		req.Tier = user.Tier
	}

	if req.Email == "" {
		req.Email = user.Email
	}

	//check if email is already used by another user and id is not the same
	var user2 model.User
	err = u.db.Where("email = ? AND id != ?", req.Email, req.Id).First(&user2).Error
	if err == nil {
		return nil, errors.New("email already used by another user")
	}

	if req.Password == "" {
		req.Password = user.PasswordHash
	} else {
		hashed, err := util.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		req.Password = hashed
	}

	updateUser := model.User{
		Username:     req.Username,
		Phonenumber:  req.Phonenumber,
		Email:        req.Email,
		PasswordHash: req.Password,
		Tier:         req.Tier,
	}

	err = u.db.Where("id = ?", req.Id).Updates(&updateUser).Error
	if err != nil {
		return nil, errors.New("failed to update user " + err.Error())
	}

	return &pb.UserResponse{
		Username:    updateUser.Username,
		Email:       updateUser.Email,
		Phonenumber: updateUser.Phonenumber,
		Tier:        updateUser.Tier,
	}, nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.UserResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("id is required")
	}

	var user model.User
	err := u.db.Where("id = ?", req.Id).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = u.db.Delete(&user).Error
	if err != nil {
		return nil, errors.New("failed to delete user")
	}

	return &pb.UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Tier:        user.Tier,
	}, nil
}

func (u *UserHandler) IsValidToken(ctx context.Context, req *pb.IsValidTokenRequest) (*pb.IsValidTokenResponse, error) {
	//validate requests
	if req.Token == "" {
		return nil, errors.New("token is required")
	}

	//validate token
	key := os.Getenv("JWT_SECRET")
	t, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	var user model.User
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	err = u.db.Where("email = ?", claims["email"]).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	res := &pb.IsValidTokenResponse{
		Valid: true,
		User: &pb.UserResponse{
			Username:    user.Username,
			Email:       user.Email,
			Phonenumber: user.Phonenumber,
			Tier:        user.Tier,
		},
	}

	return res, nil
}
