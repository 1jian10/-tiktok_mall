package main

import (
	"github.com/gin-gonic/gin"
	"mall/api"
)

func main() {
	engine := gin.Default()
	api.InitOrder(engine)

	err := engine.Run("0.0.0.0:23333")
	if err != nil {
		panic(err)
	}
}
