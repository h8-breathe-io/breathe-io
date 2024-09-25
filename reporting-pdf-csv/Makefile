protoc:
	protoc reporting.proto --proto_path proto --go_out . --go-grpc_out .
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:proto/pb --go-grpc_out=paths=source_relative:proto/pb