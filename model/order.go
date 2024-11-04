package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Currency string

	UserID   uint      `gorm:"foreignKey:UserID"`
	Products []Product `gorm:"many2many:order_products"`
}
