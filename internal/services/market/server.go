package market

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/internal/models"
	"net/http"
)

type Server struct {
}

func NewServer(dsn string) *Server {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Order{})
	return &Server{}
}

func (s *Server) run() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// User routes
	r.Post("/api/user/register", RegisterHandler)
	r.Post("/api/user/login", LoginHandler)
	r.Post("/api/user/orders", OrderHandler)
	r.Get("/api/user/orders", GetOrdersHandler)
	r.Get("/api/user/balance", GetBalanceHandler)
	r.Post("/api/user/balance/withdraw", WithdrawHandler)
	r.Get("/api/user/withdrawals", GetWithdrawalsHandler)

	http.ListenAndServe(":8000", r)
}
