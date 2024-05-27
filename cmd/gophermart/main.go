package main

import (
	"flag"
	"market/internal/services"
	"market/internal/services/server"
	"os"
	"time"
)

func main() {
	defaultRunAddress := "localhost:8080"
	defaultDatabaseURI := "postgres://postgres:root@localhost:5432"
	defaultAccrualSystemAddress := "http://localhost:8000"

	runAddress := getEnv("RUN_ADDRESS", defaultRunAddress)
	databaseURI := getEnv("DATABASE_URI", defaultDatabaseURI)
	accrualSystemAddress := getEnv("ACCRUAL_SYSTEM_ADDRESS", defaultAccrualSystemAddress)

	flag.StringVar(&runAddress, "a", runAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&databaseURI, "d", databaseURI, "адрес подключения к базе данных")
	flag.StringVar(&accrualSystemAddress, "r", accrualSystemAddress, "адрес системы расчёта начислений")
	flag.Parse()

	a, err := services.NewOrderAccrual(accrualSystemAddress, databaseURI)
	if err != nil {
		panic(err)
	}

	go a.PollAccrualSystem(5 * time.Second)

	s, err := server.NewServer(databaseURI)
	if err != nil {
		panic(err)
	}
	s.Run(runAddress)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
