package main

import (
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

	go a.PollAccrualSystem(1 * time.Microsecond)

	err = s.Run(c.RunAddress)
	if err != nil {
		panic(err)
	}
}
