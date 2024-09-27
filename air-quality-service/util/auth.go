package util

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc/metadata"
)

func CreateServiceContext() context.Context {
	token := os.Getenv("SERVICE_TOKEN")
	log.Printf("creating context with token %s", token)
	md := metadata.Pairs("auth_token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx
}
