package points_calculation

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/internal/entities"
	"market/internal/services/points_calculation/hendlers/user"
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

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// User routes
	r.Post("/api/user/register",
		func(writer http.ResponseWriter, request *http.Request) {
			user.RegisterHandler(s.db, writer, request)
		},
	)

	r.Post("/api/user/login",
		func(writer http.ResponseWriter, request *http.Request) {
			user.LoginHandler(s.db, writer, request)
		})
	//r.Post("/api/user/orders", OrderHandler)
	//r.Get("/api/user/orders", GetOrdersHandler)
	//r.Get("/api/user/balance", GetBalanceHandler)
	//r.Post("/api/user/balance/withdraw", WithdrawHandler)
	//r.Get("/api/user/withdrawals", GetWithdrawalsHandler)

	http.ListenAndServe(":8080", r)
}
