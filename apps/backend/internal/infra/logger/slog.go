package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type Config struct {
	Level  string
	Format string
}

func New(cfg Config) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	if level == slog.LevelDebug {
		opts.AddSource = true
	}

	var handler slog.Handler
	var out io.Writer = os.Stderr

	switch strings.ToLower(cfg.Format) {
	case "text":
		handler = slog.NewTextHandler(out, opts)
	default:
		handler = slog.NewJSONHandler(out, opts)
	}

	return slog.New(handler)
}

func NewGORMLogger(l *slog.Logger, level slog.Level) logger.Interface {
	return logger.New(
		slog.NewLogLogger(l.Handler(), level),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      gormLevel(level),
			Colorful:      false,
		},
	)
}

func gormLevel(l slog.Level) logger.LogLevel {
	if l <= slog.LevelInfo {
		return logger.Info
	}
	if l <= slog.LevelWarn {
		return logger.Warn
	}
	return logger.Error
}
