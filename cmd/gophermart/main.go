package main

import (
	"market/config"
	"market/internal/services/server"
)

func main() {
	c := config.NewBaseConfig()
	s, err := server.NewServer(c.DatabaseURI)
	if err != nil {
		panic(err)
	}

	err = s.Run(c.RunAddress)
	if err != nil {
		panic(err)
	}
}
