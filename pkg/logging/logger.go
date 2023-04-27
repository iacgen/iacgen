package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

func getLoggerLevel(level string) zapcore.Level {
	if zapLevel, ok := levelMap[level]; ok {
		return zapLevel
	}
	return zapcore.InfoLevel
}

func init() {
	globalLogger = newLogger("info")
}

func Init(level string) {
	globalLogger = newLogger(level)
}

func newLogger(level string) *zap.Logger {
	var err error
	zapConfig := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(getLoggerLevel(level)),
		OutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			TimeKey:    "time",
			NameKey:    "process",
			CallerKey:  "file",
			// FunctionKey:   "method",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
	}
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	return logger
}

func GetLogger() *zap.Logger {
	return globalLogger
}
