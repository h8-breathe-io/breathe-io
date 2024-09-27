package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "api-gateway/pb"
	"api-gateway/service"
	"api-gateway/util"
)

func NewSubsPaymentClient() pb.SubPaymentClient {
	addr := os.Getenv("SUBS_PAYMENT_SERVICE_URL")
	log.Printf("subs-payment service url: %s", addr)
	// Set up a connection to the server.
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.NewClient(addr, opts...)
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

	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewUserClient(conn)

	return client
}

func NewAQClient() pb.AirQualityServiceClient {
	addr := os.Getenv("AIR_QUALITY_SERVICE_URL")
	log.Printf("aq service url: %s", addr)
	// Set up a connection to the server.
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewAirQualityServiceClient(conn)

	return client
}

func NewBFClient() pb.BusinessFacilitiesClient {
	addr := os.Getenv("BUSINESS_FACILITIES_SERVICE_URL")
	log.Printf("bf service url: %s", addr)
	// Set up a connection to the server.
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewBusinessFacilitiesClient(conn)

	return client
}

func NewLocClient() pb.LocationServiceClient {
	addr := os.Getenv("LOCATION_SERVICE_URL")
	log.Printf("loc service url: %s", addr)
	// Set up a connection to the server.
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewLocationServiceClient(conn)

	return client
}

func NewReportClient() pb.ReportServiceClient {
	addr := os.Getenv("REPORT_SERVICE_URL")
	log.Printf("report service url: %s", addr)
	// Set up a connection to the server.
	opts := []grpc.DialOption{}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("filed to get certs: %v", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewReportServiceClient(conn)

	return client
}

type handler struct {
	subsPaymentCLient pb.SubPaymentClient
	userClient        pb.UserClient
	aqClient          pb.AirQualityServiceClient
	bfClient          pb.BusinessFacilitiesClient
	locClient         pb.LocationServiceClient
	reportService     pb.ReportServiceClient
}

func (h *handler) HandleCreateUserSubscription(c echo.Context) error {
	// pb definitions have json annotations, can use it directly
	var req pb.CreateUserSubcriptionReq
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := util.CreateContext(c)
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

	ctx := util.CreateContext(c)
	res, err := h.aqClient.SaveAirQualities(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)

}

func (h *handler) HandleSaveAirQualityBusiness(c echo.Context) error {
	var req pb.SaveAirQualityForBusinessReq
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := util.CreateContext(c)
	res, err := h.aqClient.SaveAirQualityForBusiness(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)

}

func (h *handler) HandleSaveAirQualitiesHistorical(c echo.Context) error {
	var req pb.SaveHistoricalAirQualitiesRequest
	err := c.Bind(&req)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid request body", err.Error())
	}

	ctx := util.CreateContext(c)
	res, err := h.aqClient.SaveHistoricalAirQualities(
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

	ctx := util.CreateContext(c)
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

	ctx := util.CreateContext(c)
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

	ctx := util.CreateContext(c)
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

	ctx := util.CreateContext(c)
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

	var req pb.UpdateBFRequest
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	req.Id = uint64(bfId)

	ctx := util.CreateContext(c)
	res, err := h.bfClient.UpdateBusinessFacility(
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

	ctx := util.CreateContext(c)
	res, err := h.bfClient.DeleteBusinessFacility(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGetLocations(c echo.Context) error {

	ctx := util.CreateContext(c)
	res, err := h.locClient.GetLocations(
		ctx,
		&emptypb.Empty{},
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGetLocation(c echo.Context) error {
	// get id param
	idParam := c.Param("id")
	bfId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid BF id")
	}

	var req pb.GetLocationRequest
	req.LocationId = uint64(bfId)

	ctx := util.CreateContext(c)
	res, err := h.locClient.GetLocation(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGetLocationRecommendation(c echo.Context) error {
	// get id param
	idParam := c.Param("id")
	bfId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid BF id")
	}

	var req pb.LocationRecommendationRequest
	req.BusinessId = uint64(bfId)

	ctx := util.CreateContext(c)
	res, err := h.locClient.GetLocationRecommendation(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGenerateReport(c echo.Context) error {

	var req pb.ReportRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	ctx := util.CreateContext(c)
	res, err := h.reportService.GenerateReport(
		ctx,
		&req,
	)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *handler) HandleGetProfile(c echo.Context) error {
	token := util.ExtractAuthToken(c)
	user, err := h.userClient.IsValidToken(context.TODO(), &pb.IsValidTokenRequest{Token: token})
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "service error", err.Error())
	}

	return c.JSON(http.StatusOK, user.User)
}

func main() {
	godotenv.Load()

	handler := handler{
		subsPaymentCLient: NewSubsPaymentClient(),
		userClient:        NewUserClient(),
		aqClient:          NewAQClient(),
		bfClient:          NewBFClient(),
		locClient:         NewLocClient(),
		reportService:     NewReportClient(),
	}

	cs := service.NewCronServices(NewAQClient(), NewLocClient())
	//declare cron services
	c := cron.New()
	//running cron job exactly every start of the
	c.AddFunc("*/15 * * * *", func() {
		cs.RenewAQData()
	})
	c.Start()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	// set error handler
	e.HTTPErrorHandler = util.ErrorHandler

	// user subs
	userSubs := e.Group("/user-subscriptions")
	userSubs.POST("", handler.HandleCreateUserSubscription)
	// callback for xendit
	// don't need auth since it uses xendit token, will be authenticated in service
	e.POST("/payment-callback", handler.HandlePaymentCallback)

	// user
	users := e.Group("/users")
	users.POST("/register", handler.HandleRegister)
	users.POST("/login", handler.HandleLogin)
	users.GET("/profile", handler.HandleGetProfile)

	// air-qualities
	aq := e.Group("/air-qualities")
	aq.POST("", handler.HandleSaveAirQualities)
	aq.POST("/business", handler.HandleSaveAirQualityBusiness)
	aq.GET("", handler.HandleGetAirQualities)
	aq.POST("/historical", handler.HandleSaveAirQualitiesHistorical)

	// business facilities
	bf := e.Group("/business-facilities")
	bf.POST("", handler.HandleAddBusinessFacility)
	bf.GET("", handler.HandleGetBusinessFacilities)
	bf.GET("/:id", handler.HandleGetBusinessFacility)
	bf.PUT("/:id", handler.HandleUpdateBusinessFacility)
	bf.DELETE("/:id", handler.HandleDeleteBusinessFacility)
	bf.GET("/:id/recommendation", handler.HandleGetLocationRecommendation)

	// locations
	l := e.Group("/locations")
	l.GET("", handler.HandleGetLocations)
	l.GET("/:id", handler.HandleGetLocation)

	// reporting
	e.POST("/reports", handler.HandleGenerateReport)

	// start server
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
