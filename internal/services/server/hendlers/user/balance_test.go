package user

import (
	"bytes"
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
	"testing"
	"time"
)

func TestBalanceHandlers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	testCases := []testCase{
		{
			name: "GetBalanceHandler - Success",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				GetBalanceHandler(gdb, writer, upLogin(request))
			}),
			method:         "GET",
			url:            "/balance",
			requestBody:    nil,
			expectedStatus: http.StatusOK,
			queryMock:      GetBalanceHandlerMockQuery,
		},
		{
			name: "WithdrawHandler - Success",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    WithdrawRequest{Order: "order123", Sum: 100},
			expectedStatus: http.StatusOK,
			queryMock:      WithdrawHandlerMockQuery,
		},
		{
			name: "WithdrawHandler - Success no order",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    WithdrawRequest{Order: "60480142", Sum: 100},
			expectedStatus: http.StatusOK,
			queryMock:      WithdrawHandlerMockQueryNoOrder,
		},
		{
			name: "WithdrawHandler - error no order",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    WithdrawRequest{Order: "order123", Sum: 100},
			expectedStatus: http.StatusInternalServerError,
			queryMock:      WithdrawHandlerMockQueryNoOrderError,
		},
		{
			name: "WithdrawHandler - bad request",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    `{"orderws": "order123", "suma": 100}`,
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "WithdrawHandler - not enough money",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    WithdrawRequest{Order: "order123", Sum: 100},
			expectedStatus: http.StatusPaymentRequired,
			queryMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
						AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 0, 0))
			},
		},
		{
			name: "WithdrawHandler - invalid sum",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    WithdrawRequest{Order: "order123", Sum: 0},
			expectedStatus: http.StatusUnprocessableEntity,
			queryMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
						AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 0, 0))
			},
		},
		{
			name: "WithdrawHandler - Over Error",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WithdrawHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/withdraw",
			requestBody:    WithdrawRequest{Order: "order123", Sum: 100},
			expectedStatus: http.StatusInternalServerError,
			queryMock:      WithdrawHandlerMockQueryError,
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

var dateMock, _ = time.Parse("2006-01-02 15:04:05.000000", "")

func GetBalanceHandlerMockQuery(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 5000000, 0))
}

func WithdrawHandlerMockQuery(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 5000000, 0))
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE number = \$1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}).
			AddRow(1, dateMock, dateMock, nil, "order123", "PROCESSING", 1000, 123))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5,"balance"=\$6,"withdrawal"=\$7 WHERE "users"."deleted_at" IS NULL AND "id" = \$8`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "balance_history" ("created_at","updated_at","deleted_at","user_id","amount","operation","order_id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
}

func WithdrawHandlerMockQueryNoOrder(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 5000000, 0))
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE number = \$1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}))
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 5000000, 0))

	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE number = \$1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"email"=$4,"password"=$5,"balance"=$6,"withdrawal"=$7 WHERE "users"."deleted_at" IS NULL AND "id" = $8`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders" ("created_at","updated_at","deleted_at","number","status","accrual","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5,"balance"=\$6,"withdrawal"=\$7 WHERE "users"."deleted_at" IS NULL AND "id" = \$8`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "balance_history" ("created_at","updated_at","deleted_at","user_id","amount","operation","order_id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
}

func WithdrawHandlerMockQueryNoOrderError(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 5000000, 0))
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE number = \$1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}))
}

func WithdrawHandlerMockQueryError(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "nUJ4D@example.com", "password", 5000000, 0))
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE number = \$1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}).
			AddRow(1, dateMock, dateMock, nil, "order123", "PROCESSING", 1000, 123))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5,"balance"=\$6,"withdrawal"=\$7 WHERE "users"."deleted_at" IS NULL AND "id" = \$8`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "balance_history" ("created_at","updated_at","deleted_at","user_id","amount","operation","order_id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	mock.ExpectRollback()
}