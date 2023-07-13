package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kripsy/gophermart/internal/auth/application"
	"go.uber.org/zap"
)

// @title Swagger API Gophermart
// @version 1.0
// @description This is a swagger server for Gophermart Auth server.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /

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
