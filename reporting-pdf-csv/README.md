# Reporting CSV

## Domain GCP
```
```

## Run Protobuf
```
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get google.golang.org/grpc
protoc -I proto/ proto/reporting.proto --go_out=. --go-grpc_out=.
```

## Dockerize
1. Write Dockerfiles for build container
```
# Step 1: Build the Go binary with gRPC support
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies, including gRPC and protobuf dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app, ensure gRPC is enabled
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Step 2: Build a small image using the Go binary from step 1
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x /app/main

# Copy .env if necessary (optional, only if needed in the container)
COPY .env .

# Expose port 50051 for gRPC
EXPOSE 50051

# Command to run the executable
CMD ["./main"]
```

4. Build, tag Docker images & push to registry
```
docker build -t gcr.io/hacktiv8-431914/bookmanagement:v1.0 .
gcloud auth login
gcloud config set project hacktiv8-431914
gcloud auth configure-docker
docker push gcr.io/hacktiv8-431914/bookmanagement:v1.0
```


