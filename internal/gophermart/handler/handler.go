package handler

import (
	"context"
	"net/http"

	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
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

func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
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
	w.Write([]byte("Hello world"))
}
