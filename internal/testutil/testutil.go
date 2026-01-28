// Package testutil provides common test utilities and helpers for Ship Shape tests.
package testutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// CaptureOutput captures stdout and stderr output during test execution.
// Returns stdout, stderr, and any error that occurred.
func CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	// Save original stdout/stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create pipes
	rOut, wOut, errOut := os.Pipe()
	rErr, wErr, errErr := os.Pipe()

	if errOut != nil || errErr != nil {
		t.Fatalf("Failed to create pipes: %v, %v", errOut, errErr)
	}

	// Capture output in goroutines
	outC := make(chan string, 1)
	errC := make(chan string, 1)

	go func() {
		var buf bytes.Buffer
		//nolint:errcheck // io.Copy error handled by returning partial output
		_, _ = io.Copy(&buf, rOut) // #nosec G104 -- errors don't prevent output capture
		outC <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		//nolint:errcheck // io.Copy error handled by returning partial output
		_, _ = io.Copy(&buf, rErr) // #nosec G104 -- errors don't prevent output capture
		errC <- buf.String()
	}()

	// Replace stdout/stderr AFTER starting goroutines
	os.Stdout = wOut
	os.Stderr = wErr

	// Run function
	fn()

	// Close writers to signal EOF to readers
	//nolint:errcheck // Close errors are not critical in test cleanup
	_ = wOut.Close() // #nosec G104 -- Close errors are not critical in test cleanup
	//nolint:errcheck // Close errors are not critical in test cleanup
	_ = wErr.Close() // #nosec G104 -- Close errors are not critical in test cleanup

	// Restore original stdout/stderr immediately
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Get captured output from goroutines with buffered channels
	stdout = <-outC
	stderr = <-errC

	// Close readers to free resources
	//nolint:errcheck // Close errors are not critical in test cleanup
	_ = rOut.Close() // #nosec G104 -- Close errors are not critical in test cleanup
	//nolint:errcheck // Close errors are not critical in test cleanup
	_ = rErr.Close() // #nosec G104 -- Close errors are not critical in test cleanup

	return stdout, stderr
}

// TempDir creates a temporary directory for testing and returns the path.
// The directory is automatically cleaned up when the test finishes.
func TempDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "shipshape-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("Failed to remove temp dir %s: %v", dir, err)
		}
	})

	return dir
}

// WriteFile writes content to a file in the given directory.
// Creates parent directories as needed.
func WriteFile(t *testing.T, dir, path, content string) string {
	t.Helper()

	fullPath := filepath.Join(dir, path)

	// Create parent directories
	// #nosec G301 -- test utilities need standard permissions
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("Failed to create directories for %s: %v", fullPath, err)
	}

	// Write file
	// #nosec G306 -- test utilities need standard permissions
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to write file %s: %v", fullPath, err)
	}

	return fullPath
}

// Chdir changes the current working directory for the duration of the test.
// The original directory is restored when the test finishes.
func Chdir(t *testing.T, dir string) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory to %s: %v", dir, err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Logf("Failed to restore directory to %s: %v", oldDir, err)
		}
	})
}

// SetEnv sets an environment variable for the duration of the test.
// The original value is restored when the test finishes.
func SetEnv(t *testing.T, key, value string) {
	t.Helper()

	oldValue, exists := os.LookupEnv(key)

	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set env var %s: %v", key, err)
	}

	t.Cleanup(func() {
		if exists {
			if err := os.Setenv(key, oldValue); err != nil {
				t.Logf("Failed to restore env var %s: %v", key, err)
			}
		} else {
			if err := os.Unsetenv(key); err != nil {
				t.Logf("Failed to unset env var %s: %v", key, err)
			}
		}
	})
}

// UnsetEnv unsets an environment variable for the duration of the test.
// The original value is restored when the test finishes.
func UnsetEnv(t *testing.T, key string) {
	t.Helper()

	oldValue, exists := os.LookupEnv(key)

	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed to unset env var %s: %v", key, err)
	}

	t.Cleanup(func() {
		if exists {
			if err := os.Setenv(key, oldValue); err != nil {
				t.Logf("Failed to restore env var %s: %v", key, err)
			}
		}
	})
}
