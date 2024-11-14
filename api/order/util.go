package order

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/service/cart/proto/cart"
	"mall/service/order/proto/order"
	"net/http"
)

func List(c *gin.Context) {
	id := c.GetUint("userid")
	resp, err := OrderClient.ListOrder(c, &order.ListOrderReq{UserId: uint32(id)})
	if err != nil {
		c.JSON(http.StatusOK, ListResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, ListResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
}

func CheckOut(c *gin.Context) {
	req := &order.ProcessOrderReq{}
	id := c.GetUint("userid")
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, CheckOutResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	req.UserId = uint32(id)

	if len(req.OrderItems) == 0 {
		GetResp, err := CartClient.GetCart(c, &cart.GetCartReq{UserId: uint32(id)})
		if err != nil {
			c.JSON(http.StatusOK, CheckOutResp{
				Status: api.Status{
					Code:     api.ERROR,
					ErrorMsg: err.Error(),
				},
			})
			return
		}
		_, _ = CartClient.EmptyCart(c, &cart.EmptyCartReq{UserId: uint32(id)})
		for _, v := range GetResp.Items {
			req.OrderItems = append(req.OrderItems, &order.OrderItem{
				ProductId: v.ProductId,
				Quantity:  int32(v.Quantity),
			})
		}
	}
	if !IsSync {
		ASyncMake(c, req)
		return
	}
	SyncMake(c, req)
}

func Charge(c *gin.Context) {
	id := c.GetUint("userid")
	req := order.MarkOrderPaidReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, ChargeResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	req.UserId = uint32(id)

	_, err := OrderClient.MarkOrderPaid(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, ChargeResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, ChargeResp{
		Status: api.Status{
			Code: api.OK,
		},
	})

}

func ASyncMake(c *gin.Context, req *order.ProcessOrderReq) {
	id := c.GetUint("userid")
	PlaceReq := order.PlaceOrderReq{
		ProductId: make([]uint32, len(req.OrderItems)),
		Quantity:  make([]int32, len(req.OrderItems)),
		UserId:    uint32(id),
	}
	for i, v := range req.OrderItems {
		PlaceReq.ProductId[i] = v.ProductId
		PlaceReq.Quantity[i] = v.Quantity
	}
	_, err := OrderClient.PlaceOrder(c, &PlaceReq)
	if err != nil {
		c.JSON(http.StatusOK, CheckOutResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	res, err := json.Marshal(&req)
	if err != nil {
		c.JSON(http.StatusOK, CheckOutResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: "json marshal:" + err.Error(),
			},
		})
	}
	err = producer.Publish("order", res)
	if err != nil {
		c.JSON(http.StatusOK, CheckOutResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: "mq publish:" + err.Error(),
			},
		})
	}
	c.JSON(http.StatusOK, CheckOutResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
}

func SyncMake(c *gin.Context, req *order.ProcessOrderReq) {
	_, err := OrderClient.ProcessOrder(c, req)
	if err != nil {
		c.JSON(http.StatusOK, CheckOutResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, CheckOutResp{
		Status: api.Status{
			Code: api.OK,
		},
	})

}
