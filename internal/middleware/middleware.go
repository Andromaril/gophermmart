package middleware

import (
	"net/http"

	"github.com/andromaril/gophermmart/internal/verification"
)

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		var err error
		cookie, err := r.Cookie("Token")
		_, err = r.Cookie("Login")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := verification.GetUserID(cookie.Value)
		if token == -1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(ow, r)
	})

}
