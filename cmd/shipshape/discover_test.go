package main

import (
	"encoding/json"
	"testing"

	"github.com/chambridge/ship-shape/internal/testutil"
	"github.com/chambridge/ship-shape/pkg/types"
)

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestDiscoverCommand(t *testing.T) {
	t.Run("discovers Go repository", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create Go repository structure
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "main_test.go", "package main\nimport \"testing\"")
		testutil.WriteFile(t, dir, "go.mod", "module example.com/app\n\ngo 1.21")

		// Run discover command
		rootCmd.SetArgs([]string{"discover", dir})

		stdout, stderr := testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Verify output contains key information
		if len(stdout) == 0 {
			t.Error("Expected output, got empty string")
		}

		// Should mention repository path
		if !contains(stdout, dir) {
			t.Errorf("Output should contain repository path %s", dir)
		}

		// Should detect Go language
		if !contains(stdout, "Go") {
			t.Error("Output should mention Go language")
		}

		// Should detect testing framework
		if !contains(stdout, "testing") {
			t.Error("Output should mention testing framework")
		}

		// Stderr should be empty (no errors)
		if len(stderr) > 0 {
			t.Logf("Unexpected stderr: %s", stderr)
		}
	})

	t.Run("discovers JavaScript repository", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create JavaScript repository structure
		packageJSON := `{
			"name": "my-app",
			"devDependencies": {
				"jest": "^29.0.0",
				"eslint": "^8.0.0"
			}
		}`
		testutil.WriteFile(t, dir, "package.json", packageJSON)
		testutil.WriteFile(t, dir, "index.js", "console.log('hello');")
		testutil.WriteFile(t, dir, "app.test.js", "test('example', () => {});")

		// Run discover command
		rootCmd.SetArgs([]string{"discover", dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Should detect JavaScript
		if !contains(stdout, "JavaScript") {
			t.Error("Output should mention JavaScript language")
		}

		// Should detect jest
		if !contains(stdout, "jest") {
			t.Error("Output should mention jest framework")
		}

		// Should detect eslint
		if !contains(stdout, "eslint") {
			t.Error("Output should mention eslint")
		}
	})

	t.Run("JSON output format", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create simple repository
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "main_test.go", "package main\nimport \"testing\"")

		// Run discover command with --json flag
		rootCmd.SetArgs([]string{"discover", "--json", dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Parse JSON output
		var repo types.Repository

		err := json.Unmarshal([]byte(stdout), &repo)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, stdout)
		}

		// Verify structure
		if repo.Path != dir {
			t.Errorf("Repository path = %s, want %s", repo.Path, dir)
		}

		if len(repo.Languages) == 0 {
			t.Error("Expected languages to be detected")
		}

		// Should detect Go
		hasGo := false
		for _, lang := range repo.Languages {
			if lang.Language == types.LanguageGo {
				hasGo = true
				break
			}
		}

		if !hasGo {
			t.Error("Go language not detected in JSON output")
		}
	})

	t.Run("handles non-existent directory", func(t *testing.T) {
		// Run discover command on non-existent directory
		rootCmd.SetArgs([]string{"discover", "/nonexistent/path/123456"})

		_, _ = testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err == nil {
				t.Error("Expected error for non-existent directory")
			}
		})
	})

	t.Run("uses current directory when no args", func(t *testing.T) {
		// Create temp directory and change to it
		tempDir := testutil.TempDir(t)
		testutil.Chdir(t, tempDir)

		// Create files in current directory
		testutil.WriteFile(t, ".", "test.go", "package test")

		// Run discover command without arguments
		rootCmd.SetArgs([]string{"discover"})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Should detect Go
		if !contains(stdout, "Go") {
			t.Error("Should detect Go in current directory")
		}
	})

	t.Run("discovers multi-language repository", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create multi-language repository
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "script.py", "print('hello')")
		testutil.WriteFile(t, dir, "app.js", "console.log('hello');")

		// Run discover command
		rootCmd.SetArgs([]string{"discover", dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Should detect all languages
		if !contains(stdout, "Go") {
			t.Error("Should detect Go")
		}

		if !contains(stdout, "Python") {
			t.Error("Should detect Python")
		}

		if !contains(stdout, "JavaScript") {
			t.Error("Should detect JavaScript")
		}
	})

	t.Run("handles empty repository", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create empty repository (just README)
		testutil.WriteFile(t, dir, "README.md", "# Empty Project")

		// Run discover command
		rootCmd.SetArgs([]string{"discover", dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Should handle gracefully
		if !contains(stdout, "None detected") {
			t.Error("Should indicate no languages/frameworks detected")
		}
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf returns the index of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
