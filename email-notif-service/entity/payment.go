package entity

import (
	"time"
)

type Payment struct {
	Id              int64      `json:"id,omitempty"`
	UserId          int64      `json:"user_id,omitempty"`
	PaymentGateway  string     `json:"payment_gateway,omitempty"`
	Amount          float32    `json:"amount,omitempty"`
	Currency        string     `json:"currency,omitempty"`
	TransactionDate *time.Time `json:"transaction_date,omitempty"` // can be null
	Status          string     `json:"status,omitempty"`
	Url             string     `json:"url,omitempty"`
}
