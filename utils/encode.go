package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(pwd, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
