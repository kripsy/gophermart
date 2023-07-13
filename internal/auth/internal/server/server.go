package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kripsy/gophermart/internal/auth/internal/handler"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
  "github.com/kripsy/gophermart/internal/auth/internal/usecase"
  httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "github.com/kripsy/gophermart/docs/auth"
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

	m.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	m.Router.Post("/api/register", h.RegisterUserHandler)

	m.Router.Get("/swagger/*", httpSwagger.WrapHandler)

	return m, nil
}
