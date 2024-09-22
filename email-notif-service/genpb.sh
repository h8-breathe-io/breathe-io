protoc subs-payment.proto --proto_path ../proto --go_out . --go-grpc_out .
protoc email-notif.proto --proto_path ../proto --go_out . --go-grpc_out .