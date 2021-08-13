package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Auth struct {
	accessTokenTTL  time.Duration
	accessTokenKey  string
	refreshTokenTTL time.Duration
	refreshTokenKey string
}

func NewAuth() *Auth {
	return &Auth{
		accessTokenTTL:  time.Minute * 5,
		accessTokenKey:  "99873dachas7d",
		refreshTokenTTL: time.Hour * 8,
		refreshTokenKey: "aklsjdfhsdfh1",
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	GUID uuid.UUID
}

//func (a *Auth) GenerateToken(username, password string) (string, error) {
//svc := NewNebulaService(a.nebulaAddress)
//
//if !svc.CheckAuthorization(username, password) {
//	return "", errors.New("неверные аутентификационные данные")
//}
//
//token := jwt.NewWithClaims(
//	jwt.SigningMethodHS256,
//	&tokenClaims{
//		jwt.StandardClaims{
//			ExpiresAt: time.Now().Add(a.tokenTTL).Unix(),
//			IssuedAt:  time.Now().Unix(),
//		},
//		username,
//		helpers.EncodeThis(password),
//	},
//)
//return token.SignedString([]byte(a.tokenKey))
//}

//func (a *Auth) ParseToken(tokenIn string) ([]string, error) {
//token, err := jwt.ParseWithClaims(tokenIn, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//		return nil, errors.New("неверный алгоритм подписи токена")
//	}
//	return []byte(a.tokenKey), nil
//})
//if err != nil {
//	return nil, err
//}
//
//claims, ok := token.Claims.(*tokenClaims)
//if !ok {
//	return nil, errors.New("bad token")
//}
//
//params := make([]string, 2)
//params[0] = claims.PNumber
//params[1] = claims.HashPasswd
//
//return params, nil
//}
