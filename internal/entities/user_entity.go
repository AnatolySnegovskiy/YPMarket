package entities

import "gorm.io/gorm"

type UserEntity struct {
	gorm.Model     `json:"gorm_._model"`
	Email          string                 `json:"email,omitempty" gorm:"not null;unique;name:email;type:varchar(255)"`
	Password       string                 `json:"password,omitempty" gorm:"not null;name:password;type:varchar(255)"`
	Balance        float64                `json:"balance,omitempty" gorm:"not null;default:0;name:balance;type:float"`
	BalanceHistory []BalanceHistoryEntity `json:"balance_history,omitempty" gorm:"foreignKey:user;references:ID"`
	Orders         []OrderEntity          `json:"orders,omitempty" gorm:"foreignKey:user;references:ID"`
}

func (UserEntity) TableName() string {
	return "users"
}
