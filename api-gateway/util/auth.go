package util

import (
	"context"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
)

func ExtractAuthToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")

	// Check if the header is in the correct format
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Print("No token found, returning empty context")
		// if not return emtpty context
		return ""
	}

	// Extract the token part from the header (after "Bearer ")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token
}

func CreateContext(c echo.Context) context.Context {
	//get token from header

	authHeader := c.Request().Header.Get("Authorization")

	// Check if the header is in the correct format
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Print("No token found, returning empty context")
		// if not return emtpty context
		return context.TODO()
	}

	// Extract the token part from the header (after "Bearer ")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Token found '%s', attaching to context", token)
	// attach token to context
	md := metadata.Pairs("auth_token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}
