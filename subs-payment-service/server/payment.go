package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"subs-payment-service/model"
	"subs-payment-service/pb"
	"subs-payment-service/service"

	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func NewPaymentServer(
	db *gorm.DB,
	emailNotifService service.EmailNotifService,
	invoiceService service.InvoiceService,
	userService service.UserService,
) *PaymentServer {
	return &PaymentServer{
		db:                db,
		emailNotifService: emailNotifService,
		invoiceService:    invoiceService,
		userService:       userService,
	}
}

type PaymentServer struct {
	db                *gorm.DB
	emailNotifService service.EmailNotifService
	invoiceService    service.InvoiceService
	userService       service.UserService
	pb.UnimplementedSubPaymentServer
}

// CompletePayment implements pb.SubPaymentServer.
func (ps *PaymentServer) CompletePayment(c context.Context, req *pb.CompletePaymentReq) (*pb.CompletePaymentResp, error) {
	// verify webhook token
	verifToken := req.CallbackToken
	if verifToken == "" || verifToken != os.Getenv("XENDIT_WEBHOOK_TOKEN") {
		log.Printf("received xendit webhook token: %s", verifToken)
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
		return nil, status.Errorf(codes.Internal, "internal server error: %s", err.Error())
	}

	// check payment not updated yet
	if payment.Status != "pending" {
		return nil, status.Errorf(http.StatusBadRequest, "payment already updated")
	}

	var transDate time.Time
	transDate, err = time.Parse("", req.PaidAt)
	if err != nil {
		transDate = time.Now()
	}

	// update payment
	payment.TransactionDate = &transDate
	if req.Status == "PAID" {
		payment.Status = "completed"
	} else {
		payment.Status = "failed"
	}
	err = ps.db.Save(&payment).Error
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "internal server error: %s", err.Error())
	}

	// notify email service
	ps.emailNotifService.NotifyPaymentSucccess(int(payment.ID))

	return &pb.CompletePaymentResp{
		Payment: &pb.Payment{
			Id:              int64(payment.ID),
			UserId:          int64(payment.UserID),
			PaymentGateway:  payment.PaymentGateway,
			Amount:          float32(payment.Amount),
			Currency:        payment.Currency,
			TransactionDate: timestamppb.New(*payment.TransactionDate),
			Status:          payment.Status,
			Url:             payment.Url,
		},
	}, nil
}

func (ps *PaymentServer) validateUserSubscriptionData(req *pb.CreateUserSubcriptionReq) error {
	// Check if UserID is provided and valid
	if req.UserId <= 0 {
		return errors.New("user_id must be a positive integer")
	}

	// Check if Tier is provided and valid (non-empty and meaningful)
	if strings.TrimSpace(req.Tier) == "" && strings.TrimSpace(req.Tier) == "business" {
		return errors.New("tier cannot be empty")
	}

	// Check if Duration is valid (positive integer)
	if req.Duration <= 0 {
		return errors.New("duration must be a positive integer")
	}

	return nil
}

func (ps *PaymentServer) rollbackUserSub(userSub *model.UserSubscription) {
	payment := userSub.Payment
	err := ps.db.Delete(&userSub).Error
	if err != nil {
		log.Printf("failed rolling back changes deleting user sub: %s", err.Error())
	}
	err = ps.db.Delete(&payment).Error
	if err != nil {
		log.Printf("failed rolling back changes deleting payment: %s", err.Error())
	}
}

// CreateUserSubcription implements pb.SubPaymentServer.
func (ps *PaymentServer) CreateUserSubcription(c context.Context, req *pb.CreateUserSubcriptionReq) (*pb.CreateUserSubcriptionResp, error) {
	err := ps.validateUserSubscriptionData(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err.Error())
	}

	// validate token and get user
	user, err := ps.userService.ValidateAndGetUser(c)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}

	// get subscribtion
	var sub model.Subscription
	err = ps.db.Where("tier=?", req.Tier).First(&sub).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "subscription tier '%s' does not exist", req.GetTier())
	}
	// TODO only allow business tier

	// generate payment and invoice
	newPayment := &model.Payment{
		UserID:         user.ID,
		PaymentGateway: "xendit",
		Amount:         sub.PricePerMonth * float64(req.Duration),
		Currency:       "IDR",
		Status:         "pending",
		// filled  after creation
		Url: "",
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
		return nil, status.Errorf(codes.Internal, "unknown: %s", err.Error())
	}
	newPayment = &newUserSub.Payment

	// generate invoice
	url, err := ps.invoiceService.GenerateInvoice(
		user,
		newUserSub,
		newPayment,
	)
	if err != nil {
		ps.rollbackUserSub(newUserSub)
		return nil, status.Errorf(codes.Internal, "failed to generate invoice: %s", err.Error())
	}
	//update url
	newPayment.Url = url
	err = ps.db.Save(&newPayment).Error
	if err != nil {
		ps.rollbackUserSub(newUserSub)
		return nil, status.Errorf(codes.Internal, "failed to update payment invoice: %s", err.Error())
	}

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

func (ps *PaymentServer) GetPaymentByID(c context.Context, req *pb.GetPaymentByIDReq) (*pb.GetPaymentByIDResp, error) {
	if req.PaymentId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "payment id must be a positive integer")
	}

	var payment model.Payment
	err := ps.db.Where("id=?", req.PaymentId).First(&payment).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment with id %d not found; %s", req.PaymentId, err.Error())
	}

	return &pb.GetPaymentByIDResp{
		Payment: &pb.Payment{
			Id:              int64(payment.ID),
			UserId:          int64(payment.UserID),
			PaymentGateway:  payment.PaymentGateway,
			Amount:          float32(payment.Amount),
			Currency:        payment.Currency,
			TransactionDate: timestamppb.New(*payment.TransactionDate),
			Status:          payment.Status,
			Url:             payment.Url,
		},
	}, nil
}
