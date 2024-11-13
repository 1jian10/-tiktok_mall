package cart

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/service/cart/proto/cart"
	"net/http"
)

func Add(c *gin.Context) {
	id := c.GetUint("userid")
	req := cart.AddItemReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, AddResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	req.UserId = uint32(id)
	_, err := CartClient.AddItem(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, AddResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, AddResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
}

func Get(c *gin.Context) {
	id := c.GetUint("userid")

	resp, err := CartClient.GetCart(c, &cart.GetCartReq{UserId: uint32(id)})
	if err != nil {
		c.JSON(http.StatusOK, GetResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, GetResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
}

func Empty(c *gin.Context) {
	id := c.GetUint("userid")
	_, err := CartClient.EmptyCart(c, &cart.EmptyCartReq{UserId: uint32(id)})
	if err != nil {
		c.JSON(http.StatusOK, EmptyResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, EmptyResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
}
