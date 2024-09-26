package handler

import (
	"context"
	"errors"
	"user-service/model"
	pb "user-service/pb/generated"
	"user-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type BusinessFacilityHandler struct {
	pb.UnimplementedBusinessFacilitiesServer
	db          *gorm.DB
	ls          pb.LocationServiceClient
	userService service.UserService
}

func NewBusinessFacilitiesHandler(db *gorm.DB, ls pb.LocationServiceClient) *BusinessFacilityHandler {
	return &BusinessFacilityHandler{
		db:          db,
		ls:          ls,
		userService: service.NewUserService(db),
	}
}

func (bf *BusinessFacilityHandler) AddBusinessFacility(ctx context.Context, req *pb.AddBFRequest) (*pb.BFResponse, error) {
	//validate requests
	// if req.UserId == 0 {
	// 	return nil, errors.New("user id is required")
	// }

	if req.CompanyType == "" {
		return nil, errors.New("company Type is required")
	}

	if req.TotalEmission == 0 {
		return nil, errors.New("total emission is required")
	}

	if req.LocationId == 0 {
		return nil, errors.New("location id is required")
	}

	// validate token and get user
	user, err := bf.userService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%d': %s", req.UserId, err.Error())
	}

	//check if location id valid
	location, err := bf.ls.GetLocation(ctx, &pb.GetLocationRequest{LocationId: req.LocationId})
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "location not found: %s", err.Error())
	}

	// create new BusinessFacility
	businessFacility := model.BusinessFacility{
		UserID:        uint64(user.ID),
		CompanyType:   req.CompanyType,
		TotalEmission: req.TotalEmission,
		LocationID:    location.LocationId,
	}

	err = bf.db.Create(&businessFacility).Error
	if err != nil {
		return nil, err
	}

	return &pb.BFResponse{
		Id:            uint64(businessFacility.ID),
		UserId:        businessFacility.UserID,
		CompanyType:   businessFacility.CompanyType,
		TotalEmission: businessFacility.TotalEmission,
		LocationId:    businessFacility.LocationID,
	}, nil
}

func (bf *BusinessFacilityHandler) GetBusinessFacilities(ctx context.Context, req *pb.GetBFRequests) (*pb.BFResponses, error) {
	//validate requests
	// if req.UserId == 0 {
	// 	return nil, errors.New("user_id is required")
	// }

	// validate token and get user
	user, err := bf.userService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}

	//get all business facilities of the user
	var businessFacilities []model.BusinessFacility
	err = bf.db.Where("user_id = ?", user.ID).Find(&businessFacilities).Error
	if err != nil {
		return nil, errors.New("failed to get business facilities")
	}

	res := &pb.BFResponses{
		BusinessFacilities: make([]*pb.BFResponse, 0, len(businessFacilities)),
	}

	for _, businessFacility := range businessFacilities {
		res.BusinessFacilities = append(res.BusinessFacilities, &pb.BFResponse{
			Id:            uint64(businessFacility.ID),
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

	// validate token and get user
	user, err := bf.userService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}

	//check if business facility exists
	var businessFacility model.BusinessFacility
	err = bf.db.Where("id = ?", req.Id).First(&businessFacility).Error
	if err != nil {
		return nil, errors.New("business facility not found")
	}

	// ensure business facility belongs to user
	if businessFacility.UserID != uint64(user.ID) {
		return nil, errors.New("business facility doesn't belong to user")
	}

	return &pb.BFResponse{
		Id:            uint64(businessFacility.ID),
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

	// if req.UserId == 0 {
	// 	req.UserId = businessFacility.UserID
	// }

	if req.CompanyType == "" {
		req.CompanyType = businessFacility.CompanyType
	}

	if req.TotalEmission == 0 {
		req.TotalEmission = businessFacility.TotalEmission
	}

	if req.LocationId == 0 {
		req.LocationId = businessFacility.LocationID
	}

	// validate token and get user
	user, err := bf.userService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}

	// ensure business facility belongs to user
	if businessFacility.UserID != uint64(user.ID) {
		return nil, errors.New("business facility doesn't belong to user")
	}

	//check if location id valid
	//check if location id valid
	_, err = bf.ls.GetLocation(ctx, &pb.GetLocationRequest{LocationId: req.LocationId})
	if err != nil {
		return nil, errors.New("location not found")
	}

	updateBusinessFacility := model.BusinessFacility{
		UserID:        uint64(user.ID),
		CompanyType:   req.CompanyType,
		TotalEmission: req.TotalEmission,
		LocationID:    req.LocationId,
	}

	err = bf.db.Where("id = ?", req.Id).Updates(&updateBusinessFacility).Error
	if err != nil {
		return nil, errors.New("failed to update business facility")
	}

	return &pb.BFResponse{
		Id:            uint64(req.Id),
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

	// validate token and get user
	user, err := bf.userService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}
	// ensure business facility belongs to user
	if businessFacility.UserID != uint64(user.ID) {
		return nil, errors.New("business facility doesn't belong to user")
	}

	err = bf.db.Delete(&businessFacility).Error
	if err != nil {
		return nil, errors.New("failed to delete business facility")
	}

	return &pb.BFResponse{
		Id:            uint64(businessFacility.ID),
		UserId:        businessFacility.UserID,
		CompanyType:   businessFacility.CompanyType,
		TotalEmission: businessFacility.TotalEmission,
		LocationId:    businessFacility.LocationID,
	}, nil
}

func (bf *BusinessFacilityHandler) GetCarbonTax(ctx context.Context, req *pb.GetCarbonTaxRequest) (*pb.GetCarbonTaxResponse, error) {
	//validate requests
	if req.Id == 0 {
		return nil, errors.New("business id is required")
	}

	//check if business facility exists
	var businessFacility model.BusinessFacility
	err := bf.db.Where("id = ?", req.Id).First(&businessFacility).Error
	if err != nil {
		return nil, errors.New("business facility not found")
	}

	//carbonTaxRate for indonesia per ton
	carbonTaxRate := 2

	//get carbon tax value for indonesia
	carbonTax := businessFacility.TotalEmission * float64(carbonTaxRate)

	return &pb.GetCarbonTaxResponse{
		Currency:  "USD",
		CarbonTax: carbonTax,
	}, nil
}
