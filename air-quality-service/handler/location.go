package handler

import (
	"air-quality-service/model"
	pb "air-quality-service/pb/generated"
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type LocationHandler struct {
	pb.UnimplementedLocationServiceServer
	db *gorm.DB
}

func NewLocationHandler(db *gorm.DB) *LocationHandler {
	return &LocationHandler{
		db: db,
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
