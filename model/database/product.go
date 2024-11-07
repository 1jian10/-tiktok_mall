package database

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"unique" json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`
	Stock       uint

	Carts      []Cart       `gorm:"many2many:cart_products" json:"-"`
	Orders     []Order      `gorm:"many2many:order_products" json:"-"`
	Categories []Categories `gorm:"many2many:categories_products" json:"categories"`
}

type Categories struct {
	gorm.Model
	Name string `gorm:"unique"`

	Products []Product `gorm:"many2many:categories_products"`
}
