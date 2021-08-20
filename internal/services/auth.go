package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

type AuthService interface {
	CreateTokenPair(string) (*tokenPair, error)
	ValidateToken(string, bool) (string, error)
	GetDataFromToken(string, bool) (string, error)
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
	uid            string
}

type RefreshTokenInfo struct {
	RefreshTokenTTL time.Duration
	token           string
	uid             string
}

func NewAuthInstance() *Auth {
	return &Auth{
		AccessTokenInfo: AccessTokenInfo{
			AccessTokenTTL: time.Second * 600,
		},
		RefreshTokenInfo: RefreshTokenInfo{
			RefreshTokenTTL: time.Hour * 72,
		},
		accessTokenKey:  []byte(os.Getenv("ACCESS_KEY")),
		refreshTokenKey: []byte(os.Getenv("REFRESH_KEY")),
	}
}

type tokenPair struct {
	AccessToken    string    `json:"access_token,omitempty" bson:"-"`
	AccessUID      string    `json:"access_token_uid,omitempty" bson:"access_token_uid,omitempty"`
	AccessExpired  time.Time `json:"access_token_expired,omitempty" bson:"access_token_expired,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty" bson:"refresh_token_hash,omitempty"`
	RefreshUID     string    `json:"refresh_token_uid,omitempty" bson:"refresh_token_uid,omitempty"`
	RefreshExpired time.Time `json:"refresh_token_expired,omitempty" bson:"refresh_token_expired,omitempty"`
}

func (a *Auth) CreateTokenPair(refreshUid string) (*tokenPair, error) {
	pair := new(tokenPair)
	refreshTokenExist := len(refreshUid) > 0

	aUid, err := generateUUID()
	if err != nil {
		return nil, err
	}

	aToken, aExpired, err := a.generateToken(aUid, a.AccessTokenInfo)
	if err != nil {
		return nil, err
	}

	pair.AccessUID = aUid.String()
	pair.AccessToken = aToken
	pair.AccessExpired = time.Unix(aExpired, 0)
	pair.RefreshUID = refreshUid

	if !refreshTokenExist {
		rUid, err := generateUUID()
		if err != nil {
			return nil, err
		}

		rToken, rExpired, err := a.generateToken(rUid, a.RefreshTokenInfo)
		if err != nil {
			return nil, err
		}

		pair.RefreshUID = rUid.String()
		pair.RefreshExpired = time.Unix(rExpired, 0)
		pair.RefreshToken = rToken
	}

	return pair, nil
}

func generateUUID() (uuid.UUID, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return [16]byte{}, err
	}
	return uid, nil
}

func (a *Auth) generateToken(uid uuid.UUID, tokenInfo interface{}) (string, int64, error) {
	var expTime int64
	var key []byte
	switch tokenInfo.(type) {
	case AccessTokenInfo:
		expTime = time.Now().Add(a.AccessTokenTTL).Unix()
		key = a.accessTokenKey
	case RefreshTokenInfo:
		expTime = time.Now().Add(a.RefreshTokenTTL).Unix()
		key = a.refreshTokenKey
	default:
		return "", -1, errors.New("invalid argument")
	}

	claims := jwt.MapClaims{}
	claims["iat"] = time.Now().Unix()
	claims["uid"] = uid
	claims["exp"] = expTime
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := at.SignedString(key)
	if err != nil {
		return "", -1, err
	}
	return token, expTime, nil
}

func (a *Auth) ValidateToken(bearerHeader string, refresh bool) (string, error) {
	tokenStr, err := parseHeader(bearerHeader)
	if err != nil {
		return "", err
	}
	uid, err := a.GetDataFromToken(tokenStr, refresh)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func parseHeader(bearerHeader string) (string, error) {
	bearerHeaderArr := strings.Split(bearerHeader, " ")
	if len(bearerHeaderArr) != 2 || bearerHeaderArr[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	token := bearerHeaderArr[1]
	if len(token) == 0 {
		return "", errors.New("invalid authorization header")
	}
	return token, nil
}

func (a *Auth) GetDataFromToken(tokenStr string, refresh bool) (string, error) {
	c := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("invalid signing method %v", token.Header))
		}
		if refresh {
			return a.refreshTokenKey, nil
		}
		return a.accessTokenKey, nil
	})
	if err != nil || !jwtToken.Valid {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	return claims["uid"].(string), nil
}
