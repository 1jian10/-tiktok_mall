package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string

	Cart   *Cart
	Orders []Order
}

type Address struct {
	gorm.Model
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}
