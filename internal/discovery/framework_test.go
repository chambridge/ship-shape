package discovery

import (
	"testing"

	"github.com/chambridge/ship-shape/internal/testutil"
	"github.com/chambridge/ship-shape/pkg/types"
)

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestFrameworkDetector_Detect(t *testing.T) {
	t.Run("detects Go built-in testing framework", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create Go test files
		testutil.WriteFile(t, dir, "main_test.go", "package main\nimport \"testing\"")
		testutil.WriteFile(t, dir, "util_test.go", "package util\nimport \"testing\"")

		walker := NewWalker(dir)
		detector := NewFrameworkDetector(dir, walker)

		frameworks, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should find Go's testing framework
		found := false
		for _, fw := range frameworks {
			if fw.Name == "testing" && fw.Language == types.LanguageGo {
				found = true

				if fw.Type != types.FrameworkTypeTest {
					t.Errorf("testing framework type = %v, want %v", fw.Type, types.FrameworkTypeTest)
				}
			}
		}

		if !found {
			t.Error("Go testing framework not detected")
		}
	})

	t.Run("detects Python unittest framework", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create Python unittest files
		testutil.WriteFile(t, dir, "test_example.py", "import unittest\n\nclass TestExample(unittest.TestCase):\n    pass")

		walker := NewWalker(dir)
		detector := NewFrameworkDetector(dir, walker)

		frameworks, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should find Python's unittest framework
		found := false
		for _, fw := range frameworks {
			if fw.Name == "unittest" && fw.Language == types.LanguagePython {
				found = true

				if fw.Type != types.FrameworkTypeTest {
					t.Errorf("unittest framework type = %v, want %v", fw.Type, types.FrameworkTypeTest)
				}
			}
		}

		if !found {
			t.Error("Python unittest framework not detected")
		}
	})

	t.Run("combines manifest and built-in frameworks", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create Go test files
		testutil.WriteFile(t, dir, "main_test.go", "package main\nimport \"testing\"")

		// Create go.mod with testify
		goMod := `module github.com/example/app

go 1.21

require github.com/stretchr/testify v1.8.4
`
		testutil.WriteFile(t, dir, "go.mod", goMod)

		walker := NewWalker(dir)
		detector := NewFrameworkDetector(dir, walker)

		frameworks, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should find both testing and testify
		hasBuiltin := false
		hasTestify := false

		for _, fw := range frameworks {
			if fw.Name == "testing" && fw.Language == types.LanguageGo {
				hasBuiltin = true
			}

			if fw.Name == "testify" && fw.Language == types.LanguageGo {
				hasTestify = true
			}
		}

		if !hasBuiltin {
			t.Error("Built-in testing framework not detected")
		}

		if !hasTestify {
			t.Error("testify from go.mod not detected")
		}
	})

	t.Run("detects JavaScript frameworks from package.json", func(t *testing.T) {
		dir := testutil.TempDir(t)

		packageJSON := `{
			"name": "my-app",
			"devDependencies": {
				"jest": "^29.0.0",
				"eslint": "^8.0.0"
			}
		}`
		testutil.WriteFile(t, dir, "package.json", packageJSON)

		walker := NewWalker(dir)
		detector := NewFrameworkDetector(dir, walker)

		frameworks, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should find jest and eslint
		hasJest := false
		hasEslint := false

		for _, fw := range frameworks {
			if fw.Name == "jest" {
				hasJest = true
			}

			if fw.Name == "eslint" {
				hasEslint = true
			}
		}

		if !hasJest {
			t.Error("jest not detected from package.json")
		}

		if !hasEslint {
			t.Error("eslint not detected from package.json")
		}
	})

	t.Run("deduplicates frameworks", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create scenario where same framework might be detected multiple ways
		// (This is hypothetical - in practice, built-in frameworks won't be in manifests)
		testutil.WriteFile(t, dir, "main_test.go", "package main")

		walker := NewWalker(dir)
		detector := NewFrameworkDetector(dir, walker)

		frameworks, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Count occurrences of each framework
		counts := make(map[string]int)
		for _, fw := range frameworks {
			key := fw.Name + string(fw.Language)
			counts[key]++
		}

		// No framework should appear more than once
		for key, count := range counts {
			if count > 1 {
				t.Errorf("Framework %s appears %d times, expected 1", key, count)
			}
		}
	})

	t.Run("handles repository with no frameworks", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create non-test files
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "README.md", "# Project")

		walker := NewWalker(dir)
		detector := NewFrameworkDetector(dir, walker)

		frameworks, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should return empty list, not error
		if len(frameworks) != 0 {
			t.Errorf("Expected 0 frameworks for repo with no test files, got %d", len(frameworks))
		}
	})
}

func TestIsGoTestFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"main_test.go", true},
		{"util_test.go", true},
		{"pkg_test.go", true},
		{"main.go", false},
		{"test.go", false},
		{"_test.go", false}, // needs prefix to be valid
		{"test", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isGoTestFile(tt.name)
			if got != tt.want {
				t.Errorf("isGoTestFile(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsPythonTestFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"test_example.py", true},
		{"test_util.py", true},
		{"example_test.py", true},
		{"util_test.py", true},
		{"test.py", false},
		{"example.py", false},
		{"test_example", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPythonTestFile(tt.name)
			if got != tt.want {
				t.Errorf("isPythonTestFile(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHasUnittestImport(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "import unittest",
			content: "import unittest\n\nclass TestExample(unittest.TestCase):\n    pass",
			want:    true,
		},
		{
			name:    "from unittest",
			content: "from unittest import TestCase\n\nclass TestExample(TestCase):\n    pass",
			want:    true,
		},
		{
			name:    "no unittest import",
			content: "import os\nimport sys\n\ndef test_example():\n    pass",
			want:    false,
		},
		{
			name:    "empty file",
			content: "",
			want:    false,
		},
		{
			name:    "unittest in comment",
			content: "# This uses unittest framework\nimport pytest",
			want:    false, // No actual unittest import
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasUnittestImport(tt.content)
			if got != tt.want {
				t.Errorf("hasUnittestImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeduplicateFrameworks(t *testing.T) {
	t.Run("removes duplicates", func(t *testing.T) {
		frameworks := []types.Framework{
			{Name: "jest", Language: types.LanguageJavaScript, Type: types.FrameworkTypeTest},
			{Name: "jest", Language: types.LanguageJavaScript, Type: types.FrameworkTypeTest},
			{Name: "pytest", Language: types.LanguagePython, Type: types.FrameworkTypeTest},
			{Name: "pytest", Language: types.LanguagePython, Type: types.FrameworkTypeTest},
		}

		result := deduplicateFrameworks(frameworks)

		if len(result) != 2 {
			t.Errorf("Expected 2 unique frameworks, got %d", len(result))
		}

		// Verify we have one jest and one pytest
		hasJest := false
		hasPytest := false

		for _, fw := range result {
			if fw.Name == "jest" {
				hasJest = true
			}

			if fw.Name == "pytest" {
				hasPytest = true
			}
		}

		if !hasJest {
			t.Error("jest not found after deduplication")
		}

		if !hasPytest {
			t.Error("pytest not found after deduplication")
		}
	})

	t.Run("keeps frameworks with same name but different languages", func(t *testing.T) {
		// Hypothetical: same tool name in different ecosystems
		frameworks := []types.Framework{
			{Name: "mock", Language: types.LanguageGo, Type: types.FrameworkTypeTest},
			{Name: "mock", Language: types.LanguagePython, Type: types.FrameworkTypeTest},
		}

		result := deduplicateFrameworks(frameworks)

		if len(result) != 2 {
			t.Errorf("Expected 2 frameworks (different languages), got %d", len(result))
		}
	})

	t.Run("handles empty list", func(t *testing.T) {
		frameworks := []types.Framework{}

		result := deduplicateFrameworks(frameworks)

		if len(result) != 0 {
			t.Errorf("Expected empty list, got %d frameworks", len(result))
		}
	})
}
