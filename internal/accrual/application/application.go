package application

import (
	"context"

	"github.com/kripsy/gophermart/internal/accrual/internal/config"
	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	"github.com/kripsy/gophermart/internal/accrual/internal/server"
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

	srv, err := server.InitServer(ctx)
	if err != nil {
		return nil, err
	}

	return &Application{
		appConfig:  cfg,
		appServer:  srv,
		appContext: ctx,
	}, nil
}
