package entities

import (
	"gorm.io/gorm"
	"time"
)

type OrderEntity struct {
	gorm.Model `json:"gorm_._model"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty" gorm:"not null;name:updated_at;type:timestamp"`
	Number     string     `json:"order,omitempty" gorm:"not null;unique;name:order;type:varchar(255)"`
	Status     string     `json:"status,omitempty" gorm:"not null;name:status;type:varchar(255)"`
	Accrual    float64    `json:"accrual,omitempty" gorm:"not null;name:accrual;type:float"`
	User       UserEntity `json:"-" gorm:"not null;foreignKey:ID;references:UserId;name:user;type:bigint"`
}

func (OrderEntity) TableName() string {
	return "orders"
}
