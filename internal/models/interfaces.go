package models

import "market/internal/entities"

type BalanceManagerInterface interface {
	GetBalance() *CurrentBalance
	Deposit(order string, sum float64) error
	Withdraw(order string, sum float64) error
	GetWithdrawals() []Withdrawals
}

type OrderModelInterface interface {
	CreateOrder(orderNumber string) error
	GetOrders() []entities.OrderEntity
	GetOrdersByStatus(status []string) []entities.OrderEntity
}

type UserModelInterface interface {
	Authenticate(email string, password string) (string, error)
	Registration(email string, password string) error
	getUserByEmail(email string) (*entities.UserEntity, error)
}
