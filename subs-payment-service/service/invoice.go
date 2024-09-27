package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"subs-payment-service/model"
)

type InvoiceService interface {
	GenerateInvoice(
		forUser *User,
		forUserSub *model.UserSubscription,
		forPayment *model.Payment,
	) (string, error)
}

type invoiceService struct {
	client          *http.Client
	hostname        string
	invoiceDuration int
}

func NewInvoiceService() InvoiceService {
	return &invoiceService{
		client:   &http.Client{},
		hostname: "https://api.xendit.co/v2/invoices",
		// default 24 hours
		invoiceDuration: 86400,
	}
}

func (is *invoiceService) BuildSuccessUrl() string {
	return os.Getenv("XENDIT_INVOICE_CALLBACK")
}

func (is *invoiceService) BuildFailureUrl() string {
	return os.Getenv("XENDIT_INVOICE_CALLBACK")
}

func (is *invoiceService) generateInvoiceDesc(
	forUser *User,
	forUserSub *model.UserSubscription,
	forPayment *model.Payment,
) string {
	return fmt.Sprintf("Invoice for Breathe.io %s Tier, %d months", forUserSub.Subscription.Tier, forUserSub.Duration)
}

func (is *invoiceService) generateInvoiceItems(forUserSub *model.UserSubscription) []map[string]any {
	items := []map[string]any{
		{
			"name":     fmt.Sprintf("Breathe.io %s Tier, %d months", forUserSub.Subscription.Tier, forUserSub.Duration),
			"quantity": forUserSub.Duration,
			"price":    forUserSub.Subscription.PricePerMonth,
			"url":      fmt.Sprintf("/user-subscriptions/%d", forUserSub.ID),
		},
	}
	return items
}

func (is *invoiceService) GenerateInvoice(
	forUser *User,
	forUserSub *model.UserSubscription,
	forPayment *model.Payment,
) (string, error) {
	data := map[string]any{
		"external_id":      fmt.Sprintf("%d", forPayment.ID),
		"amount":           forPayment.Amount,
		"description":      is.generateInvoiceDesc(forUser, forUserSub, forPayment),
		"invoice_duration": is.invoiceDuration,
		"customer": map[string]any{
			"email":         forUser.Email,
			"mobile_number": forUser.PhoneNumber,
		},
		"success_redirect_url": is.BuildSuccessUrl(),
		"failure_redirect_url": is.BuildFailureUrl(),
		"payment_methods": []string{
			"CREDIT_CARD", "BCA", "BNI", "BSI", "BRI", "MANDIRI", "PERMATA",
			"OVO", "DANA", "SHOPEEPAY", "LINKAJA", "JENIUSPAY", "DD_BRI", "DD_BCA_KLIKPAY", "QRIS",
		},
		"currency": "IDR",
		"locale":   "en",
		"items":    is.generateInvoiceItems(forUserSub),
	}
	body, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", is.hostname, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("XENDIT_API_KEY"), "")

	resp, err := is.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		bytes, err := io.ReadAll(resp.Body)
		if err == nil {
			fmt.Printf("%s\n", string(string(bytes)))
		}
		// fmt.Printf("Here\n")
		return "", fmt.Errorf("xendit status code %d", resp.StatusCode)
	}

	var respBody map[string]any
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", nil
	}

	if url, ok := respBody["invoice_url"].(string); ok {
		return url, nil
	} else {
		return "", errors.New("failed to get invoice url from resp")
	}
}
