package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"

	con "github.com/gorilla/context"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
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

	number, err := strconv.Atoi(string(byteNumber))
	if err != nil {
		l.Error("ERROR Can't get value from body.", zap.String("msg", err.Error()))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	//422 — неверный формат номера заказа;
	if !utils.LuhnValid(number) {
		l.Debug("ERROR invalid order number format.")
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
	if order.UserName != username {
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

func (h *Handler) ReadOrdersHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ReadUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) CreateWithdrawHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ReadWithdrawsTestHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Debug("TestHandler")
	w.Header().Add("Content-Type", "plain/text")
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		return
	}
}
