package middleware

import (
	"context"
	"market/internal/system"
	"net/http"
)

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userId, err := system.GetUserID(request.Header.Get("Authorization"))
		if err == nil && userId > 0 {
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), "user_id", userId)))
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
		}
	})
}
