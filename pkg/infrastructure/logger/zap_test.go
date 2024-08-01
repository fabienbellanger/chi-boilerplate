package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestZapLogLevel(t *testing.T) {
	cases := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
		"":      zapcore.DebugLevel,
	}

	env := "development"
	for level, expected := range cases {
		assert.Equal(t, expected, getZapLoggerLevel(level, env))
	}

	env = "production"
	cases[""] = zapcore.WarnLevel
	for level, expected := range cases {
		assert.Equal(t, expected, getZapLoggerLevel(level, env))
	}
}
