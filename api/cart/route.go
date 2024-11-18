package cart

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/model"
	"mall/service/cart/proto/cart"
	"mall/util"
)

func Add(c *gin.Context) {
	id := c.GetUint("userid")
	req := cart.AddItemReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	req.UserId = uint32(id)
	_, err := CartClient.AddItem(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "")
}

func Get(c *gin.Context) {
	id := c.GetUint("userid")

	resp, err := CartClient.GetCart(c, &cart.GetCartReq{UserId: uint32(id)})
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}

func Empty(c *gin.Context) {
	id := c.GetUint("userid")
	_, err := CartClient.EmptyCart(c, &cart.EmptyCartReq{UserId: uint32(id)})
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "")
}
