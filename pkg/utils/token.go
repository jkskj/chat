package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JwtSecret  指定加密密钥
var JwtSecret = []byte("chat")

type Claims struct {
	Id        uint   `json:"id"`
	UserName  string `json:"user_name" `
	Authority int    `json:"authority" `
	jwt.StandardClaims
}

// GenerateToken 根据用户的用户名和密码产生token
func GenerateToken(id uint, username string, authority int) (string, error) {
	//设置token有效时间
	notTime := GetLocalTime()
	expireTime := notTime.Add(24 * time.Hour)
	claims := Claims{
		Id:        id,
		UserName:  username,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "chat",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JwtSecret)
	return token, err
}

// ParseToken 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*Claims, error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}
