package services

import (
	"net/http"
)

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestFetchOrderAccrual(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"order": "123", "status": "PROCESSED", "accrual": 100}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/api/orders/456":
			w.WriteHeader(http.StatusNoContent)
		case "/api/orders/789":
			w.WriteHeader(http.StatusTooManyRequests)
		case "/api/orders/999":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	o := &OrderAccrual{address: server.URL}

	testCases := []struct {
		orderNumber      string
		expectedResponse *OrderAccrualResponse
		expectedError    string
	}{
		{"123", &OrderAccrualResponse{"123", "PROCESSED", 100}, ""},
		{"456", nil, "заказ не зарегистрирован в системе расчёта"},
		{"789", nil, "превышено количество запросов к сервису"},
		{"999", nil, "внутренняя ошибка сервера"},
		{"000", nil, "неизвестный код ответа: 404"},
	}

	for _, tc := range testCases {
		resp, err := o.fetchOrderAccrual(tc.orderNumber)
		assert.Equal(t, tc.expectedResponse, resp)
		if err != nil {
			assert.EqualError(t, err, tc.expectedError)
		}
	}
}
