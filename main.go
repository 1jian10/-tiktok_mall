package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()

	err := engine.Run("0.0.0.0:23333")
	if err != nil {
		panic(err)
	}
}
