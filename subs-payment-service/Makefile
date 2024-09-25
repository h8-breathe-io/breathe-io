protoc:
	protoc subs-payment.proto --proto_path proto --go_out . --go-grpc_out .
	protoc email-notif.proto --proto_path ../email-notif-service/proto --go_out . --go-grpc_out .
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
