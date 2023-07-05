package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/kripsy/gophermart/internal/auth/handler"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
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

	m.Router.HandleFunc("/test", h.TestHandler)

	return m, nil
}
