package logger

import (
	"fmt"
	"go.uber.org/zap"
)

type CustomLogger interface {
	// Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	// Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	// Fatal(msg string, fields ...interface{})
	// Panic(msg string, fields ...interface{})
}

type CustomLog struct {
	inner *zap.Logger
}

func toZapFileds(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		if f, ok := field.(zap.Field); ok {
			zapFields[i] = f
		} else {
			// Handle the case where the field is not of type zap.Field
			zapFields[i] = zap.Any(fmt.Sprintf("field%d", i), field)
		}
	}
	return zapFields
}

func (c *CustomLog) Info(msg string, fields ...interface{}) {
	zapFields := toZapFileds(fields...)
	c.inner.Info(msg, zapFields...)
}

func (c *CustomLog) Error(msg string, fields ...interface{}) {
	zapFields := toZapFileds(fields...)
	c.inner.Error(msg, zapFields...)
}
