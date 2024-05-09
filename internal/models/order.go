package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model `json:"gorm_._model"`
	Number     string `json:"number,omitempty" gorm:"not null;unique;name:number;type:varchar(255)"`
	Status     string `json:"status,omitempty" gorm:"not null;default:'NEW';name:status;type:varchar(255)"` // NEW, PROCESSING, INVALID, PROCESSED
}

func NewOrderModel() *Order {
	return &Order{}
}
