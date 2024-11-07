package main

import (
	"mall/model/database"
)
import "gorm.io/gorm"
import "gorm.io/driver/mysql"

func main() {
	dsn := "root:2549124159f@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	err = db.SetupJoinTable(&database.Cart{}, "Products", &database.CartProducts{})
	if err != nil {
		panic(err)
	}
	err = db.SetupJoinTable(&database.Order{}, "Products", &database.OrderProducts{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.User{}, &database.Categories{}, &database.Product{}, &database.Cart{}, &database.Order{}, &database.Address{})
	if err != nil {
		panic(err)
	}
}
