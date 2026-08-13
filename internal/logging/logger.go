package logging

import (
	"context"
	"log/slog"
)

type Logger struct {
	inner *slog.Logger
}

func New(inner *slog.Logger) *Logger {
	return &Logger{inner: inner}
}

func (l *Logger) withRequestID(ctx context.Context, args []any) []any {
	if requestID := RequestID(ctx); requestID != "" {
		return append([]any{"request_id", requestID}, args...)
	}

	return args
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.inner.Info(msg, l.withRequestID(ctx, args)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.inner.Warn(msg, l.withRequestID(ctx, args)...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.inner.Error(msg, l.withRequestID(ctx, args)...)
}
