package main

import (
	"github.com/gin-gonic/gin"
	"mall/api/cart"
	"mall/api/file"
	"mall/api/order"
	Product "mall/api/product"
	"mall/api/user"
)

func main() {
	engine := gin.Default()

	order.Init(engine)
	cart.Init(engine)
	user.Init(engine)
	Product.Init(engine)
	file.Init(engine)

	err := engine.Run("0.0.0.0:23333")
	if err != nil {
		panic(err)
	}
}
