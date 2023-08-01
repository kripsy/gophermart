package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/kripsy/gophermart/internal/gophermart/internal/models"
	"github.com/kripsy/gophermart/internal/gophermart/internal/storage"

	"github.com/kripsy/gophermart/internal/gophermart/internal/handler"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"github.com/kripsy/gophermart/internal/gophermart/internal/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Router *chi.Mux
}

func InitServer(ctx context.Context, publicKey string, channelForRequestToAccrual chan models.ResponseOrder) (*Server, error) {
	m := &Server{
		Router: chi.NewRouter(),
	}
	l := logger.LoggerFromContext(ctx)
	h, err := handler.InitHandler(ctx, publicKey)
	if err != nil {
		l.Error("Error in Init server", zap.String("msg", err.Error()))
		return nil, err
	}
	mw := middleware.InitMiddleware(ctx, publicKey)

	m.Router.Use(mw.JWTMiddleware)
	store := storage.GetStorage()
	m.Router.Post("/api/user/orders", h.CreateOrderHandler(store, channelForRequestToAccrual)) // загрузка пользователем номера заказа для расчёта;
	m.Router.Get("/api/user/orders", h.ReadOrdersHandler(store))                               // получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	m.Router.Get("/api/user/balance", h.ReadUserBalanceHandler(store))                         // получение текущего баланса счёта баллов лояльности пользователя;
	m.Router.Post("/api/user/balance/withdraw", h.CreateWithdrawHandler(store))                // запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	m.Router.Get("/api/user/withdrawals", h.ReadWithdrawsHandler(store))                       // получение информации о выводе средств с накопительного счёта пользователем.

	m.Router.HandleFunc("/test", h.TestHandler)

	return m, nil
}
