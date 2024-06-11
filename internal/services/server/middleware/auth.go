package middleware

import (
	"context"
	"market/internal/system"
	"net/http"
)

type UserContextKey string

const UserIDContextKey UserContextKey = "userID"

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userID, err := system.GetUserID(request.Header.Get("Authorization"))
		if err == nil && userID > 0 {
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), UserIDContextKey, userID)))
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
		}
	})
}
