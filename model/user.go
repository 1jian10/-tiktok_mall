package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string

	Cart    *Cart
	Address *Address
}

type Address struct {
	gorm.Model
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32

	UserID uint `gorm:"foreignKey:UserID"`
}
