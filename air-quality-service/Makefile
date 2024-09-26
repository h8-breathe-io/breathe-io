include .env

PROJECT_ID=breathe-io
SERVICE_NAME=air-quality-service
IMAGE_NAME=gcr.io/$(PROJECT_ID)/$(SERVICE_NAME)
DB_HOST=34.101.36.131
USER_SERVICE_URL=user-service-43931901018.asia-southeast2.run.app:443
SUBS_PAYMENT_SERVICE_URL=subs-payment-service-43931901018.asia-southeast2.run.app:443
EMAIL_NOTIF_URL=email-notif-service-43931901018.asia-southeast2.run.app:443

protoc:
	protoc --proto_path=pb pb/*.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
	protoc --proto_path=../user-service/pb business_facilities.proto --go_out=paths=source_relative:pb/generated --go-grpc_out=paths=source_relative:pb/generated
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
		--set-env-vars DB_PASS=$(DB_PASS) \
		--set-env-vars OPEN_WEATHER_API_KEY=$(OPEN_WEATHER_API_KEY) \
		--set-env-vars OPEN_WEATHER_BASE_URL=$(OPEN_WEATHER_BASE_URL) \
		--set-env-vars EMAIL_NOTIF_URL=$(EMAIL_NOTIF_URL) \
		--set-env-vars USER_SERVICE_URL=$(USER_SERVICE_URL) \
		--set-env-vars BUSINESS_FACILITY_SERVICE_URL=$(USER_SERVICE_URL)

cloud:protoc build_push cloud_run