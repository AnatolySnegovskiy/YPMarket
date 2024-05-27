package services

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"market/internal/models"
	db2 "market/internal/system/db"
	"net/http"
	"time"
)

type OrderAccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type OrderAccrual struct {
	address string
	db      *gorm.DB
}

func NewOrderAccrual(address string, dsn string) (*OrderAccrual, error) {
	db, err := db2.Init(dsn)
	return &OrderAccrual{
		address: address,
		db:      db,
	}, err
}

func (o *OrderAccrual) PollAccrualSystem(interval time.Duration) {
	ticker := time.NewTicker(interval)
	balanceModel := models.NewBalanceModel(o.db, 0)
	orderModel := models.NewOrderModel(o.db, 0)

	for {
		select {
		case <-ticker.C:
			orders := orderModel.GetOrdersByStatus([]string{"REGISTERED", "PROCESSING"})
			for _, order := range orders {
				res, err := o.fetchOrderAccrual(order.Number)

				if err != nil {
					fmt.Println(err)
					continue
				}

				order.Accrual = res.Accrual
				order.Status = res.Status

				if res.Status == "PROCESSED" {
					err := balanceModel.Deposit("", res.Accrual)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}

			if len(orders) != 0 {
				o.db.Save(orders)
			}
		}
	}
}

func (o *OrderAccrual) fetchOrderAccrual(orderNumber string) (*OrderAccrualResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 10, // Таймаут запроса
	}

	url := fmt.Sprintf("%s/api/orders/%s", o.address, orderNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response OrderAccrualResponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	case http.StatusNoContent:
		return nil, fmt.Errorf("заказ не зарегистрирован в системе расчёта")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("превышено количество запросов к сервису")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	default:
		return nil, fmt.Errorf("неизвестный код ответа: %d", resp.StatusCode)
	}
}
