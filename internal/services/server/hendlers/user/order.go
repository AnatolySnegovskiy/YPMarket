package user

import (
	"encoding/json"
	"gorm.io/gorm"
	"io"
	"market/internal/models"
	"net/http"
)

func CreateOrderHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	om := models.NewOrderModel(db, getUserID(request))
	err = om.CreateOrder(string(body))
	if err != nil {
		switch err.Error() {
		case "already exists current user":
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte("Already exists current user"))
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		case "already exist":
			http.Error(writer, "Already exist other user", http.StatusConflict)
			return
		case "invalid request":
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			return
		case "invalid format":
			http.Error(writer, "Invalid format", http.StatusUnprocessableEntity)
			return
		default:
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writer.WriteHeader(http.StatusAccepted)
}

func GetOrdersHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	orders := models.NewOrderModel(db, getUserID(request)).GetOrders()
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(orders)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func GetWithdrawalsHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	withdrawals := models.NewBalanceModel(db, getUserID(request)).GetWithdrawals()
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(withdrawals)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
