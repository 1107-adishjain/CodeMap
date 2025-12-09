package middlewares

import (
	"codemap/backend/internal/helper"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "user_ID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		claims, err := helper.VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// Set user ID in context for downstream handlers
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
