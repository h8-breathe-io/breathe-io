package server

import (
	"bytes"
	"context"
	"email-notif-service/constants"
	"email-notif-service/pb"
	"email-notif-service/service"
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
) *EmailNotifServer {
	// load templates
	paymentTmpl := template.Must(template.ParseFiles(
		filepath.Join(constants.ROOT_TEMPLATE_PATH, "payment.html"),
	))
	return &EmailNotifServer{
		paymentService: paymentService,
		userService:    userService,
		emailService:   emailService,
		paymentTmpl:    paymentTmpl,
	}
}

type EmailNotifServer struct {
	paymentService service.SubsPaymentService
	userService    service.UserService
	emailService   service.EmailService
	paymentTmpl    *template.Template
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
