package entity

import "time"

type User struct {
	ID          int
	Username    string
	Email       string
	PhoneNumber string
	Tier        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
