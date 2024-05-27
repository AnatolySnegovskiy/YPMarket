package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
	"market/internal/services/server/hendlers/user"
	iMiddleware "market/internal/services/server/middleware"
	db2 "market/internal/system/db"
	"net/http"
)

type Server struct {
	db *gorm.DB
}

func NewServer(dsn string) (*Server, error) {
	db, err := db2.Init(dsn)
	return &Server{
		db: db,
	}, err
}

func (s *Server) Run(runAddress string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/api/user/register",
		func(writer http.ResponseWriter, request *http.Request) {
			user.RegisterHandler(s.db, writer, request)
		},
	)

	r.Post("/api/user/login",
		func(writer http.ResponseWriter, request *http.Request) {
			user.LoginHandler(s.db, writer, request)
		})
	r.With(iMiddleware.JwtAuthMiddleware).Post("/api/user/orders",
		func(writer http.ResponseWriter, request *http.Request) {
			user.CreateOrderHandler(s.db, writer, request)
		})
	r.With(iMiddleware.JwtAuthMiddleware).Get("/api/user/orders",
		func(writer http.ResponseWriter, request *http.Request) {
			user.GetOrdersHandler(s.db, writer, request)
		})
	r.With(iMiddleware.JwtAuthMiddleware).Get("/api/user/balance",
		func(writer http.ResponseWriter, request *http.Request) {
			user.GetBalanceHandler(s.db, writer, request)
		})
	r.Post("/api/user/balance/withdraw",
		func(writer http.ResponseWriter, request *http.Request) {
			user.WithdrawHandler(s.db, writer, request)
		})
	r.Get("/api/user/withdrawals",
		func(writer http.ResponseWriter, request *http.Request) {
			user.GetWithdrawalsHandler(s.db, writer, request)
		})

	http.ListenAndServe(runAddress, r)
}
