include .env

PROJECT_ID=breathe-io
SERVICE_NAME=api-gateway
IMAGE_NAME=gcr.io/$(PROJECT_ID)/$(SERVICE_NAME)
USER_SERVICE_URL=user-service-43931901018.asia-southeast2.run.app:443
SUBS_PAYMENT_SERVICE_URL=subs-payment-service-43931901018.asia-southeast2.run.app:443
AIR_QUALITY_SERVICE_URL=air-quality-service-43931901018.asia-southeast2.run.app:443
REPORT_SERVICE_URL=reporting-43931901018.asia-southeast2.run.app:443

protoc:
	protoc subs-payment.proto --proto_path ../subs-payment-service/proto --go_out . --go-grpc_out .
	protoc email-notif.proto --proto_path ../email-notif-service/proto --go_out . --go-grpc_out .
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../air-quality-service/pb air_quality.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../air-quality-service/pb location.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../user-service/pb business_facilities.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb
	protoc --proto_path=../reporting-pdf-csv/proto reporting.proto --go_out=paths=source_relative:pb --go-grpc_out=paths=source_relative:pb

build_push:
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)


cloud_run:
	gcloud run deploy $(SERVICE_NAME) --image $(IMAGE_NAME) --platform managed --region asia-southeast2 --allow-unauthenticated --port 8080 \
		--set-env-vars BUSINESS_FACILITIES_SERVICE_URL=$(USER_SERVICE_URL) \
		--set-env-vars LOCATION_SERVICE_URL=$(AIR_QUALITY_SERVICE_URL) \
		--set-env-vars SUBS_PAYMENT_SERVICE_URL=$(SUBS_PAYMENT_SERVICE_URL) \
		--set-env-vars AIR_QUALITY_SERVICE_URL=$(AIR_QUALITY_SERVICE_URL) \
		--set-env-vars USER_SERVICE_URL=$(USER_SERVICE_URL) \
		--set-env-vars REPORT_SERVICE_URL=$(REPORT_SERVICE_URL) \
		--set-env-vars AQ_SERVICE_TOKEN=$(AQ_SERVICE_TOKEN)

cloud:protoc build_push cloud_run