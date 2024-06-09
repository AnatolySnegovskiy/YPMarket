package user

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"io"
	"market/internal/services/server/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

var dateMock, _ = time.Parse("2006-01-02 15:04:05.000000", "")

func upLogin(request *http.Request) *http.Request {
	return request.WithContext(context.WithValue(request.Context(), UserIDContextKey, 123))
}

func RunTestCases(t *testing.T, testCases []testCase, mock sqlmock.Sqlmock) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			if tc.requestBody != nil {
				body, _ := json.Marshal(tc.requestBody)
				req.Body = io.NopCloser(bytes.NewReader(body))
			}

			rr := httptest.NewRecorder()
			tc.queryMock(mock)
			tc.handler.ServeHTTP(rr, req)
			if assert.Equal(t, tc.expectedStatus, rr.Code) == false {
				t.Errorf(rr.Body.String())
			}
		})
	}
}
