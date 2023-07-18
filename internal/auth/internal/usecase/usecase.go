// This is package of business logic level.
// Here realized logic for register, login user.

package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kripsy/gophermart/internal/auth/internal/config"
	"github.com/kripsy/gophermart/internal/auth/internal/db"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/models"
	"github.com/kripsy/gophermart/internal/auth/internal/utils"

	"go.uber.org/zap"
)

type UseCase struct {
	ctx context.Context
	db  *db.DB
	cfg *config.Config
}

func InitUseCases(ctx context.Context, db *db.DB, cfg *config.Config) (*UseCase, error) {
	uc := &UseCase{
		ctx: ctx,
		db:  db,
		cfg: cfg,
	}
	return uc, nil
}

// RegisterUser get context, username, password and return token, expired time, error.
// At the first step we check is user exists. If exists - return error conflict.
// If user not exists we get new user ID.
// After register new user we generate new jwt token.
func (uc *UseCase) RegisterUser(ctx context.Context, username, password string) (string, time.Time, error) {
	l := logger.LoggerFromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
	defer cancel()
	isUserExists, err := uc.db.IsUserExists(ctx, username)
	if err != nil {
		l.Error("error check isUserExists in RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}

	if isUserExists {
		l.Debug("user already exists")
		userExistsError := models.NewUserExistsError(username)
		return "", time.Time{}, userExistsError
	}
	l.Debug("user not exists, continue get new ID")

	ctx, cancel = context.WithTimeout(ctx, 400*time.Millisecond)
	defer cancel()

	newID, err := uc.db.GetNextUserID(ctx)
	if err != nil {
		l.Error("error get newID in RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	l.Debug("new ID", zap.Int("msg", newID))

	ctx, cancel = context.WithTimeout(ctx, 400*time.Millisecond)
	defer cancel()

	hash, err := utils.GetHash(ctx, password)
	if err != nil {
		l.Error("error GetHash in RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}

	err = uc.db.RegisterUser(ctx, username, hash, newID)
	if err != nil {
		var ue *models.UserExistsError
		if errors.As(err, &ue) {
			l.Debug("register dublicate user", zap.String("msg", username))
			return "", time.Time{}, err
		} else {
			l.Error("error usecase RegisterUser", zap.String("msg", err.Error()))
			return "", time.Time{}, err
		}
	}

	l.Debug("registered user", zap.String("msg", username))

	l.Debug("generate new token", zap.String("msg", username))

	tokenString, expAt, err := utils.BuildJWTString(ctx, newID, username, uc.cfg.PrivateKey, uc.cfg.TokenExp)

	if err != nil {
		l.Error("error BuildJWTString", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	return tokenString, expAt, nil
}

func (uc *UseCase) LoginUser(ctx context.Context, username, password string) (string, time.Time, error) {
	l := logger.LoggerFromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
	defer cancel()

	l.Debug("start LoginUser in UseCase")

	userID, hashPassword, err := uc.db.GetUserHashPassword(ctx, username)
	if err != nil {
		l.Error("error GetUserHashPassword in LoginUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}

	err = utils.IsPasswordCorrect(ctx, []byte(password), []byte(hashPassword))
	if err != nil {
		l.Debug("password incorrect", zap.String("msg", err.Error()))
		return "", time.Time{}, models.NewUserLoginError(username)
	}

	tokenString, expAt, err := utils.BuildJWTString(ctx, userID, username, uc.cfg.PrivateKey, uc.cfg.TokenExp)
	if err != nil {
		l.Error("error BuildJWTString", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	l.Debug("success login user", zap.String("msg", username))
	return tokenString, expAt, nil
}
