package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null" json:"username"`
	Email        string `gorm:"unique;not null" json:"email"`
	Phonenumber  string `gorm:"unique;not null" json:"phone_number"`
	PasswordHash string `gorm:"not null" json:"password"`
	Tier         string `gorm:"type:varchar(20);not null;check:tier in ('free', 'business')" json:"tier"`
}
