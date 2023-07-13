package middleware

import (
	"net/http"

	"github.com/gorilla/context"
)

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		// TODO add AuthMiddleware

		context.Set(r, "userName", "userName")
		//context.Set(r, "userName", nil)
		h.ServeHTTP(ow, r)
	}
}
