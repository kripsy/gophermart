// This is package of business logic level.
// Here realized logic for register, login user.

package usecase

import (
	"context"
	"time"

	"github.com/kripsy/gophermart/internal/auth/internal/db"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/models"
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
		userExistsError := models.NewUserExistsError(username, err)
		l.Debug("user already exists")
		return "", time.Time{}, userExistsError
	}

	token, expTime, err := uc.db.RegisterUser(ctx, username, password)
	if err != nil {
		l.Error("error usecase RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}

	return "", time.Time{}, nil
}
