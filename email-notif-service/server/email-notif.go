package server

import (
	"bytes"
	"context"
	"email-notif-service/constants"
	"email-notif-service/pb"
	"email-notif-service/service"
	"email-notif-service/util"
	"fmt"
	"html/template"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewEmailNotifServer(
	paymentService service.SubsPaymentService,
	userService service.UserService,
	emailService service.EmailService,
	airQualityService service.AirQualityService,
) *EmailNotifServer {
	// load templates
	paymentTmpl := template.Must(template.ParseFiles(
		filepath.Join(constants.ROOT_TEMPLATE_PATH, "payment.html"),
	))
	registerTmpl := template.Must(template.ParseFiles(
		filepath.Join(constants.ROOT_TEMPLATE_PATH, "register.html"),
	))
	aqAlertTmpl := template.Must(template.ParseFiles(
		filepath.Join(constants.ROOT_TEMPLATE_PATH, "aq-alert.html"),
	))
	return &EmailNotifServer{
		paymentService:    paymentService,
		userService:       userService,
		emailService:      emailService,
		airQualityService: airQualityService,
		paymentTmpl:       paymentTmpl,
		registerTmpl:      registerTmpl,
		aqAlertTmpl:       aqAlertTmpl,
		bfService:         service.NewBusinessFacilityService(),
	}
}

type EmailNotifServer struct {
	paymentService    service.SubsPaymentService
	userService       service.UserService
	emailService      service.EmailService
	airQualityService service.AirQualityService
	paymentTmpl       *template.Template
	registerTmpl      *template.Template
	aqAlertTmpl       *template.Template
	bfService         service.BusinessFacilityService
	pb.UnimplementedEmailNotifServiceServer
}

type PaymentEmailData struct {
	Username string
	Amount   string
}

func (es *EmailNotifServer) NotifyPaymentComplete(c context.Context, req *pb.NotifyPaymentCompleteReq) (*pb.NotifyPaymentCompleteResp, error) {
	// get payment
	payment, err := es.paymentService.GetPaymentByID(int(req.PaymentId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get payment: %s", err.Error())
	}
	// get user
	user, err := es.userService.GetUserByID(int(payment.UserId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %s", err.Error())
	}

	buf := bytes.NewBuffer([]byte{})

	// data for template
	data := PaymentEmailData{
		Username: user.Username,
		Amount:   fmt.Sprintf("%.0f", payment.Amount),
	}

	// generate html for email
	err = es.paymentTmpl.Execute(buf, &data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to generate email: %s", err.Error())
	}

	// create email req
	emailReq := service.SendEmailRequest{
		From: service.Address{Email: constants.SENDER_EMAIL, Name: "Breathe.io"},
		To: []service.Address{
			{Email: user.Email, Name: user.Username},
		},
		Subject:  "Payment Completed - Breathe.io",
		Text:     buf.String(),
		Html:     buf.String(),
		Category: "Payment",
	}

	err = es.emailService.SendEmail(&emailReq)
	if err != nil {
		return nil, err
	}

	return &pb.NotifyPaymentCompleteResp{
		Status: "OK",
		Email:  user.Email,
	}, nil
}

type RegisterEmailData struct {
	Username string
}

func (es *EmailNotifServer) NotifyRegister(c context.Context, req *pb.NotifyRegisterReq) (*pb.NotifyRegisterResp, error) {

	// get user
	user, err := es.userService.GetUserByID(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %s", err.Error())
	}

	buf := bytes.NewBuffer([]byte{})

	// data for template
	data := PaymentEmailData{
		Username: user.Username,
	}

	// generate html for email
	err = es.registerTmpl.Execute(buf, &data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to generate email: %s", err.Error())
	}

	// create email req
	emailReq := service.SendEmailRequest{
		From: service.Address{Email: constants.SENDER_EMAIL, Name: "Breathe.io"},
		To: []service.Address{
			{Email: user.Email, Name: user.Username},
		},
		Subject:  "Welcome to Breathe.io",
		Text:     buf.String(),
		Html:     buf.String(),
		Category: "Registration",
	}

	err = es.emailService.SendEmail(&emailReq)
	if err != nil {
		return nil, err
	}

	return &pb.NotifyRegisterResp{
		Status: "OK",
		Email:  user.Email,
	}, nil
}

type AQAlertEmailData struct {
	Username string
	Subject  string
	Location string
	Aqi      string
}

func GetLocString(loc *pb.Location) string {
	var locStr string
	if loc != nil {
		if loc.LocationName == "Location" {
			locStr = fmt.Sprintf("Lat: %f, Long: %f", loc.Latitude, loc.Longitude)
		} else {
			locStr = loc.LocationName
		}

	}
	return locStr
}

func (es *EmailNotifServer) NotifyAirQuality(c context.Context, req *pb.NotifyAirQualityReq) (*pb.NotifyAirQualityResp, error) {
	// get user
	user, err := es.userService.GetUserByID(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %s", err.Error())
	}

	// get air quality
	aq, err := es.airQualityService.GetAirQualityByID(int(req.AirQualityId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %s", err.Error())
	}

	// if <=100 still considered healthy
	if aq.Aqi <= 3 {
		return &pb.NotifyAirQualityResp{Status: "Not Sent", Email: user.Email}, nil
	}

	data := AQAlertEmailData{
		Username: user.Username,
		Subject:  "Current Location",
		Location: GetLocString(aq.Location),
		Aqi:      fmt.Sprintf("%d", aq.Aqi),
	}

	buf := bytes.NewBuffer([]byte{})
	// generate html for email
	err = es.aqAlertTmpl.Execute(buf, &data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to generate email: %s", err.Error())
	}

	// create email req
	emailReq := service.SendEmailRequest{
		From: service.Address{Email: constants.SENDER_EMAIL, Name: "Breathe.io"},
		To: []service.Address{
			{Email: user.Email, Name: user.Username},
		},
		Subject:  "Unhealthy Air Quality in Your Area",
		Text:     buf.String(),
		Html:     buf.String(),
		Category: "AQ Alert",
	}

	err = es.emailService.SendEmail(&emailReq)
	if err != nil {
		return nil, err
	}

	return &pb.NotifyAirQualityResp{
		Status: "OK",
		Email:  user.Email,
	}, nil
}

func (es *EmailNotifServer) NotifyAirQualityBusiness(c context.Context, req *pb.NotifyAirQualityBusinessReq) (*pb.NotifyAirQualityBusinessResp, error) {

	ctx := util.CreateServiceContext()
	//get business
	bf, err := es.bfService.GetBusinessFacilityByID(ctx, int(req.BusinessId))
	if err != nil {
		return nil, err
	}

	// get user
	user, err := es.userService.GetUserByID(int(bf.UserID))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %s", err.Error())
	}

	// get air quality
	aq, err := es.airQualityService.GetAirQualityByID(int(req.AirQualityId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get aq: %s", err.Error())
	}

	data := AQAlertEmailData{
		Username: user.Username,
		Subject:  bf.CompanyType,
		Location: GetLocString(aq.Location),
		Aqi:      fmt.Sprintf("%d", aq.Aqi),
	}

	buf := bytes.NewBuffer([]byte{})
	// generate html for email
	err = es.aqAlertTmpl.Execute(buf, &data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to generate email: %s", err.Error())
	}

	// create email req
	emailReq := service.SendEmailRequest{
		From: service.Address{Email: constants.SENDER_EMAIL, Name: "Breathe.io"},
		To: []service.Address{
			{Email: user.Email, Name: user.Username},
		},
		Subject:  "Unhealthy Air Quality in Your Area",
		Text:     buf.String(),
		Html:     buf.String(),
		Category: "AQ Alert",
	}

	err = es.emailService.SendEmail(&emailReq)
	if err != nil {
		return nil, err
	}

	return &pb.NotifyAirQualityBusinessResp{
		Status: "OK",
		Email:  user.Email,
	}, nil
}
