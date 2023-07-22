package application

import (
	"context"

	"github.com/kripsy/gophermart/internal/gophermart/internal/config"
	"github.com/kripsy/gophermart/internal/gophermart/internal/db"
	"github.com/kripsy/gophermart/internal/gophermart/internal/etl"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/server"
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

	//<<<<<<< HEAD
	_, err = db.InitDB(ctx, cfg.DatabaseAddress, cfg.MigrationsPath)
	if err != nil {
		l.Error("error init DB", zap.String("msg", err.Error()))
		return nil, err
	}

	//srv, err := server.InitServer(ctx)
	//=======
	srv, err := server.InitServer(ctx, cfg.PublicKey)
	//>>>>>>> dev
	if err != nil {
		return nil, err
	}

	etl.InitETL(ctx, cfg.AccrualAddress)

	return &Application{
		appConfig:  cfg,
		appServer:  srv,
		appContext: ctx,
	}, nil
}
