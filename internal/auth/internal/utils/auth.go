package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	common "github.com/kripsy/gophermart/internal/common/auth"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func BuildJWTString(ctx context.Context, userID int, username, privateKey string, expTime time.Duration) (string, time.Time, error) {
	privateKeyPEM, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	l := logger.LoggerFromContext(ctx)
	expAt := time.Now().Add(expTime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, common.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
		UserID:   userID,
		Username: username,
	})

	tokenString, err := token.SignedString(privateKeyPEM)
	if err != nil {
		l.Error("failed in BuildJWTString", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	return tokenString, expAt, nil
}

func AddToken(w http.ResponseWriter, token string, expTime time.Time) error {
	w.Header().Add("Authorization", "Bearer "+token)

	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	}
	http.SetCookie(w, cookie)
	return nil
}

func GetHash(ctx context.Context, password string) (string, error) {
	l := logger.LoggerFromContext(ctx)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		l.Error("error GetHash", zap.String("msg", err.Error()))
		return "", err
	}
	return string(bytes), nil

}

func IsPasswordCorrect(ctx context.Context, password, hashPassowrd []byte) error {
	l := logger.LoggerFromContext(ctx)

	err := bcrypt.CompareHashAndPassword(hashPassowrd, password)

	if err != nil {
		l.Error("error compare password and hash", zap.String("msg", err.Error()))
		return err
	}
	return nil
}
