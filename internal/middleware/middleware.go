package middleware

import (
	"context"
	"net/http"

	"github.com/andromaril/gophermmart/internal/verification"
)

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		cookie, err := r.Cookie("Token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			//wef := fmt.Sprintf("%v", err)
			//w.Write([]byte(wef))
			return
		}
		token := verification.GetUserId(cookie.Value)
		if token == -1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(ow, r)
	})

}

func AuthMiddlewareContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		login := ow.Header().Get("Authorization")
		//r.Header.Set("Authorization", login)
		//ow.Write([]byte(login))
		ctx := context.WithValue(r.Context(), "Authorization", login)
		h.ServeHTTP(ow, r.WithContext(ctx))
		//h.ServeHTTP(ow, r)
	})

}
