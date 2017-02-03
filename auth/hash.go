package auth

import (
	"crypto/base64"
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func encrypt(str string) (string, error) {
	passwd := []byte(str)
	hash, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswd(str, hash string) error {
	return bcrypt.CompareHasAndPassword([]byte(hash), []byte(str))
}

func genHash(id string) string {
	idBytes := ([]byte)(id)
	idSha := sha256.New()
	idSha.Write(idBytes)
	shaStr := base64.URLEncoding.EncodeToString(idSha.Sum(nil))
	return shaStr
}
