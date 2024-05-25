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

	err = models.NewOrderModel(db, request.Context().Value("user_id").(int)).CreateOrder(string(body))
	if err != nil {
		switch err.Error() {
		case "already exists current user":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Already exists current user"))
			break
		case "already exist":
			http.Error(writer, "Already exist other user", http.StatusConflict)
			break
		case "invalid request":
			http.Error(writer, "Invalid request", http.StatusBadRequest)
			break
		case "invalid format":
			http.Error(writer, "Invalid format", http.StatusUnprocessableEntity)
		default:
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetOrdersHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	orders := models.NewOrderModel(db, request.Context().Value("user_id").(int)).GetOrders()
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(orders)
}
