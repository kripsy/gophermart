package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/kripsy/gophermart/internal/auth/internal/logger"
)

func RegisterUser(ctx context.Context, username, password string) (string, time.Time, error) {
	l := logger.LoggerFromContext(ctx)
	l.Debug("usecase RegisterUser")
	token := ""
	expTime := time.Now().Add(100500 * time.Hour)

	return token, expTime, fmt.Errorf("not implemented yet")
}
