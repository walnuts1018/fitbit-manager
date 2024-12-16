package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	slog.Warn("Log Level Change is not supported")
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	duration := time.Since(begin)
	sql, rowsAffected := fc()

	defaultAttrs := []slog.Attr{
		slog.String("sql", sql),
		slog.Int64("rows_affected", rowsAffected),
		slog.Duration("duration", duration),
	}

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "gorm trace", append(defaultAttrs, slog.Any("error", err))...)
		return
	} else {
		slog.LogAttrs(ctx, slog.LevelDebug, "gorm trace", defaultAttrs...)
	}
}
