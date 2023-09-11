package auth

import (
	"context"
	"log"
	"net/http"
)

type ContextKey string

// Middleware for authentication.
func AuthMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.Proto, r.Method, r.Host, r.RequestURI)

		token := r.URL.Query()["bearer"]

		if len(token) == 1 {
			uid := ValidateToken(token[0])

			if uid == 0 {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			contextKey := ContextKey("userId")

			ctx := context.WithValue(r.Context(), contextKey, uid)
			f(w, r.WithContext(ctx))

		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}
