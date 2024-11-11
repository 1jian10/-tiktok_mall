package auth

import "github.com/dgrijalva/jwt-go"

type MyClaims struct {
	Userid uint `json:"user_id"`
	jwt.StandardClaims
}
