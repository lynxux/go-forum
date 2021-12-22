package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

/*
jwt token
*/

const TokenExpireDuration = time.Hour * 2

var mySigningSecret = []byte("link")
var ErrorInvalidToken = errors.New("Invalid Token!")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenToken(userID int64, username string) (string, error) {
	// 创建一个自己声明的数据
	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的secret签名获得完整的加密字符串
	return token.SignedString(mySigningSecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var myclaims = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, myclaims, func(token *jwt.Token) (interface{}, error) {
		return mySigningSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return myclaims, nil
	}
	return nil, ErrorInvalidToken
}
