package discovery

import (
	"os"
	"path/filepath"

	"github.com/chambridge/ship-shape/pkg/types"
)

// FrameworkDetector detects testing frameworks and development tools in a repository.
type FrameworkDetector struct {
	rootPath string
	walker   *Walker
}

// NewFrameworkDetector creates a new framework detector.
func NewFrameworkDetector(rootPath string, walker *Walker) *FrameworkDetector {
	return &FrameworkDetector{
		rootPath: rootPath,
		walker:   walker,
	}
}

// Detect analyzes the repository and returns all detected frameworks.
// It combines manifest-based detection with file-based detection for built-in frameworks.
func (d *FrameworkDetector) Detect() ([]types.Framework, error) {
	var frameworks []types.Framework

	// Parse dependency manifests
	parser := NewManifestParser(d.rootPath)

	manifestFrameworks, err := parser.ParseAll()
	if err == nil {
		frameworks = append(frameworks, manifestFrameworks...)
	}

	// Detect built-in frameworks
	builtinFrameworks := d.detectBuiltinFrameworks()
	frameworks = append(frameworks, builtinFrameworks...)

	// Deduplicate frameworks by name
	frameworks = deduplicateFrameworks(frameworks)

	return frameworks, nil
}

// detectBuiltinFrameworks detects frameworks that don't require package manager entries.
func (d *FrameworkDetector) detectBuiltinFrameworks() []types.Framework {
	var frameworks []types.Framework

	// Detect Go's built-in testing package
	if d.hasGoTestFiles() {
		frameworks = append(frameworks, types.Framework{
			Name:     "testing",
			Language: types.LanguageGo,
			Type:     types.FrameworkTypeTest,
		})
	}

	// Detect Python's built-in unittest
	if d.hasPythonUnittestFiles() {
		frameworks = append(frameworks, types.Framework{
			Name:     "unittest",
			Language: types.LanguagePython,
			Type:     types.FrameworkTypeTest,
		})
	}

	return frameworks
}

// hasGoTestFiles checks if the repository contains Go test files.
func (d *FrameworkDetector) hasGoTestFiles() bool {
	hasTestFiles := false

	// Ignore errors - we're just checking for existence
	_, _ = d.walker.Walk(func(fi FileInfo) error { //nolint:errcheck // Intentionally checking existence only
		if filepath.Ext(fi.Name) == ".go" && isGoTestFile(fi.Name) {
			hasTestFiles = true
			return filepath.SkipAll // Stop walking once we find one
		}

		return nil
	})

	return hasTestFiles
}

// hasPythonUnittestFiles checks if the repository contains Python unittest files.
func (d *FrameworkDetector) hasPythonUnittestFiles() bool {
	hasUnittestFiles := false

	// Ignore errors - we're just checking for existence
	_, _ = d.walker.Walk(func(fi FileInfo) error { //nolint:errcheck // Intentionally checking existence only
		if filepath.Ext(fi.Name) == ".py" {
			// Check if file contains unittest imports
			// For now, just check for test_*.py or *_test.py naming
			if isPythonTestFile(fi.Name) {
				// Read file to check for unittest imports
				data, err := os.ReadFile(fi.Path) //nolint:gosec // Reading source files from repository
				if err == nil && hasUnittestImport(string(data)) {
					hasUnittestFiles = true
					return filepath.SkipAll
				}
			}
		}

		return nil
	})

	return hasUnittestFiles
}

// isGoTestFile checks if a filename is a Go test file.
func isGoTestFile(name string) bool {
	return len(name) > 8 && name[len(name)-8:] == "_test.go"
}

// isPythonTestFile checks if a filename follows Python test naming conventions.
func isPythonTestFile(name string) bool {
	// Must have .py extension
	if len(name) < 8 || name[len(name)-3:] != ".py" {
		return false
	}

	// Check for test_*.py or *_test.py patterns
	return (len(name) > 8 && name[:5] == "test_") || (len(name) > 8 && name[len(name)-8:] == "_test.py")
}

// hasUnittestImport checks if Python code contains unittest imports.
func hasUnittestImport(content string) bool {
	// Simple check for unittest import
	// More sophisticated: would parse the AST
	return len(content) > 0 && (containsSubstring(content, "import unittest") ||
		containsSubstring(content, "from unittest"))
}

// containsSubstring is a simple substring check.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}

// deduplicateFrameworks removes duplicate frameworks, keeping the first occurrence.
func deduplicateFrameworks(frameworks []types.Framework) []types.Framework {
	seen := make(map[string]bool)

	var result []types.Framework

	for _, fw := range frameworks {
		key := fw.Name + string(fw.Language)
		if !seen[key] {
			seen[key] = true

			result = append(result, fw)
		}
	}

	return result
}
