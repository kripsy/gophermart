// This is package of business logic level.
// Here realized logic for register, login user.

package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/kripsy/gophermart/internal/auth/internal/db"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"go.uber.org/zap"
)

type UseCase struct {
	ctx context.Context
	db  *db.DB
}

func InitUseCases(ctx context.Context, db *db.DB) (*UseCase, error) {
	uc := &UseCase{
		ctx: ctx,
		db:  db,
	}
	return uc, nil
}

func (uc *UseCase) RegisterUser(ctx context.Context, username, password string) (string, time.Time, error) {
	l := logger.LoggerFromContext(ctx)
	isUserExists, err := uc.db.IsUserExists(ctx, username)
	if err != nil {
		l.Error("error check isUserExists in RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}

	if isUserExists {
		l.Debug("user already exists")
		return "", time.Time{}, fmt.Errorf("user already exists")
	}

	token, expTime, err := uc.db.RegisterUser(ctx, username, password)
	if err != nil {
		l.Error("error usecase RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}

	return "", time.Time{}, nil
}
