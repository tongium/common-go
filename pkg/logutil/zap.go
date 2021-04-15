package logutil

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var quiet bool = os.Getenv("LOG_QUIET") == "true"

// Get global logger, create new from environment if not exist.
//
// Make sure you call New() before call this or set LOG_LEVEL and LOG_ENCODING.
func Logger() *zap.Logger {
	if logger == nil {
		if err := NewFromEnv(); err != nil && !quiet {
			fmt.Println("ERROR: can not init global logger due", err.Error())
		}
	}

	return logger
}

// Create global logger from environment
//
// LOG_LEVEL (default is error), LOG_ENCODING (default is json)
func NewFromEnv() error {
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok && !quiet {
		fmt.Println("WARNING: LOG_LEVEL is nil")
	}

	encoding, _ := os.LookupEnv("LOG_ENCODING")

	return New(level, encoding)
}

// Create global logger with stackdriver format
//
// level: debug, info, warn, error, panic, and fatal (default is error)
//
// encoding: json and console (default is json)
func New(level, encoding string) error {
	encodeLevel := zapcore.LowercaseLevelEncoder
	var development bool

	if strings.ToLower(encoding) == "console" {
		encodeLevel = zapcore.CapitalColorLevelEncoder
		development = true
		encoding = "console"
	} else {
		encoding = "json"
	}

	zapLogger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLoggerLevelByString(level)),
		Development: development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: encoding,
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

	if err != nil {
		return err
	}

	// set global logger
	logger = zapLogger
	return nil
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
