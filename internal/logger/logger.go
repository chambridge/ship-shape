// Package logger provides structured logging for Ship Shape.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Level represents logging severity levels.
type Level int

const (
	// LevelDebug is the debug level (most verbose).
	LevelDebug Level = iota
	// LevelInfo is the info level (default).
	LevelInfo
	// LevelWarn is the warning level.
	LevelWarn
	// LevelError is the error level (least verbose).
	LevelError
)

// Config holds logger configuration.
type Config struct {
	// Output is where to write logs (defaults to os.Stderr)
	Output io.Writer
	// Format is the output format ("text" or "json")
	Format string
	// Level is the minimum log level to output
	Level Level
	// NoColor disables colored output for text format
	NoColor bool
}

// Logger wraps slog.Logger with additional Ship Shape-specific functionality.
type Logger struct {
	*slog.Logger
	config Config
}

var (
	// Default is the default logger instance
	Default *Logger
)

func init() {
	// Initialize with default configuration
	Default = New(Config{
		Level:   LevelInfo,
		Format:  "text",
		Output:  os.Stderr,
		NoColor: false,
	})
}

// New creates a new logger with the given configuration.
func New(cfg Config) *Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}

	// Convert our Level to slog.Level
	var slogLevel slog.Level

	switch cfg.Level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	// Create handler based on format
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(cfg.Output, opts)
	default: // "text"
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
		config: cfg,
	}
}

// SetDefault sets the default logger instance.
func SetDefault(l *Logger) {
	Default = l
	slog.SetDefault(l.Logger)
}

// Debug logs a debug message with optional key-value pairs.
func Debug(msg string, args ...any) {
	Default.Debug(msg, args...)
}

// Info logs an info message with optional key-value pairs.
func Info(msg string, args ...any) {
	Default.Info(msg, args...)
}

// Warn logs a warning message with optional key-value pairs.
func Warn(msg string, args ...any) {
	Default.Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs.
func Error(msg string, args ...any) {
	Default.Error(msg, args...)
}

// DebugContext logs a debug message with context and optional key-value pairs.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Default.DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context and optional key-value pairs.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Default.InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and optional key-value pairs.
func WarnContext(ctx context.Context, msg string, args ...any) {
	Default.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context and optional key-value pairs.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Default.ErrorContext(ctx, msg, args...)
}

// With creates a new logger with the given attributes added to every log entry.
func With(args ...any) *Logger {
	return &Logger{
		Logger: Default.With(args...),
		config: Default.config,
	}
}

// WithGroup creates a new logger with the given group name.
func WithGroup(name string) *Logger {
	return &Logger{
		Logger: Default.WithGroup(name),
		config: Default.config,
	}
}
