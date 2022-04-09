package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
}

// zap logger delegates all calls to the underlying zap.Logger
type zapLogger struct {
	logger *zap.Logger
}

// Info logs an info msg with fields
func (l zapLogger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an error msg with fields
func (l zapLogger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l zapLogger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l zapLogger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// With creates a child log, and optionally adds some context fields to that log.
func (l zapLogger) With(fields ...zapcore.Field) Logger {
	return zapLogger{logger: l.logger.With(fields...)}
}
