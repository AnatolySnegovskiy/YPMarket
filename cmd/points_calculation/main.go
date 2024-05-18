package main

import "market/internal/services/points_calculation"

func main() {
	s, err := points_calculation.NewServer("postgres://postgres:root@localhost:5432")
	if err != nil {
		panic(err)
	}
	s.Run()
}
