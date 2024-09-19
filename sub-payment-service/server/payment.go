package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"sub-payment-service/entity"
	"sub-payment-service/model"
	"sub-payment-service/pb"
	"sub-payment-service/service"
	"sub-payment-service/util"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PaymentServer struct {
	db                *gorm.DB
	emailNotifService service.EmailNotifService
	pb.UnimplementedSubPaymentServer
}

// CompletePayment implements pb.SubPaymentServer.
func (ps *PaymentServer) CompletePayment(context.Context, *pb.CompletePaymentReq) (*pb.CompletePaymentResp, error) {
	panic("unimplemented")
}

// CreateUserSubcription implements pb.SubPaymentServer.
func (ps *PaymentServer) CreateUserSubcription(context.Context, *pb.CreateUserSubcriptionReq) (*pb.CreateUserSubcriptionResp, error) {
	panic("unimplemented")
}

// GetUserSubcriptions implements pb.SubPaymentServer.
func (ps *PaymentServer) GetUserSubcriptions(context.Context, *pb.GetUserSubcriptionsReq) (*pb.GetUserSubcriptionsResp, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedSubPaymentServer implements pb.SubPaymentServer.
func (ps *PaymentServer) mustEmbedUnimplementedSubPaymentServer() {
	panic("unimplemented")
}

func NewPaymentServer(db *gorm.DB, emailNotifService service.EmailNotifService) *PaymentServer {
	return &PaymentServer{
		db:                db,
		emailNotifService: emailNotifService,
	}
}

// webhook payload structure
type WebhookPayload struct {
	ID            string  `json:"id"`
	ExternalID    string  `json:"external_id"`
	PaymentMethod string  `json:"payment_method"`
	PaidAmount    float64 `json:"paid_amount"`
	Status        string  `json:"status"`
}

type PaymentResp struct {
	ID            uint    `json:"payment_id"`
	PaymentMethod string  `json:"payment_method"`
	PaidAmount    float64 `json:"paid_amount"`
}

func (ph *PaymentServer) GetUserForPayment(p *model.Payment) *entity.User {
	// TODO
	return nil
}

func (ph *PaymentServer) HandlePaymentSuccess(c echo.Context) error {
	// verify webhook token
	verifToken := c.Request().Header.Get("x-callback-token")
	if verifToken == "" || verifToken != os.Getenv("XENDIT_WEBHOOK_TOKEN") {
		return util.NewAppError(http.StatusUnauthorized, "invalid webhook token", "")
	}

	// parse body
	var reqBody WebhookPayload
	err := c.Bind(&reqBody)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "bad request", err.Error())
	}
	if reqBody.PaymentMethod == "" || reqBody.PaidAmount <= 0 {
		return util.NewAppError(http.StatusBadRequest, "payment method and paid amount cannot be empty", "")
	}

	// get corresponding payment
	paymentId, err := strconv.Atoi(reqBody.ExternalID)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid payment id", "")
	}
	var payment model.Payment
	err = ph.db.Where("id=?", paymentId).Select("*").First(&payment).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return util.NewAppError(http.StatusNotFound, "payment not found", "")
	} else if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// check payment not updated yet
	if payment.Status != "Unpaid" {
		return util.NewAppError(http.StatusBadRequest, "payment already updated", "")
	}

	// update payment
	payment.PaymentMethod = reqBody.PaymentMethod
	payment.TotalPayment = reqBody.PaidAmount
	if reqBody.Status == "PAID" {
		payment.Status = "Completed"
	} else {
		payment.Status = reqBody.Status
	}
	err = ph.db.Save(&payment).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// send email is successful
	// get user
	if user := ph.GetUserForPayment(&payment); user != nil {
		if payment.Status == "Completed" {
			// TODO
			ph.emailNotifService.NotifyPaymentSucccess()
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "payment updated successfully",
		"payment": PaymentResp{
			ID:            payment.ID,
			PaymentMethod: payment.PaymentMethod,
			PaidAmount:    payment.TotalPayment,
		},
	})
}
