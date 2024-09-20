package main

import (
	"air-quality-service/config"
	"air-quality-service/handler"
	"air-quality-service/service"
	"air-quality-service/util"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := config.CreateDBInstance()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = util.ErrorHandler

	// jwt middleware
	// config := echojwt.Config{
	// 	NewClaimsFunc: func(c echo.Context) jwt.Claims {
	// 		return new(auth.JwtAppClaims)
	// 	},
	// 	SigningKey: []byte(os.Getenv("JWT_KEY")),
	// }
	// jwtAuth := echojwt.WithConfig(config)

	// instantiate dependencies
	airQualityService := service.NewAirQualityService()

	// payments, for call backs by xendit
	airQuality := handler.NewAirQualityHandler(db, airQualityService)
	aqs := e.Group("/air-qualities")
	//authorization commented, waiting for working
	// aqs.Use(jwtAuth)
	aqs.POST("", airQuality.FetchAirQualityData)
	aqs.GET("", airQuality.GetAirQualities)

	// swagger docs
	// e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// start server
	log.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
