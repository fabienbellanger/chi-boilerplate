package pkg

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigAMQP(t *testing.T) {
	viper.Set("AMQP_URL", "amqp://guest:guest@localhost:5672/")

	c := NewConfigAMQP()

	assert.Equal(t, c.URL, "amqp://guest:guest@localhost:5672/")
}

func TestNewConfigPprof(t *testing.T) {
	viper.Set("PPROF_ENABLE", true)
	viper.Set("PPROF_BASICAUTH_USERNAME", "john")
	viper.Set("PPROF_BASICAUTH_PASSWORD", "test")

	c := NewConfigPprof()

	assert.Equal(t, c.Enable, true)
	assert.Equal(t, c.BasicAuthUsername, "john")
	assert.Equal(t, c.BasicAuthPassword, "test")
}

func TestNewConfigCORS(t *testing.T) {
	// TODO: Implement
}

func TestNewConfigJWT(t *testing.T) {
	// TODO: Implement
}
