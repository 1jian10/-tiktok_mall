package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Currency string
	Paid     string
	Cost     float32

	UserID    uint      `gorm:"foreignKey:UserID;index"`
	Products  []Product `gorm:"many2many:order_products"`
	AddressID uint      `gorm:"foreignKey:AddressID;index"`
	Address   *Address
}

type Address struct {
	gorm.Model
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

type OrderProducts struct {
	OrderID   uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  uint
}
