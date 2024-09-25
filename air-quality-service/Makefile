protoc:
	protoc --proto_path=pb pb/*.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc --proto_path=../user-service/pb business_facilities.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc email-notif.proto --proto_path ../email-notif-service/proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
run:
	@go run main.go