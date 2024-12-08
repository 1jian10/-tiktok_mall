package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string

	Cart     *Cart
	Orders   []Order
	Messages []Message
}

type Message struct {
	gorm.Model
	Message string `gorm:"message"`

	UserID uint `gorm:"primaryKey;index"`
}
