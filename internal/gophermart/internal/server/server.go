package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kripsy/gophermart/cmd/gophermart/middleware"
	"github.com/kripsy/gophermart/internal/gophermart/internal/handler"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"go.uber.org/zap"
)

type Server struct {
	Router *chi.Mux
}

func middlewares(h http.HandlerFunc) http.HandlerFunc {
	return middleware.AuthMiddleware(h)
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

	m.Router.Post("/api/user/orders", middlewares(h.CreateOrderHandler)) // загрузка пользователем номера заказа для расчёта;
	m.Router.Get("/api/user/orders", middlewares(h.ReadOrdersHandler))   // получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	m.Router.Get("/api/user/balance", h.ReadUserBalanceHandler)          // получение текущего баланса счёта баллов лояльности пользователя;
	m.Router.Post("/api/user/balance/withdraw", h.CreateWithdrawHandler) // запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	m.Router.Get("/api/user/withdrawals", h.ReadWithdrawsTestHandler)    // получение информации о выводе средств с накопительного счёта пользователем.
	m.Router.HandleFunc("/test", h.TestHandler)

	return m, nil
}
