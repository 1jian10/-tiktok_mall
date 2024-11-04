package model

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model

	UserID   uint      `gorm:"foreignKey:UserID"`
	Products []Product `gorm:"many2many:cart_products"`
}
