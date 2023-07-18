package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	con "github.com/gorilla/context"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"github.com/kripsy/gophermart/internal/gophermart/internal/storage"
	"github.com/kripsy/gophermart/internal/gophermart/internal/utils"
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
	username := con.Get(r, "username")

	//401 — пользователь не аутентифицирован;
	if username == nil {
		l.Error("ERROR User is Unauthorized")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	//400 — неверный формат запроса;
	byteNumber, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("ERROR Can't get value from body.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	number, err := strconv.ParseInt(string(byteNumber), 10, 64)
	if err != nil {
		l.Error("ERROR Can't get value from body.", zap.String("msg", err.Error()))
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
	order, err := getStorage.PutOrder(h.ctx, username, number)

	//500 — внутренняя ошибка сервера.
	if err != nil {
		l.Error("ERROR DB.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	//409 — номер заказа уже был загружен другим пользователем;
	if order.Username != username {
		l.Error("ERROR the order number has already been uploaded by another user.")
		rw.WriteHeader(http.StatusConflict)
		return
	}

	//200 — номер заказа уже был загружен этим пользователем;
	if !order.ProcessedAt.Time.IsZero() {
		l.Error("ERROR the order number has already been uploaded by another user.")
		rw.WriteHeader(http.StatusOK)
		return
	}

	//202 — новый номер заказа принят в обработку;
	rw.WriteHeader(http.StatusAccepted)

	// TODO Здесь я буду передавать в канал объект ордер в горутину которая будет ходить в сервис начислений.
}

func (h *Handler) ReadOrdersHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Info("ReadOrdersHandler")
	username := con.Get(r, "username")
	rw.Header().Set("Content-Type", "application/json")

	//401 — пользователь не аутентифицирован;
	if username == nil {
		l.Error("ERROR User is Unauthorized")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	getStorage := storage.GetStorage()
	orders, err := getStorage.GetOrders(h.ctx, username)

	//500 — внутренняя ошибка сервера.
	if err != nil {
		l.Error("ERROR DB.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 204 — заказ не зарегистрирован в системе расчёта.
	if len(orders) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	rw.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(orders); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		l.Error("ERROR encoding response.", zap.String("msg", err.Error()))
		return
	}
}

func (h *Handler) ReadUserBalanceHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Info("ReadUserBalanceHandler")
	username := con.Get(r, "username")

	//401 — пользователь не аутентифицирован;
	if username == nil {
		l.Error("ERROR User is Unauthorized")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	getStorage := storage.GetStorage()
	balance, err := getStorage.GetBalance(h.ctx, username)

	// 204 — заказ не зарегистрирован в системе расчёта.
	var e *models.ResponseBalanceError
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
	enc := json.NewEncoder(rw)
	if err := enc.Encode(balance); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		l.Error("ERROR encoding response.", zap.String("msg", err.Error()))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateWithdrawHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Info("CreateOrderHandler")
	username := con.Get(r, "username")

	//401 — пользователь не аутентифицирован;
	if username == nil {
		l.Error("ERROR User is Unauthorized")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req models.RequestWithdraw
	dec := json.NewDecoder(r.Body)

	//400 — неверный формат запроса;
	if err := dec.Decode(&req); err != nil {
		l.Error("ERROR Can't decode request JSON body.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	number, err := strconv.ParseInt(req.Number, 10, 64)
	if err != nil {
		l.Error("ERROR Can't get value from body.", zap.String("msg", err.Error()))
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
	err = getStorage.PutWithdraw(h.ctx, username, number, req.Accrual)

	//402 — на счету недостаточно средств;
	var e *models.ResponseBalanceError
	if errors.As(err, &e) {
		l.Error("ERROR there are not enough funds in the account.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusPaymentRequired)
		return
	}

	//500 — внутренняя ошибка сервера.
	if err != nil {
		l.Error("ERROR DB.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	//200 — успешная обработка запроса;
	rw.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) ReadWithdrawsHandler(rw http.ResponseWriter, r *http.Request) {

	l := logger.LoggerFromContext(h.ctx)
	l.Info("ReadWithdrawsHandler")
	username := con.Get(r, "username")

	//401 — пользователь не аутентифицирован;
	if username == nil {
		l.Error("ERROR User is Unauthorized")
		rw.WriteHeader(http.StatusUnauthorized)
	}

	getStorage := storage.GetStorage()
	withdraws, err := getStorage.GetWithdraws(h.ctx, username)

	//204 — нет ни одного списания.
	var e *models.ResponseBalanceError
	if errors.As(err, &e) {
		l.Error("ERROR the order is not registered in the payment system.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusNoContent)
	}

	//500 — внутренняя ошибка сервера.
	if err != nil {
		l.Error("ERROR DB.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(withdraws); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		l.Error("ERROR encoding response.", zap.String("msg", err.Error()))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
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
