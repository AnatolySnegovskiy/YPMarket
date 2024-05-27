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

func NewBalanceModel(db *gorm.DB, userId int) *OrderModel {
	u := &entities.UserEntity{}
	db.First(u, userId)
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

func (m *OrderModel) Withdraw(order int, sum int) error {
	if sum <= 0 {
		return fmt.Errorf("invalid sum")
	}
	if m.UserEntity.Balance < float64(sum) {
		return fmt.Errorf("not enough money")
	}

	m.UserEntity.Balance -= float64(sum)
	m.UserEntity.Withdrawal += float64(sum)
	historyEntity := entities.BalanceHistoryEntity{}
	historyEntity.Amount = float64(sum)
	historyEntity.Operation = withdrawOperation
	historyEntity.User = *m.UserEntity
	historyEntity.Order = m.UserEntity.Orders[order]
	m.UserEntity.BalanceHistory = append(m.UserEntity.BalanceHistory, historyEntity)
	return m.DB.Save(m.UserEntity).Error
}

func (m *OrderModel) Deposit(order int, sum int) error {
	if sum <= 0 {
		return fmt.Errorf("invalid sum")
	}

	m.UserEntity.Balance += float64(sum)
	historyEntity := entities.BalanceHistoryEntity{}
	historyEntity.Amount = float64(sum)
	historyEntity.Operation = depositOperation
	historyEntity.User = *m.UserEntity
	historyEntity.Order = m.UserEntity.Orders[order]
	m.UserEntity.BalanceHistory = append(m.UserEntity.BalanceHistory, historyEntity)
	return m.DB.Save(m.UserEntity).Error
}

func (m *OrderModel) GetWithdrawals() []Withdrawals {
	var withdrawals []Withdrawals
	m.DB.Model(&entities.BalanceHistoryEntity{}).Select("sum(amount) as sum, processed_at as processed_at, orders.number as order").
		Joins("JOIN orders ON balance_history.order_id = orders.id").
		Where("user_id = ? AND operation = ?", m.UserEntity.ID, "withdraw").Find(withdrawals)
	return withdrawals
}
