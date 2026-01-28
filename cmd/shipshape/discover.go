// Ship Shape - Discover Command
// Copyright (c) 2026 Ship Shape Contributors
// Licensed under Apache License 2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chambridge/ship-shape/internal/discovery"
	"github.com/chambridge/ship-shape/internal/logger"
	"github.com/chambridge/ship-shape/pkg/types"
	"github.com/spf13/cobra"
)

var (
	discoverJSON bool
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover [directory]",
	Short: "Discover languages and frameworks in a repository",
	Long: `Analyzes a repository to discover programming languages, testing frameworks,
and development tools.

The discover command scans the repository and identifies:
  • Programming languages and their distribution
  • Testing frameworks (Jest, pytest, testify, etc.)
  • Coverage tools (nyc, c8, coverage.py, etc.)
  • Linters and formatters (eslint, prettier, black, etc.)
  • Build tools and task runners

Example:
  shipshape discover .
  shipshape discover /path/to/repo
  shipshape discover --json > repo-context.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDiscover,
}

func init() {
	rootCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().BoolVar(&discoverJSON, "json", false, "output in JSON format")
}

func runDiscover(_ *cobra.Command, args []string) error {
	// Determine target directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Verify directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	logger.Info("Discovering repository context", "directory", dir)

	// Create walker
	walker := discovery.NewWalker(dir)

	// Count total files for progress reporting
	totalFiles, err := walker.CountFiles()
	if err != nil {
		return fmt.Errorf("failed to count files: %w", err)
	}

	logger.Debug("Repository scan", "total_files", totalFiles)

	// Detect languages
	logger.Debug("Detecting languages...")
	languageDetector := discovery.NewLanguageDetector(walker)

	languages, err := languageDetector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect languages: %w", err)
	}

	logger.Debug("Languages detected", "count", len(languages))

	// Detect frameworks
	logger.Debug("Detecting frameworks...")
	frameworkDetector := discovery.NewFrameworkDetector(dir, walker)

	frameworks, err := frameworkDetector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect frameworks: %w", err)
	}

	logger.Debug("Frameworks detected", "count", len(frameworks))

	// Build repository context
	repo := types.Repository{
		Path:          dir,
		Languages:     languages,
		Frameworks:    frameworks,
		TotalFiles:    totalFiles,
		ExcludedPaths: walker.ExcludePatterns,
	}

	// Output results
	if discoverJSON {
		return outputJSON(&repo)
	}

	return outputText(&repo)
}

func outputJSON(repo *types.Repository) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(repo); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

//nolint:gocognit,gocyclo,nestif // Output formatting can be complex but is readable
func outputText(repo *types.Repository) error {
	fmt.Printf("Repository: %s\n", repo.Path)
	fmt.Printf("Total Files: %d\n\n", repo.TotalFiles)

	// Languages section
	if len(repo.Languages) > 0 {
		fmt.Println("Languages:")
		for _, lang := range repo.Languages {
			primary := ""
			if lang.IsPrimary {
				primary = " (primary)"
			}

			fmt.Printf("  • %s: %.1f%% (%d files)%s\n",
				lang.Language, lang.Percentage, lang.FileCount, primary)
		}

		fmt.Println()
	} else {
		fmt.Println("Languages: None detected")
		fmt.Println()
	}

	// Frameworks section
	if len(repo.Frameworks) > 0 {
		// Group frameworks by type
		frameworksByType := make(map[types.FrameworkType][]types.Framework)
		for _, fw := range repo.Frameworks {
			frameworksByType[fw.Type] = append(frameworksByType[fw.Type], fw)
		}

		fmt.Println("Frameworks & Tools:")

		// Test frameworks
		if frameworks, ok := frameworksByType[types.FrameworkTypeTest]; ok {
			fmt.Println("  Testing:")
			for _, fw := range frameworks {
				fmt.Printf("    • %s (%s)\n", fw.Name, fw.Language)
			}
		}

		// Coverage tools
		if frameworks, ok := frameworksByType[types.FrameworkTypeCoverage]; ok {
			fmt.Println("  Coverage:")
			for _, fw := range frameworks {
				fmt.Printf("    • %s (%s)\n", fw.Name, fw.Language)
			}
		}

		// Linters
		if frameworks, ok := frameworksByType[types.FrameworkTypeLint]; ok {
			fmt.Println("  Linting:")
			for _, fw := range frameworks {
				fmt.Printf("    • %s (%s)\n", fw.Name, fw.Language)
			}
		}

		// Formatters
		if frameworks, ok := frameworksByType[types.FrameworkTypeFormat]; ok {
			fmt.Println("  Formatting:")
			for _, fw := range frameworks {
				fmt.Printf("    • %s (%s)\n", fw.Name, fw.Language)
			}
		}

		// Build tools
		if frameworks, ok := frameworksByType[types.FrameworkTypeBuild]; ok {
			fmt.Println("  Build:")
			for _, fw := range frameworks {
				fmt.Printf("    • %s (%s)\n", fw.Name, fw.Language)
			}
		}

		// Other tools
		if frameworks, ok := frameworksByType[types.FrameworkTypeOther]; ok {
			fmt.Println("  Other:")
			for _, fw := range frameworks {
				fmt.Printf("    • %s (%s)\n", fw.Name, fw.Language)
			}
		}

		fmt.Println()
	} else {
		fmt.Println("Frameworks & Tools: None detected")
		fmt.Println()
	}

	return nil
}
