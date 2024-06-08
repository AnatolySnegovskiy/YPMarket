package server

import (
	"github.com/stretchr/testify/assert"
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
