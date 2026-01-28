package discovery

import (
	"testing"

	"github.com/chambridge/ship-shape/internal/testutil"
	"github.com/chambridge/ship-shape/pkg/types"
)

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestManifestParser_ParsePackageJSON(t *testing.T) {
	t.Run("detects jest test framework", func(t *testing.T) {
		dir := testutil.TempDir(t)

		packageJSON := `{
			"name": "my-app",
			"devDependencies": {
				"jest": "^29.0.0"
			}
		}`
		testutil.WriteFile(t, dir, "package.json", packageJSON)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parsePackageJSON()

		if err != nil {
			t.Fatalf("parsePackageJSON() error = %v", err)
		}

		if len(frameworks) != 1 {
			t.Fatalf("Expected 1 framework, got %d", len(frameworks))
		}

		fw := frameworks[0]
		if fw.Name != "jest" {
			t.Errorf("Framework name = %q, want %q", fw.Name, "jest")
		}

		if fw.Type != types.FrameworkTypeTest {
			t.Errorf("Framework type = %v, want %v", fw.Type, types.FrameworkTypeTest)
		}

		if fw.Language != types.LanguageJavaScript {
			t.Errorf("Framework language = %v, want %v", fw.Language, types.LanguageJavaScript)
		}
	})

	t.Run("detects multiple tools", func(t *testing.T) {
		dir := testutil.TempDir(t)

		packageJSON := `{
			"name": "my-app",
			"devDependencies": {
				"jest": "^29.0.0",
				"eslint": "^8.0.0",
				"prettier": "^2.8.0",
				"nyc": "^15.1.0"
			}
		}`
		testutil.WriteFile(t, dir, "package.json", packageJSON)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parsePackageJSON()

		if err != nil {
			t.Fatalf("parsePackageJSON() error = %v", err)
		}

		if len(frameworks) != 4 {
			t.Fatalf("Expected 4 frameworks, got %d", len(frameworks))
		}

		// Verify each type is present
		frameworkTypes := make(map[types.FrameworkType]int)
		for _, fw := range frameworks {
			frameworkTypes[fw.Type]++
		}

		if frameworkTypes[types.FrameworkTypeTest] != 1 {
			t.Errorf("Expected 1 test framework, got %d", frameworkTypes[types.FrameworkTypeTest])
		}

		if frameworkTypes[types.FrameworkTypeLint] != 1 {
			t.Errorf("Expected 1 lint framework, got %d", frameworkTypes[types.FrameworkTypeLint])
		}

		if frameworkTypes[types.FrameworkTypeFormat] != 1 {
			t.Errorf("Expected 1 format framework, got %d", frameworkTypes[types.FrameworkTypeFormat])
		}

		if frameworkTypes[types.FrameworkTypeCoverage] != 1 {
			t.Errorf("Expected 1 coverage framework, got %d", frameworkTypes[types.FrameworkTypeCoverage])
		}
	})

	t.Run("handles missing package.json", func(t *testing.T) {
		dir := testutil.TempDir(t)

		parser := NewManifestParser(dir)
		_, err := parser.parsePackageJSON()

		if err == nil {
			t.Error("Expected error for missing package.json")
		}
	})
}

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestManifestParser_ParseGoMod(t *testing.T) {
	t.Run("detects testify", func(t *testing.T) {
		dir := testutil.TempDir(t)

		goMod := `module github.com/example/app

go 1.21

require (
	github.com/stretchr/testify v1.8.4
)
`
		testutil.WriteFile(t, dir, "go.mod", goMod)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parseGoMod()

		if err != nil {
			t.Fatalf("parseGoMod() error = %v", err)
		}

		if len(frameworks) != 1 {
			t.Fatalf("Expected 1 framework, got %d", len(frameworks))
		}

		fw := frameworks[0]
		if fw.Name != "testify" {
			t.Errorf("Framework name = %q, want %q", fw.Name, "testify")
		}

		if fw.Language != types.LanguageGo {
			t.Errorf("Framework language = %v, want %v", fw.Language, types.LanguageGo)
		}

		if fw.Type != types.FrameworkTypeTest {
			t.Errorf("Framework type = %v, want %v", fw.Type, types.FrameworkTypeTest)
		}
	})

	t.Run("detects multiple test libraries", func(t *testing.T) {
		dir := testutil.TempDir(t)

		goMod := `module github.com/example/app

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/golang/mock v1.6.0
	github.com/onsi/ginkgo/v2 v2.11.0
)
`
		testutil.WriteFile(t, dir, "go.mod", goMod)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parseGoMod()

		if err != nil {
			t.Fatalf("parseGoMod() error = %v", err)
		}

		if len(frameworks) != 3 {
			t.Fatalf("Expected 3 frameworks, got %d", len(frameworks))
		}

		// Verify all frameworks are test type
		for _, fw := range frameworks {
			if fw.Type != types.FrameworkTypeTest {
				t.Errorf("Framework %s type = %v, want %v", fw.Name, fw.Type, types.FrameworkTypeTest)
			}
		}
	})

	t.Run("handles missing go.mod", func(t *testing.T) {
		dir := testutil.TempDir(t)

		parser := NewManifestParser(dir)
		_, err := parser.parseGoMod()

		if err == nil {
			t.Error("Expected error for missing go.mod")
		}
	})
}

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestManifestParser_ParsePyprojectToml(t *testing.T) {
	t.Run("detects pytest", func(t *testing.T) {
		dir := testutil.TempDir(t)

		pyproject := `[tool.pytest.ini_options]
testpaths = ["tests"]

[tool.coverage.run]
source = ["src"]
`
		testutil.WriteFile(t, dir, "pyproject.toml", pyproject)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parsePyprojectToml()

		if err != nil {
			t.Fatalf("parsePyprojectToml() error = %v", err)
		}

		// Should find pytest and coverage
		if len(frameworks) < 1 {
			t.Fatalf("Expected at least 1 framework, got %d", len(frameworks))
		}

		// Check pytest is detected
		found := false
		for _, fw := range frameworks {
			if fw.Name == "pytest" {
				found = true

				if fw.Language != types.LanguagePython {
					t.Errorf("pytest language = %v, want %v", fw.Language, types.LanguagePython)
				}

				if fw.Type != types.FrameworkTypeTest {
					t.Errorf("pytest type = %v, want %v", fw.Type, types.FrameworkTypeTest)
				}
			}
		}

		if !found {
			t.Error("pytest not detected")
		}
	})

	t.Run("detects multiple Python tools", func(t *testing.T) {
		dir := testutil.TempDir(t)

		pyproject := `[tool.pytest.ini_options]
testpaths = ["tests"]

[tool.black]
line-length = 100

[tool.ruff]
line-length = 100
`
		testutil.WriteFile(t, dir, "pyproject.toml", pyproject)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parsePyprojectToml()

		if err != nil {
			t.Fatalf("parsePyprojectToml() error = %v", err)
		}

		if len(frameworks) < 2 {
			t.Fatalf("Expected at least 2 frameworks, got %d", len(frameworks))
		}
	})
}

//nolint:gocognit // Table-driven tests can be complex but are still readable
func TestManifestParser_ParseRequirementsTxt(t *testing.T) {
	t.Run("detects pytest from requirements.txt", func(t *testing.T) {
		dir := testutil.TempDir(t)

		requirements := `pytest==7.4.0
pytest-cov==4.1.0
black==23.7.0
`
		testutil.WriteFile(t, dir, "requirements.txt", requirements)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parseRequirementsTxt()

		if err != nil {
			t.Fatalf("parseRequirementsTxt() error = %v", err)
		}

		if len(frameworks) != 3 {
			t.Fatalf("Expected 3 frameworks, got %d", len(frameworks))
		}

		// Verify pytest is detected
		found := false
		for _, fw := range frameworks {
			if fw.Name == "pytest" {
				found = true

				if fw.Language != types.LanguagePython {
					t.Errorf("pytest language = %v, want %v", fw.Language, types.LanguagePython)
				}

				if fw.Type != types.FrameworkTypeTest {
					t.Errorf("pytest type = %v, want %v", fw.Type, types.FrameworkTypeTest)
				}
			}
		}

		if !found {
			t.Error("pytest not detected")
		}
	})

	t.Run("handles version specifiers", func(t *testing.T) {
		dir := testutil.TempDir(t)

		requirements := `pytest>=7.0.0
coverage~=6.5
black>22.0
ruff!=0.0.280
`
		testutil.WriteFile(t, dir, "requirements.txt", requirements)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parseRequirementsTxt()

		if err != nil {
			t.Fatalf("parseRequirementsTxt() error = %v", err)
		}

		if len(frameworks) != 4 {
			t.Fatalf("Expected 4 frameworks, got %d", len(frameworks))
		}
	})

	t.Run("ignores comments and empty lines", func(t *testing.T) {
		dir := testutil.TempDir(t)

		requirements := `# Test dependencies
pytest==7.4.0

# Code quality
black==23.7.0
# ruff==0.0.280  (commented out)
`
		testutil.WriteFile(t, dir, "requirements.txt", requirements)

		parser := NewManifestParser(dir)
		frameworks, err := parser.parseRequirementsTxt()

		if err != nil {
			t.Fatalf("parseRequirementsTxt() error = %v", err)
		}

		// Should only find pytest and black (not the commented ruff)
		if len(frameworks) != 2 {
			t.Fatalf("Expected 2 frameworks, got %d", len(frameworks))
		}
	})
}

func TestManifestParser_ParseAll(t *testing.T) {
	t.Run("combines all manifest sources", func(t *testing.T) {
		dir := testutil.TempDir(t)

		// Create multiple manifest files
		packageJSON := `{
			"devDependencies": {
				"jest": "^29.0.0"
			}
		}`
		testutil.WriteFile(t, dir, "package.json", packageJSON)

		goMod := `module github.com/example/app
require github.com/stretchr/testify v1.8.4
`
		testutil.WriteFile(t, dir, "go.mod", goMod)

		requirements := `pytest==7.4.0
`
		testutil.WriteFile(t, dir, "requirements.txt", requirements)

		parser := NewManifestParser(dir)
		frameworks, err := parser.ParseAll()

		if err != nil {
			t.Fatalf("ParseAll() error = %v", err)
		}

		// Should find frameworks from all sources
		if len(frameworks) < 3 {
			t.Fatalf("Expected at least 3 frameworks, got %d", len(frameworks))
		}

		// Verify we have frameworks from each language
		langs := make(map[types.Language]bool)
		for _, fw := range frameworks {
			langs[fw.Language] = true
		}

		if !langs[types.LanguageJavaScript] {
			t.Error("No JavaScript frameworks detected")
		}

		if !langs[types.LanguageGo] {
			t.Error("No Go frameworks detected")
		}

		if !langs[types.LanguagePython] {
			t.Error("No Python frameworks detected")
		}
	})

	t.Run("handles repository with no manifests", func(t *testing.T) {
		dir := testutil.TempDir(t)

		parser := NewManifestParser(dir)
		frameworks, err := parser.ParseAll()

		if err != nil {
			t.Fatalf("ParseAll() error = %v", err)
		}

		// Should return empty list, not error
		if len(frameworks) != 0 {
			t.Errorf("Expected 0 frameworks, got %d", len(frameworks))
		}
	})
}
