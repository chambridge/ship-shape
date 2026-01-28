package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/chambridge/ship-shape/pkg/types"
)

// ManifestParser parses dependency manifests to detect frameworks and tools.
type ManifestParser struct {
	rootPath string
}

// NewManifestParser creates a new manifest parser.
func NewManifestParser(rootPath string) *ManifestParser {
	return &ManifestParser{
		rootPath: rootPath,
	}
}

// ParseAll finds and parses all dependency manifests in the repository.
func (p *ManifestParser) ParseAll() ([]types.Framework, error) {
	var frameworks []types.Framework

	// Parse package.json (JavaScript/TypeScript)
	if pkgFrameworks, err := p.parsePackageJSON(); err == nil {
		frameworks = append(frameworks, pkgFrameworks...)
	}

	// Parse go.mod (Go)
	if goFrameworks, err := p.parseGoMod(); err == nil {
		frameworks = append(frameworks, goFrameworks...)
	}

	// Parse pyproject.toml (Python)
	if pyFrameworks, err := p.parsePyprojectToml(); err == nil {
		frameworks = append(frameworks, pyFrameworks...)
	}

	// Parse requirements.txt (Python)
	if reqFrameworks, err := p.parseRequirementsTxt(); err == nil {
		frameworks = append(frameworks, reqFrameworks...)
	}

	return frameworks, nil
}

// PackageJSON represents a simplified package.json structure.
type PackageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

// parsePackageJSON parses package.json and extracts framework information.
func (p *ManifestParser) parsePackageJSON() ([]types.Framework, error) {
	path := filepath.Join(p.rootPath, "package.json")

	data, err := os.ReadFile(path) //nolint:gosec // Reading manifest files from repository root
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var frameworks []types.Framework

	// Combine all dependencies
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}

	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	// Known test frameworks
	testFrameworks := map[string]string{
		"jest":       "jest",
		"mocha":      "mocha",
		"vitest":     "vitest",
		"jasmine":    "jasmine",
		"@jest/core": "jest",
	}

	// Known coverage tools
	coverageTools := map[string]string{
		"nyc":      "nyc",
		"c8":       "c8",
		"istanbul": "istanbul",
	}

	// Known linters
	linters := map[string]string{
		"eslint":                    "eslint",
		"tslint":                    "tslint",
		"@typescript-eslint/parser": "eslint",
	}

	// Known formatters
	formatters := map[string]string{
		"prettier": "prettier",
	}

	// Detect frameworks
	for dep, version := range allDeps {
		if name, ok := testFrameworks[dep]; ok {
			lang := types.LanguageJavaScript
			if strings.Contains(dep, "typescript") || hasTypeScriptFiles(p.rootPath) {
				lang = types.LanguageTypeScript
			}

			frameworks = append(frameworks, types.Framework{
				Name:        name,
				Language:    lang,
				Type:        types.FrameworkTypeTest,
				Version:     version,
				ConfigFiles: []string{"package.json"},
			})
		}

		if name, ok := coverageTools[dep]; ok {
			frameworks = append(frameworks, types.Framework{
				Name:        name,
				Language:    types.LanguageJavaScript,
				Type:        types.FrameworkTypeCoverage,
				Version:     version,
				ConfigFiles: []string{"package.json"},
			})
		}

		if name, ok := linters[dep]; ok {
			frameworks = append(frameworks, types.Framework{
				Name:        name,
				Language:    types.LanguageJavaScript,
				Type:        types.FrameworkTypeLint,
				Version:     version,
				ConfigFiles: []string{"package.json"},
			})
		}

		if name, ok := formatters[dep]; ok {
			frameworks = append(frameworks, types.Framework{
				Name:        name,
				Language:    types.LanguageJavaScript,
				Type:        types.FrameworkTypeFormat,
				Version:     version,
				ConfigFiles: []string{"package.json"},
			})
		}
	}

	return frameworks, nil
}

// parseGoMod parses go.mod and extracts framework information.
func (p *ManifestParser) parseGoMod() ([]types.Framework, error) {
	path := filepath.Join(p.rootPath, "go.mod")

	data, err := os.ReadFile(path) //nolint:gosec // Reading manifest files from repository root
	if err != nil {
		return nil, err
	}

	content := string(data)

	var frameworks []types.Framework

	// Check for testify (popular Go testing library)
	if strings.Contains(content, "github.com/stretchr/testify") {
		frameworks = append(frameworks, types.Framework{
			Name:        "testify",
			Language:    types.LanguageGo,
			Type:        types.FrameworkTypeTest,
			ConfigFiles: []string{"go.mod"},
		})
	}

	// Check for gomock
	if strings.Contains(content, "github.com/golang/mock") || strings.Contains(content, "go.uber.org/mock") {
		frameworks = append(frameworks, types.Framework{
			Name:        "gomock",
			Language:    types.LanguageGo,
			Type:        types.FrameworkTypeTest,
			ConfigFiles: []string{"go.mod"},
		})
	}

	// Check for ginkgo
	if strings.Contains(content, "github.com/onsi/ginkgo") {
		frameworks = append(frameworks, types.Framework{
			Name:        "ginkgo",
			Language:    types.LanguageGo,
			Type:        types.FrameworkTypeTest,
			ConfigFiles: []string{"go.mod"},
		})
	}

	// Note: Go's built-in testing package doesn't appear in go.mod
	// We'll detect it by looking for *_test.go files

	return frameworks, nil
}

// parsePyprojectToml parses pyproject.toml (simplified version).
func (p *ManifestParser) parsePyprojectToml() ([]types.Framework, error) {
	path := filepath.Join(p.rootPath, "pyproject.toml")

	data, err := os.ReadFile(path) //nolint:gosec // Reading manifest files from repository root
	if err != nil {
		return nil, err
	}

	content := string(data)

	var frameworks []types.Framework

	// Simple string matching for common frameworks
	// TODO: Use proper TOML parser (github.com/pelletier/go-toml/v2)

	if strings.Contains(content, "pytest") {
		frameworks = append(frameworks, types.Framework{
			Name:        "pytest",
			Language:    types.LanguagePython,
			Type:        types.FrameworkTypeTest,
			ConfigFiles: []string{"pyproject.toml"},
		})
	}

	if strings.Contains(content, "coverage") || strings.Contains(content, "pytest-cov") {
		frameworks = append(frameworks, types.Framework{
			Name:        "coverage.py",
			Language:    types.LanguagePython,
			Type:        types.FrameworkTypeCoverage,
			ConfigFiles: []string{"pyproject.toml"},
		})
	}

	if strings.Contains(content, "black") {
		frameworks = append(frameworks, types.Framework{
			Name:        "black",
			Language:    types.LanguagePython,
			Type:        types.FrameworkTypeFormat,
			ConfigFiles: []string{"pyproject.toml"},
		})
	}

	if strings.Contains(content, "ruff") {
		frameworks = append(frameworks, types.Framework{
			Name:        "ruff",
			Language:    types.LanguagePython,
			Type:        types.FrameworkTypeLint,
			ConfigFiles: []string{"pyproject.toml"},
		})
	}

	return frameworks, nil
}

// parseRequirementsTxt parses requirements.txt.
func (p *ManifestParser) parseRequirementsTxt() ([]types.Framework, error) {
	path := filepath.Join(p.rootPath, "requirements.txt")

	data, err := os.ReadFile(path) //nolint:gosec // Reading manifest files from repository root
	if err != nil {
		return nil, err
	}

	content := string(data)

	var frameworks []types.Framework

	// Parse line by line
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Extract package name (before ==, >=, ~=, etc.)
		pkg := strings.FieldsFunc(line, func(r rune) bool {
			return r == '=' || r == '>' || r == '<' || r == '!' || r == '~'
		})[0]

		pkg = strings.TrimSpace(pkg)

		switch pkg {
		case "pytest":
			frameworks = append(frameworks, types.Framework{
				Name:        "pytest",
				Language:    types.LanguagePython,
				Type:        types.FrameworkTypeTest,
				ConfigFiles: []string{"requirements.txt"},
			})
		case "coverage", "pytest-cov":
			frameworks = append(frameworks, types.Framework{
				Name:        "coverage.py",
				Language:    types.LanguagePython,
				Type:        types.FrameworkTypeCoverage,
				ConfigFiles: []string{"requirements.txt"},
			})
		case "black":
			frameworks = append(frameworks, types.Framework{
				Name:        "black",
				Language:    types.LanguagePython,
				Type:        types.FrameworkTypeFormat,
				ConfigFiles: []string{"requirements.txt"},
			})
		case "pylint", "flake8", "ruff":
			frameworks = append(frameworks, types.Framework{
				Name:        pkg,
				Language:    types.LanguagePython,
				Type:        types.FrameworkTypeLint,
				ConfigFiles: []string{"requirements.txt"},
			})
		}
	}

	return frameworks, nil
}

// hasTypeScriptFiles checks if the repository contains TypeScript files.
func hasTypeScriptFiles(rootPath string) bool {
	// Check for tsconfig.json
	if _, err := os.Stat(filepath.Join(rootPath, "tsconfig.json")); err == nil {
		return true
	}

	// TODO: Walk directory to check for .ts/.tsx files
	return false
}
