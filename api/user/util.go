package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/middleware/auth"
	"mall/service/user/proto/user"
	"net/http"
)

func Register(c *gin.Context) {
	var req user.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, RegisterResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	if _, err := UserClient.Register(c, &req); err != nil {
		c.JSON(http.StatusOK, RegisterResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, RegisterResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
	return
}

func Login(c *gin.Context) {
	var req user.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, LoginResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
	}

	resp, err := UserClient.Login(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, LoginResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}

	c.Set("userid", uint(resp.UserId))
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusOK, LoginResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, LoginResp{
		Status: api.Status{
			Code: api.OK,
		},
		Token: token,
	})
}

func Logout(c *gin.Context) {
	token := c.GetHeader("authorization")

	log.Debug("token:" + token)
	auth.DeleteToken(token)
	c.JSON(http.StatusOK, LogOutResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
}

func Info(c *gin.Context) {
	id := c.GetUint("userid")
	log.Debug("userid:" + fmt.Sprint(id))
	resp, err := UserClient.Info(c, &user.InfoReq{UserId: uint32(id)})
	if err != nil {
		c.JSON(http.StatusOK, InfoResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, InfoResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
}

func Delete(c *gin.Context) {
	id := c.GetUint("userid")

	_, err := UserClient.Delete(c, &user.DeleteReq{UserId: uint32(id)})
	if err != nil {
		c.JSON(http.StatusOK, DeleteResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, DeleteResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
}

func Update(c *gin.Context) {
	id := c.GetUint("userid")
	req := user.UpdateReq{Password: ""}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, UpdateResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	req.UserId = uint32(id)
	_, err := UserClient.Update(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, UpdateResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UpdateResp{
		Status: api.Status{
			Code: api.OK,
		},
	})
}
