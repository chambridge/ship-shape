package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Level
	}{
		{
			name: "default config",
			config: Config{
				Level:  LevelInfo,
				Format: "text",
			},
			want: LevelInfo,
		},
		{
			name: "debug level",
			config: Config{
				Level:  LevelDebug,
				Format: "text",
			},
			want: LevelDebug,
		},
		{
			name: "error level",
			config: Config{
				Level:  LevelError,
				Format: "json",
			},
			want: LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.config.Output = buf

			logger := New(tt.config)
			if logger == nil {
				t.Fatal("New() returned nil")
			}

			if logger.config.Level != tt.want {
				t.Errorf("New() level = %v, want %v", logger.config.Level, tt.want)
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		logFunc       func(string, ...any)
		name          string
		expectedLevel string
		logLevel      Level
		shouldLog     bool
	}{
		{
			name:          "debug logs at debug level",
			logLevel:      LevelDebug,
			logFunc:       func(msg string, args ...any) { Debug(msg, args...) },
			shouldLog:     true,
			expectedLevel: "DEBUG",
		},
		{
			name:          "debug does not log at info level",
			logLevel:      LevelInfo,
			logFunc:       func(msg string, args ...any) { Debug(msg, args...) },
			shouldLog:     false,
			expectedLevel: "DEBUG",
		},
		{
			name:          "info logs at info level",
			logLevel:      LevelInfo,
			logFunc:       func(msg string, args ...any) { Info(msg, args...) },
			shouldLog:     true,
			expectedLevel: "INFO",
		},
		{
			name:          "info does not log at warn level",
			logLevel:      LevelWarn,
			logFunc:       func(msg string, args ...any) { Info(msg, args...) },
			shouldLog:     false,
			expectedLevel: "INFO",
		},
		{
			name:          "warn logs at warn level",
			logLevel:      LevelWarn,
			logFunc:       func(msg string, args ...any) { Warn(msg, args...) },
			shouldLog:     true,
			expectedLevel: "WARN",
		},
		{
			name:          "error logs at all levels",
			logLevel:      LevelDebug,
			logFunc:       func(msg string, args ...any) { Error(msg, args...) },
			shouldLog:     true,
			expectedLevel: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			cfg := Config{
				Output: buf,
				Format: "text",
				Level:  tt.logLevel,
			}

			logger := New(cfg)
			SetDefault(logger)

			tt.logFunc("test message", "key", "value")

			output := buf.String()

			if tt.shouldLog {
				verifyLogOutput(t, output, tt.expectedLevel)
			} else {
				verifyNoLogOutput(t, output)
			}
		})
	}
}

func verifyLogOutput(t *testing.T, output, expectedLevel string) {
	t.Helper()

	if output == "" {
		t.Error("Expected log output, got empty string")
	}

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log to contain 'test message', got: %s", output)
	}

	if !strings.Contains(output, expectedLevel) {
		t.Errorf("Expected log to contain level '%s', got: %s", expectedLevel, output)
	}

	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log to contain 'key=value', got: %s", output)
	}
}

func verifyNoLogOutput(t *testing.T, output string) {
	t.Helper()

	if output != "" {
		t.Errorf("Expected no log output, got: %s", output)
	}
}

func TestJSONFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := Config{
		Output: buf,
		Format: "json",
		Level:  LevelInfo,
	}

	logger := New(cfg)
	SetDefault(logger)

	Info("test message", "key", "value", "number", 42)

	output := buf.String()
	if output == "" {
		t.Fatal("Expected JSON log output, got empty string")
	}

	// Parse JSON to validate structure
	var logEntry map[string]any
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log: %v\nOutput: %s", err, output)
	}

	// Verify expected fields
	if msg, ok := logEntry["msg"].(string); !ok || msg != "test message" {
		t.Errorf("Expected msg='test message', got: %v", logEntry["msg"])
	}

	if level, ok := logEntry["level"].(string); !ok || level != "INFO" {
		t.Errorf("Expected level='INFO', got: %v", logEntry["level"])
	}

	if key, ok := logEntry["key"].(string); !ok || key != "value" {
		t.Errorf("Expected key='value', got: %v", logEntry["key"])
	}

	if number, ok := logEntry["number"].(float64); !ok || number != 42 {
		t.Errorf("Expected number=42, got: %v", logEntry["number"])
	}
}

func TestContextLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := Config{
		Output: buf,
		Format: "text",
		Level:  LevelDebug,
	}

	logger := New(cfg)
	SetDefault(logger)

	ctx := context.Background()

	// Test each context-aware function
	DebugContext(ctx, "debug message")
	InfoContext(ctx, "info message")
	WarnContext(ctx, "warn message")
	ErrorContext(ctx, "error message")

	output := buf.String()

	expectedMessages := []string{"debug message", "info message", "warn message", "error message"}
	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected log to contain '%s', got: %s", msg, output)
		}
	}
}

func TestWith(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := Config{
		Output: buf,
		Format: "text",
		Level:  LevelInfo,
	}

	logger := New(cfg)
	SetDefault(logger)

	// Create logger with additional attributes
	loggerWithAttrs := With("component", "test", "version", "1.0")
	SetDefault(loggerWithAttrs)

	Info("test message")

	output := buf.String()

	if !strings.Contains(output, "component=test") {
		t.Errorf("Expected log to contain 'component=test', got: %s", output)
	}

	if !strings.Contains(output, "version=1.0") {
		t.Errorf("Expected log to contain 'version=1.0', got: %s", output)
	}
}

func TestWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := Config{
		Output: buf,
		Format: "json",
		Level:  LevelInfo,
	}

	logger := New(cfg)
	SetDefault(logger)

	// Create logger with group
	loggerWithGroup := WithGroup("request")
	SetDefault(loggerWithGroup)

	Info("test message", "method", "GET", "path", "/api/v1")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected JSON log output, got empty string")
	}

	// Parse JSON to validate structure
	var logEntry map[string]any
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log: %v\nOutput: %s", err, output)
	}

	// Verify group exists
	request, ok := logEntry["request"].(map[string]any)
	if !ok {
		t.Errorf("Expected 'request' group in log entry, got: %v", logEntry)

		return
	}

	if method, ok := request["method"].(string); !ok || method != "GET" {
		t.Errorf("Expected request.method='GET', got: %v", request["method"])
	}

	if path, ok := request["path"].(string); !ok || path != "/api/v1" {
		t.Errorf("Expected request.path='/api/v1', got: %v", request["path"])
	}
}

func TestSetDefault(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	cfg1 := Config{Output: buf1, Format: "text", Level: LevelInfo}
	cfg2 := Config{Output: buf2, Format: "text", Level: LevelDebug}

	logger1 := New(cfg1)
	SetDefault(logger1)

	Info("message from logger1")

	// Switch default logger
	logger2 := New(cfg2)
	SetDefault(logger2)

	Debug("message from logger2")

	// Verify first logger received first message
	output1 := buf1.String()
	if !strings.Contains(output1, "message from logger1") {
		t.Errorf("Expected logger1 to contain 'message from logger1', got: %s", output1)
	}

	// Verify second logger received second message
	output2 := buf2.String()
	if !strings.Contains(output2, "message from logger2") {
		t.Errorf("Expected logger2 to contain 'message from logger2', got: %s", output2)
	}

	// Verify first logger did NOT receive second message
	if strings.Contains(output1, "message from logger2") {
		t.Errorf("Expected logger1 NOT to contain 'message from logger2', got: %s", output1)
	}
}

func TestDefaultLogger(t *testing.T) {
	// Test that Default logger is initialized
	if Default == nil {
		t.Fatal("Default logger should be initialized")
	}

	// Test that Default logger can log
	buf := &bytes.Buffer{}
	cfg := Config{
		Output: buf,
		Format: "text",
		Level:  LevelInfo,
	}

	logger := New(cfg)
	SetDefault(logger)

	Info("test default logger")

	output := buf.String()
	if !strings.Contains(output, "test default logger") {
		t.Errorf("Expected default logger to work, got: %s", output)
	}
}

func TestNilOutput(t *testing.T) {
	cfg := Config{
		Output: nil, // Should default to os.Stderr
		Format: "text",
		Level:  LevelInfo,
	}

	logger := New(cfg)
	if logger == nil {
		t.Fatal("New() with nil output should not return nil logger")
	}

	if logger.config.Output == nil {
		t.Error("Logger output should be set to default when nil is provided")
	}
}

func TestLevelConversion(t *testing.T) {
	tests := []struct {
		name      string
		level     Level
		wantSlog  slog.Level
		shouldLog bool
	}{
		{
			name:      "debug level",
			level:     LevelDebug,
			wantSlog:  slog.LevelDebug,
			shouldLog: true,
		},
		{
			name:      "info level",
			level:     LevelInfo,
			wantSlog:  slog.LevelInfo,
			shouldLog: true,
		},
		{
			name:      "warn level",
			level:     LevelWarn,
			wantSlog:  slog.LevelWarn,
			shouldLog: true,
		},
		{
			name:      "error level",
			level:     LevelError,
			wantSlog:  slog.LevelError,
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			cfg := Config{
				Output: buf,
				Format: "text",
				Level:  tt.level,
			}

			logger := New(cfg)
			SetDefault(logger)

			// Log at the configured level
			switch tt.level {
			case LevelDebug:
				Debug("test")
			case LevelInfo:
				Info("test")
			case LevelWarn:
				Warn("test")
			case LevelError:
				Error("test")
			}

			output := buf.String()
			if tt.shouldLog && output == "" {
				t.Errorf("Expected log output at level %v, got empty string", tt.level)
			}
		})
	}
}
