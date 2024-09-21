package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "api-gateway/pb"
	"api-gateway/util"
)

func NewSubsPaymentClient() pb.SubPaymentClient {
	addr := os.Getenv("SUBS_PAYMENT_SERVICE_URL")
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewSubPaymentClient(conn)

	return client
}

type handler struct {
	subsPaymentCLient pb.SubPaymentClient
}

func (h *handler) HandleCreateUserSubscription(c echo.Context) error {
	// pb definitions have json annotations, can use it directly
	var req pb.CreateUserSubcriptionReq
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	// TODO get user_id from context
	// userId, ok := c.Get("user_id").(string)
	// if !ok {
	// 	return util.NewAppError(http.StatusInternalServerError, "internal server error", "user id not set in context")
	// }
	// for now hardcode
	req.UserId = 1

	res, err := h.subsPaymentCLient.CreateUserSubcription(
		context.TODO(),
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandlePaymentCallback(c echo.Context) error {
	// pb definitions have json annotations, can use it directly
	var req pb.CompletePaymentReq
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	res, err := h.subsPaymentCLient.CompletePayment(
		context.TODO(),
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func main() {
	godotenv.Load()

	handler := handler{
		subsPaymentCLient: NewSubsPaymentClient(),
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// TODO middleware to handle JWT token
	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Parse and get token from `Authorization` header
			// 2. use RegisterLoginClient to call login function
			// 3. get result of login call, return err if failed
			// 4. if token valid, get user_id from the claims in the JWT
			// 5. Attach user id to context
			// 		c.Set("user_id", userId)
			// 6. handlers can then get the user id and pass it on to grpc calls that need it
			return next(c)
		}
	}

	// user subs
	userSubs := e.Group("/user-subscriptions")
	userSubs.Use(authMiddleware)
	userSubs.POST("", handler.HandleCreateUserSubscription)
	// callback for xendit
	// don't need auth since it uses xendit token, will be authenticated in service
	e.POST("/payment-callback", handler.HandlePaymentCallback)

	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))))
}
