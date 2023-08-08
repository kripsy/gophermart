package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
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

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return publicKeyPEM, nil
	})

	if err != nil || !token.Valid {
		return Claims{}, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}

// GetTokenFromBearer returns token from header.
// Header should start from "Baerer ", otherwise return empty string and error.
func GetTokenFromBearer(bearerString string) (string, error) {
	splitString := strings.Split(bearerString, "Bearer ")
	fmt.Printf("len splitString %d\n", len(splitString))
	if len(splitString) < 2 {
		fmt.Printf("bearer string not valid")
		return "", fmt.Errorf("bearer string not valid")
	}
	tokenString := splitString[1]
	if tokenString == "" {
		fmt.Printf("tokenString is empty")
		return "", fmt.Errorf("tokenString is empty")
	}
	return tokenString, nil
}

func GetToken(w http.ResponseWriter, r *http.Request) (string, error) {
	var token string

	tokenString := r.Header.Get("Authorization")
	if tokenString != "" {
		fmt.Printf("get token from header: %s\n", tokenString)
		token, _ = GetTokenFromBearer(tokenString)
		fmt.Printf("token %s\n", token)
	}
	if token != "" {
		return token, nil
	}

	// if we continue - it means that in header isn't token. Try find it in cookie
	cookieToken, err := r.Cookie("token")
	if err != nil {
		return "", errors.Wrap(err, "cannot get token from cookie")
	}
	token = cookieToken.Value
	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("token not found")

}
