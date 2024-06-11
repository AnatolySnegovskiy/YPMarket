package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/internal/entities"
)

func Init(dsn string) (*gorm.DB, error) {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	err := db.AutoMigrate(&entities.UserEntity{}, &entities.BalanceHistoryEntity{}, &entities.OrderEntity{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
