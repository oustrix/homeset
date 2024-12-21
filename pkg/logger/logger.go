package logger

import (
	"io"
	"log/slog"
	"strings"
)

// Config provide data to Configure.
type Config struct {
	Writer io.Writer
	Level  string
}

// Configure configures global slog logger.
func Configure(config Config) {
	var lvl slog.Level

	switch strings.ToLower(config.Level) {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(config.Writer, &slog.HandlerOptions{
		Level: lvl,
	})

	slog.SetDefault(slog.New(handler))
}
