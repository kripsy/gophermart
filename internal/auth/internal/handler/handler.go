package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/models"
	"github.com/kripsy/gophermart/internal/auth/internal/usecase"
	"github.com/kripsy/gophermart/internal/auth/internal/utils"
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

func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Debug("TestHandler")
	w.Header().Add("Content-Type", "plain/text")
	w.Write([]byte("Hello world"))
}

// RegisterUserHandler accepts a username and password in json format.
// If we have success register new user, we insert token into cookie `token` and header `Authorization`.
func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Debug("RegisterUserHandler")
	body, err := io.ReadAll(r.Body)
	err = r.Body.Close()
	if err != nil {
		l.Debug("error close body", zap.String("msg", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user, err := models.InitNewUser(body)
	if err != nil {
		l.Debug("error init model of user from request", zap.String("msg", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	l.Debug("init new user from body", zap.String("msg", user.Username))

	token, expTime, err := usecase.RegisterUser(h.ctx, user.Username, user.Password)
	if err != nil {
		l.Debug("error register user", zap.String("msg", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	utils.AddToken(w, token, expTime)
	w.Write([]byte("Hello world"))
}
