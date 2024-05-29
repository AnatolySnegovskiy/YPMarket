package entities

import (
	"gorm.io/gorm"
	"time"
)

type OrderEntity struct {
	gorm.Model `json:"-"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" gorm:"not null;name:updated_at;type:timestamp"`
	Number     string    `json:"order,omitempty" gorm:"not null;unique;name:order;type:varchar(255)"`
	Status     string    `json:"status,omitempty" gorm:"not null;name:status;type:varchar(255)"`
	Accrual    float64   `json:"accrual,omitempty" gorm:"not null;name:accrual;type:float;default:0"`
	UserID     uint      `json:"-" gorm:"not null;name:user_id;type:bigint"`
}

func (OrderEntity) TableName() string {
	return "orders"
}
