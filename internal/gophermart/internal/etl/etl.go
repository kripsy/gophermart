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

	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"github.com/kripsy/gophermart/internal/gophermart/internal/storage"
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
	fmt.Println("restore")
	getStorage := storage.GetStorage()
	newOrders, _ := getStorage.GetNewOrders(ctx)
	processingOrders, _ := getStorage.GetProcessingOrders(ctx)
	// TODO Добавить обработку ошибок

	for _, order := range newOrders {
		ch1 <- order
	}

	for _, order := range processingOrders {
		ch2 <- order
	}
}

func registeringNewOrder(ctx context.Context, ch1 chan models.ResponseOrder, ch2 chan models.ResponseOrder, accrualAddress string) {
	getStorage := storage.GetStorage()
	fmt.Println("registeringNewOrder")
	for {
		order := <-ch1
		u, err := url.Parse(accrualAddress + "/api/orders")
		if err != nil {
			log.Fatal(err)
		}

		body := fmt.Sprintf(`{"order": "%s"}`, order.Number)
		jsonBody := []byte(body)
		r := bytes.NewReader(jsonBody)
		resp, err := http.Post(u.String(), "application/json", r)
		if err != nil {
			log.Fatalln(err)
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
	getStorage := storage.GetStorage()
	fmt.Println("getAccrualForOrder")
	for {
		order := <-ch2
		u, err := url.Parse(fmt.Sprintf(accrualAddress+"/api/orders/%s", order.Number))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(u.String())
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		accrual := &models.ResponseAccrual{}
		err = json.Unmarshal(body, accrual)

		if resp.StatusCode == http.StatusOK && accrual.Status == models.StatusProcessed {
			getStorage.UpdateStatusOrder(ctx, order.Number, models.StatusProcessed, accrual.Accrual)
			//ch2 <- order
			//} else {
			//	ch1 <- order
		}
	}
}
