package user

import (
	"encoding/json"
	"gorm.io/gorm"
	"market/internal/models"
	"net/http"
)

type WithdrawRequest struct {
	Order int `json:"order"`
	Sum   int `json:"sum"`
}

func GetBalanceHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	balance := models.NewBalanceModel(db, request.Context().Value("user_id").(int)).GetBalance()
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
	err = models.NewBalanceModel(db, request.Context().Value("user_id").(int)).Withdraw(withdrawRequest.Order, withdrawRequest.Sum)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
