# List of services
SERVICES = air-quality-service api-gateway email-notif-service subs-payment-service user-service

# Target to build each service by running make in its folder
.PHONY: all build_services docker_compose

compose: build_services docker_compose

# runs make file in each service
build_services:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(MAKE) -C $$service protoc; \
	done

# Target to run docker-compose
docker_compose:
	@echo "Running docker-compose..."
	docker compose up --build
