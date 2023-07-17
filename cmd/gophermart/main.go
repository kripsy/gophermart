package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kripsy/gophermart/internal/gophermart/application"
	"go.uber.org/zap"
)

func main() {

	application, err := application.NewApp(context.Background())
	//etlApplication, err := etl.NewEtlApp(context.Background())
	//
	//go restore()

	if err != nil {
		fmt.Println("Error init application: ", err.Error())
		os.Exit(1)
	}
	logger := application.GetAppLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Error logger.Sync()", zap.String("msg", err.Error()))
		}
	}(logger)
	loggerLevel, runAddress, dbURI := application.GetAppConfig()
	logger.Info("LOGGER_LEVEL", zap.String("msg", loggerLevel))
	logger.Info("RUN_ADDRESS", zap.String("msg", runAddress))
	logger.Info("DATABASE_URI", zap.String("msg", dbURI))

	//go restore()

	err = http.ListenAndServe(runAddress, application.GetAppServer().Router)
	if err != nil {
		logger.Error("Error ListenAndServe", zap.String("msg", err.Error()))
		os.Exit(1)
	}
}

//func restore() {
//	ctx := context.Background()
//	fmt.Println("restore")
//	getStorage := storage.GetStorage()
//	NewOrders, err := getStorage.GetNewOrders(ctx)
//	fmt.Println(NewOrders, err)
//
//	ProcessingOrders, err := getStorage.GetProcessingOrders(ctx)
//	fmt.Println(ProcessingOrders, err)
//
//}
