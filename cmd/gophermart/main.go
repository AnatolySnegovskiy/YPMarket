package main

import (
	"fmt"
	"market/config"
	"market/internal/services"
	"market/internal/services/server"
	"time"
)

func main() {
	c := config.NewBaseConfig()
	s, err := server.NewServer(c.DatabaseURI)
	if err != nil {
		panic(err)
	}

	a, err := services.NewOrderAccrual(c.AccrualSystemAddress, c.DatabaseURI)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go a.PollAccrualSystem(ticker)

	fmt.Println("Server started")
	err = s.Run(c.RunAddress)
	if err != nil {
		panic(err)
	}
}
