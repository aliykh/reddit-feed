package log

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogType uint8

const (
	ZapLogger LogType = iota
	Mock
)

type Factory struct {
	logger Logger
	tr     func()
}

func NewFactory(t LogType, levelInfo string) *Factory {
	// for now, zap log should be enough

	zapLoggerMaker := func() *Factory {
		logger, tr := newZapLogger(parseLevel(levelInfo))
		return &Factory{
			logger: logger,
			tr:     tr,
		}
	}

	switch t {

	case ZapLogger:
		return zapLoggerMaker()

	case Mock:
		logger := MockLogger{}
		return &Factory{
			logger: logger,
			tr:     func() {},
		}

	default:
		return zapLoggerMaker()
	}

}

func (f *Factory) Default() Logger {
	return f.logger
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (f Factory) For(ctx context.Context) Logger {

	if span := opentracing.SpanFromContext(ctx); span != nil {

		logger := spanLogger{span: span, logger: f.logger}

		if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
			logger.spanFields = []zapcore.Field{
				zap.String("trace_id", jaegerCtx.TraceID().String()),
				zap.String("span_id", jaegerCtx.SpanID().String()),
			}
		}

		return logger
	}
	return f.Default()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (f Factory) With(fields ...zapcore.Field) Factory {
	return Factory{logger: f.logger.With(fields...)}
}


func parseLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

const (
	// LevelDebug ...
	LevelDebug = "debug"
	// LevelInfo ...
	LevelInfo = "info"
	// LevelWarn ...
	LevelWarn = "warn"
	// LevelError ...
	LevelError = "error"
	// LevelPanic ...
	LevelPanic = "panic"
	// LevelFatal ...
	LevelFatal = "fatal"
)