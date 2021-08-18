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
	CreateTokenPair() (*tokenPair, error)
	RefreshTokens(string) (*tokenPair, error)
	ValidateToken(string) (string, error)
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
			AccessTokenTTL: time.Minute * 15,
		},
		RefreshTokenInfo: RefreshTokenInfo{
			RefreshTokenTTL: time.Hour * 24,
		},
		accessTokenKey:  []byte(os.Getenv("ACCESS_KEY")),
		refreshTokenKey: []byte(os.Getenv("REFRESH_KEY")),
	}
}

type tokenPair struct {
	AccessToken    string    `json:"access_token,omitempty" bson:"access_token,omitempty"`
	AccessUID      string    `json:"access_uid" bson:"access_uid"`
	AccessExpired  time.Time `json:"access_expired,omitempty" bson:"expires_at,omitempty"`
	RefreshToken   string    `json:"refresh_token" bson:"refresh_token"`
	RefreshUID     string    `json:"refresh_uid" bson:"refresh_uid"`
	RefreshExpired time.Time `json:"refresh_expired,omitempty" bson:"refresh_expires,omitempty"`
}

func (a *Auth) CreateTokenPair() (*tokenPair, error) {
	rUid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	aUid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	aToken, aExpired, err := a.generateToken(aUid, a.AccessTokenInfo)
	if err != nil {
		return nil, err
	}
	rToken, rExpired, err := a.generateToken(rUid, a.RefreshTokenInfo)
	if err != nil {
		return nil, err
	}

	pair := &tokenPair{
		AccessToken:    aToken,
		AccessUID:      aUid.String(),
		AccessExpired:  time.Unix(aExpired, 0),
		RefreshToken:   rToken,
		RefreshUID:     rUid.String(),
		RefreshExpired: time.Unix(rExpired, 0),
	}
	return pair, nil
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

func (a *Auth) RefreshTokens(rToken string) (*tokenPair, error) {
	//pair := new(tokenPair)
	return nil, errors.New("method not implemented")
}

func (a *Auth) ValidateToken(bearerHeader string) (string, error) {
	tokenStr, err := parseHeader(bearerHeader)
	if err != nil {
		return "", err
	}
	uid, err := a.GetDataFromToken(tokenStr)
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

func (a *Auth) GetDataFromToken(tokenStr string) (string, error) {
	c := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("invalid signing method %v", token.Header))
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
