package utils

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   int
	Username string
}

func BuildJWTString(ctx context.Context, userID int, username, secretKey string, expTime time.Duration) (string, time.Time, error) {
	l := logger.LoggerFromContext(ctx)
	expAt := time.Now().Add(expTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
		UserID:   userID,
		Username: username,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		l.Error("failed in BuildJWTString", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	return tokenString, expAt, nil
}
