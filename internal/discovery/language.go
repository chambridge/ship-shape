package discovery

import (
	"strings"

	"github.com/chambridge/ship-shape/pkg/types"
)

// ExtensionMap maps file extensions to programming languages.
// This is a simplified version; production would use go-enry for more accuracy.
var ExtensionMap = map[string]types.Language{
	// Go
	".go": types.LanguageGo,

	// Python
	".py":    types.LanguagePython,
	".pyw":   types.LanguagePython,
	".pyx":   types.LanguagePython,
	".pyi":   types.LanguagePython,
	".ipynb": types.LanguagePython, // Jupyter notebooks

	// JavaScript/TypeScript
	".js":  types.LanguageJavaScript,
	".jsx": types.LanguageJavaScript,
	".mjs": types.LanguageJavaScript,
	".cjs": types.LanguageJavaScript,
	".ts":  types.LanguageTypeScript,
	".tsx": types.LanguageTypeScript,
	".mts": types.LanguageTypeScript,
	".cts": types.LanguageTypeScript,

	// Java
	".java": types.LanguageJava,

	// Rust
	".rs": types.LanguageRust,

	// C#
	".cs":     types.LanguageCSharp,
	".cshtml": types.LanguageCSharp,
	".csx":    types.LanguageCSharp,

	// Ruby
	".rb":   types.LanguageRuby,
	".rake": types.LanguageRuby,
}

// LanguageDetector detects languages in a repository.
type LanguageDetector struct {
	walker *Walker
}

// NewLanguageDetector creates a new language detector.
func NewLanguageDetector(walker *Walker) *LanguageDetector {
	return &LanguageDetector{
		walker: walker,
	}
}

// Detect analyzes the repository and returns language statistics.
func (d *LanguageDetector) Detect() ([]types.LanguageStats, error) {
	// Count files by language
	langCounts := make(map[types.Language]int)
	totalFiles := 0

	_, err := d.walker.Walk(func(fi FileInfo) error {
		// Determine language from extension
		lang := d.detectLanguage(fi.Ext, fi.Name)
		if lang != types.LanguageUnknown {
			langCounts[lang]++
			totalFiles++
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Convert to LanguageStats
	var stats []types.LanguageStats

	for lang, count := range langCounts {
		percentage := 0.0
		if totalFiles > 0 {
			percentage = (float64(count) / float64(totalFiles)) * 100.0
		}

		stats = append(stats, types.LanguageStats{
			Language:   lang,
			FileCount:  count,
			Percentage: percentage,
			IsPrimary:  percentage > 10.0, // >10% threshold for primary languages
		})
	}

	// Sort by percentage (descending)
	sortLanguageStats(stats)

	return stats, nil
}

// detectLanguage determines the language from file extension and name.
func (d *LanguageDetector) detectLanguage(ext, name string) types.Language {
	// Check extension map
	if lang, ok := ExtensionMap[strings.ToLower(ext)]; ok {
		return lang
	}

	// Special case: files without extensions
	switch strings.ToLower(name) {
	case "gemfile", "rakefile":
		return types.LanguageRuby
	case "makefile":
		// Not a programming language
		return types.LanguageUnknown
	}

	return types.LanguageUnknown
}

// sortLanguageStats sorts language statistics by percentage (descending).
func sortLanguageStats(stats []types.LanguageStats) {
	// Simple bubble sort for small slices
	for i := 0; i < len(stats); i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].Percentage < stats[j].Percentage {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}
}
