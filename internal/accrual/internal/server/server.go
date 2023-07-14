package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/kripsy/gophermart/internal/accrual/internal/handler"
	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	"go.uber.org/zap"
)

type Server struct {
	Router *chi.Mux
}

func InitServer(ctx context.Context) (*Server, error) {
	m := &Server{
		Router: chi.NewRouter(),
	}
	l := logger.LoggerFromContext(ctx)
	h, err := handler.InitHandler(ctx)
	if err != nil {
		l.Error("Error in Init server", zap.String("msg", err.Error()))
		return nil, err
	}

	m.Router.Get("/api/orders/{number}", h.ReadOrdersHandler) // получение информации о расчёте начислений баллов лояльности;
	m.Router.Post("/api/orders", h.CreateOrderHandler)        // регистрация нового совершённого заказа;
	m.Router.Post("/api/goods", h.CreateGoodsHandler)         // регистрация информации о новой механике вознаграждения за товар.
	m.Router.HandleFunc("/test", h.TestHandler)

	return m, nil
}
