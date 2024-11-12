package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	mlog "mall/log"
	"mall/middleware/auth"
	"mall/service/user/proto/user"
	"net/http"
)

var UserClient user.UserServiceClient

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
		group.POST("/Logout", Logout)
		group.POST("/Info", auth.ParseToken, Info)
	}
}

var log *mlog.Log

func Register(c *gin.Context) {
	var req user.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	resp, _ := UserClient.Register(c, &req)
	c.JSON(http.StatusOK, resp)
	return
}

func Login(c *gin.Context) {
	var req user.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	resp, _ := UserClient.Login(c, &req)
	if resp.UserId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Success": false,
			"Token":   "",
		})
		return
	}

	c.Set("userid", uint(resp.UserId))
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Success": false,
			"Token":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	t := auth.Token{}
	if err := c.ShouldBind(&t); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	log.Debug("token:" + t.Token)
	auth.DeleteToken(t.Token)
	c.JSON(http.StatusOK, gin.H{})
}

func Info(c *gin.Context) {
	id := c.GetUint("userid")
	log.Debug("userid:" + fmt.Sprint(id))
	resp, _ := UserClient.Info(c, &user.InfoReq{UserId: uint32(id)})
	c.JSON(http.StatusOK, resp)
}
