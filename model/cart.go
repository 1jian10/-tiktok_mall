package model

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model

	UserID   uint      `gorm:"foreignKey:UserID;index"`
	Products []Product `gorm:"many2many:cart_products"`
}

type CartProducts struct {
	CartID    uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  uint
}
