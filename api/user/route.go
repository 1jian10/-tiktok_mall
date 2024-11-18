package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/middleware/auth"
	"mall/model"
	"mall/service/user/proto/user"
	"mall/util"
)

func Register(c *gin.Context) {
	var req user.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	resp, err := UserClient.Register(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	log.Debug("register id:" + fmt.Sprint(resp.UserId))
	util.Response(c, model.OK, "")
	return
}

func Login(c *gin.Context) {
	var req user.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}

	resp, err := UserClient.Login(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}

	c.Set("userid", uint(resp.UserId))
	token, err := auth.GetToken(c)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}

	util.Response(c, model.OK, "", gin.H{"token": token})
}

func Logout(c *gin.Context) {
	token := c.GetHeader("authorization")

	log.Debug("token:" + token)
	auth.DeleteToken(token)
	util.Response(c, model.OK, "")
}

func Info(c *gin.Context) {
	id := c.GetUint("userid")
	token := c.GetHeader("authorization")
	log.Debug("userid:" + fmt.Sprint(id))
	resp, err := UserClient.Info(c, &user.InfoReq{UserId: uint32(id)})
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	auth.DeleteToken(token)
	util.Response(c, model.OK, "", resp)
}

func Message(c *gin.Context) {
	id := c.GetUint("userid")

	resp, err := UserClient.GetMessage(c, &user.GetMessageReq{UserId: uint32(id)})
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}

func Delete(c *gin.Context) {
	id := c.GetUint("userid")

	_, err := UserClient.Delete(c, &user.DeleteReq{UserId: uint32(id)})
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "")
}

func Update(c *gin.Context) {
	id := c.GetUint("userid")
	req := user.UpdateReq{Password: ""}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	req.UserId = uint32(id)
	_, err := UserClient.Update(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "")
}
