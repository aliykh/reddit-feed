package log

import "go.uber.org/zap/zapcore"

type Field = zapcore.Field

type MockLogger struct{}

func (m MockLogger) Debug(msg string, fields ...Field)       {}
func (m MockLogger) Error(msg string, fields ...Field)       {}
func (m MockLogger) Info(msg string, fields ...Field)        {}
func (m MockLogger) Fatal(msg string, fields ...Field)       {}
func (m MockLogger) With(fields ...zapcore.Field) Logger { return nil }
