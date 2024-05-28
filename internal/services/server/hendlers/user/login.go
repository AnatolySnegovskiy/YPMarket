package user

import (
	"encoding/json"
	"gorm.io/gorm"
	"market/internal/models"
	"market/internal/services/server/middleware"
	"net/http"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func LoginHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	var loginRequest LoginRequest

	err := json.NewDecoder(request.Body).Decode(&loginRequest)

	if err != nil || loginRequest.Login == "" || loginRequest.Password == "" {
		http.Error(writer, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	u := models.NewUserModel(db)
	token, err := u.Authenticate(loginRequest.Login, loginRequest.Password)

	if err != nil {
		switch err.Error() {
		case "invalid login or password":
			http.Error(writer, "Invalid login or password", http.StatusUnauthorized)
			return
		default:
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(token))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func getUserID(request *http.Request) int {
	return request.Context().Value(middleware.UserIDContextKey).(int)
}
