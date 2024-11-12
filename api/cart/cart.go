package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/middleware/auth"
	"mall/service/cart/proto/cart"
	"net/http"
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
		group.POST("/Get", Get)
		group.POST("/Empty", Empty)
	}
}

func Add(c *gin.Context) {
	id := c.GetUint("userid")
	req := cart.AddItemReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	req.UserId = uint32(id)
	_, err := CartClient.AddItem(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"OK": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func Get(c *gin.Context) {
	id := c.GetUint("userid")

	resp, _ := CartClient.GetCart(c, &cart.GetCartReq{UserId: uint32(id)})
	c.JSON(http.StatusOK, resp)
}

func Empty(c *gin.Context) {
	id := c.GetUint("userid")
	_, err := CartClient.EmptyCart(c, &cart.EmptyCartReq{UserId: uint32(id)})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"OK": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"OK": true})
}
