package middleware

import (
	"net/http"

	"github.com/gorilla/context"
	commonAuth "github.com/kripsy/gophermart/internal/common/auth"
	commonUtils "github.com/kripsy/gophermart/internal/common/utils"
	"github.com/kripsy/gophermart/internal/gophermart/internal/logger"
	"go.uber.org/zap"
)

// this middleware for try get jwt from cookie.
//  1. if URL not protected - generate/update new jwt and set in cookie or pass if jwt is valid
//  2. if URL is protected:
//     2.1. if jwt valid and URL protected - pass
//     2.2. if jwt invalid/empty - generate/update new jwt

func (m *Middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.LoggerFromContext(m.ctx)
		protectedURL := []string{
			"/api/user/orders",
			"/api/user/balance",
			"/api/user/balance/withdraw",
			"/api/user/withdrawals",
		}
		l.Debug("Start JWTMiddleware")

		// check if current URL is protected
		isURLProtected := commonUtils.StingContains(protectedURL, r.URL.Path)
		l.Debug("URL protected value", zap.Bool("msg", isURLProtected))

		// try get token from header

		tokenString, err := commonAuth.GetToken(w, r)

		// if token empty and url is protected -  return 401
		if err != nil && isURLProtected {
			l.Debug("error split bearer token", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := commonAuth.Decrypt(tokenString, m.PublicKey)
		if err != nil && isURLProtected {
			l.Debug("error decrypt token", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims.Username == "" && isURLProtected {
			l.Debug("user in token is empty", zap.String("msg", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		context.Set(r, "username", claims.Username)
		next.ServeHTTP(w, r)
	})
}
