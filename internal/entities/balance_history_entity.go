package entities

import "gorm.io/gorm"

type BalanceHistoryEntity struct {
	gorm.Model `json:"gorm_._model"`
	User       UserEntity  `json:"user,omitempty" gorm:"not null;foreignKey:ID;references:UserId;name:user;type:bigint"`
	Amount     float64     `json:"amount,omitempty" gorm:"not null;name:amount;type:float"`
	Operation  string      `json:"operation,omitempty" gorm:"not null;name:operation;type:varchar(255)"`
	Order      OrderEntity `json:"order,omitempty" gorm:"foreignKey:ID;references:OrderID;name:order;type:bigint"`
}

func (BalanceHistoryEntity) TableName() string {
	return "balance_history"
}