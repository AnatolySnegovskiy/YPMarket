package middleware

import (
	"context"
	"market/internal/system"
	"net/http"
)

const userIDKey string = "user_id"

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userID, err := system.GetUserID(request.Header.Get("Authorization"))
		if err == nil && userID > 0 {
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), userIDKey, userID)))
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
		}
	})
}
