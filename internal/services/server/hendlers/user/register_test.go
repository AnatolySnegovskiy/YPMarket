package user

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"regexp"
	"testing"
)

func TestRegisterHandlers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	testCases := []testCase{
		{
			name: "RegisterHandler - nojson",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				RegisterHandler(gdb, writer, request)
			}),
			method:         "POST",
			url:            "/register",
			requestBody:    `nojson`,
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "RegisterHandler - nologin",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				RegisterHandler(gdb, writer, request)
			}),
			method:         "POST",
			url:            "/register",
			requestBody:    `{"login": "", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},
		{
			name: "RegisterHandler - nopass",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				RegisterHandler(gdb, writer, request)
			}),
			method:         "POST",
			url:            "/register",
			requestBody:    `{"login": "123", "password": ""}`,
			expectedStatus: http.StatusBadRequest,
			queryMock: func(mock sqlmock.Sqlmock) {
			},
		},

		{
			name: "RegisterHandler - StatusConflict",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				RegisterHandler(gdb, writer, request)
			}),
			method:         "POST",
			url:            "/register",
			requestBody:    RegisterRequest{Login: "login", Password: "password"},
			expectedStatus: http.StatusConflict,
			queryMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("login", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
						AddRow(123, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))
			},
		},

		{
			name: "RegisterHandler - Success",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				RegisterHandler(gdb, writer, request)
			}),
			method:         "POST",
			url:            "/register",
			requestBody:    RegisterRequest{Login: "login", Password: "TE7AcasnAMewDfIjqVJJX"},
			expectedStatus: http.StatusOK,
			queryMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("login", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}))
				registration(mock)
			},
		},

		{
			name: "RegisterHandler - StatusInternalServerError",
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				RegisterHandler(gdb, writer, request)
			}),
			method:         "POST",
			url:            "/register",
			requestBody:    RegisterRequest{Login: "login", Password: "TE7AcasnAMewDfIjqVJJX"},
			expectedStatus: http.StatusInternalServerError,
			queryMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("login", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}))
				registrationFatal(mock)
			},
		},
	}

	RunTestCases(t, testCases, mock)
}

func registrationFatal(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password","balance","withdrawal") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()
}

func registration(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password","balance","withdrawal") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password", "balance", "withdrawal"}).
			AddRow(1, dateMock, dateMock, nil, "login@example.com", "$2a$10$2/cm/mpDH7sLoYHResqdvukbGA.6WcHkEFYDmSAhFIwjMsLdxyIxe", 5000000, 0))
}
