package util

import (
	"github.com/gin-gonic/gin"
	"mall/model"
	"net/http"
)

// Response 封装响应，统一http返回格式
func Response(c *gin.Context, code int, msg string, data ...interface{}) {
	response := model.Response{
		Status: model.Status{
			Code:     uint(code),
			ErrorMsg: msg,
		},
	}
	if len(data) != 0 {
		response.Data = data[0]
	}
	c.JSON(http.StatusOK, response)
}
