package application

import (
	"context"

	"github.com/kripsy/gophermart/internal/auth/internal/config"
	"github.com/kripsy/gophermart/internal/auth/internal/db"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/server"
	"github.com/kripsy/gophermart/internal/auth/internal/usecase"
	"go.uber.org/zap"
)

type Application struct {
	appConfig  *config.Config
	appServer  *server.Server
	appContext context.Context
}

func (a *Application) GetAppServer() *server.Server {
	return a.appServer
}

func (a *Application) GetAppLogger() *zap.Logger {
	return logger.LoggerFromContext(a.appContext)
}

func (a *Application) GetAppConfig() (string, string, string) {
	return a.appConfig.LoggerLevel, a.appConfig.RunAddress, a.appConfig.DatabaseAddress
}

func NewApp(ctx context.Context) (*Application, error) {
	cfg := config.InitConfig()

	l, err := logger.InitLogger(cfg.LoggerLevel)
	ctx = logger.ContextWithLogger(ctx, l)
	if err != nil {
		return nil, err
	}

	db, err := db.InitDB(ctx, cfg.DatabaseAddress, cfg.MigrationsPath)
	if err != nil {
		l.Error("error init DB", zap.String("msg", err.Error()))
		return nil, err
	}

	uc, err := usecase.InitUseCases(ctx, db, cfg)
	if err != nil {
		l.Error("error init usecase", zap.String("msg", err.Error()))
		return nil, err
	}

	srv, err := server.InitServer(ctx, uc)
	if err != nil {
		return nil, err
	}

	return &Application{
		appConfig:  cfg,
		appServer:  srv,
		appContext: ctx,
	}, nil
}
