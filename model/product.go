package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	// 模糊查询 name like %abc% 索引失效
	Name      string  `gorm:"unique;index" json:"name"`
	ImagePath string  `json:"image_path"`
	FilePath  string  `json:"file_path"`
	Price     float32 `json:"price"`
	Stock     uint

	Carts      []Cart       `gorm:"many2many:cart_products" json:"-"`
	Orders     []Order      `gorm:"many2many:order_products" json:"-"`
	Categories []Categories `gorm:"many2many:categories_products" json:"categories"`
}

type Categories struct {
	gorm.Model
	Name string `gorm:"unique;index"`

	Products []Product `gorm:"many2many:categories_products"`
}
