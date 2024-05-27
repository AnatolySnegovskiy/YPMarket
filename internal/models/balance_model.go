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

type CurrentBalance struct {
	Current   float64 `json:"current,omitempty"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
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

func (o *OrderModel) Withdraw(order int, sum int) error {
	if sum <= 0 {
		return fmt.Errorf("invalid sum")
	}
	if o.UserEntity.Balance < float64(sum) {
		return fmt.Errorf("not enough money")
	}
	
	o.UserEntity.Balance -= float64(sum)
	o.UserEntity.Withdrawal += float64(sum)
	historyEntity := entities.BalanceHistoryEntity{}
	historyEntity.Amount = float64(sum)
	historyEntity.Operation = "withdraw"
	o.UserEntity.BalanceHistory = append(o.UserEntity.BalanceHistory, historyEntity)
	return o.DB.Save(o.UserEntity).Error
}
