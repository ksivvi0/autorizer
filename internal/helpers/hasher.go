package helpers

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func GetHash(data string) (string, error) {
	raw, err := bcrypt.GenerateFromPassword([]byte(data), 11)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func CheckHash(want string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(want))
	return err == nil
}

func EncodeString(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeString(data string) (string, error) {
	rawStr, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(rawStr), nil
}
