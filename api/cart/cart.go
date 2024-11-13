package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/middleware/auth"
	"mall/service/cart/proto/cart"
)

var CartClient cart.CartServiceClient
var log *mlog.Log

func Init(engine *gin.Engine) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "cart.rpc",
		},
	})
	CartClient = cart.NewCartServiceClient(conn.Conn())
	log = mlog.NewLog("CartAPI")
	group := engine.Group("/Cart", auth.ParseToken)
	{
		group.POST("/Add", Add)
		group.GET("/Get", Get)
		group.PUT("/Empty", Empty)
	}
}
