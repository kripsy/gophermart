package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kripsy/gophermart/internal/accrual/application"
	"go.uber.org/zap"
)

func main() {

	application, err := application.NewApp(context.Background())

	if err != nil {
		fmt.Println("Error init application: ", err.Error())
		os.Exit(1)
	}
	logger := application.GetAppLogger()
	defer logger.Sync()
	loggerLevel, runAddress, dbURI := application.GetAppConfig()
	logger.Info("LOGGER_LEVEL", zap.String("msg", loggerLevel))
	logger.Info("RUN_ADDRESS", zap.String("msg", runAddress))
	logger.Info("DATABASE_URI", zap.String("msg", dbURI))
	err = http.ListenAndServe(runAddress, application.GetAppServer().Router)
	if err != nil {
		logger.Error("Error ListenAndServe", zap.String("msg", err.Error()))
		os.Exit(1)
	}
}
