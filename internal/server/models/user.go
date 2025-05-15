package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email               string `gorm:"unique;not null"`
	Username            string `gorm:"not null"`
	PublicKey           string `gorm:"default:''"`
	EncryptedPrivateKey string `gorm:"default:''"`
	SecurityToken       string `gorm:"default:''"`
	OtpVerified         bool   `gorm:"default:false"`
}
