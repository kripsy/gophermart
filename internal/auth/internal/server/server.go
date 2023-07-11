package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/kripsy/gophermart/internal/auth/internal/handler"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/usecase"
	"go.uber.org/zap"
)

type Server struct {
	Router *chi.Mux
}

func InitServer(ctx context.Context, uc *usecase.UseCase) (*Server, error) {
	m := &Server{
		Router: chi.NewRouter(),
	}
	l := logger.LoggerFromContext(ctx)
	h, err := handler.InitHandler(ctx, uc)
	if err != nil {
		l.Error("Error in Init server", zap.String("msg", err.Error()))
		return nil, err
	}

	m.Router.HandleFunc("/test", h.TestHandler)
	m.Router.Post("/api/register", h.RegisterUserHandler)

	return m, nil
}
