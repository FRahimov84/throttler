package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerCtxKey string

const newKey = LoggerCtxKey("abrakadabra")

func FromCtx(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(newKey).(*zap.Logger)
	if ok {
		return l
	}
	return nil
}

func ToCtx(parent context.Context, l *zap.Logger) context.Context {
	return context.WithValue(parent, newKey, l)
}

func InitLogger(filename string) *zap.Logger {

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
