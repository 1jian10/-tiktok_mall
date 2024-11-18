package user

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/middleware/auth"
	"mall/service/user/proto/user"
)

var UserClient user.UserServiceClient
var log *mlog.Log

func Init(engine *gin.Engine) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:4379"},
			Key:   "user.rpc",
		},
	})
	UserClient = user.NewUserServiceClient(conn.Conn())
	log = mlog.NewLog("UserAPI")
	group := engine.Group("/User")
	{
		group.POST("/Register", Register)
		group.POST("/Login", Login)
		group.GET("/Logout", Logout)
		group.GET("/Info", auth.ParseToken, Info)
		group.GET("/Message", auth.ParseToken, Message)
		group.DELETE("/Delete", auth.ParseToken, Delete)
		group.PUT("/Update", auth.ParseToken, Update)
	}
}
