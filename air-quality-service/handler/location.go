package handler

import (
	"air-quality-service/model"
	pb "air-quality-service/pb/generated"
	"air-quality-service/service"
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type LocationHandler struct {
	pb.UnimplementedLocationServiceServer
	db                      *gorm.DB
	userService             service.UserService
	businessFacilityService service.BusinessFacilityService
}

func NewLocationHandler(db *gorm.DB, userService service.UserService, businessFacilityService service.BusinessFacilityService) *LocationHandler {
	return &LocationHandler{
		db:                      db,
		userService:             userService,
		businessFacilityService: businessFacilityService,
	}
}

func (lh *LocationHandler) AddLocation(ctx context.Context, req *pb.AddLocationRequest) (*pb.LocationResponse, error) {
	//validate request
	// check if location name is empty
	if req.LocationName == "" {
		return nil, errors.New("location name is required")
	}

	// check if latitude is off limit
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, errors.New("latitude value is off limit")
	}

	// check if longitude is off limit
	if req.Longitude == -180 || req.Longitude > 180 {
		return nil, errors.New("longitude value is off limit")
	}

	// save location to db
	location := model.Location{
		LocationName: req.LocationName,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
	}
	err := lh.db.Create(&location).Error
	if err != nil {
		return nil, errors.New("failed to create location")
	}

	return &pb.LocationResponse{
		LocationId:   uint64(location.ID),
		LocationName: location.LocationName,
		Latitude:     location.Latitude,
		Longitude:    location.Longitude,
	}, nil
}

func (lh *LocationHandler) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.LocationResponse, error) {
	//validate ID if valid
	if req.LocationId == 0 {
		return nil, errors.New("location id is required")
	}

	var location model.Location
	err := lh.db.Where("id=?", req.LocationId).First(&location).Error
	if err != nil {
		return nil, errors.New("location not found")
	}

	return &pb.LocationResponse{
		LocationId:   uint64(location.ID),
		LocationName: location.LocationName,
		Latitude:     location.Latitude,
		Longitude:    location.Longitude,
	}, nil
}

func (lh *LocationHandler) GetLocations(ctx context.Context, req *emptypb.Empty) (*pb.GetLocationsResponse, error) {
	var locations []model.Location
	err := lh.db.Find(&locations).Error
	if err != nil {
		return nil, errors.New("failed to get locations")
	}

	var locationResponses []*pb.LocationResponse
	for _, location := range locations {
		locationResponses = append(locationResponses, &pb.LocationResponse{
			LocationId:   uint64(location.ID),
			LocationName: location.LocationName,
			Latitude:     location.Latitude,
			Longitude:    location.Longitude,
		})
	}

	return &pb.GetLocationsResponse{
		Locations: locationResponses,
	}, nil
}

func (lh *LocationHandler) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.LocationResponse, error) {
	//search if exists
	var location model.Location
	err := lh.db.Where("id=?", req.LocationId).First(&location).Error
	if err != nil {
		return nil, errors.New("location not found")
	}

	//validate request
	// check if location name is empty
	if req.LocationName == "" {
		return nil, errors.New("location name is required")
	}

	// check if latitude is off limit
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, errors.New("latitude value is off limit")
	}

	// check if longitude is off limit
	if req.Longitude == -180 || req.Longitude > 180 {
		return nil, errors.New("longitude value is off limit")
	}

	// update location to db
	updateLocation := model.Location{
		LocationName: req.LocationName,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
	}
	err = lh.db.Where("id=?", req.LocationId).Updates(&updateLocation).Error
	if err != nil {
		return nil, errors.New("failed to update location")
	}

	return &pb.LocationResponse{
		LocationId:   uint64(location.ID),
		LocationName: updateLocation.LocationName,
		Latitude:     updateLocation.Latitude,
		Longitude:    updateLocation.Longitude,
	}, nil
}

func (lh *LocationHandler) DeleteLocation(ctx context.Context, req *pb.DeleteLocationRequest) (*pb.LocationResponse, error) {
	//validate ID if valid
	if req.LocationId == 0 {
		return nil, errors.New("location id is required")
	}

	//search if exists
	var location model.Location
	err := lh.db.Where("id=?", req.LocationId).First(&location).Error
	if err != nil {
		return nil, errors.New("location not found")
	}

	// delete location to db
	err = lh.db.Where("id=?", req.LocationId).Delete(&location).Error
	if err != nil {
		return nil, errors.New("failed to delete location")
	}

	return &pb.LocationResponse{
		LocationId:   uint64(location.ID),
		LocationName: location.LocationName,
		Latitude:     location.Latitude,
		Longitude:    location.Longitude,
	}, nil
}

func (lh *LocationHandler) GetLocationRecommendation(ctx context.Context, req *pb.LocationRecommendationRequest) (*pb.LocationRecommendationResponse, error) {
	//validate the request
	if req.BusinessId == 0 {
		return nil, errors.New("business facility id is required")
	}

	// validate token and get user
	user, err := lh.userService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}

	//get initial location where business facility is located
	// need to attach token to outgoing context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("failed to extract metadata")
	}

	// Step 2: Create a new outgoing context with the extracted metadata
	outCtx := metadata.NewOutgoingContext(ctx, md)
	businessFacility, err := lh.businessFacilityService.GetBusinessFacilityByID(outCtx, int(req.BusinessId))
	if err != nil {
		return nil, errors.New("business facility not found" + err.Error())
	}

	// ensure business facility belongs to user
	if businessFacility.UserID != uint64(user.ID) {
		return nil, errors.New("business facility doesn't belong to user")
	}

	//find the location based on location ID provided
	var initialLocation model.Location
	err = lh.db.Where("id=?", businessFacility.LocationID).First(&initialLocation).Error
	if err != nil {
		return nil, errors.New("location not found" + err.Error())
	}

	//get all locations
	var candidateLocations []model.Location
	err = lh.db.Preload("AirQualities").Find(&candidateLocations).Error
	if err != nil {
		return nil, errors.New("failed to get locations")
	}

	var bestLocations []model.Location
	var bestScore float64

	fmt.Println("===============================================")
	for i, candidateLocation := range candidateLocations {
		candidateLocationScore := 0.0
		for _, airQuality := range candidateLocation.AirQualities {
			// Calculate feasibility score based on various pollutants
			score := (1.0 / float64(airQuality.AQI+1)) * 0.4 // Weight for AQI
			score += (1.0 / (airQuality.CO + 1)) * 0.1       // Weight for CO
			score += (1.0 / (airQuality.NO2 + 1)) * 0.1      // Weight for NO2
			score += (1.0 / (airQuality.O3 + 1)) * 0.1       // Weight for O3
			score += (1.0 / (airQuality.SO2 + 1)) * 0.1      // Weight for SO2
			score += (1.0 / (airQuality.PM25 + 1)) * 0.1     // Weight for PM2.5
			score += (1.0 / (airQuality.PM10 + 1)) * 0.05    // Weight for PM10
			score += (1.0 / (airQuality.NH3 + 1)) * 0.05     // Weight for NH3
			candidateLocationScore += score                  // Aggregate scores for all air qualities
		}
		fmt.Printf("Calculating air quality score for ID: %d, air quality score is %f \n", candidateLocations[i].ID, candidateLocationScore)
		// Update best score and location if current location is better
		if candidateLocationScore > bestScore {
			bestScore = candidateLocationScore
			bestLocations = []model.Location{candidateLocations[i]}
		} else if candidateLocationScore == bestScore {
			bestLocations = append(bestLocations, candidateLocations[i]) // Add to list of best locations
		}
	}
	fmt.Println("===============================================")

	return &pb.LocationRecommendationResponse{
		InitialLocation: &pb.LocationResponse{
			LocationId:   uint64(initialLocation.ID),
			LocationName: initialLocation.LocationName,
			Latitude:     initialLocation.Latitude,
			Longitude:    initialLocation.Longitude,
		},
		RecommendedLocations: func() []*pb.LocationResponse {
			var recommendedLocations []*pb.LocationResponse
			for _, location := range bestLocations {
				recommendedLocations = append(recommendedLocations, &pb.LocationResponse{
					LocationId:   uint64(location.ID),
					LocationName: location.LocationName,
					Latitude:     location.Latitude,
					Longitude:    location.Longitude,
				})
			}
			return recommendedLocations
		}(),
	}, nil
}
