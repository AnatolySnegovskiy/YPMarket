package services

import (
	"context"
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

func (o *OrderAccrual) PollAccrualSystem(ticker *time.Ticker, ctx context.Context) {
	balanceModel := models.NewBalanceModel(o.db, 0)
	orderModel := models.NewOrderModel(o.db, 0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			orders := orderModel.GetOrdersByStatus([]string{"REGISTERED", "PROCESSING", "NEW"})
			for _, order := range orders {
				res, err := o.fetchOrderAccrual(order.Number)

				if err != nil {
					fmt.Println(err)
					continue
				}

				order.Accrual = res.Accrual
				order.Status = res.Status
				o.db.Save(order)

				if res.Status == "PROCESSED" {
					err := balanceModel.Deposit(order.Number, res.Accrual)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		}
	}
}

func (o *OrderAccrual) fetchOrderAccrual(orderNumber string) (*OrderAccrualResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 2,
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
		return nil, fmt.Errorf("the order is not registered in the payment system")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("the number of requests to the service has been exceeded\n")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error")
	default:
		return nil, fmt.Errorf("unknown response code: %d", resp.StatusCode)
	}
}
