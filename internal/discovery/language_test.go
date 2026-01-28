package discovery

import (
	"testing"

	"github.com/chambridge/ship-shape/internal/testutil"
	"github.com/chambridge/ship-shape/pkg/types"
)

func TestLanguageDetector_Detect(t *testing.T) {
	t.Run("single language - go", func(t *testing.T) {
		dir := testutil.TempDir(t)
		testutil.WriteFile(t, dir, "main.go", "package main")
		testutil.WriteFile(t, dir, "util.go", "package util")
		testutil.WriteFile(t, dir, "helper.go", "package helper")

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(stats) != 1 {
			t.Fatalf("Detect() returned %d languages, want 1", len(stats))
		}

		lang := stats[0]
		if lang.Language != types.LanguageGo {
			t.Errorf("Language = %v, want %v", lang.Language, types.LanguageGo)
		}

		if lang.FileCount != 3 {
			t.Errorf("FileCount = %d, want 3", lang.FileCount)
		}

		if lang.Percentage != 100.0 {
			t.Errorf("Percentage = %.1f, want 100.0", lang.Percentage)
		}

		if !lang.IsPrimary {
			t.Error("Language should be marked as primary (100%)")
		}
	})

	t.Run("multi-language distribution", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// 6 Go files (60%)
		testutil.WriteFile(t, dir, "main.go", "")
		testutil.WriteFile(t, dir, "util.go", "")
		testutil.WriteFile(t, dir, "helper.go", "")
		testutil.WriteFile(t, dir, "cmd/app/main.go", "")
		testutil.WriteFile(t, dir, "pkg/lib/lib.go", "")
		testutil.WriteFile(t, dir, "internal/core/core.go", "")

		// 3 Python files (30%)
		testutil.WriteFile(t, dir, "script.py", "")
		testutil.WriteFile(t, dir, "utils.py", "")
		testutil.WriteFile(t, dir, "config.py", "")

		// 1 JavaScript file (10%)
		testutil.WriteFile(t, dir, "index.js", "")

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(stats) != 3 {
			t.Fatalf("Detect() returned %d languages, want 3", len(stats))
		}

		// Should be sorted by percentage (descending)
		if stats[0].Language != types.LanguageGo {
			t.Errorf("Primary language = %v, want %v", stats[0].Language, types.LanguageGo)
		}

		if stats[1].Language != types.LanguagePython {
			t.Errorf("Second language = %v, want %v", stats[1].Language, types.LanguagePython)
		}

		if stats[2].Language != types.LanguageJavaScript {
			t.Errorf("Third language = %v, want %v", stats[2].Language, types.LanguageJavaScript)
		}

		// Check percentages
		if stats[0].Percentage != 60.0 {
			t.Errorf("Go percentage = %.1f, want 60.0", stats[0].Percentage)
		}

		if stats[1].Percentage != 30.0 {
			t.Errorf("Python percentage = %.1f, want 30.0", stats[1].Percentage)
		}

		if stats[2].Percentage != 10.0 {
			t.Errorf("JavaScript percentage = %.1f, want 10.0", stats[2].Percentage)
		}

		// Check IsPrimary (>10% threshold)
		if !stats[0].IsPrimary {
			t.Error("Go should be primary (60%)")
		}

		if !stats[1].IsPrimary {
			t.Error("Python should be primary (30%)")
		}

		if stats[2].IsPrimary {
			t.Error("JavaScript should NOT be primary (exactly 10%)")
		}
	})

	t.Run("primary language threshold", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// 9 Go files (90%)
		for i := 0; i < 9; i++ {
			testutil.WriteFile(t, dir, "file"+string(rune('0'+i))+".go", "")
		}

		// 1 Python file (10%)
		testutil.WriteFile(t, dir, "script.py", "")

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Go: 90% - should be primary
		goStats := findLanguage(stats, types.LanguageGo)
		if goStats == nil {
			t.Fatal("Go not found in stats")
		}

		if !goStats.IsPrimary {
			t.Error("Go should be primary (90%)")
		}

		// Python: 10% - should NOT be primary (threshold is >10%, not >=10%)
		pyStats := findLanguage(stats, types.LanguagePython)
		if pyStats == nil {
			t.Fatal("Python not found in stats")
		}

		if pyStats.IsPrimary {
			t.Error("Python should NOT be primary (exactly 10%, threshold is >10%)")
		}
	})

	t.Run("ignores unknown extensions", func(t *testing.T) {
		dir := testutil.TempDir(t)

		testutil.WriteFile(t, dir, "main.go", "")
		testutil.WriteFile(t, dir, "README.md", "")       // Unknown
		testutil.WriteFile(t, dir, "config.yaml", "")     // Unknown
		testutil.WriteFile(t, dir, "data.json", "")       // Unknown
		testutil.WriteFile(t, dir, "Makefile", "")        // Unknown
		testutil.WriteFile(t, dir, "build.sh", "")        // Unknown

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should only detect Go
		if len(stats) != 1 {
			t.Fatalf("Detect() returned %d languages, want 1", len(stats))
		}

		if stats[0].Language != types.LanguageGo {
			t.Errorf("Language = %v, want %v", stats[0].Language, types.LanguageGo)
		}

		if stats[0].FileCount != 1 {
			t.Errorf("FileCount = %d, want 1", stats[0].FileCount)
		}
	})

	t.Run("detects all supported languages", func(t *testing.T) {
		dir := testutil.TempDir(t)

		testutil.WriteFile(t, dir, "main.go", "")
		testutil.WriteFile(t, dir, "script.py", "")
		testutil.WriteFile(t, dir, "app.js", "")
		testutil.WriteFile(t, dir, "component.tsx", "")
		testutil.WriteFile(t, dir, "Main.java", "")
		testutil.WriteFile(t, dir, "lib.rs", "")
		testutil.WriteFile(t, dir, "Program.cs", "")
		testutil.WriteFile(t, dir, "app.rb", "")

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should detect 8 languages
		if len(stats) != 8 {
			t.Fatalf("Detect() returned %d languages, want 8", len(stats))
		}

		// Verify all expected languages are present
		expectedLangs := []types.Language{
			types.LanguageGo,
			types.LanguagePython,
			types.LanguageJavaScript,
			types.LanguageTypeScript,
			types.LanguageJava,
			types.LanguageRust,
			types.LanguageCSharp,
			types.LanguageRuby,
		}

		for _, expectedLang := range expectedLangs {
			found := false
			for _, stat := range stats {
				if stat.Language == expectedLang {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Language %v not detected", expectedLang)
			}
		}
	})

	t.Run("handles special files", func(t *testing.T) {
		dir := testutil.TempDir(t)

		testutil.WriteFile(t, dir, "Gemfile", "source 'https://rubygems.org'")
		testutil.WriteFile(t, dir, "Rakefile", "task :default => :test")

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Should detect Ruby
		if len(stats) != 1 {
			t.Fatalf("Detect() returned %d languages, want 1", len(stats))
		}

		if stats[0].Language != types.LanguageRuby {
			t.Errorf("Language = %v, want %v", stats[0].Language, types.LanguageRuby)
		}

		if stats[0].FileCount != 2 {
			t.Errorf("FileCount = %d, want 2 (Gemfile + Rakefile)", stats[0].FileCount)
		}
	})

	t.Run("empty repository", func(t *testing.T) {
		dir := testutil.TempDir(t)

		walker := NewWalker(dir)
		detector := NewLanguageDetector(walker)

		stats, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(stats) != 0 {
			t.Errorf("Detect() returned %d languages for empty repo, want 0", len(stats))
		}
	})
}

func TestDetectLanguage(t *testing.T) {
	detector := &LanguageDetector{}

	tests := []struct {
		name string
		ext  string
		filename string
		want types.Language
	}{
		// Go
		{
			name: "go file",
			ext:  ".go",
			filename: "main.go",
			want: types.LanguageGo,
		},

		// Python
		{
			name: "python file",
			ext:  ".py",
			filename: "script.py",
			want: types.LanguagePython,
		},
		{
			name: "python pyi file",
			ext:  ".pyi",
			filename: "types.pyi",
			want: types.LanguagePython,
		},
		{
			name: "jupyter notebook",
			ext:  ".ipynb",
			filename: "analysis.ipynb",
			want: types.LanguagePython,
		},

		// JavaScript
		{
			name: "javascript file",
			ext:  ".js",
			filename: "index.js",
			want: types.LanguageJavaScript,
		},
		{
			name: "jsx file",
			ext:  ".jsx",
			filename: "Component.jsx",
			want: types.LanguageJavaScript,
		},

		// TypeScript
		{
			name: "typescript file",
			ext:  ".ts",
			filename: "app.ts",
			want: types.LanguageTypeScript,
		},
		{
			name: "tsx file",
			ext:  ".tsx",
			filename: "Component.tsx",
			want: types.LanguageTypeScript,
		},

		// Java
		{
			name: "java file",
			ext:  ".java",
			filename: "Main.java",
			want: types.LanguageJava,
		},

		// Rust
		{
			name: "rust file",
			ext:  ".rs",
			filename: "lib.rs",
			want: types.LanguageRust,
		},

		// C#
		{
			name: "csharp file",
			ext:  ".cs",
			filename: "Program.cs",
			want: types.LanguageCSharp,
		},

		// Ruby
		{
			name: "ruby file",
			ext:  ".rb",
			filename: "app.rb",
			want: types.LanguageRuby,
		},
		{
			name: "Gemfile",
			ext:  "",
			filename: "Gemfile",
			want: types.LanguageRuby,
		},
		{
			name: "Rakefile",
			ext:  "",
			filename: "Rakefile",
			want: types.LanguageRuby,
		},

		// Unknown
		{
			name: "unknown extension",
			ext:  ".txt",
			filename: "README.txt",
			want: types.LanguageUnknown,
		},
		{
			name: "Makefile",
			ext:  "",
			filename: "Makefile",
			want: types.LanguageUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detector.detectLanguage(tt.ext, tt.filename)
			if got != tt.want {
				t.Errorf("detectLanguage(%q, %q) = %v, want %v", tt.ext, tt.filename, got, tt.want)
			}
		})
	}
}

// Helper function to find a language in stats
func findLanguage(stats []types.LanguageStats, lang types.Language) *types.LanguageStats {
	for i := range stats {
		if stats[i].Language == lang {
			return &stats[i]
		}
	}

	return nil
}
