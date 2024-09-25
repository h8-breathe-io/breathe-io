package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "api-gateway/pb"
	"api-gateway/util"
)

func NewSubsPaymentClient() pb.SubPaymentClient {
	addr := os.Getenv("SUBS_PAYMENT_SERVICE_URL")
	log.Printf("subs-payment service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewSubPaymentClient(conn)

	return client
}

func NewUserClient() pb.UserClient {
	addr := os.Getenv("USER_SERVICE_URL")
	log.Printf("user service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewUserClient(conn)

	return client
}

func NewAQClient() pb.AirQualityServiceClient {
	addr := os.Getenv("AIR_QUALITY_SERVICE_URL")
	log.Printf("user service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewAirQualityServiceClient(conn)

	return client
}

func NewBFClient() pb.BusinessFacilitiesClient {
	addr := os.Getenv("BUSINESS_FACILITIES_SERVICE_URL")
	log.Printf("user service url: %s", addr)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewBusinessFacilitiesClient(conn)

	return client
}

type handler struct {
	subsPaymentCLient pb.SubPaymentClient
	userClient        pb.UserClient
	aqClient          pb.AirQualityServiceClient
	bfClient          pb.BusinessFacilitiesClient
}

func (h *handler) createContext(c echo.Context) context.Context {
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

func (h *handler) HandleCreateUserSubscription(c echo.Context) error {
	// pb definitions have json annotations, can use it directly
	var req pb.CreateUserSubcriptionReq
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := h.createContext(c)
	// forward
	req.UserId = 1
	res, err := h.subsPaymentCLient.CreateUserSubcription(
		ctx,
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

	// verify webhook token
	verifToken := c.Request().Header.Get("x-callback-token")
	if verifToken == "" {
		return util.NewAppError(http.StatusUnauthorized, "invalid webhook token", "")
	}
	// set token from header
	req.CallbackToken = verifToken

	res, err := h.subsPaymentCLient.CompletePayment(
		context.TODO(),
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) HandleRegister(c echo.Context) error {
	var req pb.RegisterRequest
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	res, err := h.userClient.Register(
		context.TODO(),
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)

}

func (h *handler) HandleLogin(c echo.Context) error {
	var req pb.LoginRequest
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	res, err := h.userClient.Login(
		context.TODO(),
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)

}

func (h *handler) HandleSaveAirQualities(c echo.Context) error {
	var req pb.SaveAirQualitiesRequest
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := h.createContext(c)
	res, err := h.aqClient.SaveAirQualities(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)

}

func (h *handler) HandleGetAirQualities(c echo.Context) error {
	locId := c.QueryParam("locId")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	locIdInt, err := strconv.Atoi(locId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid location id")
	}

	var req pb.GetAirQualitiesRequest
	req.LocationId = uint64(locIdInt)
	req.StartDate = startDate
	req.EndDate = endDate

	ctx := h.createContext(c)
	res, err := h.aqClient.GetAirQualities(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)

}

func (h *handler) HandleAddBusinessFacility(c echo.Context) error {
	var req pb.AddBFRequest
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := h.createContext(c)
	res, err := h.bfClient.AddBusinessFacility(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGetBusinessFacilities(c echo.Context) error {
	var req pb.GetBFRequests
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := h.createContext(c)
	res, err := h.bfClient.GetBusinessFacilities(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGetBusinessFacility(c echo.Context) error {
	// get id param
	idParam := c.Param("id")
	bfId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid BF id")
	}

	var req pb.GetBFRequest
	req.Id = uint64(bfId)

	ctx := h.createContext(c)
	res, err := h.bfClient.GetBusinessFacility(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleUpdateBusinessFacility(c echo.Context) error {
	// get id param
	idParam := c.Param("id")
	bfId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid BF id")
	}

	var req pb.GetBFRequest
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	req.Id = uint64(bfId)

	ctx := h.createContext(c)
	res, err := h.bfClient.GetBusinessFacility(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleDeleteBusinessFacility(c echo.Context) error {
	// get id param
	idParam := c.Param("id")
	bfId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid BF id")
	}

	var req pb.DeleteBFRequest
	req.Id = uint64(bfId)

	ctx := h.createContext(c)
	res, err := h.bfClient.DeleteBusinessFacility(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func main() {
	godotenv.Load()

	handler := handler{
		subsPaymentCLient: NewSubsPaymentClient(),
		userClient:        NewUserClient(),
		aqClient:          NewAQClient(),
		bfClient:          NewBFClient(),
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	// set error handler
	e.HTTPErrorHandler = util.ErrorHandler

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

	// user
	users := e.Group("/users")
	users.POST("/register", handler.HandleRegister)
	users.POST("/login", handler.HandleLogin)

	// air-qualities
	aq := e.Group("/air-qualities")
	aq.POST("", handler.HandleSaveAirQualities)
	aq.GET("", handler.HandleGetAirQualities)

	// business facilities
	bf := e.Group("/business-facilities")
	bf.POST("", handler.HandleAddBusinessFacility)
	bf.GET("", handler.HandleGetBusinessFacilities)
	bf.GET("/:id", handler.HandleGetBusinessFacility)
	bf.PUT("/:id", handler.HandleUpdateBusinessFacility)
	bf.DELETE("/:id", handler.HandleDeleteBusinessFacility)
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))))
}
