package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/service/cart/proto/cart"
	"net/http"
)

var CartClient cart.CartServiceClient

func Init(engine *gin.Engine) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "cart.rpc",
		},
	})
	CartClient = cart.NewCartServiceClient(conn.Conn())
	mlog.SetName("CartAPI")
	group := engine.Group("/Cart")
	{
		group.POST("/Add", Add)
		group.GET("/Get", Get)
		group.GET("/Empty", Empty)
	}
}

func Add(c *gin.Context) {
	//v:=c.Value("userid")
	id := uint32(1)
	req := cart.AddItemReq{}
	if err := c.ShouldBind(&req); err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	req.UserId = id
	_, err := CartClient.AddItem(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"OK": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func Get(c *gin.Context) {
	//v:=c.Value("userid")
	id := uint32(1)

	resp, _ := CartClient.GetCart(c, &cart.GetCartReq{UserId: id})
	c.JSON(http.StatusOK, resp)
}

func Empty(c *gin.Context) {
	//v:=c.Value("userid")
	id := uint32(1)
	_, err := CartClient.EmptyCart(c, &cart.EmptyCartReq{UserId: id})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"OK": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"OK": true})
}
