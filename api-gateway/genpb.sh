protoc subs-payment.proto --proto_path ../proto --go_out . --go-grpc_out .
protoc email-notif.proto --proto_path ../proto --go_out . --go-grpc_out .
protoc user.proto --proto_path ../user-service/pb --go_out pb --go-grpc_out pb
protoc air_quality.proto --proto_path ../air-quality-service/pb --go_out pb --go-grpc_out pb 
protoc location.proto --proto_path ../air-quality-service/pb --go_out pb --go-grpc_out pb 
