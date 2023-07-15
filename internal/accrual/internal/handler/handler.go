package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	"github.com/kripsy/gophermart/internal/accrual/internal/models"
	"github.com/kripsy/gophermart/internal/accrual/internal/storage"
	"github.com/kripsy/gophermart/internal/accrual/internal/utils"
	"go.uber.org/zap"
)

type Handler struct {
	ctx context.Context
}

func InitHandler(ctx context.Context) (*Handler, error) {
	h := &Handler{
		ctx: ctx,
	}
	return h, nil
}

func (h *Handler) CreateOrderHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Info("CreateOrderHandler")

	var req models.RequestOrder
	dec := json.NewDecoder(r.Body)

	//400 — неверный формат запроса;
	if err := dec.Decode(&req); err != nil {
		l.Error("ERROR decode json.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	number, err := strconv.ParseInt(req.Number, 10, 64)
	if err != nil {
		l.Error("ERROR decode json.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	//422 — неверный формат номера заказа;
	if !utils.LuhnValid(number) {
		l.Error("ERROR invalid order number format.")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	getStorage := storage.GetStorage()
	order, err := getStorage.PutOrder(h.ctx, number)

	//500 — внутренняя ошибка сервера.
	if err != nil {
		l.Error("ERROR DB.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	//409 — заказ уже принят в обработку;
	if !order.ProcessedAt.Time.IsZero() {
		l.Error("ERROR the order number has already been uploaded.")
		rw.WriteHeader(http.StatusConflict)
		return
	}

	//202 — новый номер заказа принят в обработку;
	rw.WriteHeader(http.StatusAccepted)

	// TODO Здесь будем передавать в канал объект ордер для горутины которая будет выполнять начисления и менять статус.

}

func (h *Handler) ReadOrdersHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) CreateGoodsHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Debug("TestHandler")
	w.Header().Add("Content-Type", "plain/text")
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		l.Error("Error w.Write([]byte", zap.String("msg", err.Error()))
	}
}
