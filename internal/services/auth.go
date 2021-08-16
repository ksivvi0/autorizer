package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Auth struct {
	AccessTokenInfo
	RefreshTokenInfo
}

type AccessTokenInfo struct {
	AccessTokenTTL time.Duration
	AccessTokenKey []byte
	token          string
}

type RefreshTokenInfo struct {
	RefreshTokenTTL time.Duration
	RefreshTokenKey []byte
	token           string
}

func NewAuth() *Auth {
	return &Auth{
		AccessTokenInfo: AccessTokenInfo{
			AccessTokenTTL: time.Minute * 5,
			AccessTokenKey: []byte("99873dachas7d"), //insecure
		},
		RefreshTokenInfo: RefreshTokenInfo{
			RefreshTokenTTL: time.Hour * 24,
			RefreshTokenKey: []byte("aklsjdfhsdfh1"), //insecure
		},
	}
}

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type tokenClaims struct {
	jwt.StandardClaims
	uuid uuid.UUID
}

func (a *Auth) GetTokenPair() (*tokenPair, error) {
	accessToken, err := a.generateToken(a.AccessTokenInfo)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateToken(a.RefreshTokenInfo)
	if err != nil {
		return nil, err
	}

	pair := &tokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return pair, nil
}

func (a *Auth) generateToken(tokenInfo interface{}) (string, error) {
	switch tokenInfo.(type) {
	case AccessTokenInfo:
		_uuid, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		claims := jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			tokenClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(a.AccessTokenTTL).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				uuid: _uuid,
			},
		)
		aToken, err := claims.SignedString(a.AccessTokenKey)
		if err != nil {
			return "", err
		}
		return aToken, nil

	case RefreshTokenInfo:
		_uuid, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		claims := jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			tokenClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(a.RefreshTokenTTL).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				uuid: _uuid,
			},
		)
		rToken, err := claims.SignedString(a.RefreshTokenKey)
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
