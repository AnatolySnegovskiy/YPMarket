package user

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"io"
	"market/internal/models"
	"net/http"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(db *gorm.DB, writer http.ResponseWriter, request *http.Request) {
	db.WithContext(request.Context())
	var registerRequest RegisterRequest
	bodyBytes, err := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.NewDecoder(request.Body).Decode(&registerRequest)

	if err != nil || registerRequest.Login == "" || registerRequest.Password == "" {
		http.Error(writer, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	u := models.NewUserModel(db)
	err = u.Registration(registerRequest.Login, registerRequest.Password)

	if err != nil {
		switch err.Error() {
		case "user already exists":
			http.Error(writer, "User already exists", http.StatusConflict)
		default:
			http.Error(writer, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	clonedRequest := request.Clone(request.Context())
	clonedRequest.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	LoginHandler(db, writer, clonedRequest)
}
