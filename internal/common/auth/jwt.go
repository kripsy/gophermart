package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   int
	Username string
}

func Decrypt(tokenString, publicKey string) (Claims, error) {
	publicKeyPEM, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	claims := Claims{}
	if err != nil {
		return Claims{}, fmt.Errorf("validate: parse key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return publicKeyPEM, nil
	})

	if err != nil || !token.Valid {
		return Claims{}, fmt.Errorf("validate: invalid")
	}
	return claims, nil
}
