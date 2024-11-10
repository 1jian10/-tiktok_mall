package main

import (
	"github.com/gin-gonic/gin"
	"mall/api/order"
)

func main() {
	engine := gin.Default()
	order.Init(engine)

	err := engine.Run("0.0.0.0:23333")
	if err != nil {
		panic(err)
	}
}
