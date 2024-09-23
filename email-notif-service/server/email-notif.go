package server

import (
	"context"
	"email-notif-service/constants"
	"email-notif-service/pb"
	"email-notif-service/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewEmailNotifServer(
	paymentService service.SubsPaymentService,
	userService service.UserService,
	emailService service.EmailService,
) *EmailNotifServer {
	return &EmailNotifServer{
		paymentService: paymentService,
		userService:    userService,
		emailService:   emailService,
	}
}

type EmailNotifServer struct {
	paymentService service.SubsPaymentService
	userService    service.UserService
	emailService   service.EmailService
	pb.UnimplementedEmailNotifServiceServer
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

	// create email req
	emailReq := service.SendEmailRequest{
		From: service.Address{Email: constants.SENDER_EMAIL, Name: "Breathe.io"},
		To: []service.Address{
			{Email: user.Email, Name: user.Username},
		},
		Subject:  "Payment Completed - Breathe.io",
		Text:     "Congrats for sending test email with Mailtrap!",
		Html:     "",
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
