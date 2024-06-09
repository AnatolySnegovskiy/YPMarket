package services

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"regexp"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestNewOrderAccrual(t *testing.T) {
	_, err := NewOrderAccrual("", "")
	assert.Error(t, err)
}

func TestPollAccrualSystem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"order": "123", "status": "PROCESSED", "accrual": 100}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)
	o := &OrderAccrual{address: server.URL, db: gdb}

	ticker := time.NewTicker(1 * time.Second)

	var dateMock, _ = time.Parse("2006-01-02 15:04:05.000000", "")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE status IN ($1,$2,$3) AND "orders"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}).
			AddRow(123, dateMock, dateMock, nil, "123", "PROCESSED", 1000, 123))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE number = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}).
			AddRow(123, dateMock, dateMock, nil, "123", "PROCESSED", 1000, 123))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"email"=$4,"password"=$5,"balance"=$6,"withdrawal"=$7 WHERE "users"."deleted_at" IS NULL AND "id" = $8`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "balance_history" ("created_at","updated_at","deleted_at","user_id","amount","operation","order_id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second+time.Millisecond*500)
	defer cancel()
	o.PollAccrualSystem(ticker, ctx)
}

func TestPollAccrualSystemError1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`"order": "123", "status": "PROCESSED", "accrual": 100}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)
	o := &OrderAccrual{address: server.URL, db: gdb}

	ticker := time.NewTicker(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second+time.Millisecond*500)
	defer cancel()
	o.PollAccrualSystem(ticker, ctx)
}

func TestPollAccrualSystemError2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"order": "123", "status": "PROCESSED", "accrual": 100}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)
	o := &OrderAccrual{address: server.URL, db: gdb}

	ticker := time.NewTicker(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second+time.Millisecond*500)
	defer cancel()
	o.PollAccrualSystem(ticker, ctx)
}

func TestFetchOrderAccrual(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"order": "123", "status": "PROCESSED", "accrual": 100}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/api/orders/333":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`"order2": "123", "status": "PROCESSED", "accrual": 100}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/api/orders/456":
			w.WriteHeader(http.StatusNoContent)
		case "/api/orders/789":
			w.WriteHeader(http.StatusTooManyRequests)
		case "/api/orders/999":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	o := &OrderAccrual{address: server.URL}

	testCases := []struct {
		orderNumber      string
		expectedResponse *OrderAccrualResponse
		expectedError    string
	}{
		{"123", &OrderAccrualResponse{"123", "PROCESSED", 100}, ""},
		{"333", nil, "invalid character ':' after top-level value"},
		{"456", nil, "the order is not registered in the payment system"},
		{"789", nil, "the number of requests to the service has been exceeded\n"},
		{"999", nil, "internal server error"},
		{"000", nil, "unknown response code: 404"},
	}

	for _, tc := range testCases {
		resp, err := o.fetchOrderAccrual(tc.orderNumber)
		assert.Equal(t, tc.expectedResponse, resp)
		if err != nil {
			assert.EqualError(t, err, tc.expectedError)
		}
	}
}
