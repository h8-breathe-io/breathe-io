package util

import (
	"context"
	"os"

	"google.golang.org/grpc/metadata"
)

func CreateServiceContext() context.Context {
	token := os.Getenv("SERVICE_TOKEN")
	md := metadata.Pairs("auth_token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}
