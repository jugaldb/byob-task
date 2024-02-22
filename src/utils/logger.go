package utils

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StandardLogger enforces specific log message formats.
type StandardLogger struct {
	*zap.SugaredLogger
}

// IntegerLevelEncoder returns custom encoder for level field.
func IntegerLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt8((int8(l) + 3) * 10)
}

var appLogger *StandardLogger

// NewLogger creates a new application logger.
func NewLogger(config *Config) *StandardLogger {
	var cfg zap.Config
	outputLevel, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		outputLevel = zapcore.InfoLevel
	}

	if config.DGN != "local" {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stdout"}
		cfg.InitialFields = map[string]any{"name": "fam.service.halfblood"}
		cfg.EncoderConfig.EncodeLevel = IntegerLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.TimeKey = "time"
		cfg.Level = zap.NewAtomicLevelAt(outputLevel)
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &StandardLogger{SugaredLogger: logger.Sugar()}
}

func SetAppLogger(l *StandardLogger) *StandardLogger {
	appLogger = l
	return appLogger
}

func GetAppLogger() *StandardLogger {
	if appLogger == nil {
		panic("Logger not initialized")
	}
	return appLogger
}

func (l *StandardLogger) Printf(format string, v ...any) {
	if strings.Contains(format, "failed") {
		l.Errorf(format, v)
	} else {
		l.Infof(format, v)
	}
}
