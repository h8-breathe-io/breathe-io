include .env

PROJECT_ID=breathe-io
SERVICE_NAME=user-service
IMAGE_NAME=gcr.io/$(PROJECT_ID)/$(SERVICE_NAME)
DB_HOST=34.101.36.131
EMAIL_NOTIF_URL=email-notif-service-43931901018.asia-southeast2.run.app:8080
LOCATION_GRPC_ADDR=air-quality-service-43931901018.asia-southeast2.run.app:8080

protoc:
	protoc --proto_path=pb pb/*.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc email-notif.proto --proto_path ../email-notif-service/proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated

run:
	@go run main.go

build_push:
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)

cloud_run:
	gcloud run deploy $(SERVICE_NAME) --image $(IMAGE_NAME) --platform managed --region asia-southeast2 --allow-unauthenticated --port 8080 \
		--set-env-vars DB_HOST=$(DB_HOST) \
		--set-env-vars DB_PORT=5432 \
		--set-env-vars DB_NAME=postgres \
		--set-env-vars DB_USER=postgres \
		--set-env-vars JWT_SECRET=$(JWT_SECRET) \
		--set-env-vars DB_PASS=$(DB_PASS) \
		--set-env-vars EMAIL_NOTIF_URL=$(EMAIL_NOTIF_URL) \
		--set-env-vars LOCATION_GRPC_ADDR=$(LOCATION_GRPC_ADDR)


cloud:protoc build_push cloud_run