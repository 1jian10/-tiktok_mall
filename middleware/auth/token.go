package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	mlog "mall/log"
	"mall/model"
	"net/http"
	"strconv"
	"time"
)

var log *mlog.Log

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	log = mlog.NewLog("Auth", mlog.Info)
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		log.Warn("redis ping fail:" + err.Error())
		RDB = nil
	}
}

var Key = []byte("1-10 mall key")
var RDB *redis.Client

func GetToken(c *gin.Context) (string, error) {
	id := c.GetUint("userid")
	res, err := RDB.Get(c, "token:id:"+strconv.FormatUint(uint64(id), 10)).Result()
	if err == nil {
		return res, nil
	}

	claims := MyClaims{
		Userid: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(Key)
	if err == nil {
		b, err := json.Marshal(&claims)
		if err != nil {
			return str, err
		}
		RDB.Set(c, "token:"+str, string(b), time.Hour*24*7)
		RDB.Set(c, "token:id:"+strconv.FormatUint(uint64(id), 10), str, time.Hour*24*7)
	}
	return str, err
}

func DeleteToken(token string) {
	c := context.Background()
	if RDB == nil {
		log.Warn("redis is not connect")
		return
	}
	str, err := RDB.Get(c, "token:"+token).Result()
	if err != nil {
		log.Warn(err.Error())
		return
	}
	var m MyClaims
	err = json.Unmarshal([]byte(str), &m)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	RDB.Del(c, "token:id:"+strconv.FormatUint(uint64(m.Userid), 10))
	RDB.Del(c, "token:"+token)
}

func ParseToken(c *gin.Context) {
	var t Token
	t.Token = c.GetHeader("Authorization")
	log.Debug("token:" + t.Token)
	t.Ctx = c
	if RDB != nil {
		t.withRedis()
		return
	}
	t.direct()
}

func (t Token) withRedis() {
	c := t.Ctx
	res, err := RDB.Get(c, "token:"+t.Token).Result()
	if err != nil {
		log.Info("get token fail:" + err.Error())
		c.JSON(http.StatusOK, AuthResp{
			Status: model.Status{
				Code:     model.FORBIDDEN,
				ErrorMsg: "you need login to use it",
			},
		})
		c.Abort()
		return
	}
	log.Debug("get from redis" + fmt.Sprint(res))

	var m MyClaims
	err = json.Unmarshal([]byte(res), &m)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, AuthResp{
			Status: model.Status{
				Code:     model.ERROR,
				ErrorMsg: "json unmarshal failed",
			},
		})
		c.Abort()
		return
	}
	log.Debug("unmarshal:" + fmt.Sprint(m))
	c.Set("userid", m.Userid)
	c.Next()
	return
}

func (t Token) direct() {
	c := t.Ctx
	token, err := jwt.ParseWithClaims(t.Token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if err != nil {
		log.Info("direct parse token fail:" + err.Error())
		c.JSON(http.StatusOK, AuthResp{
			Status: model.Status{
				Code:     model.FORBIDDEN,
				ErrorMsg: "parse token failed",
			},
		})
		c.Abort()
		return
	}
	if !token.Valid {
		log.Info("token is invalid")
		c.JSON(http.StatusOK, AuthResp{
			Status: model.Status{
				Code:     model.FORBIDDEN,
				ErrorMsg: "token is invalid",
			},
		})
		c.Abort()
		return
	}
	c.Set("userid", token.Claims.(*MyClaims).Userid)
	c.Next()
	return
}
