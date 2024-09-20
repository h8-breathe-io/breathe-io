package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sub-payment-service/model"
	"sub-payment-service/pb"
	"sub-payment-service/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func NewPaymentServer(
	db *gorm.DB,
	emailNotifService service.EmailNotifService,
	invoiceService service.InvoiceService) *PaymentServer {
	return &PaymentServer{
		db:                db,
		emailNotifService: emailNotifService,
		invoiceService:    invoiceService,
	}
}

type PaymentServer struct {
	db                *gorm.DB
	emailNotifService service.EmailNotifService
	invoiceService    service.InvoiceService
	pb.UnimplementedSubPaymentServer
}

// CompletePayment implements pb.SubPaymentServer.
func (ps *PaymentServer) CompletePayment(c context.Context, req *pb.CompletePaymentReq) (*pb.CompletePaymentResp, error) {
	// verify webhook token
	verifToken := req.CallbackToken
	if verifToken == "" || verifToken != os.Getenv("XENDIT_WEBHOOK_TOKEN") {
		return nil, status.Errorf(codes.InvalidArgument, "invalid webhook token")
	}

	// validate req
	if req.PaymentMethod == "" || req.PaidAmount <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "payment method and paid amount cannot be empty")
	}

	// get corresponding payment
	paymentId, err := strconv.Atoi(req.ExternalId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment id: %s", req.ExternalId)
	}
	var payment model.Payment
	err = ps.db.Where("id=?", paymentId).Select("*").First(&payment).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "payment not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Unknown, "internal server error: %s", err.Error())
	}

	// check payment not updated yet
	if payment.Status != "pending" {
		return nil, status.Errorf(http.StatusBadRequest, "payment already updated")
	}

	// update payment
	payment.PaymentGateway = req.PaymentMethod
	payment.Amount = float64(req.PaidAmount)
	if req.Status == "PAID" {
		payment.Status = "completed"
	} else {
		payment.Status = "failed"
	}
	err = ps.db.Save(&payment).Error
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "internal server error: %s", err.Error())
	}

	// send email is successful
	// get user
	// if user := ps.GetUserForPayment(&payment); user != nil {
	// 	if payment.Status == "Completed" {
	// 		// TODO
	// 		ps.emailNotifService.NotifyPaymentSucccess()
	// 	}
	// }

	return &pb.CompletePaymentResp{}, nil
}

func (ps *PaymentServer) validateUserSubscriptionData(req *pb.CreateUserSubcriptionReq) error {
	// Check if UserID is provided and valid
	if req.UserId <= 0 {
		return errors.New("user_id must be a positive integer")
	}

	// Check if Tier is provided and valid (non-empty and meaningful)
	if strings.TrimSpace(req.Tier) == "" {
		return errors.New("tier cannot be empty")
	}

	// Check if Duration is valid (positive integer)
	if req.Duration <= 0 {
		return errors.New("duration must be a positive integer")
	}

	return nil
}

// CreateUserSubcription implements pb.SubPaymentServer.
func (ps *PaymentServer) CreateUserSubcription(c context.Context, req *pb.CreateUserSubcriptionReq) (*pb.CreateUserSubcriptionResp, error) {
	err := ps.validateUserSubscriptionData(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err.Error())
	}

	// check user exists?
	//TODO

	// get subscribtion
	var sub model.Subscription
	err = ps.db.Where("tier=?", req.Tier).First(&sub).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "subscription tier '%s' does not exist", req.GetTier())
	}
	// TODO only allow business

	// generate payment and invoice
	newPayment := &model.Payment{
		UserID:         int(req.UserId),
		PaymentGateway: "xendit",
		Amount:         sub.PricePerMonth * float64(req.Duration),
		Currency:       "IDR",
		Status:         "pending",
		// TODO generate invoice
		Url: "test",
	}

	//create new model
	newUserSub := &model.UserSubscription{
		UserID:         uint(req.UserId),
		SubscriptionID: sub.ID,
		Duration:       int(req.Duration),
		Payment:        *newPayment,
	}
	err = ps.db.Save(newUserSub).Error
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "unknown: %s", err.Error())
	}
	newPayment = &newUserSub.Payment

	return &pb.CreateUserSubcriptionResp{
		Id:     int64(newUserSub.ID),
		UserId: int64(newUserSub.UserID),
		Subscription: &pb.Subscription{
			Id:            int64(sub.ID),
			Tier:          sub.Tier,
			PricePerMonth: float32(sub.PricePerMonth),
		},
		Duration: int64(newUserSub.Duration),
		Payment: &pb.Payment{
			Id:             int64(newPayment.ID),
			UserId:         int64(newPayment.UserID),
			PaymentGateway: newPayment.PaymentGateway,
			Amount:         float32(newPayment.Amount),
			Currency:       newPayment.Currency,
			Status:         newPayment.Status,
			Url:            newPayment.Url,
		},
	}, nil
}

// GetUserSubcriptions implements pb.SubPaymentServer.
func (ps *PaymentServer) GetUserSubcriptions(context.Context, *pb.GetUserSubcriptionsReq) (*pb.GetUserSubcriptionsResp, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedSubPaymentServer implements pb.SubPaymentServer.
func (ps *PaymentServer) mustEmbedUnimplementedSubPaymentServer() {
	panic("unimplemented")
}
