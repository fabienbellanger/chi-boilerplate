package logger

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fabienbellanger/goutils"
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

// getLoggerOutputs returns an array with the log outputs.
// Outputs can be stdout and/or file.
func getLoggerOutputs(logOutputs []string, appName, filePath string) (outputs []string, err error) {
	if goutils.StringInSlice("file", logOutputs) {
		logPath := path.Clean(filePath)
		_, err := os.Stat(logPath)
		if err != nil {
			return nil, err
		}

		if appName == "" {
			return nil, errors.New("no APP_NAME variable defined")
		}

		outputs = append(outputs, fmt.Sprintf("%s/%s.log",
			logPath,
			appName))
	}
	if goutils.StringInSlice("stdout", logOutputs) {
		outputs = append(outputs, "stdout")
	}
	return
}

func (c *CustomLog) Info(msg string, fields ...interface{}) {
	zapFields := toZapFileds(fields...)
	c.inner.Info(msg, zapFields...)
}

func (c *CustomLog) Error(msg string, fields ...interface{}) {
	zapFields := toZapFileds(fields...)
	c.inner.Error(msg, zapFields...)
}
