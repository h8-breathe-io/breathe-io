package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Phonenumber  string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Tier         string    `gorm:"type:varchar(20);not null;check:tier in ('free', 'business')"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
