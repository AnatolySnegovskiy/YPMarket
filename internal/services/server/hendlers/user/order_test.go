package user

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestOrderHandlers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	testCases := []testCase{
		{
			name: "CreateOrderHandler - StatusBadRequest",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				CreateOrderHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/create_order",
			requestBody:    strings.NewReader(`newOrder`),
			expectedStatus: http.StatusBadRequest,
			queryMock:      login,
		},
		{
			name: "CreateOrderHandler - StatusUnprocessableEntity",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				CreateOrderHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/create_order",
			requestBody:    123123,
			expectedStatus: http.StatusUnprocessableEntity,
			queryMock:      login,
		},
		{
			name: "CreateOrderHandler - StatusConflict",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				CreateOrderHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/create_order",
			requestBody:    60480142,
			expectedStatus: http.StatusConflict,
			queryMock:      CreateOrderStatusConflictMockQuery,
		},
		{
			name: "CreateOrderHandler - Success",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				CreateOrderHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/create_order",
			requestBody:    60480142,
			expectedStatus: http.StatusAccepted,
			queryMock:      CreateOrderHandlerMockQuery,
		},
		{
			name: "CreateOrderHandler - StatusInternalServerError",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				CreateOrderHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/create_order",
			requestBody:    60480142,
			expectedStatus: http.StatusInternalServerError,
			queryMock:      StatusInternalServerErrorOrderHandlerMockQuery,
		},

		{
			name: "GetOrdersHandler - Success",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				GetOrdersHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/getorder",
			requestBody:    "",
			expectedStatus: http.StatusOK,
			queryMock: func(mock sqlmock.Sqlmock) {
				login(mock)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL AND "users"."id" = $2 ORDER BY "users"."id" LIMIT $3`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
						AddRow(123, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE "orders"."user_id" = $1 AND "orders"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}).
						AddRow(1, dateMock, dateMock, nil, "order123", "PROCESSING", 1000, 123))
			},
		},

		{
			name: "GetWithdrawalsHandler - Success",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				GetWithdrawalsHandler(gdb, writer, upLogin(request))
			}),
			method:         "POST",
			url:            "/getorder",
			requestBody:    "",
			expectedStatus: http.StatusOK,
			queryMock: func(mock sqlmock.Sqlmock) {
				login(mock)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT sum(balance_history.amount) as sum, balance_history.updated_at as processed_at, orders.number as order FROM "balance_history" LEFT JOIN orders ON balance_history.order_id = orders.id WHERE (orders.user_id = $1 AND balance_history.operation = $2) AND "balance_history"."deleted_at" IS NULL GROUP BY balance_history.updated_at, orders.number ORDER BY balance_history.updated_at desc`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"sum", "processed_at", "order"}).
						AddRow(100, dateMock, 12546445))
			},
		},
	}

	RunTestCases(t, testCases, mock)
}

func CreateOrderStatusConflictMockQuery(mock sqlmock.Sqlmock) {
	login(mock)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE number = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}).
			AddRow(123, dateMock, dateMock, nil, "order123", "PROCESSING", 1000, 1))
}

func CreateOrderHandlerMockQuery(mock sqlmock.Sqlmock) {
	login(mock)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE number = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
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
}

func StatusInternalServerErrorOrderHandlerMockQuery(mock sqlmock.Sqlmock) {
	login(mock)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE number = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "number", "status", "accrual", "user_id"}))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"email"=$4,"password"=$5,"balance"=$6,"withdrawal"=$7 WHERE "users"."deleted_at" IS NULL AND "id" = $8`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders" ("created_at","updated_at","deleted_at","number","status","accrual","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
}

func login(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(123, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))
}
