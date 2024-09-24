// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: air_quality.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetAirQualitiesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LocationId uint64 `protobuf:"varint,1,opt,name=location_id,json=locationId,proto3" json:"location_id,omitempty"`
}

func (x *GetAirQualitiesRequest) Reset() {
	*x = GetAirQualitiesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAirQualitiesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAirQualitiesRequest) ProtoMessage() {}

func (x *GetAirQualitiesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAirQualitiesRequest.ProtoReflect.Descriptor instead.
func (*GetAirQualitiesRequest) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{0}
}

func (x *GetAirQualitiesRequest) GetLocationId() uint64 {
	if x != nil {
		return x.LocationId
	}
	return 0
}

type GetAirQualitiesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AirQualities []*AirQuality `protobuf:"bytes,1,rep,name=air_qualities,json=airQualities,proto3" json:"air_qualities,omitempty"`
}

func (x *GetAirQualitiesResponse) Reset() {
	*x = GetAirQualitiesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAirQualitiesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAirQualitiesResponse) ProtoMessage() {}

func (x *GetAirQualitiesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAirQualitiesResponse.ProtoReflect.Descriptor instead.
func (*GetAirQualitiesResponse) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{1}
}

func (x *GetAirQualitiesResponse) GetAirQualities() []*AirQuality {
	if x != nil {
		return x.AirQualities
	}
	return nil
}

type SaveAirQualitiesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *SaveAirQualitiesRequest) Reset() {
	*x = SaveAirQualitiesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveAirQualitiesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveAirQualitiesRequest) ProtoMessage() {}

func (x *SaveAirQualitiesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveAirQualitiesRequest.ProtoReflect.Descriptor instead.
func (*SaveAirQualitiesRequest) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{2}
}

func (x *SaveAirQualitiesRequest) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *SaveAirQualitiesRequest) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type SaveAirQualitiesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *SaveAirQualitiesResponse) Reset() {
	*x = SaveAirQualitiesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveAirQualitiesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveAirQualitiesResponse) ProtoMessage() {}

func (x *SaveAirQualitiesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveAirQualitiesResponse.ProtoReflect.Descriptor instead.
func (*SaveAirQualitiesResponse) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{3}
}

func (x *SaveAirQualitiesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type AirQuality struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	LocationId int64                  `protobuf:"varint,2,opt,name=location_id,json=locationId,proto3" json:"location_id,omitempty"`
	Aqi        int64                  `protobuf:"varint,3,opt,name=aqi,proto3" json:"aqi,omitempty"`
	Co         float64                `protobuf:"fixed64,4,opt,name=co,proto3" json:"co,omitempty"`
	No         float64                `protobuf:"fixed64,5,opt,name=no,proto3" json:"no,omitempty"`
	No2        float64                `protobuf:"fixed64,6,opt,name=no2,proto3" json:"no2,omitempty"`
	O3         float64                `protobuf:"fixed64,7,opt,name=o3,proto3" json:"o3,omitempty"`
	So2        float64                `protobuf:"fixed64,8,opt,name=so2,proto3" json:"so2,omitempty"`
	Pm25       float64                `protobuf:"fixed64,9,opt,name=pm25,proto3" json:"pm25,omitempty"`
	Pm10       float64                `protobuf:"fixed64,10,opt,name=pm10,proto3" json:"pm10,omitempty"`
	Nh3        float64                `protobuf:"fixed64,11,opt,name=nh3,proto3" json:"nh3,omitempty"`
	FetchTime  *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=fetch_time,json=fetchTime,proto3" json:"fetch_time,omitempty"`
}

func (x *AirQuality) Reset() {
	*x = AirQuality{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AirQuality) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AirQuality) ProtoMessage() {}

func (x *AirQuality) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AirQuality.ProtoReflect.Descriptor instead.
func (*AirQuality) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{4}
}

func (x *AirQuality) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AirQuality) GetLocationId() int64 {
	if x != nil {
		return x.LocationId
	}
	return 0
}

func (x *AirQuality) GetAqi() int64 {
	if x != nil {
		return x.Aqi
	}
	return 0
}

func (x *AirQuality) GetCo() float64 {
	if x != nil {
		return x.Co
	}
	return 0
}

func (x *AirQuality) GetNo() float64 {
	if x != nil {
		return x.No
	}
	return 0
}

func (x *AirQuality) GetNo2() float64 {
	if x != nil {
		return x.No2
	}
	return 0
}

func (x *AirQuality) GetO3() float64 {
	if x != nil {
		return x.O3
	}
	return 0
}

func (x *AirQuality) GetSo2() float64 {
	if x != nil {
		return x.So2
	}
	return 0
}

func (x *AirQuality) GetPm25() float64 {
	if x != nil {
		return x.Pm25
	}
	return 0
}

func (x *AirQuality) GetPm10() float64 {
	if x != nil {
		return x.Pm10
	}
	return 0
}

func (x *AirQuality) GetNh3() float64 {
	if x != nil {
		return x.Nh3
	}
	return 0
}

func (x *AirQuality) GetFetchTime() *timestamppb.Timestamp {
	if x != nil {
		return x.FetchTime
	}
	return nil
}

type GetAirQualityByIDReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetAirQualityByIDReq) Reset() {
	*x = GetAirQualityByIDReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAirQualityByIDReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAirQualityByIDReq) ProtoMessage() {}

func (x *GetAirQualityByIDReq) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAirQualityByIDReq.ProtoReflect.Descriptor instead.
func (*GetAirQualityByIDReq) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{5}
}

func (x *GetAirQualityByIDReq) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetAirQualityByIDResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AirQuality *AirQuality `protobuf:"bytes,1,opt,name=air_quality,json=airQuality,proto3" json:"air_quality,omitempty"`
}

func (x *GetAirQualityByIDResp) Reset() {
	*x = GetAirQualityByIDResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_air_quality_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAirQualityByIDResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAirQualityByIDResp) ProtoMessage() {}

func (x *GetAirQualityByIDResp) ProtoReflect() protoreflect.Message {
	mi := &file_air_quality_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAirQualityByIDResp.ProtoReflect.Descriptor instead.
func (*GetAirQualityByIDResp) Descriptor() ([]byte, []int) {
	return file_air_quality_proto_rawDescGZIP(), []int{6}
}

func (x *GetAirQualityByIDResp) GetAirQuality() *AirQuality {
	if x != nil {
		return x.AirQuality
	}
	return nil
}

var File_air_quality_proto protoreflect.FileDescriptor

var file_air_quality_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x10, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x39, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72,
	0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x22, 0x5c, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69,
	0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0d,
	0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74,
	0x79, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74,
	0x79, 0x52, 0x0c, 0x61, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x22,
	0x53, 0x0a, 0x17, 0x53, 0x61, 0x76, 0x65, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74,
	0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x22, 0x34, 0x0a, 0x18, 0x53, 0x61, 0x76, 0x65, 0x41, 0x69, 0x72, 0x51,
	0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x98, 0x02, 0x0a, 0x0a, 0x41,
	0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x71,
	0x69, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x61, 0x71, 0x69, 0x12, 0x0e, 0x0a, 0x02,
	0x63, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x02, 0x63, 0x6f, 0x12, 0x0e, 0x0a, 0x02,
	0x6e, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x02, 0x6e, 0x6f, 0x12, 0x10, 0x0a, 0x03,
	0x6e, 0x6f, 0x32, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6e, 0x6f, 0x32, 0x12, 0x0e,
	0x0a, 0x02, 0x6f, 0x33, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x02, 0x6f, 0x33, 0x12, 0x10,
	0x0a, 0x03, 0x73, 0x6f, 0x32, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x73, 0x6f, 0x32,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6d, 0x32, 0x35, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04,
	0x70, 0x6d, 0x32, 0x35, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6d, 0x31, 0x30, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x04, 0x70, 0x6d, 0x31, 0x30, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x68, 0x33, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6e, 0x68, 0x33, 0x12, 0x39, 0x0a, 0x0a, 0x66, 0x65,
	0x74, 0x63, 0x68, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x66, 0x65, 0x74, 0x63,
	0x68, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x26, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x51,
	0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x56, 0x0a,
	0x15, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x42, 0x79,
	0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x12, 0x3d, 0x0a, 0x0b, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75,
	0x61, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x69,
	0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x41,
	0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x0a, 0x61, 0x69, 0x72, 0x51, 0x75,
	0x61, 0x6c, 0x69, 0x74, 0x79, 0x32, 0xd2, 0x02, 0x0a, 0x11, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61,
	0x6c, 0x69, 0x74, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x68, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x28,
	0x2e, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x61, 0x69, 0x72, 0x5f, 0x71,
	0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6b, 0x0a, 0x10, 0x53, 0x61, 0x76, 0x65, 0x41, 0x69, 0x72,
	0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x29, 0x2e, 0x61, 0x69, 0x72, 0x5f,
	0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x61, 0x76,
	0x65, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69,
	0x74, 0x79, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x41, 0x69, 0x72, 0x51,
	0x75, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x66, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c,
	0x69, 0x74, 0x79, 0x42, 0x79, 0x49, 0x44, 0x12, 0x26, 0x2e, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75,
	0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x69,
	0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x1a,
	0x27, 0x2e, 0x61, 0x69, 0x72, 0x5f, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x69, 0x72, 0x51, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79,
	0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x38, 0x2d, 0x62, 0x72, 0x65, 0x61,
	0x74, 0x68, 0x65, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x72, 0x65, 0x61, 0x74, 0x68, 0x65, 0x2d, 0x69,
	0x6f, 0x2f, 0x61, 0x69, 0x72, 0x2d, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_air_quality_proto_rawDescOnce sync.Once
	file_air_quality_proto_rawDescData = file_air_quality_proto_rawDesc
)

func file_air_quality_proto_rawDescGZIP() []byte {
	file_air_quality_proto_rawDescOnce.Do(func() {
		file_air_quality_proto_rawDescData = protoimpl.X.CompressGZIP(file_air_quality_proto_rawDescData)
	})
	return file_air_quality_proto_rawDescData
}

var file_air_quality_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_air_quality_proto_goTypes = []any{
	(*GetAirQualitiesRequest)(nil),   // 0: air_quality_grpc.GetAirQualitiesRequest
	(*GetAirQualitiesResponse)(nil),  // 1: air_quality_grpc.GetAirQualitiesResponse
	(*SaveAirQualitiesRequest)(nil),  // 2: air_quality_grpc.SaveAirQualitiesRequest
	(*SaveAirQualitiesResponse)(nil), // 3: air_quality_grpc.SaveAirQualitiesResponse
	(*AirQuality)(nil),               // 4: air_quality_grpc.AirQuality
	(*GetAirQualityByIDReq)(nil),     // 5: air_quality_grpc.GetAirQualityByIDReq
	(*GetAirQualityByIDResp)(nil),    // 6: air_quality_grpc.GetAirQualityByIDResp
	(*timestamppb.Timestamp)(nil),    // 7: google.protobuf.Timestamp
}
var file_air_quality_proto_depIdxs = []int32{
	4, // 0: air_quality_grpc.GetAirQualitiesResponse.air_qualities:type_name -> air_quality_grpc.AirQuality
	7, // 1: air_quality_grpc.AirQuality.fetch_time:type_name -> google.protobuf.Timestamp
	4, // 2: air_quality_grpc.GetAirQualityByIDResp.air_quality:type_name -> air_quality_grpc.AirQuality
	0, // 3: air_quality_grpc.AirQualityService.GetAirQualities:input_type -> air_quality_grpc.GetAirQualitiesRequest
	2, // 4: air_quality_grpc.AirQualityService.SaveAirQualities:input_type -> air_quality_grpc.SaveAirQualitiesRequest
	5, // 5: air_quality_grpc.AirQualityService.GetAirQualityByID:input_type -> air_quality_grpc.GetAirQualityByIDReq
	1, // 6: air_quality_grpc.AirQualityService.GetAirQualities:output_type -> air_quality_grpc.GetAirQualitiesResponse
	3, // 7: air_quality_grpc.AirQualityService.SaveAirQualities:output_type -> air_quality_grpc.SaveAirQualitiesResponse
	6, // 8: air_quality_grpc.AirQualityService.GetAirQualityByID:output_type -> air_quality_grpc.GetAirQualityByIDResp
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_air_quality_proto_init() }
func file_air_quality_proto_init() {
	if File_air_quality_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_air_quality_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetAirQualitiesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_air_quality_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetAirQualitiesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_air_quality_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SaveAirQualitiesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_air_quality_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SaveAirQualitiesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_air_quality_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*AirQuality); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_air_quality_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*GetAirQualityByIDReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_air_quality_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GetAirQualityByIDResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_air_quality_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_air_quality_proto_goTypes,
		DependencyIndexes: file_air_quality_proto_depIdxs,
		MessageInfos:      file_air_quality_proto_msgTypes,
	}.Build()
	File_air_quality_proto = out.File
	file_air_quality_proto_rawDesc = nil
	file_air_quality_proto_goTypes = nil
	file_air_quality_proto_depIdxs = nil
}
