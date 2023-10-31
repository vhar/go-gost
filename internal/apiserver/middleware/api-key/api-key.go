package apikey

import (
	"go-gost/internal/apiserver/response"
	"net/http"
)

func Authorize(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			authorize := r.Header.Get("Authorization")

			if key != "" {
				if authorize == "" || authorize != key {
					response.ErrorResponse(w, "Authorization required", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
