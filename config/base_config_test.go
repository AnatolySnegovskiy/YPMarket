package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewBaseConfig(t *testing.T) {
	t.Run("Default values", func(t *testing.T) {
		config := NewBaseConfig()
		assert.Equal(t, "localhost:8000", config.RunAddress, "expected default address")
		assert.Equal(t, "postgres://postgres:root@localhost:5432", config.DatabaseURI, "expected default database uri")
		assert.Equal(t, "http://localhost:8080", config.AccrualSystemAddress, "expected default accrual system address")
	})

	t.Run("ENV", func(t *testing.T) {
		resetVars()
		os.Setenv("RUN_ADDRESS", "localhost:1000")
		os.Setenv("DATABASE_URI", "postgres")
		os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:100")
		config := NewBaseConfig()
		assert.Equal(t, "localhost:1000", config.RunAddress, "expected default address")
		assert.Equal(t, "postgres", config.DatabaseURI, "expected default database uri")
		assert.Equal(t, "http://localhost:100", config.AccrualSystemAddress, "expected default accrual system address")
	})

	t.Run("Flag", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-a", "localhost:0000", "-d", "postgres2", "-r", "http://localhost2:8080"}
		config := NewBaseConfig()
		assert.Equal(t, "localhost:0000", config.RunAddress, "expected default address")
		assert.Equal(t, "postgres2", config.DatabaseURI, "expected default database uri")
		assert.Equal(t, "http://localhost2:8080", config.AccrualSystemAddress, "expected default accrual system address")
	})
}

func resetVars() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd"}
	os.Clearenv()
}
