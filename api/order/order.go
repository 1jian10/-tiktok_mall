package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/middleware/auth"
	"mall/model/api"
	"mall/service/cart/proto/cart"
	"mall/service/order/proto/order"
	"net/http"
	"strconv"
)

var OrderClient order.OrderServiceClient
var CartClient cart.CartServiceClient
var log *mlog.Log

func Init(engine *gin.Engine) {
	OrderConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "order.rpc",
		},
	})
	CartConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "cart.rpc",
		},
	})
	log = mlog.NewLog("OrderAPI")
	OrderClient = order.NewOrderServiceClient(OrderConn.Conn())
	CartClient = cart.NewCartServiceClient(CartConn.Conn())

	group := engine.Group("/Order", auth.ParseToken)
	{
		group.POST("/CheckOut", CheckOut)
		group.POST("/Charge", Charge)
		group.POST("/List", List)
	}

}

func List(c *gin.Context) {
	id := c.GetUint("userid")
	resp, _ := OrderClient.ListOrder(c, &order.ListOrderReq{UserId: uint32(id)})
	c.JSON(http.StatusOK, resp.Orders)
}

func CheckOut(c *gin.Context) {
	req := api.CheckOutReq{}
	id := c.GetUint("userid")
	//v = c.Value("email")
	email := ""
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if len(req.ProductID) == 0 {
		GetResp, _ := CartClient.GetCart(c, &cart.GetCartReq{
			UserId: uint32(id),
		})
		_, _ = CartClient.EmptyCart(c, &cart.EmptyCartReq{UserId: uint32(id)})
		for _, v := range GetResp.Items {
			req.Quantity = append(req.Quantity, v.Quantity)
			req.ProductID = append(req.ProductID, uint(v.ProductId))
		}
	}
	Make(c, req, uint32(id), email)
}

func Charge(c *gin.Context) {
	id := c.GetUint("userid")
	req := api.ChargeReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	_, err := OrderClient.MarkOrderPaid(context.Background(), &order.MarkOrderPaidReq{
		UserId:  uint32(id),
		OrderId: strconv.Itoa(int(req.OrderId)), //maybe bug
	})
	if err != nil {
		c.JSON(http.StatusOK, api.ChargeResp{Success: false})
	} else {
		c.JSON(http.StatusOK, api.ChargeResp{Success: true})
	}
}

func Make(c *gin.Context, req api.CheckOutReq, id uint32, email string) {
	PlaceReq := order.PlaceOrderReq{
		Items:  make([]*order.CartItem, len(req.ProductID)),
		UserId: id,
	}
	for i, v := range req.ProductID {
		PlaceReq.Items[i] = &order.CartItem{
			ProductId: uint32(v),
			Quantity:  req.Quantity[i],
		}
	}
	PlaceResp, _ := OrderClient.PlaceOrder(c, &PlaceReq)
	if PlaceResp.Success == "Yes" {
		c.JSON(http.StatusOK, &api.CheckOutResp{
			Success: true,
		})
		ProcessReq := &order.ProcessOrderReq{
			UserId:       id,
			UserCurrency: "default",
			Address: &order.Address{
				StreetAddress: req.Address.StreetAddress,
				City:          req.Address.City,
				State:         req.Address.State,
				Country:       req.Address.Country,
				ZipCode:       int32(req.ZipCode),
			},
			Email:      email,
			OrderItems: make([]*order.OrderItem, len(req.ProductID)),
		}
		for i, v := range req.ProductID {
			ProcessReq.OrderItems[i] = &order.OrderItem{
				Item: &order.CartItem{
					ProductId: uint32(v),
					Quantity:  req.Quantity[i],
				},
				Cost: 0,
			}
		}
		_, _ = OrderClient.ProcessOrder(c, ProcessReq)
	} else {
		c.JSON(http.StatusOK, api.CheckOutResp{
			Success: false,
		})
	}
}
