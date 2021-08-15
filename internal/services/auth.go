package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Auth struct {
	accessTokenInfo
	refreshTokenInfo
}

type accessTokenInfo struct {
	accessTokenTTL time.Duration
	accessTokenKey string
	token          string
}

type refreshTokenInfo struct {
	refreshTokenTTL time.Duration
	refreshTokenKey string
	token           string
}

func NewAuth() *Auth {
	return &Auth{
		accessTokenInfo: accessTokenInfo{
			accessTokenTTL: time.Minute * 5,
			accessTokenKey: "99873dachas7d", //insecure
		},
		refreshTokenInfo: refreshTokenInfo{
			refreshTokenTTL: time.Hour * 24,
			refreshTokenKey: "aklsjdfhsdfh1", //insecure
		},
	}
}

type tokenPair struct {
	accessToken  string
	refreshToken string
}

type tokenClaims struct {
	jwt.StandardClaims
	uuid uuid.UUID
}

func (a *Auth) GetTokenPair() (*tokenPair, error) {
	newUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	accessToken, err := a.generateToken(newUuid, a.accessTokenInfo)
	refreshToken, err := a.generateToken(newUuid, a.refreshTokenInfo)

	pair := &tokenPair{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
	return pair, nil
}

func (a *Auth) generateToken(uuid uuid.UUID, tokenInfo interface{}) (string, error) {
	switch tokenInfo.(type) {
	case accessTokenInfo:
		claims := jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			tokenClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(a.accessTokenTTL).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				uuid: uuid,
			},
		)
		aToken, err := claims.SignedString(a.accessTokenKey)
		if err != nil {
			return "", err
		}
		return aToken, nil
	case refreshTokenInfo:
		claims := jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			tokenClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(a.accessTokenTTL).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				uuid: uuid,
			},
		)
		rToken, err := claims.SignedString(a.accessTokenKey)
		if err != nil {
			return "", err
		}
		return rToken, nil
	default:
		return "", errors.New("invalid argument")
	}
}

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
