include .env

PROJECT_ID=breathe-io
SERVICE_NAME=reporting
IMAGE_NAME=gcr.io/$(PROJECT_ID)/$(SERVICE_NAME)
DB_HOST=34.101.36.131
USER_SERVICE_URL=user-service-43931901018.asia-southeast2.run.app:443

protoc:
	protoc reporting.proto --proto_path proto --go_out . --go-grpc_out .
	protoc --proto_path=../user-service/pb user.proto --go_out=paths=source_relative:proto/pb --go-grpc_out=paths=source_relative:proto/pb


build_push:
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)


cloud_run:
	gcloud run deploy $(SERVICE_NAME) --image $(IMAGE_NAME) --platform managed --region asia-southeast2 --allow-unauthenticated --port 8080 \
		--set-env-vars DB_HOST=$(DB_HOST) \
		--set-env-vars DB_PORT=5432 \
		--set-env-vars DB_NAME=postgres \
		--set-env-vars DB_USER=postgres \
		--set-env-vars DB_PASSWORD=$(DB_PASSWORD) \
		--set-env-vars TOKEN_GPT=$(TOKEN_GPT) \
		--set-env-vars USER_SERVICE_URL=$(USER_SERVICE_URL)

cloud:protoc build_push cloud_run