package accrual

import (
	"market/config"
	"market/internal/services"
	"time"
)

func main() {
	c := config.NewBaseConfig()
	a, err := services.NewOrderAccrual(c.AccrualSystemAddress, c.DatabaseURI)
	if err != nil {
		panic(err)
	}

	go a.PollAccrualSystem(5 * time.Second)
}
