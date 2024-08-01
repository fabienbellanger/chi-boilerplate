package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger = zap.Logger

// NewZapLogger creates a new custom Zap logger.
func NewZapLogger() (*ZapLogger, error) {
	// Logs outputs
	outputs, err := getLoggerOutputs(viper.GetStringSlice("LOG_OUTPUTS"), viper.GetString("APP_NAME"), viper.GetString("LOG_PATH"))
	if err != nil {
		return nil, err
	}

	// Level
	level := getZapLoggerLevel(viper.GetString("LOG_LEVEL"), viper.GetString("APP_ENV"))

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      outputs,
		ErrorOutputPaths: outputs,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return zap.NewProduction()
	}

	return logger, nil
}

// getZapLoggerLevel returns the minimum log level.
// If nothing is specified in the environment variable LOG_LEVEL,
// The level is DEBUG in development mode and WARN in others cases.
func getZapLoggerLevel(l string, env string) (level zapcore.Level) {
	switch l {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		if env == "development" {
			level = zapcore.DebugLevel
		} else {
			level = zapcore.WarnLevel
		}
	}
	return
}

// toZapFileds converts fields to zap.Field.
// TODO: Add tests
func toZapFileds(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		if f, ok := field.(zap.Field); ok {
			zapFields[i] = f
		} /* else {
			// Handle the case where the field is not of type zap.Field
			zapFields[i] = zap.Any(fmt.Sprintf("field%d", i), field)
		}*/
	}
	return zapFields
}
