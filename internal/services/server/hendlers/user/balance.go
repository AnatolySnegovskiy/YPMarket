package user

import (
	"encoding/json"
	"gorm.io/gorm"
	"market/internal/models"
	"net/http"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func GetBalanceHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	balance := models.NewBalanceModel(db, getUserID(request)).GetBalance()
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(balance)
}

func WithdrawHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	var withdrawRequest WithdrawRequest
	err := json.NewDecoder(request.Body).Decode(&withdrawRequest)
	if err != nil {
		http.Error(writer, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = models.NewBalanceModel(db, getUserID(request)).Withdraw(withdrawRequest.Order, withdrawRequest.Sum)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
