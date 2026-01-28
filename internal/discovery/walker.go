// Package discovery handles repository discovery and context analysis.
package discovery

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DefaultExcludePatterns are directory patterns excluded from analysis.
var DefaultExcludePatterns = []string{
	// Version control
	".git",
	".svn",
	".hg",

	// Dependencies
	"node_modules",
	"vendor",
	"venv",
	".venv",
	"env",
	".env",
	"__pycache__",
	".tox",

	// Build outputs
	"dist",
	"build",
	"target",
	"out",
	"bin",
	".next",
	".nuxt",

	// IDE/Editor
	".idea",
	".vscode",
	".vs",
	"*.swp",
	"*.swo",
	".DS_Store",

	// Test coverage
	"coverage",
	".coverage",
	"htmlcov",
	".nyc_output",

	// Misc
	"tmp",
	"temp",
	".cache",
}

// Walker provides file system traversal with exclusion patterns.
type Walker struct {
	// Root is the starting directory for traversal
	Root string

	// ExcludePatterns are directory/file patterns to skip
	ExcludePatterns []string

	// IncludeHidden includes hidden files/directories (starting with .)
	IncludeHidden bool
}

// FileInfo contains information about a discovered file.
type FileInfo struct {
	// Path is the absolute path to the file
	Path string

	// RelPath is the path relative to the walker root
	RelPath string

	// Name is the file name with extension
	Name string

	// Ext is the file extension (including the dot, e.g., ".go")
	Ext string

	// IsDir indicates if this is a directory
	IsDir bool

	// Size is the file size in bytes
	Size int64
}

// NewWalker creates a new file system walker with default exclusions.
func NewWalker(root string) *Walker {
	return &Walker{
		Root:            root,
		ExcludePatterns: DefaultExcludePatterns,
		IncludeHidden:   false,
	}
}

// Walk traverses the file system and calls fn for each file.
// Directories matching exclusion patterns are skipped.
// Returns the total number of files processed and any error.
func (w *Walker) Walk(fn func(FileInfo) error) (int, error) {
	fileCount := 0

	err := filepath.WalkDir(w.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Return error for debugging
			return err
		}

		// Get relative path from root
		relPath, err := filepath.Rel(w.Root, path)
		if err != nil {
			return nil
		}

		// Check exclusions
		if w.shouldExclude(relPath, d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		// Skip directories (we only process files)
		if d.IsDir() {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Build FileInfo
		fileInfo := FileInfo{
			Path:    path,
			RelPath: relPath,
			Name:    d.Name(),
			Ext:     filepath.Ext(d.Name()),
			IsDir:   d.IsDir(),
			Size:    info.Size(),
		}

		// Call the callback
		if err := fn(fileInfo); err != nil {
			return err
		}

		fileCount++

		return nil
	})

	return fileCount, err
}

// shouldExclude checks if a path should be excluded based on patterns.
func (w *Walker) shouldExclude(relPath string, _ bool) bool {
	// Don't exclude root directory (relPath = ".")
	if relPath == "." {
		return false
	}

	// Skip hidden files/directories if not included
	if !w.IncludeHidden && strings.HasPrefix(filepath.Base(relPath), ".") {
		// Allow some common dotfiles
		base := filepath.Base(relPath)
		if !isAllowedDotfile(base) {
			return true
		}
	}

	// Check each exclusion pattern
	for _, pattern := range w.ExcludePatterns {
		// Check if any path component matches the pattern
		// nolint: stringsseq
		pathParts := strings.Split(relPath, string(os.PathSeparator))
		for _, part := range pathParts {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}

			// Direct string match for common patterns
			if part == pattern {
				return true
			}
		}
	}

	return false
}

// isAllowedDotfile checks if a dotfile is allowed (not excluded).
func isAllowedDotfile(name string) bool {
	allowed := []string{
		".gitignore",
		".gitattributes",
		".editorconfig",
		".prettierrc",
		".eslintrc",
		".pylintrc",
		".go-version",
		".python-version",
		".ruby-version",
		".nvmrc",
	}

	for _, a := range allowed {
		if name == a || strings.HasPrefix(name, a+".") {
			return true
		}
	}

	return false
}

// CountFiles returns the total number of files that would be processed.
// This is useful for progress reporting.
func (w *Walker) CountFiles() (int, error) {
	count := 0
	_, err := w.Walk(func(FileInfo) error {
		count++
		return nil
	})

	return count, err
}
