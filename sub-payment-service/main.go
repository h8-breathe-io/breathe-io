package main

import (
	"fmt"
	"log"
	"os"

	"sub-payment-service/auth"
	"sub-payment-service/config"
	"sub-payment-service/handler"
	"sub-payment-service/service"

	// _ "h8-p2-finalproj-app/docs"

	"sub-payment-service/util"

	"github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//	@title			H8 P2 Final Project App
//	@version		1.0
//	@description	Hacktiv8 Phase 2 Final Project

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securitydefinitions.basic	BasicAuth
//	@tokenUrl					https://localhost:8080/users/login
//	@scope.read					Grants read access
//	@scope.write				Grants write access

func main() {
	db := config.CreateDBInstance()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = util.ErrorHandler

	// jwt middleware
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtAppClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	}
	jwtAuth := echojwt.WithConfig(config)

	// instantiate dependencies
	emailNotifService := service.NewEmailNotifService()

	// payments, for call backs by xendit
	payment := handler.NewPaymentHandler(db, emailNotifService)
	payments := e.Group("/payments")
	payments.Use(jwtAuth)
	payments.POST("/callback", payment.HandlePaymentSuccess)

	// swagger docs
	// e.GET("/swagger/*", echoSwagger.WrapHandler)

	// start server
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
