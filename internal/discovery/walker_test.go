package discovery

import (
	"path/filepath"
	"testing"

	"github.com/chambridge/ship-shape/internal/testutil"
)

func TestWalker_Walk(t *testing.T) {
	// Create temp directory with test structure
	dir := testutil.TempDir(t)

	// Create test file structure
	testutil.WriteFile(t, dir, "main.go", "package main")
	testutil.WriteFile(t, dir, "README.md", "# Project")
	testutil.WriteFile(t, dir, "internal/pkg/util.go", "package pkg")
	testutil.WriteFile(t, dir, "cmd/app/main.go", "package main")

	// Create excluded directories
	testutil.WriteFile(t, dir, "node_modules/lib/index.js", "module.exports = {}")
	testutil.WriteFile(t, dir, "vendor/pkg/dep.go", "package dep")
	testutil.WriteFile(t, dir, ".git/config", "[core]")

	walker := NewWalker(dir)

	var files []FileInfo
	count, err := walker.Walk(func(fi FileInfo) error {
		files = append(files, fi)
		return nil
	})

	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should find 4 files (main.go, README.md, internal/pkg/util.go, cmd/app/main.go)
	// Should exclude node_modules, vendor, .git
	expectedCount := 4
	if count != expectedCount {
		t.Errorf("Walk() processed %d files, want %d", count, expectedCount)
	}

	if len(files) != expectedCount {
		t.Errorf("Walk() collected %d files, want %d", len(files), expectedCount)
	}

	// Verify main.go is included
	found := false
	for _, f := range files {
		if f.Name == "main.go" && f.RelPath == "main.go" {
			found = true

			if f.Ext != ".go" {
				t.Errorf("main.go extension = %q, want %q", f.Ext, ".go")
			}

			if f.IsDir {
				t.Error("main.go should not be marked as directory")
			}
		}
	}

	if !found {
		t.Error("main.go not found in walked files")
	}

	// Verify excluded files are not included
	for _, f := range files {
		if contains(f.RelPath, "node_modules") {
			t.Errorf("File from node_modules should be excluded: %s", f.RelPath)
		}

		if contains(f.RelPath, "vendor") {
			t.Errorf("File from vendor should be excluded: %s", f.RelPath)
		}

		if contains(f.RelPath, ".git") {
			t.Errorf("File from .git should be excluded: %s", f.RelPath)
		}
	}
}

func TestWalker_ExcludePatterns(t *testing.T) {
	dir := testutil.TempDir(t)

	// Create files in various excluded directories
	testutil.WriteFile(t, dir, "src/main.go", "package main")
	testutil.WriteFile(t, dir, "build/output.txt", "build output")
	testutil.WriteFile(t, dir, "dist/bundle.js", "bundled code")
	testutil.WriteFile(t, dir, "target/classes/App.class", "compiled")
	testutil.WriteFile(t, dir, "__pycache__/module.pyc", "cached")

	walker := NewWalker(dir)

	var files []FileInfo
	_, err := walker.Walk(func(fi FileInfo) error {
		files = append(files, fi)
		return nil
	})

	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should only find src/main.go
	if len(files) != 1 {
		t.Errorf("Walk() found %d files, want 1", len(files))
		for _, f := range files {
			t.Logf("  - %s", f.RelPath)
		}
	}

	if len(files) > 0 && files[0].Name != "main.go" {
		t.Errorf("Found file %q, want %q", files[0].Name, "main.go")
	}
}

func TestWalker_HiddenFiles(t *testing.T) {
	dir := testutil.TempDir(t)

	testutil.WriteFile(t, dir, "main.go", "package main")
	testutil.WriteFile(t, dir, ".gitignore", "*.log")       // Allowed dotfile
	testutil.WriteFile(t, dir, ".eslintrc.json", "{}")      // Allowed dotfile
	testutil.WriteFile(t, dir, ".hidden/secret.txt", "ssh") // Hidden directory

	t.Run("exclude hidden by default", func(t *testing.T) {
		walker := NewWalker(dir)

		var files []FileInfo
		_, err := walker.Walk(func(fi FileInfo) error {
			files = append(files, fi)
			return nil
		})

		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// Should find: main.go, .gitignore, .eslintrc.json (3 files)
		// Should exclude: .hidden/secret.txt
		if len(files) != 3 {
			t.Errorf("Walk() found %d files, want 3", len(files))
			for _, f := range files {
				t.Logf("  - %s", f.RelPath)
			}
		}

		// Verify .hidden/secret.txt is excluded
		for _, f := range files {
			if contains(f.RelPath, ".hidden") {
				t.Errorf("Hidden directory file should be excluded: %s", f.RelPath)
			}
		}
	})

	t.Run("include hidden when enabled", func(t *testing.T) {
		walker := NewWalker(dir)
		walker.IncludeHidden = true

		var files []FileInfo
		_, err := walker.Walk(func(fi FileInfo) error {
			files = append(files, fi)
			return nil
		})

		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// Should find all files including hidden
		if len(files) < 4 {
			t.Errorf("Walk() found %d files, want at least 4", len(files))
		}
	})
}

func TestWalker_CustomExclusions(t *testing.T) {
	dir := testutil.TempDir(t)

	testutil.WriteFile(t, dir, "src/main.go", "package main")
	testutil.WriteFile(t, dir, "test/helper.go", "package test")
	testutil.WriteFile(t, dir, "docs/README.md", "# Docs")

	walker := NewWalker(dir)
	// Add custom exclusion for test directory
	walker.ExcludePatterns = append(walker.ExcludePatterns, "test")

	var files []FileInfo
	_, err := walker.Walk(func(fi FileInfo) error {
		files = append(files, fi)
		return nil
	})

	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should exclude test/helper.go
	for _, f := range files {
		if contains(f.RelPath, "test") {
			t.Errorf("File from excluded 'test' directory found: %s", f.RelPath)
		}
	}

	// Should include src and docs
	foundSrc := false
	foundDocs := false

	for _, f := range files {
		if f.Name == "main.go" {
			foundSrc = true
		}

		if f.Name == "README.md" {
			foundDocs = true
		}
	}

	if !foundSrc {
		t.Error("src/main.go should be included")
	}

	if !foundDocs {
		t.Error("docs/README.md should be included")
	}
}

func TestWalker_CountFiles(t *testing.T) {
	dir := testutil.TempDir(t)

	testutil.WriteFile(t, dir, "file1.go", "")
	testutil.WriteFile(t, dir, "file2.go", "")
	testutil.WriteFile(t, dir, "file3.go", "")

	walker := NewWalker(dir)

	count, err := walker.CountFiles()
	if err != nil {
		t.Fatalf("CountFiles() error = %v", err)
	}

	if count != 3 {
		t.Errorf("CountFiles() = %d, want 3", count)
	}
}

func TestWalker_RelativePaths(t *testing.T) {
	dir := testutil.TempDir(t)

	testutil.WriteFile(t, dir, "pkg/util.go", "package pkg")

	walker := NewWalker(dir)

	var files []FileInfo
	_, err := walker.Walk(func(fi FileInfo) error {
		files = append(files, fi)
		return nil
	})

	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	file := files[0]

	// RelPath should be relative to root
	expectedRelPath := filepath.Join("pkg", "util.go")
	if file.RelPath != expectedRelPath {
		t.Errorf("RelPath = %q, want %q", file.RelPath, expectedRelPath)
	}

	// Path should be absolute
	if !filepath.IsAbs(file.Path) {
		t.Errorf("Path %q is not absolute", file.Path)
	}
}

func TestWalker_FileExtensions(t *testing.T) {
	dir := testutil.TempDir(t)

	testutil.WriteFile(t, dir, "main.go", "")
	testutil.WriteFile(t, dir, "script.py", "")
	testutil.WriteFile(t, dir, "README", "") // No extension
	testutil.WriteFile(t, dir, "config.yaml", "")

	walker := NewWalker(dir)

	extensions := make(map[string]int)
	_, err := walker.Walk(func(fi FileInfo) error {
		extensions[fi.Ext]++
		return nil
	})

	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	tests := []struct {
		ext   string
		count int
	}{
		{".go", 1},
		{".py", 1},
		{"", 1}, // README with no extension
		{".yaml", 1},
	}

	for _, tt := range tests {
		if got := extensions[tt.ext]; got != tt.count {
			t.Errorf("Extension %q count = %d, want %d", tt.ext, got, tt.count)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return filepath.ToSlash(s) != "" && (s == substr || filepath.Dir(s) == substr ||
		filepath.Base(filepath.Dir(s)) == substr ||
		filepath.Base(filepath.Dir(filepath.Dir(s))) == substr)
}
