package middleware

import (
	"github.com/stretchr/testify/assert"
	"market/internal/system"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJwtAuthMiddleware(t *testing.T) {
	server := httptest.NewServer(JwtAuthMiddleware(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userID := request.Context().Value(UserIDContextKey).(int)
		assert.Equal(t, 123, userID)
	})))
	defer server.Close()
	token, _ := system.CreateToken(123)

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("GET", server.URL, nil)
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
