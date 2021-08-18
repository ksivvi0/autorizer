package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type AuthService interface {
	CreateTokenPair() (*tokenPair, error)
	RefreshTokens(rToken string) (*tokenPair, error)
	ValidateToken(token string, refresh bool) (string, error)
}

type Auth struct {
	AccessTokenInfo
	RefreshTokenInfo
	accessTokenKey  []byte
	refreshTokenKey []byte
}

type AccessTokenInfo struct {
	AccessTokenTTL time.Duration
	token          string
	uuid           string
}

type RefreshTokenInfo struct {
	RefreshTokenTTL time.Duration
	token           string
	uuid            string
}

func NewAuthInstance() *Auth {
	return &Auth{
		AccessTokenInfo: AccessTokenInfo{
			AccessTokenTTL: time.Minute * 5,
		},
		RefreshTokenInfo: RefreshTokenInfo{
			RefreshTokenTTL: time.Hour * 24,
		},
		accessTokenKey:  []byte("kdfjjhsdfpw"),    //insecure, read from config or environment
		refreshTokenKey: []byte("asasddhnkjasl8"), //insecure, read from config or environment
	}
}

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshToken string `json:"refresh_token"`
	RefreshUUID  string `json:"refresh_uuid"`
}

type tokenClaims struct {
	jwt.StandardClaims
	uuid uuid.UUID
}

func (a *Auth) CreateTokenPair() (*tokenPair, error) {
	rUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	aUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	atInfo, err := a.generateToken(aUuid, a.AccessTokenInfo)
	if err != nil {
		return nil, err
	}
	rtInfo, err := a.generateToken(rUuid, a.RefreshTokenInfo)
	if err != nil {
		return nil, err
	}

	pair := &tokenPair{
		AccessToken:  atInfo,
		AccessUUID:   aUuid.String(),
		RefreshToken: rtInfo,
		RefreshUUID:  rUuid.String(),
	}
	return pair, nil
}

func (a *Auth) generateToken(uuid uuid.UUID, tokenInfo interface{}) (string, error) {
	switch tokenInfo.(type) {
	case AccessTokenInfo:
		claims := jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			tokenClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(a.AccessTokenTTL).Unix(),
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

	case RefreshTokenInfo:
		claims := jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			tokenClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(a.RefreshTokenTTL).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				uuid: uuid,
			},
		)
		rToken, err := claims.SignedString(a.refreshTokenKey)
		if err != nil {
			return "", err
		}
		return rToken, nil
	default:
		return "", errors.New("invalid argument")
	}
}

func (a *Auth) RefreshTokens(rToken string) (*tokenPair, error) {
	//pair := new(tokenPair)
	return nil, errors.New("method not implemented")
}

func (a *Auth) ValidateToken(tokenIn string, refresh bool) (string, error) {

	//TODO: get information about token in DB
	token, err := jwt.ParseWithClaims(tokenIn, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New(fmt.Sprintf("invalid signing method %v", token.Header))
		}
		if !refresh {
			return a.accessTokenKey, nil
		}
		return a.refreshTokenKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("bad token")
	}

	return claims.uuid.String(), nil
}
