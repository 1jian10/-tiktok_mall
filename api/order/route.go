package order

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/model"
	"mall/service/cart/proto/cart"
	"mall/service/order/proto/order"
	"mall/util"
)

func List(c *gin.Context) {
	id := c.GetUint("userid")
	resp, err := OrderClient.ListOrder(c, &order.ListOrderReq{UserId: uint32(id)})
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.ERROR, "", resp)
}

func CheckOut(c *gin.Context) {
	req := &order.ProcessOrderReq{}
	id := c.GetUint("userid")
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	req.UserId = uint32(id)

	if len(req.OrderItems) == 0 {
		GetResp, err := CartClient.GetCart(c, &cart.GetCartReq{UserId: uint32(id)})
		if err != nil {
			util.Response(c, model.ERROR, "")
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
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	req.UserId = uint32(id)

	_, err := OrderClient.MarkOrderPaid(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "")

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
	resp, err := OrderClient.PlaceOrder(c, &PlaceReq)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	req.OrderId = resp.OrderId
	res, err := json.Marshal(&req)
	if err != nil {
		util.Response(c, model.ERROR, "json marshal:"+err.Error())
	}
	err = producer.Publish("order", res)
	if err != nil {
		util.Response(c, model.ERROR, "mq publish:"+err.Error())
	}
	util.Response(c, model.OK, "", resp)
}

func SyncMake(c *gin.Context, req *order.ProcessOrderReq) {
	resp, err := OrderClient.ProcessOrder(c, req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}
