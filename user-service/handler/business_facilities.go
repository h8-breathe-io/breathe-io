package handler

import (
	"context"
	"errors"
	"user-service/model"
	pb "user-service/pb/generated"

	"gorm.io/gorm"
)

type BusinessFacilityHandler struct {
	pb.UnimplementedBusinessFacilitiesServer
	db *gorm.DB
}

func NewBusinessFacilitiesHandler(db *gorm.DB) *BusinessFacilityHandler {
	return &BusinessFacilityHandler{
		db: db,
	}
}

func (bf *BusinessFacilityHandler) AddBusinessFacility(ctx context.Context, req *pb.AddBFRequest) (*pb.BFResponse, error) {
	//validate requests
	if req.UserId == 0 {
		return nil, errors.New("user id is required")
	}

	if req.CompanyType == "" {
		return nil, errors.New("company Type is required")
	}

	if req.TotalEmission == 0 {
		return nil, errors.New("total emission is required")
	}

	if req.LocationId == 0 {
		return nil, errors.New("password is required")
	}

	// check if user id valid

	var user model.User
	err := bf.db.Where("id = ?", req.UserId).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	//check if location id valid
	//need location service
	// var location model.Location
	// err = bf.db.Where("id = ?", req.LocationId).First(&location).Error
	// if err != nil {
	// 	return nil, errors.New("location not found")
	// }

	// create new BusinessFacility
	businessFacility := model.BusinessFacility{
		UserID:        req.UserId,
		CompanyType:   req.CompanyType,
		TotalEmission: req.TotalEmission,
		LocationID:    req.LocationId,
	}

	err = bf.db.Create(&businessFacility).Error
	if err != nil {
		return nil, err
	}

	return &pb.BFResponse{
		UserId:        businessFacility.UserID,
		CompanyType:   businessFacility.CompanyType,
		TotalEmission: businessFacility.TotalEmission,
		LocationId:    businessFacility.LocationID,
	}, nil
}

func (bf *BusinessFacilityHandler) GetBusinessFacilities(ctx context.Context, req *pb.GetBFRequests) (*pb.BFResponses, error) {
	//validate requests
	if req.UserId == 0 {
		return nil, errors.New("user_id is required")
	}

	//check if user exists

	var user model.User
	err := bf.db.Where("id = ?", req.UserId).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	//get all business facilities of the user
	var businessFacilities []model.BusinessFacility
	err = bf.db.Where("user_id = ?", req.UserId).Find(&businessFacilities).Error
	if err != nil {
		return nil, errors.New("failed to get business facilities")
	}

	res := &pb.BFResponses{
		BusinessFacilities: make([]*pb.BFResponse, 0, len(businessFacilities)),
	}

	for _, businessFacility := range businessFacilities {
		res.BusinessFacilities = append(res.BusinessFacilities, &pb.BFResponse{
			UserId:        businessFacility.UserID,
			CompanyType:   businessFacility.CompanyType,
			TotalEmission: businessFacility.TotalEmission,
			LocationId:    businessFacility.LocationID,
		})
	}

	return res, nil
}

func (bf *BusinessFacilityHandler) GetBusinessFacility(ctx context.Context, req *pb.GetBFRequest) (*pb.BFResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("id is required")
	}

	//check if business facility exists
	var businessFacility model.BusinessFacility
	err := bf.db.Where("id = ?", req.Id).First(&businessFacility).Error
	if err != nil {
		return nil, errors.New("business facility not found")
	}

	return &pb.BFResponse{
		UserId:        businessFacility.UserID,
		CompanyType:   businessFacility.CompanyType,
		TotalEmission: businessFacility.TotalEmission,
		LocationId:    businessFacility.LocationID,
	}, nil
}

func (bf *BusinessFacilityHandler) UpdateBusinessFacility(ctx context.Context, req *pb.UpdateBFRequest) (*pb.BFResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("id is required")
	}

	var businessFacility model.BusinessFacility
	err := bf.db.Where("id = ?", req.Id).First(&businessFacility).Error
	if err != nil {
		return nil, errors.New("business facility not found")
	}

	if req.UserId == 0 {
		req.UserId = businessFacility.UserID
	}

	if req.CompanyType == "" {
		req.CompanyType = businessFacility.CompanyType
	}

	if req.TotalEmission == 0 {
		req.TotalEmission = businessFacility.TotalEmission
	}

	if req.LocationId == 0 {
		req.LocationId = businessFacility.LocationID
	}

	//check if user id valid
	var user model.User
	err = bf.db.Where("id = ?", req.UserId).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	//check if location id valid
	//need location service
	// var location model.Location
	// err = bf.db.Where("id = ?", req.LocationId).First(&location).Error
	// if err != nil {
	// 	return nil, errors.New("location not found")
	// }

	updateBusinessFacility := model.BusinessFacility{
		UserID:        req.UserId,
		CompanyType:   req.CompanyType,
		TotalEmission: req.TotalEmission,
		LocationID:    req.LocationId,
	}

	err = bf.db.Where("id = ?", req.Id).Updates(&updateBusinessFacility).Error
	if err != nil {
		return nil, errors.New("failed to update business facility")
	}

	return &pb.BFResponse{
		UserId:        updateBusinessFacility.UserID,
		CompanyType:   updateBusinessFacility.CompanyType,
		TotalEmission: updateBusinessFacility.TotalEmission,
		LocationId:    updateBusinessFacility.LocationID,
	}, nil
}

func (bf *BusinessFacilityHandler) DeleteBusinessFacility(ctx context.Context, req *pb.DeleteBFRequest) (*pb.BFResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("id is required")
	}

	var businessFacility model.BusinessFacility
	err := bf.db.Where("id = ?", req.Id).First(&businessFacility).Error
	if err != nil {
		return nil, errors.New("business facility not found")
	}

	err = bf.db.Delete(&businessFacility).Error
	if err != nil {
		return nil, errors.New("failed to delete business facility")
	}

	return &pb.BFResponse{
		UserId:        businessFacility.UserID,
		CompanyType:   businessFacility.CompanyType,
		TotalEmission: businessFacility.TotalEmission,
		LocationId:    businessFacility.LocationID,
	}, nil
}
