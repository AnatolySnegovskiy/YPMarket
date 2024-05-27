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

func NewOrderModel(db *gorm.DB, userID int) *OrderModel {
	u := &entities.UserEntity{}
	if userID != 0 {
		db.First(u, userID)
	}
	return &OrderModel{
		DB:         db,
		UserEntity: u,
	}
}

func (m *OrderModel) CreateOrder(orderNumber string) error {
	num, err := strconv.Atoi(orderNumber)

	if err != nil {
		return fmt.Errorf("invalid request")
	}

	if luhn.Valid(num) {
		return fmt.Errorf("invalid format")
	}

	var existingOrder entities.OrderEntity
	result := m.DB.Where("number = ?", orderNumber).First(&existingOrder)

	if result.RowsAffected == 1 && existingOrder.User.ID == m.UserEntity.ID {
		return fmt.Errorf("already exists current user")
	}

	if result.RowsAffected == 1 {
		return fmt.Errorf("already exists")
	}

	m.UserEntity.Orders = append(m.UserEntity.Orders, entities.OrderEntity{
		Number: orderNumber,
		Status: StatusNew,
	})

	return m.DB.Save(m.UserEntity).Error
}

func (m *OrderModel) GetOrders() []entities.OrderEntity {
	return m.UserEntity.Orders
}

func (m *OrderModel) GetOrdersByStatus(status []string) []entities.OrderEntity {
	var orders []entities.OrderEntity
	m.DB.Where("status IN ?", status).Find(&orders)
	return orders
}
