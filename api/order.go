package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/model/api"
	"mall/service/cart/proto/cart"
	"mall/service/order/proto/order"
	"net/http"
)

var OrderClient order.OrderServiceClient
var CartClient cart.CartServiceClient

func InitOrder(engine *gin.Engine) {
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
	OrderClient = order.NewOrderServiceClient(OrderConn.Conn())
	CartClient = cart.NewCartServiceClient(CartConn.Conn())
	mlog.SetName("OrderAPI")

	group := engine.Group("/order")
	{
		group.POST("/CheckOut", CheckOut)
	}

}

func CheckOut(c *gin.Context) {
	req := api.CheckOutReq{}
	//v := c.Value("userid")
	id := uint32(1)
	//v = c.Value("email")
	email := ""
	if err := c.ShouldBind(&req); err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if len(req.ProductID) == 0 {
		GetResp, _ := CartClient.GetCart(c, &cart.GetCartReq{
			UserId: id,
		})
		for _, v := range GetResp.Items {
			req.Quantity = append(req.Quantity, v.Quantity)
			req.ProductID = append(req.ProductID, uint(v.ProductId))
		}
	}
	orderMake(c, req, id, email)

}

func orderMake(c *gin.Context, req api.CheckOutReq, id uint32, email string) {
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
