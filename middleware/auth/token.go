package auth

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	mlog "mall/log"
	"mall/model/auth"
	"net/http"
	"time"
)

type Token struct {
	Token string `json:"token;binding:required"`
	Ctx   *gin.Context
}

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		mlog.Warn("redis ping fail:" + err.Error())
		RDB = nil
	}
	mlog.SetName("Auth")
}

var Key = []byte("1-10 mall key")
var RDB *redis.Client

func GetToken(c *gin.Context) (string, error) {
	id := c.GetUint("userid")
	claims := auth.MyClaims{
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
	}
	return str, err
}

func DeleteToken(token string) {
	if RDB == nil {
		return
	}
	err := RDB.Del(context.Background(), "token:"+token).Err()
	if err != nil {
		mlog.Warn(err.Error())
		return
	}
}

func ParseToken(c *gin.Context) {
	var t Token
	if err := c.ShouldBind(&t); err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		c.Abort()
		return
	}
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
		mlog.Info("get token fail:" + err.Error())
		c.JSON(http.StatusForbidden, gin.H{})
		c.Abort()
		return
	}

	var m auth.MyClaims
	err = json.Unmarshal([]byte(res), &m)
	if err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		c.Abort()
		return
	}
	c.Set("userid", m.Userid)
	c.Next()
	return
}

func (t Token) direct() {
	c := t.Ctx
	token, err := jwt.ParseWithClaims(t.Token, &auth.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if err != nil {
		mlog.Info("direct parse token fail:" + err.Error())
		c.JSON(http.StatusForbidden, gin.H{})
		c.Abort()
		return
	}
	if !token.Valid {
		mlog.Info("token is invalid")
		c.JSON(http.StatusForbidden, gin.H{})
		c.Abort()
		return
	}
	c.Set("userid", token.Claims.(*auth.MyClaims).Userid)
	c.Next()
	return
}
