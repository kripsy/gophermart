package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) CreateOrderHandler(store storage.Store) http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
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

		order, err := store.PutOrder(h.ctx, number)

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
	return http.HandlerFunc(fn)
}

func (h *Handler) ReadOrdersHandler(store storage.Store) http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
		l := logger.LoggerFromContext(h.ctx)
		l.Info("ReadOrdersHandler")

		urlNumber := chi.URLParam(r, "number")

		//400 — неверный формат запроса;
		if urlNumber == "" {
			l.Error("ERROR number is empty.")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		number, err := strconv.ParseInt(urlNumber, 10, 64)
		if err != nil {
			l.Error("ERROR invalid order number format.", zap.String("msg", err.Error()))
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		//422 — неверный формат номера заказа;
		if !utils.LuhnValid(number) {
			l.Error("ERROR invalid order number format.")
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		// TODO посмотри мое разделение server/handler/usecase/database.
		// В твоем случае у тебя раздутые handler, сложно читать.
		// Да и за бизнес-логику (неверный формат номера) лучше пусть отвечает слой usecase.
		order, err := store.GetOrder(h.ctx, number)

		// 204 — заказ не зарегистрирован в системе расчёта.
		var e *models.AccrualError
		if errors.As(err, &e) {
			l.Error("ERROR the order is not registered in the payment system.", zap.String("msg", err.Error()))
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		//500 — внутренняя ошибка сервера.
		if err != nil {
			l.Error("ERROR DB.", zap.String("msg", err.Error()))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := models.ResponseOrder{}

		resp.Number = strconv.FormatInt(order.Number, 10)
		resp.Status = order.Status
		resp.Accrual = order.Accrual

		enc := json.NewEncoder(rw)
		rw.Header().Set("Content-Type", "application/json")
		if err := enc.Encode(resp); err != nil {
			l.Error("error encoding response", zap.String("msg", err.Error()))
			return
		}
	}
	return http.HandlerFunc(fn)
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
