package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

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

func IsPasswordCorrect(ctx context.Context, password, hashPassowrd []byte) (bool, error) {
	l := logger.LoggerFromContext(ctx)

	err := bcrypt.CompareHashAndPassword(hashPassowrd, password)

	if err != nil {
		l.Error("error compare password and hash", zap.String("msg", err.Error()))
		return false, err
	}
	return true, nil
}
