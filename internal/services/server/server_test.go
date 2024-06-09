package server

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	s, _ := NewServer("")
	assert.NotNil(t, s)
	done := make(chan error)

	go func() {
		done <- s.Run("127.0.0.1:9090")
	}()

	select {
	case err := <-done:
		assert.Nil(t, err)
	case <-time.After(1 * time.Second):
		assert.NotNil(t, done)
	}
}

func TestServer_APIRequests(t *testing.T) {
	s, _ := NewServer("")
	assert.NotNil(t, s)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	gdb.Logger = gdb.Logger.LogMode(logger.Silent)

	s.db = gdb
	done := make(chan error)
	go func() {
		done <- s.Run("127.0.0.1:9999")
	}()

	time.Sleep(100 * time.Millisecond)
	client := &http.Client{}

	baseURL := "http://127.0.0.1:9999"
	url := fmt.Sprintf("%s/api/user/register", baseURL)
	Req, _ := http.NewRequest("POST", url, nil)
	Resp, err := client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, Resp.StatusCode)

	url = fmt.Sprintf("%s/api/user/login", baseURL)
	Req, _ = http.NewRequest("POST", url, nil)
	Resp, err = client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, Resp.StatusCode)

	url = fmt.Sprintf("%s/api/user/orders", baseURL)
	Req, _ = http.NewRequest("POST", url, nil)
	Resp, err = client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, Resp.StatusCode)

	url = fmt.Sprintf("%s/api/user/balance", baseURL)
	Req, _ = http.NewRequest("GET", url, nil)
	Resp, err = client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, Resp.StatusCode)

	url = fmt.Sprintf("%s/api/user/balance/withdraw", baseURL)
	Req, _ = http.NewRequest("POST", url, nil)
	Resp, err = client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, Resp.StatusCode)

	url = fmt.Sprintf("%s/api/user/withdrawals", baseURL)
	Req, _ = http.NewRequest("GET", url, nil)
	Resp, err = client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, Resp.StatusCode)

	url = fmt.Sprintf("%s/api/user/orders", baseURL)
	Req, _ = http.NewRequest("GET", url, nil)
	Resp, err = client.Do(Req)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	defer Resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, Resp.StatusCode)

	select {
	case err := <-done:
		assert.Nil(t, err)
	case <-time.After(1 * time.Second):
		assert.NotNil(t, done)
	}
}
