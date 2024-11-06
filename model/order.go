package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Currency string
	Paid     string
	Cost     float32

	UserID    uint      `gorm:"foreignKey:UserID"`
	Products  []Product `gorm:"many2many:order_products"`
	AddressID uint      `gorm:"foreignKey:AddressID"`
	Address   *Address
}

type OrderProducts struct {
	OrderID   uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  uint
}
