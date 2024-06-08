package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestUserHandlers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	testCases := []testCase{
		{
			name: "Login - StatusUnauthorized",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				LoginHandler(gdb, writer, request)
			}),
			method:         "GET",
			url:            "/login",
			requestBody:    LoginRequest{Login: "login", Password: "password"},
			expectedStatus: http.StatusUnauthorized,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "Login - StatusBadRequest no login",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				LoginHandler(gdb, writer, request)
			}),
			method:         "GET",
			url:            "/login",
			requestBody:    LoginRequest{Login: "", Password: "123"},
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "Login - StatusBadRequest no password",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				LoginHandler(gdb, writer, request)
			}),
			method:         "GET",
			url:            "/login",
			requestBody:    LoginRequest{Login: "aaa", Password: ""},
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "Login - StatusBadRequest bad request",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				LoginHandler(gdb, writer, request)
			}),
			method:         "GET",
			url:            "/login",
			requestBody:    `{"gas": "aaa", "pas": ""}`,
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "Login - StatusBadRequest bad request",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				LoginHandler(gdb, writer, request)
			}),
			method:         "GET",
			url:            "/login",
			requestBody:    LoginRequest{Login: "login", Password: "TE7AcasnAMewDfIjqVJJX"},
			expectedStatus: http.StatusOK,
			queryMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
						AddRow(1, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))
			},
		},
	}

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

type mockResponseWriter struct {
	writeError error
}

func (m *mockResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (m *mockResponseWriter) Write(p []byte) (int, error) {
	return 0, m.writeError
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
}

func TestLoginHandler_WriteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)
	rw := &mockResponseWriter{writeError: errors.New("mock write error")}

	req := httptest.NewRequest("POST", "http://example.com/login", strings.NewReader(`{"login": "test", "password": "password"}`))
	ctx := context.Background()
	req = req.WithContext(ctx)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(1, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))

	LoginHandler(gdb, rw, req)
	
	assert.NotNil(t, rw.writeError, "Expected write error to be not nil")
	assert.Equal(t, "mock write error", rw.writeError.Error(), "Expected write error message to be 'mock write error'")
}
