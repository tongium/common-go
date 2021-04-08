package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func GetZapLogger(level string) (logger *zap.Logger, err error) {
	logger, err = zap.Config{
		Level:       zap.NewAtomicLevelAt(getLoggerLevelByString(level)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "trace",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     zapcore.OmitKey,
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddCallerSkip(1))

	return
}
