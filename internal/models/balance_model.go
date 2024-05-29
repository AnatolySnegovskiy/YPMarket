package models

import (
	"fmt"
	"gorm.io/gorm"
	"market/internal/entities"
)

type BalanceModel struct {
	*gorm.DB
	*entities.UserEntity
}

const withdrawOperation = "withdraw"
const depositOperation = "deposit"

type CurrentBalance struct {
	Current   float64 `json:"current,omitempty"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
}

type Withdrawals struct {
	ProcessedAt string `json:"processed_at,omitempty"`
	Order       string `json:"order,omitempty"`
	Sum         int    `json:"sum,omitempty"`
}

func NewBalanceModel(db *gorm.DB, userID int) *OrderModel {
	u := &entities.UserEntity{}
	if userID != 0 {
		db.First(u, userID)
	}
	return &OrderModel{
		DB:         db,
		UserEntity: u,
	}
}

func (m *OrderModel) GetBalance() *CurrentBalance {
	var balance CurrentBalance
	balance.Current = m.UserEntity.Balance
	balance.Withdrawn = m.UserEntity.Withdrawal
	return &balance
}

func (m *OrderModel) Withdraw(order string, sum float64) error {
	if sum <= 0 {
		return fmt.Errorf("invalid sum")
	}
	if m.UserEntity.Balance < sum {
		return fmt.Errorf("not enough money")
	}

	orderEntity := entities.OrderEntity{}
	m.DB.Model(&entities.OrderEntity{}).Where("number = ?", order).First(&orderEntity)

	m.UserEntity.Balance -= sum
	m.UserEntity.Withdrawal += sum
	historyEntity := entities.BalanceHistoryEntity{}
	historyEntity.Amount = sum
	historyEntity.Operation = withdrawOperation
	historyEntity.OrderID = orderEntity.ID
	m.UserEntity.BalanceHistory = append(m.UserEntity.BalanceHistory, historyEntity)
	return m.DB.Save(m.UserEntity).Error
}

func (m *OrderModel) Deposit(order string, sum float64) error {
	if sum <= 0 {
		return fmt.Errorf("invalid sum")
	}
	orderEntity := entities.OrderEntity{}
	m.DB.Model(&entities.OrderEntity{}).Where("number = ?", order).First(&orderEntity)

	if orderEntity.UserID != m.UserEntity.ID {
		return fmt.Errorf("invalid order")
	}

	m.DB.Model(&entities.UserEntity{}).Where("id = ?", orderEntity.UserID).First(m.UserEntity)
	m.UserEntity.Balance += sum
	historyEntity := entities.BalanceHistoryEntity{}
	historyEntity.Amount = sum
	historyEntity.Operation = depositOperation
	historyEntity.OrderID = orderEntity.ID
	m.UserEntity.BalanceHistory = append(m.UserEntity.BalanceHistory, historyEntity)
	return m.DB.Save(m.UserEntity).Error
}

func (m *OrderModel) GetWithdrawals() []Withdrawals {
	var withdrawals []Withdrawals
	m.DB.Model(&entities.BalanceHistoryEntity{}).Select("sum(balance_history.amount) as sum, balance_history.updated_at as processed_at, orders.number as order").
		Joins("JOIN orders ON balance_history.order_id = orders.id").
		Where("orders.user_id = ? AND balance_history.operation = ?", m.UserEntity.ID, "withdraw").Group("order").Find(withdrawals)
	return withdrawals
}
