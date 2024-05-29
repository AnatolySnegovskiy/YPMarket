package entities

import "gorm.io/gorm"

type BalanceHistoryEntity struct {
	gorm.Model `json:"gorm_._model"`
	UserID     uint    `json:"user_id,omitempty" gorm:"not null;name:user_id;type:bigint"`
	Amount     float64 `json:"amount,omitempty" gorm:"not null;name:amount;type:float"`
	Operation  string  `json:"operation,omitempty" gorm:"not null;name:operation;type:varchar(255)"`
	OrderID    uint    `json:"order_id,omitempty" gorm:"name:order_id;type:bigint"`
}

func (BalanceHistoryEntity) TableName() string {
	return "balance_history"
}
