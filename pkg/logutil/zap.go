package logutil

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const JSONFormat = "json"
const TextFormat = "console"

var logger *zap.Logger

func Logger() *zap.Logger {
	return logger
}

func New(level, format string) error {
	if format != JSONFormat && format != TextFormat {
		return fmt.Errorf("log format must is json or console but got '%s'", format)
	}

	encodeLevel := zapcore.LowercaseLevelEncoder
	if format == TextFormat {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLogger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLoggerLevelByString(level)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: format,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "timestamp",
			LevelKey:         "severity",
			NameKey:          "logger",
			CallerKey:        "trace",
			FunctionKey:      zapcore.OmitKey,
			MessageKey:       "message",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      encodeLevel,
			EncodeTime:       zapcore.RFC3339TimeEncoder,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	logger = zapLogger

	return err
}

func getLoggerLevelByString(level string) zapcore.Level {
	var loggerLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		loggerLevel = zap.DebugLevel
	case "info":
		loggerLevel = zap.InfoLevel
	case "warn":
		loggerLevel = zap.WarnLevel
	case "error":
		loggerLevel = zap.ErrorLevel
	case "panic":
		loggerLevel = zap.PanicLevel
	case "fatal":
		loggerLevel = zap.FatalLevel
	default:
		loggerLevel = zap.ErrorLevel
	}

	return loggerLevel
}
