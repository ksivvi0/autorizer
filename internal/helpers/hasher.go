package helpers

import "golang.org/x/crypto/bcrypt"

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
