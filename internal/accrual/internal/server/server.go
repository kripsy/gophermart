package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/kripsy/gophermart/internal/accrual/internal/handler"
	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	"github.com/kripsy/gophermart/internal/accrual/internal/storage"
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
	store := storage.GetStorage()
	m.Router.Get("/api/orders/{number}", h.ReadOrdersHandler(store)) // получение информации о расчёте начислений баллов лояльности;
	m.Router.Post("/api/orders", h.CreateOrderHandler(store))        // регистрация нового совершённого заказа;
	m.Router.HandleFunc("/test", h.TestHandler)

	return m, nil
}
