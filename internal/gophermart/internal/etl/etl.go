package etl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"github.com/kripsy/gophermart/internal/gophermart/internal/storage"
	"go.uber.org/zap"
)

func InitETL(ctx context.Context, accrualAddress string, channelForRequestToAccrual chan models.ResponseOrder, channelForResponseFromAccrual chan models.ResponseOrder) {
	go restore(ctx, channelForRequestToAccrual, channelForResponseFromAccrual, accrualAddress)
	go registeringNewOrder(ctx, channelForRequestToAccrual, channelForResponseFromAccrual, accrualAddress)
	go getAndStoreAccrualForOrder(ctx, channelForResponseFromAccrual, accrualAddress)
}

func restore(ctx context.Context, channelForRequestToAccrual chan models.ResponseOrder, channelForResponseFromAccrual chan models.ResponseOrder, accrualAddress string) {
	l := logger.LoggerFromContext(ctx)
	l.Info("restore")
	getStorage := storage.GetStorage()
	newOrders, err := getStorage.GetNewOrders(ctx)
	if err != nil {
		l.Error("ERROR Can't get new Orders.", zap.String("msg", err.Error()))
	}
	processingOrders, err := getStorage.GetProcessingOrders(ctx)
	if err != nil {
		l.Error("ERROR Can't get processing Orders.", zap.String("msg", err.Error()))
	}

	for _, order := range newOrders {
		channelForRequestToAccrual <- order
	}

	for _, order := range processingOrders {
		channelForResponseFromAccrual <- order
	}
}

func registeringNewOrder(ctx context.Context, channelForRequestToAccrual chan models.ResponseOrder, channelForResponseFromAccrual chan models.ResponseOrder, accrualAddress string) {
	l := logger.LoggerFromContext(ctx)
	l.Info("registeringNewOrder")
	getStorage := storage.GetStorage()
	for {
		l.Info("Waiting from the channelForRequestToAccrual")
		order := <-channelForRequestToAccrual
		l.Info("Received from the channelForRequestToAccrual")

		u, err := url.Parse(accrualAddress + "/api/orders")
		if err != nil {
			l.Error("ERROR Can't parse url.", zap.String("msg", err.Error()))
		}

		body := fmt.Sprintf(`{"order": "%s"}`, order.Number)
		jsonBody := []byte(body)
		r := bytes.NewReader(jsonBody)
		resp, err := http.Post(u.String(), "application/json", r)
		if err != nil {
			channelForRequestToAccrual <- order
			l.Info("ERROR Can't get accrual.", zap.String("msg", err.Error()))
		}

		if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusConflict {
			order, _ := getStorage.UpdateStatusOrder(ctx, order.Number, models.StatusProcessing, 0)
			l.Info("Trying to send to the channelForResponseFromAccrual")
			channelForResponseFromAccrual <- order
			l.Info("Sent to the channelForResponseFromAccrual")
		} else {
			l.Info("Trying to send to the channelForRequestToAccrual")
			channelForRequestToAccrual <- order
			l.Info("Sent to the channelForRequestToAccrual")

		}
		err = resp.Body.Close()
		if err != nil {
			l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
		}
	}
}

func getAndStoreAccrualForOrder(ctx context.Context, channelForResponseFromAccrual chan models.ResponseOrder, accrualAddress string) {
	l := logger.LoggerFromContext(ctx)
	l.Info("getAndStoreAccrualForOrder")
	getStorage := storage.GetStorage()
	for {
		l.Info("Waiting from the channelForResponseFromAccrual")
		order := <-channelForResponseFromAccrual
		l.Info("Received from the channelForResponseFromAccrual")
		u, err := url.Parse(fmt.Sprintf(accrualAddress+"/api/orders/%s", order.Number))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(u.String())
		if err != nil {
			l.Error("ERROR Can't get accrual.", zap.String("msg", err.Error()))
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			l.Error("ERROR Can't get body from request.", zap.String("msg", err.Error()))
		}

		accrual := &models.ResponseAccrual{}
		err = json.Unmarshal(body, accrual)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode == http.StatusOK && accrual.Status == models.StatusProcessed {
			_, err := getStorage.UpdateStatusOrder(ctx, order.Number, models.StatusProcessed, accrual.Accrual)
			if err != nil {
				return
			}
		}
		err = resp.Body.Close()
		if err != nil {
			l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
		}
	}
}
