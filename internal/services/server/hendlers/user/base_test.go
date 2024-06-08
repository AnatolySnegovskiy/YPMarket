package user

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"market/internal/services/server/middleware"
	"net/http"
)

type testCase struct {
	name           string
	handler        http.HandlerFunc
	method         string
	url            string
	requestBody    interface{}
	expectedStatus int
	queryMock      func(mock sqlmock.Sqlmock)
}

const UserIDContextKey middleware.UserContextKey = "userID"

func upLogin(request *http.Request) *http.Request {
	return request.WithContext(context.WithValue(request.Context(), UserIDContextKey, 123))
}
