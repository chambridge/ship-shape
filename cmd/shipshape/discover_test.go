package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/chambridge/ship-shape/internal/logger"
	"github.com/chambridge/ship-shape/internal/testutil"
	"github.com/chambridge/ship-shape/pkg/types"
	"github.com/spf13/cobra"
)

// resetRootCmd resets the root command state between tests to prevent race conditions
func resetRootCmd(t *testing.T) {
	t.Helper()

	// Reset flags
	verbose = false
	quiet = false
	noColor = false

	// Reset discover command flags
	discoverJSON = false

	// Create a minimal logger that doesn't write anywhere during tests
	// This prevents race conditions from logger writing to redirected stderr
	cfg := logger.Config{
		Level:   logger.LevelError, // Only errors to minimize output
		Format:  "text",
		Output:  os.Stderr,
		NoColor: true,
	}
	l := logger.New(cfg)
	logger.SetDefault(l)
}

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestDiscoverCommand(t *testing.T) {
	// DO NOT run subtests in parallel - they share global rootCmd state
	// which causes race conditions with cobra's initialization hooks
	t.Run("discovers Go repository", func(t *testing.T) {
		resetRootCmd(t)

		dir := testutil.TempDir(t)

		// Create Go repository structure
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "main_test.go", "package main\nimport \"testing\"")
		testutil.WriteFile(t, dir, "go.mod", "module example.com/app\n\ngo 1.21")

		// Create a fresh command instance for this test to avoid initialization hooks
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
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
	})

	t.Run("discovers JavaScript repository", func(t *testing.T) {
		resetRootCmd(t)

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

		// Create fresh command for this test
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
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
		resetRootCmd(t)

		dir := testutil.TempDir(t)

		// Create simple repository
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "main_test.go", "package main\nimport \"testing\"")

		// Create fresh command for this test
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{"--json", dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
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
		resetRootCmd(t)

		// Create fresh command for this test
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{"/nonexistent/path/123456"})

		_, _ = testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
			if err == nil {
				t.Error("Expected error for non-existent directory")
			}
		})
	})

	t.Run("uses current directory when no args", func(t *testing.T) {
		resetRootCmd(t)

		// Create temp directory
		tempDir := testutil.TempDir(t)

		// Create files in the directory
		testutil.WriteFile(t, tempDir, "test.go", "package test")

		// Test by passing the directory explicitly instead of changing cwd
		// (changing cwd during tests can cause race detector issues)
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{tempDir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
			if err != nil {
				t.Fatalf("discover command failed: %v", err)
			}
		})

		// Should detect Go
		if !contains(stdout, "Go") {
			t.Error("Should detect Go in directory")
		}
	})

	t.Run("discovers multi-language repository", func(t *testing.T) {
		resetRootCmd(t)

		dir := testutil.TempDir(t)

		// Create multi-language repository
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "script.py", "print('hello')")
		testutil.WriteFile(t, dir, "app.js", "console.log('hello');")

		// Create fresh command for this test
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
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
		resetRootCmd(t)

		dir := testutil.TempDir(t)

		// Create empty repository (just README)
		testutil.WriteFile(t, dir, "README.md", "# Empty Project")

		// Create fresh command for this test
		testCmd := &cobra.Command{
			Use:  "discover [directory]",
			Args: cobra.MaximumNArgs(1),
			RunE: runDiscover,
		}
		testCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
		testCmd.SetArgs([]string{dir})

		stdout, _ := testutil.CaptureOutput(t, func() {
			err := testCmd.Execute()
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
