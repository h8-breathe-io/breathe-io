PROJECT_ID=breathe-io
SERVICE_NAME=user-service
IMAGE_NAME=gcr.io/$(PROJECT_ID)/$(SERVICE_NAME)

include .env

protoc:
	protoc --proto_path=pb pb/*.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc email-notif.proto --proto_path ../email-notif-service/proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated

run:
	@go run main.go

build_push:
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)

cloud_run:
	gcloud run deploy $(SERVICE_NAME) --image $(IMAGE_NAME) --platform managed --region asia-southeast2 --allow-unauthenticated --port 8080 --set-env-vars DB_HOST=34.101.36.131,DB_PORT=5432,DB_NAME=postgres,DB_USER=postgres


cloud:protoc build_push cloud_run