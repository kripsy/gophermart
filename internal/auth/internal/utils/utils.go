package utils

import (
	"net/http"
	"time"
)

func AddToken(w http.ResponseWriter, token string, expTime time.Time) error {
	w.Header().Add("Authorization", "Bearer "+token)

	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	}
	http.SetCookie(w, cookie)
	return nil
}
