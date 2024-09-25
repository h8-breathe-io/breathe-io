protoc:
	protoc subs-payment.proto --proto_path ../subs-payment-service/proto --go_out . --go-grpc_out .
	protoc email-notif.proto --proto_path ../email-notif-service/proto --go_out . --go-grpc_out .
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../air-quality-service/pb air_quality.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../air-quality-service/pb location.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../user-service/pb business_facilities.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
