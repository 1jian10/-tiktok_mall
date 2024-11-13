package order

import (
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/middleware/auth"
	"mall/service/cart/proto/cart"
	"mall/service/order/proto/order"
)

var OrderClient order.OrderServiceClient
var CartClient cart.CartServiceClient
var log *mlog.Log
var producer *nsq.Producer

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
	config := nsq.NewConfig()
	p, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		panic(err)
	}
	producer = p

	group := engine.Group("/Order", auth.ParseToken)
	{
		group.POST("/CheckOut", CheckOut)
		group.POST("/Charge", Charge)
		group.POST("/List", List)
	}
}
