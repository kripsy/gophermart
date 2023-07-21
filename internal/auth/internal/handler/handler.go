package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	models "github.com/kripsy/gophermart/internal/auth/internal/models"
	"github.com/kripsy/gophermart/internal/auth/internal/usecase"
	"github.com/kripsy/gophermart/internal/auth/internal/utils"

	"go.uber.org/zap"
)

type Handler struct {
	ctx context.Context
	uc  *usecase.UseCase
}

func InitHandler(ctx context.Context, uc *usecase.UseCase) (*Handler, error) {
	h := &Handler{
		ctx: ctx,
		uc:  uc,
	}
	return h, nil
}

func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerFromContext(h.ctx)
	l.Debug("TestHandler")
	w.Header().Add("Content-Type", "plain/text")
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		l.Error("Error w.Write([]byte", zap.String("msg", err.Error()))
	}
}

// ShowAccount godoc
// @Summary      Register
// @Description  Register new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user   body      models.User  true  "User register data"
// @Success      200
// @Failure      400
// @Failure      409
// @Failure      500
// @Router       /api/register [post]
// RegisterUserHandler accepts a username and password in json format.
// If we have success register new user, we insert token into cookie `token` and header `Authorization`.
func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	l := logger.LoggerFromContext(h.ctx)
	isUniqueError := false
	l.Debug("RegisterUserHandler")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("error read from body", zap.String("msg", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		l.Debug("error close body", zap.String("msg", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user, err := models.InitNewUser(body)

	if err != nil {
		l.Debug("error init model of user from request", zap.String("msg", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Debug("init new user from body", zap.String("msg", user.Username))

	token, expTime, err := h.uc.RegisterUser(h.ctx, user.Username, user.Password)
	if err != nil {
		var ue *models.UserExistsError
		if errors.As(err, &ue) {
			isUniqueError = true
		} else {
			l.Error("error register user", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if isUniqueError {
		w.WriteHeader(http.StatusConflict)
	} else {
		err := utils.AddToken(w, token, expTime)
		if err != nil {
			l.Error("error AddToken in register user", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// ShowAccount godoc
// @Summary      Login
// @Description  Login as user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user   body      models.User  true  "User login data"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      500
// @Router       /api/login [post]
// LoginUserHandler accepts a username and password in json format.
// If we have success  user login, we insert token into cookie `token` and header `Authorization`.
func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	l := logger.LoggerFromContext(h.ctx)

	l.Debug("LoginUserHandler")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("error read from body", zap.String("msg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		l.Debug("error close body", zap.String("msg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := models.InitNewUser(body)
	if err != nil {
		l.Debug("error init model of user from request", zap.String("msg", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Debug("init new user from body in LoginUserHandler", zap.String("msg", user.Username))

	token, expTime, err := h.uc.LoginUser(h.ctx, user.Username, user.Password)

	if err != nil {
		var userLoginError *models.UserLoginError
		if errors.As(err, &userLoginError) {
			l.Error("error login user", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		l.Error("error login user", zap.String("msg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = utils.AddToken(w, token, expTime)
	if err != nil {
		l.Error("error AddToken in register user", zap.String("msg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
