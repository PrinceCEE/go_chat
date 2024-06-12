package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(t *TokenClaims) (string, error) {
	key := []byte(os.Getenv("JWT_KEY"))
	t.IssuedAt = jwt.NewNumericDate(time.Now())
	t.Issuer = "escrowAPI"
	t.Subject = t.UserID
	t.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	return token.SignedString(key)
}

func VerifyToken(tokenStr string) (*TokenClaims, error) {
	key := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*TokenClaims), nil
}
