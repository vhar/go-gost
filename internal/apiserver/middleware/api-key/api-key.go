package apikey

import (
	"go-gost/internal/apiserver/response"
	"net/http"
)

func Authorize(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			authorize := r.Header.Get("Authorization")

			if key != "" {
				if authorize == "" || authorize != key {
					response.ErrorResponse(w, "Authorization required", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
