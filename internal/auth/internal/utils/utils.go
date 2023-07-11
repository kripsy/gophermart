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

func GetHash(ctx context.Context, password string) ([]byte, error) {
	l := logger.LoggerFromContext(ctx)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		l.Error("error GetHash", zap.String("msg", err.Error()))
		return nil, err
	}
	return bytes, nil

}
