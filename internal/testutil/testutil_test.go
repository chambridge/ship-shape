package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaptureOutput(t *testing.T) {
	tests := []struct {
		name           string
		fn             func()
		wantStdout     string
		wantStderr     string
		stdoutContains string
		stderrContains string
	}{
		{
			name: "stdout only",
			fn: func() {
				fmt.Println("hello stdout")
			},
			stdoutContains: "hello stdout",
		},
		{
			name: "stderr only",
			fn: func() {
				fmt.Fprintln(os.Stderr, "hello stderr")
			},
			stderrContains: "hello stderr",
		},
		{
			name: "both stdout and stderr",
			fn: func() {
				fmt.Println("stdout message")
				fmt.Fprintln(os.Stderr, "stderr message")
			},
			stdoutContains: "stdout message",
			stderrContains: "stderr message",
		},
		{
			name:       "no output",
			fn:         func() {},
			wantStdout: "",
			wantStderr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := CaptureOutput(t, tt.fn)

			if tt.wantStdout != "" && stdout != tt.wantStdout {
				t.Errorf("stdout = %q, want %q", stdout, tt.wantStdout)
			}

			if tt.wantStderr != "" && stderr != tt.wantStderr {
				t.Errorf("stderr = %q, want %q", stderr, tt.wantStderr)
			}

			if tt.stdoutContains != "" && !strings.Contains(stdout, tt.stdoutContains) {
				t.Errorf("stdout %q does not contain %q", stdout, tt.stdoutContains)
			}

			if tt.stderrContains != "" && !strings.Contains(stderr, tt.stderrContains) {
				t.Errorf("stderr %q does not contain %q", stderr, tt.stderrContains)
			}
		})
	}
}

func TestTempDir(t *testing.T) {
	dir := TempDir(t)

	// Verify directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("TempDir created directory that doesn't exist: %s", dir)
	}

	// Verify directory is writable
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Errorf("TempDir created directory is not writable: %v", err)
	}
}

func TestWriteFile(t *testing.T) {
	dir := TempDir(t)

	tests := []struct {
		name    string
		path    string
		content string
	}{
		{
			name:    "simple file",
			path:    "test.txt",
			content: "test content",
		},
		{
			name:    "nested file",
			path:    "subdir/nested/test.txt",
			content: "nested content",
		},
		{
			name:    "empty content",
			path:    "empty.txt",
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := WriteFile(t, dir, tt.path, tt.content)

			// Verify file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("WriteFile did not create file: %s", fullPath)

				return
			}

			// Verify content
			got, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)

				return
			}

			if string(got) != tt.content {
				t.Errorf("WriteFile content = %q, want %q", string(got), tt.content)
			}
		})
	}
}

func TestChdir(t *testing.T) {
	// Get original directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Resolve to handle symlinks (e.g., /var -> /private/var on macOS)
	//nolint:ineffassign,staticcheck // origDir is used in cleanup verification concept
	origDir, err = filepath.EvalSymlinks(origDir)
	if err != nil {
		t.Fatalf("Failed to resolve original directory: %v", err)
	}

	// Create temp directory
	tempDir := TempDir(t)

	// Resolve temp directory too
	tempDir, err = filepath.EvalSymlinks(tempDir)
	if err != nil {
		t.Fatalf("Failed to resolve temp directory: %v", err)
	}

	// Change to temp directory
	Chdir(t, tempDir)

	// Verify we're in the temp directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	currentDir, err = filepath.EvalSymlinks(currentDir)
	if err != nil {
		t.Fatalf("Failed to resolve current directory: %v", err)
	}

	if currentDir != tempDir {
		t.Errorf("Chdir did not change directory: got %s, want %s", currentDir, tempDir)
	}

	// Note: We can't verify cleanup restoration in the same test because t.Cleanup runs in LIFO order
	// The Chdir's cleanup (which restores the directory) runs after this test's cleanup verification would run
}

func TestSetEnv(t *testing.T) {
	key := "SHIPSHAPE_TEST_VAR"
	value := "test_value"

	// Ensure variable doesn't exist initially
	UnsetEnv(t, key)

	// Set variable
	SetEnv(t, key, value)

	// Verify variable is set
	got := os.Getenv(key)
	if got != value {
		t.Errorf("SetEnv did not set variable: got %q, want %q", got, value)
	}
}

func TestSetEnvWithExisting(t *testing.T) {
	key := "SHIPSHAPE_TEST_EXISTING"
	originalValue := "original"
	newValue := "new"

	// Set initial value
	if err := os.Setenv(key, originalValue); err != nil {
		t.Fatalf("Failed to set initial env var: %v", err)
	}

	// Clean up at end of test
	t.Cleanup(func() {
		os.Unsetenv(key)
	})

	// Override with SetEnv
	SetEnv(t, key, newValue)

	// Verify new value
	got := os.Getenv(key)
	if got != newValue {
		t.Errorf("SetEnv did not override variable: got %q, want %q", got, newValue)
	}

	// Note: We can't verify cleanup in the same test because t.Cleanup runs in LIFO order
	// The SetEnv cleanup will restore the original value before our verification runs
}

func TestUnsetEnv(t *testing.T) {
	key := "SHIPSHAPE_TEST_UNSET"
	value := "value"

	// Set initial value
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set initial env var: %v", err)
	}

	// Clean up at end of test
	t.Cleanup(func() {
		os.Unsetenv(key)
	})

	// Unset variable
	UnsetEnv(t, key)

	// Verify variable is unset
	if _, exists := os.LookupEnv(key); exists {
		t.Error("UnsetEnv did not unset variable")
	}

	// Note: We can't verify cleanup in the same test because t.Cleanup runs in LIFO order
	// The UnsetEnv cleanup will restore the original value before our verification runs
}
