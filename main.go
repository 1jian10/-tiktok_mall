package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
	"mall/model"
)

func main() {
	dsn := "root:2549124159f@tcp(127.0.0.1:3306)/tiktok_mall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = db.AutoMigrate(&model.Product{}, &model.User{}, &model.Categories{}, &model.Address{}, &model.Order{}, &model.Cart{})
	if err != nil {
		slog.Error(err.Error())
	}
	return
}
