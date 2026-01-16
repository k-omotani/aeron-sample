package logging

import (
	"log/slog"
	"os"
)

// Config defines logging configuration
type Config struct {
	Level  slog.Level
	Format string // "json" or "text"
}

// NewLogger creates a configured slog.Logger
func NewLogger(cfg *Config) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// DefaultConfig returns development-friendly defaults
func DefaultConfig() *Config {
	return &Config{
		Level:  slog.LevelDebug,
		Format: "text",
	}
}

// ParseLevel parses a string log level
func ParseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
