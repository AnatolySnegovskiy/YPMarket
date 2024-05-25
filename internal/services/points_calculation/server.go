package points_calculation

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/internal/entities"
	"market/internal/services/points_calculation/hendlers/user"
	iMiddleware "market/internal/services/points_calculation/middleware"
	"net/http"
)

type Server struct {
	db *gorm.DB
}

func NewServer(dsn string) (*Server, error) {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	err := db.AutoMigrate(&entities.UserEntity{}, &entities.BalanceHistoryEntity{})
	if err != nil {
		return nil, err
	}
	return &Server{
		db: db,
	}, nil
}

func (s *Server) Run() {
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
	//r.Get("/api/user/balance", GetBalanceHandler)
	//r.Post("/api/user/balance/withdraw", WithdrawHandler)
	//r.Get("/api/user/withdrawals", GetWithdrawalsHandler)

	http.ListenAndServe(":8080", r)
}
