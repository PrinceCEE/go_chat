package utils

import (
	"encoding/json"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(pwd, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}

func ReadJSON(r io.Reader, dst any) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(dst)
}
