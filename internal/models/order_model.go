package models

import (
	"fmt"
	"github.com/theplant/luhn"
	"gorm.io/gorm"
	"market/internal/entities"
	"strconv"
)

type OrderModel struct {
	*gorm.DB
	*entities.UserEntity
}

const StatusNew = "NEW"

func NewOrderModel(db *gorm.DB, userId int) *OrderModel {
	u := &entities.UserEntity{}
	db.First(u, userId)
	return &OrderModel{
		DB:         db,
		UserEntity: u,
	}
}

func (o *OrderModel) CreateOrder(orderNumber string) error {
	num, err := strconv.Atoi(orderNumber)

	if err != nil {
		return fmt.Errorf("invalid request")
	}

	if luhn.Valid(num) {
		return fmt.Errorf("invalid format")
	}

	var existingOrder entities.OrderEntity
	result := o.DB.Where("number = ?", orderNumber).First(&existingOrder)

	if result.RowsAffected == 1 && existingOrder.User.ID == o.UserEntity.ID {
		return fmt.Errorf("already exists current user")
	}

	if result.RowsAffected == 1 {
		return fmt.Errorf("already exists")
	}

	o.UserEntity.Orders = append(o.UserEntity.Orders, entities.OrderEntity{
		Number: orderNumber,
		Status: StatusNew,
	})

	return o.DB.Save(o.UserEntity).Error
}

func (o *OrderModel) GetOrders() []entities.OrderEntity {
	return o.UserEntity.Orders
}
