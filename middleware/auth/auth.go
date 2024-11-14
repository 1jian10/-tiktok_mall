package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mall/api"
)

type Token struct {
	Token string `header:"Authorization" binding:"required"`
	Ctx   *gin.Context
}

type MyClaims struct {
	Userid uint `json:"user_id"`
	jwt.StandardClaims
}

type AuthResp struct {
	Status api.Status `json:"status"`
}
