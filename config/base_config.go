package config

import (
	"flag"
	"os"
)

type BaseConfig struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func NewBaseConfig() BaseConfig {
	defaultRunAddress := "localhost:8080"
	defaultDatabaseURI := "postgres://postgres:root@localhost:5432"
	defaultAccrualSystemAddress := "http://localhost:8080"

	runAddress := getEnv("RUN_ADDRESS", defaultRunAddress)
	databaseURI := getEnv("DATABASE_URI", defaultDatabaseURI)
	accrualSystemAddress := getEnv("ACCRUAL_SYSTEM_ADDRESS", defaultAccrualSystemAddress)

	flag.StringVar(&runAddress, "a", runAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&databaseURI, "d", databaseURI, "адрес подключения к базе данных")
	flag.StringVar(&accrualSystemAddress, "r", accrualSystemAddress, "адрес системы расчёта начислений")
	flag.Parse()

	return BaseConfig{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualSystemAddress,
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
