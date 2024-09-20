package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID          int     `gorm:"not null"`
	PaymentGateway  string  `gorm:"not null"`
	Amount          float64 `gorm:"not null"`
	Currency        string  `gorm:"not null"`
	TransactionDate *time.Time
	Status          string `gorm:"not null"`
	Url             string `gorm:"not null"`
}

type Subscription struct {
	gorm.Model
	Tier          string  `gorm:"type:varchar(50);unique;not null"` // e.g business
	PricePerMonth float64 `gorm:"type:decimal(10,2);not null"`
}

type UserSubscription struct {
	gorm.Model
	UserID         uint       `gorm:"not null"`       // Reference to users table
	SubscriptionID uint       `gorm:"not null"`       // Reference to subscriptions table
	Duration       int        `gorm:"not null"`       // Duration of the subscription
	EndDate        *time.Time `gorm:"type:timestamp"` // End date of the subscription
	PaymentID      uint       `gorm:"not null"`       // Reference to payments table
	Payment        Payment
}
