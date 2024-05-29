package entities

import "gorm.io/gorm"

type UserEntity struct {
	gorm.Model     `json:"gorm_._model"`
	Email          string                 `json:"email,omitempty" gorm:"not null;unique;name:email;type:varchar(255)"`
	Password       string                 `json:"password,omitempty" gorm:"not null;name:password;type:varchar(255)"`
	Balance        float64                `json:"balance,omitempty" gorm:"not null;default:0;name:balance;type:float"`
	Withdrawal     float64                `json:"withdrawal,omitempty" gorm:"not null;default:0;name:withdrawal;type:float"`
	BalanceHistory []BalanceHistoryEntity `json:"balance_history,omitempty" gorm:"foreignKey:UserID;"`
	Orders         []OrderEntity          `json:"orders,omitempty" gorm:"foreignKey:UserID;"`
}

func (UserEntity) TableName() string {
	return "users"
}
