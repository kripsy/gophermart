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

var ch1 chan models.ResponseOrder
var ch2 chan models.ResponseOrder

func GetChan() chan models.ResponseOrder {
	return ch2
}

func InitETL(ctx context.Context, accrualAddress string) {
	go restore(ctx, ch1, ch2, accrualAddress)
	go registeringNewOrder(ctx, ch1, ch2, accrualAddress)
	go getAndStoreAccrualForOrder(ctx, ch2, accrualAddress)
}

func restore(ctx context.Context, ch1 chan models.ResponseOrder, ch2 chan models.ResponseOrder, accrualAddress string) {
	l := logger.LoggerFromContext(ctx)
	l.Info("restore")
	getStorage := storage.GetStorage()
	newOrders, err := getStorage.GetNewOrders(ctx)
	if err != nil {
		l.Error("ERROR Can't get new Orders.", zap.String("msg", err.Error()))
	}
	processingOrders, err := getStorage.GetProcessingOrders(ctx)
	// TODO Добавить обработку ошибок
	if err != nil {
		l.Error("ERROR Can't get processing Orders.", zap.String("msg", err.Error()))
	}

	for _, order := range newOrders {
		ch1 <- order
	}

	for _, order := range processingOrders {
		ch2 <- order
	}
}

func registeringNewOrder(ctx context.Context, ch1 chan models.ResponseOrder, ch2 chan models.ResponseOrder, accrualAddress string) {
	l := logger.LoggerFromContext(ctx)
	l.Info("registeringNewOrder")
	getStorage := storage.GetStorage()
	for {
		order := <-ch1
		u, err := url.Parse(accrualAddress + "/api/orders")
		if err != nil {
			l.Error("ERROR Can't parse url.", zap.String("msg", err.Error()))
		}

		body := fmt.Sprintf(`{"order": "%s"}`, order.Number)
		jsonBody := []byte(body)
		r := bytes.NewReader(jsonBody)
		resp, err := http.Post(u.String(), "application/json", r)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		}(resp.Body)
		if err != nil {
			l.Error("ERROR Can't get accrual.", zap.String("msg", err.Error()))
		}

		if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusConflict {
			order, _ := getStorage.UpdateStatusOrder(ctx, order.Number, models.StatusProcessing, 0)
			ch2 <- order
		} else {
			ch1 <- order
		}
	}
}

func getAndStoreAccrualForOrder(ctx context.Context, ch2 chan models.ResponseOrder, accrualAddress string) {
	l := logger.LoggerFromContext(ctx)
	l.Info("getAndStoreAccrualForOrder")
	getStorage := storage.GetStorage()
	for {
		order := <-ch2
		u, err := url.Parse(fmt.Sprintf(accrualAddress+"/api/orders/%s", order.Number))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(u.String())
		if err != nil {
			l.Error("ERROR Can't get accrual.", zap.String("msg", err.Error()))
		}

		body, err := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				l.Error("ERROR Can't close body.", zap.String("msg", err.Error()))
			}
		}(resp.Body)

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
			//ch2 <- order
			//} else {
			//	ch1 <- order
		}
	}
}
