package services

import (
	"errors"
	"testing"
)

var a *Auth = NewAuthInstance()

func TestAuth_CreateTokenPair(t *testing.T) {
	pair, err := a.CreateTokenPair()
	if err != nil {
		t.Error(err)
	}
	aToken, err := a.GetDataFromToken(pair.AccessToken, false)
	if err != nil {
		t.Error(err)
	}
	rToken, err := a.GetDataFromToken(pair.RefreshToken, true)
	if err != nil {
		t.Error(err)
	}
	if len(aToken) == 0 || len(rToken) == 0 {
		t.Error(errors.New("получен пустой токен"))
	}
}
