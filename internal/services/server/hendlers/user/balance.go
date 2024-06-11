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
	_ = json.NewEncoder(writer).Encode(balance)
	writer.WriteHeader(http.StatusOK)
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
		switch err.Error() {
		case "not enough money":
			http.Error(writer, "not enough money", http.StatusPaymentRequired)
			return
		case "invalid sum":
			http.Error(writer, "invalid sum", http.StatusUnprocessableEntity)
			return
		default:
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writer.WriteHeader(http.StatusOK)
}
