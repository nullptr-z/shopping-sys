package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// TODO: 添加 Refresh Token，刷新 Access Token

var TokenExpireDuration = time.Hour * 2 // 两小时
// var secret = []byte("门前大桥下")            // 加盐

type JWTClaims interface {
	jwt.StandardClaims
}

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

// 生成 Token
func GenToken(claims *CustomClaims) (string, error) {
	claims.NotBefore = time.Now().Unix()          //生效时间
	claims.ExpiresAt = int64(TokenExpireDuration) // 过期时间
	claims.Issuer = viper.GetString("jwt.issuer") // 签发人
	var secret = viper.GetString("jwt.signs_key") // 加盐

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// 解析 Token
func ParseTOken(tokenStr string) (*CustomClaims, error) {
	var claims = new(CustomClaims)
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		var secret = viper.GetString("jwt.signs_key")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验 Token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
