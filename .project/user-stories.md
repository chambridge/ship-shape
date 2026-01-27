# Ship Shape - User Stories with Acceptance Criteria

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Draft
**Author**: Senior Software Engineer

---

## Table of Contents

1. [Epic Overview](#epic-overview)
2. [Epic 1: Repository Discovery & Context Analysis](#epic-1-repository-discovery--context-analysis)
3. [Epic 2: Multi-Language Analysis](#epic-2-multi-language-analysis)
4. [Epic 3: Monorepo Support](#epic-3-monorepo-support)
5. [Epic 4: Test Quality Analysis](#epic-4-test-quality-analysis)
6. [Epic 5: Coverage Analysis](#epic-5-coverage-analysis)
7. [Epic 6: Tool Adoption & Recommendations](#epic-6-tool-adoption--recommendations)
8. [Epic 7: Scoring & Assessment](#epic-7-scoring--assessment)
9. [Epic 8: Quality Gates](#epic-8-quality-gates)
10. [Epic 9: CI/CD Integration (GitHub Actions)](#epic-9-cicd-integration-github-actions)
11. [Epic 10: Reporting & Dashboards](#epic-10-reporting--dashboards)
12. [Epic 11: Historical Tracking & Trends](#epic-11-historical-tracking--trends)
13. [Epic 12: Organization-Wide Features](#epic-12-organization-wide-features)

---

## Epic Overview

### Epic Prioritization

| Epic ID | Epic Name | Priority | Target Release | Open Source Tools |
|---------|-----------|----------|----------------|-------------------|
| EPIC-1 | Repository Discovery & Context Analysis | P0 (Critical) | v0.1.0 | go-enry, linguist |
| EPIC-2 | Multi-Language Analysis | P0 (Critical) | v0.2.0 | Language-specific parsers |
| EPIC-3 | Monorepo Support | P0 (Critical) | v0.2.0 | jq, yq, go-glob |
| EPIC-4 | Test Quality Analysis | P0 (Critical) | v0.3.0 | AST parsers, tree-sitter |
| EPIC-5 | Coverage Analysis | P0 (Critical) | v0.3.0 | gocov, coverage.py, c8 |
| EPIC-6 | Tool Adoption & Recommendations | P1 (High) | v0.4.0 | YAML database |
| EPIC-7 | Scoring & Assessment | P0 (Critical) | v0.4.0 | Custom Go implementation |
| EPIC-8 | Quality Gates | P0 (Critical) | v0.5.0 | Exit codes, JSON output |
| EPIC-9 | CI/CD Integration (GitHub Actions) | P0 (Critical) | v0.5.0 | GitHub Actions |
| EPIC-10 | Reporting & Dashboards | P1 (High) | v0.6.0 | templ, htmx |
| EPIC-11 | Historical Tracking & Trends | P1 (High) | v0.7.0 | SQLite, gonum/plot |
| EPIC-12 | Organization-Wide Features | P2 (Medium) | v0.8.0 | SQLite, REST API |

### Story Point Reference

- **1 point**: < 4 hours, single file/function, minimal testing
- **2 points**: 4-8 hours, 2-3 files, standard testing
- **3 points**: 1-2 days, multiple files, integration testing
- **5 points**: 2-3 days, complex logic, extensive testing
- **8 points**: 3-5 days, new subsystem, comprehensive testing
- **13 points**: 1-2 weeks, major feature, full test coverage

---

## Epic 1: Repository Discovery & Context Analysis

**Epic Goal**: Automatically discover and understand repository structure, languages, frameworks, and infrastructure before performing any analysis.

### Story SS-001: File Extension-Based Language Detection

**As a** developer
**I want** Ship Shape to automatically detect all programming languages in my repository
**So that** I receive relevant, language-specific analysis without manual configuration

**Priority**: P0 (Critical)
**Story Points**: 5
**Dependencies**: None

#### Acceptance Criteria

**AC-001.1**: File Extension Analysis
- **Given** a repository with multiple file types
- **When** Ship Shape runs discovery
- **Then** it identifies all languages by file extensions (.go, .py, .js, .ts, .java, .rs, etc.)
- **And** calculates percentage distribution of each language
- **And** determines file count per language

**AC-001.2**: Primary Language Identification
- **Given** language distribution is calculated
- **When** any language represents >10% of codebase
- **Then** that language is marked as "primary"
- **And** primary languages receive full analysis depth
- **And** secondary languages (<10%) receive basic analysis

**AC-001.3**: Ignore Common Non-Code Files
- **Given** a repository contains various file types
- **When** Ship Shape analyzes file extensions
- **Then** it excludes vendor directories (node_modules, vendor, target, etc.)
- **And** excludes binary files
- **And** excludes generated code directories
- **And** excludes documentation-only files

#### Technical Requirements

**TR-001.1**: Use go-enry for language detection
```go
import "github.com/go-enry/go-enry/v2"

func detectLanguages(repoPath string) (map[string]*LanguageInfo, error) {
    // Walk repository tree
    // Use enry.GetLanguage() for each file
    // Calculate distributions
    // Identify primary languages (>10%)
}
```

**TR-001.2**: File walking with exclusions
```go
// Use filepath.Walk with exclusion patterns
excludePatterns := []string{
    "node_modules/", "vendor/", "target/",
    ".git/", "dist/", "build/", "__pycache__/",
}
```

**TR-001.3**: Performance requirement
- Must analyze 10,000 files in <5 seconds
- Use parallel file processing with goroutines
- Stream results to avoid memory issues

#### Open Source Tools
- **go-enry/go-enry**: GitHub's language detection library (Apache 2.0)
- **filepath.Walk**: Go standard library

#### Test Requirements
- Unit tests with mock file systems
- Integration tests with real multi-language repos
- Performance benchmarks (10k files in 5s)
- Test cases: Python-only, Go-only, polyglot, monorepo

---

### Story SS-002: Configuration File-Based Framework Detection

**As a** developer
**I want** Ship Shape to detect testing frameworks and build systems from configuration files
**So that** it understands my project's infrastructure without manual input

**Priority**: P0 (Critical)
**Story Points**: 8
**Dependencies**: SS-001

#### Acceptance Criteria

**AC-002.1**: Package Manager Detection
- **Given** a repository with dependency files
- **When** Ship Shape scans for configuration files
- **Then** it detects package.json (npm/yarn/pnpm)
- **And** detects requirements.txt, setup.py, pyproject.toml (Python)
- **And** detects go.mod, go.sum (Go)
- **And** detects pom.xml, build.gradle (Java)
- **And** detects Cargo.toml (Rust)
- **And** detects Gemfile (Ruby)
- **And** detects *.csproj (C#)

**AC-002.2**: Test Framework Detection from Dependencies
- **Given** dependency files are parsed
- **When** Ship Shape analyzes dependencies
- **Then** it identifies test frameworks (pytest, jest, vitest, testing, junit, etc.)
- **And** identifies coverage tools (coverage.py, c8, gocov, jacoco)
- **And** identifies quality tools (eslint, pylint, golangci-lint)
- **And** identifies formatters (black, prettier, gofmt)

**AC-002.3**: Build System Detection
- **Given** a repository structure
- **When** Ship Shape looks for build files
- **Then** it detects Makefile, CMakeLists.txt
- **And** detects go modules, npm scripts
- **And** detects Maven, Gradle
- **And** detects Cargo

**AC-002.4**: CI Configuration Detection
- **Given** a repository with CI setup
- **When** Ship Shape scans .github/workflows/*.yml
- **Then** it detects GitHub Actions workflows
- **And** extracts test job configurations
- **And** identifies CI-based testing commands

#### Technical Requirements

**TR-002.1**: Configuration parsers
```go
// Parse package.json
func parsePackageJSON(path string) (*PackageConfig, error)

// Parse pyproject.toml
func parsePyprojectToml(path string) (*PythonConfig, error)

// Parse go.mod
func parseGoMod(path string) (*GoConfig, error)

// Parse GitHub Actions workflows
func parseGitHubWorkflow(path string) (*WorkflowConfig, error)
```

**TR-002.2**: Use standard library and minimal dependencies
- encoding/json for JSON parsing
- github.com/pelletier/go-toml/v2 for TOML (MIT license)
- gopkg.in/yaml.v3 for YAML (Apache 2.0)

**TR-002.3**: Framework version detection
- Extract version constraints from dependencies
- Identify if using latest vs. outdated versions
- Flag deprecated frameworks

#### Open Source Tools
- **go-toml**: TOML parser (MIT)
- **yaml.v3**: YAML parser (Apache 2.0)
- **encoding/json**: Go standard library

#### Test Requirements
- Unit tests for each parser function
- Test with various package manager formats
- Test version extraction accuracy
- Test malformed configuration file handling

---

### Story SS-003: Monorepo Structure Detection

**As a** developer working in a monorepo
**I want** Ship Shape to detect workspace configurations and package boundaries
**So that** it analyzes each package independently and provides aggregate scores

**Priority**: P0 (Critical)
**Story Points**: 8
**Dependencies**: SS-001, SS-002

#### Acceptance Criteria

**AC-003.1**: npm Workspaces Detection
- **Given** a repository with package.json containing "workspaces" field
- **When** Ship Shape analyzes the root package.json
- **Then** it identifies all workspace patterns (e.g., "packages/*", "apps/*")
- **And** resolves workspace patterns to actual package directories
- **And** creates a package manifest with paths and names

**AC-003.2**: Yarn/pnpm Workspaces Detection
- **Given** pnpm-workspace.yaml or yarn workspaces config
- **When** Ship Shape parses workspace configuration
- **Then** it identifies all packages in the workspace
- **And** handles both wildcard patterns and explicit paths

**AC-003.3**: Go Multi-Module Workspaces
- **Given** a Go repository with go.work file
- **When** Ship Shape parses go.work
- **Then** it identifies all module paths
- **And** treats each module as separate package
- **And** identifies shared dependencies

**AC-003.4**: Maven/Gradle Multi-Module Projects
- **Given** a Java project with multi-module structure
- **When** Ship Shape parses pom.xml or settings.gradle
- **Then** it identifies all submodules
- **And** maps module dependencies

**AC-003.5**: Python Monorepo Detection (Heuristic)
- **Given** multiple Python packages without explicit workspace config
- **When** Ship Shape finds multiple setup.py or pyproject.toml files
- **Then** it treats each as a separate package
- **And** attempts to identify shared code areas

**AC-003.6**: Lerna/Nx/Turborepo Detection
- **Given** a JavaScript monorepo using task runners
- **When** Ship Shape finds lerna.json, nx.json, or turbo.json
- **Then** it uses these configs to identify packages
- **And** respects package boundaries defined in config

#### Technical Requirements

**TR-003.1**: Workspace resolver
```go
type MonorepoDetector struct {
    detectors []MonorepoDetector
}

func (md *MonorepoDetector) Detect(repoPath string) (*MonorepoInfo, error) {
    // Try npm workspaces
    if info := detectNPMWorkspaces(repoPath); info != nil {
        return info, nil
    }

    // Try pnpm workspaces
    if info := detectPNPMWorkspaces(repoPath); info != nil {
        return info, nil
    }

    // Try Go workspaces
    if info := detectGoWorkspaces(repoPath); info != nil {
        return info, nil
    }

    // ... other detectors

    return nil, nil // Not a monorepo
}
```

**TR-003.2**: Glob pattern resolution
- Use github.com/gobwas/glob for workspace pattern matching
- Handle nested wildcards (packages/**/src)
- Resolve symlinks correctly

**TR-003.3**: Package graph construction
- Build dependency graph between packages
- Identify shared infrastructure (root-level configs)
- Detect circular dependencies

#### Open Source Tools
- **gobwas/glob**: Pattern matching (MIT)
- **yaml.v3**: YAML parsing for workspace configs
- Standard library: encoding/json, path/filepath

#### Test Requirements
- Test each monorepo type independently
- Test complex workspace patterns
- Test nested monorepos
- Test package dependency resolution
- Performance: Resolve 100 packages in <2 seconds

---

### Story SS-004: Directory Structure Mapping

**As a** developer
**I want** Ship Shape to understand my project's directory structure
**So that** it knows where source code, tests, and configs are located

**Priority**: P0 (Critical)
**Story Points**: 5
**Dependencies**: SS-001

#### Acceptance Criteria

**AC-004.1**: Standard Directory Detection
- **Given** a repository with conventional structure
- **When** Ship Shape maps directories
- **Then** it identifies source directories (src/, lib/, pkg/)
- **And** identifies test directories (test/, tests/, __tests__/)
- **And** identifies config directories (config/, configs/, .config/)
- **And** identifies documentation directories (docs/, doc/)

**AC-004.2**: Language-Specific Conventions
- **Given** different language projects
- **When** Ship Shape analyzes structure
- **Then** it applies Go conventions (cmd/, pkg/, internal/)
- **And** applies Python conventions (src/, tests/, setup.py)
- **And** applies Java conventions (src/main/java/, src/test/java/)
- **And** applies JavaScript conventions (src/, test/, dist/)

**AC-004.3**: Lines of Code Calculation
- **Given** source and test directories are identified
- **When** Ship Shape calculates metrics
- **Then** it counts total lines of code (excluding blanks/comments)
- **And** counts test lines separately
- **And** calculates test-to-source ratio
- **And** excludes vendor/generated code

**AC-004.4**: Test File Pattern Matching
- **Given** identified test directories
- **When** Ship Shape looks for test files
- **Then** it finds *_test.go (Go)
- **And** finds test_*.py, *_test.py (Python)
- **And** finds *.test.js, *.spec.js, *.test.ts (JavaScript/TypeScript)
- **And** finds *Test.java, *Tests.java (Java)
- **And** finds *_test.rs (Rust)

#### Technical Requirements

**TR-004.1**: Structure analyzer
```go
type StructureAnalyzer struct {
    languageRules map[string]*DirectoryRules
}

type DirectoryRules struct {
    SourcePatterns []string
    TestPatterns   []string
    ConfigPatterns []string
}

func (sa *StructureAnalyzer) AnalyzeStructure(repoPath string, languages []*LanguageInfo) (*StructureInfo, error) {
    // Apply language-specific rules
    // Map directories to categories
    // Calculate LOC metrics
}
```

**TR-004.2**: LOC counter (exclude comments/blanks)
- Use simple line-by-line parsing
- Language-specific comment detection
- Stream processing for large files

**TR-004.3**: Exclusion patterns
```go
var excludePatterns = []string{
    "**/node_modules/**",
    "**/vendor/**",
    "**/.git/**",
    "**/dist/**",
    "**/build/**",
    "**/__pycache__/**",
    "**/target/**",
    "**/*.min.js",
    "**/*.generated.*",
}
```

#### Open Source Tools
- **filepath.Walk**: Go standard library
- **strings**: Standard library for line processing

#### Test Requirements
- Test with various language project structures
- Test LOC counting accuracy
- Test exclusion pattern effectiveness
- Benchmark LOC counting performance (100k lines in <1s)

---

## Epic 2: Multi-Language Analysis

**Epic Goal**: Provide deep, language-specific analysis for Python, JavaScript/TypeScript, Go, Java, Rust, Ruby, C#, and C/C++.

### Story SS-010: Go Test Analysis with AST Parsing

**As a** Go developer
**I want** Ship Shape to analyze my Go tests for best practices
**So that** I can improve test quality and maintainability

**Priority**: P0 (Critical)
**Story Points**: 13
**Dependencies**: SS-001, SS-004

#### Acceptance Criteria

**AC-010.1**: Table-Driven Test Detection
- **Given** a Go test file with table-driven tests
- **When** Ship Shape parses the AST
- **Then** it identifies table-driven test patterns
- **And** counts test cases per table
- **And** rewards high test case coverage

**AC-010.2**: Parallel Test Detection
- **Given** Go tests using t.Parallel()
- **When** Ship Shape analyzes test functions
- **Then** it identifies parallel tests
- **And** calculates percentage of parallel tests
- **And** flags missing t.Parallel() in suitable tests

**AC-010.3**: Subtests Detection
- **Given** tests using t.Run()
- **When** Ship Shape parses test structure
- **Then** it identifies subtest usage
- **And** validates subtest naming consistency
- **And** checks for proper test organization

**AC-010.4**: Test Helper Detection
- **Given** test files with helper functions
- **When** Ship Shape analyzes code
- **Then** it identifies t.Helper() usage
- **And** flags helpers missing t.Helper() calls
- **And** validates helper function naming (testHelper*)

**AC-010.5**: Assertion Quality
- **Given** Go tests with assertions
- **When** Ship Shape examines test logic
- **Then** it identifies assertion libraries (testify, gomega)
- **And** flags deep.Equal usage (prefer cmp or testify)
- **And** counts assertions per test

**AC-010.6**: Benchmark Detection
- **Given** benchmark functions (BenchmarkXxx)
- **When** Ship Shape scans for benchmarks
- **Then** it counts benchmark coverage
- **And** validates b.ReportAllocs() usage
- **And** checks for proper benchmark structure

#### Technical Requirements

**TR-010.1**: AST parsing with go/ast
```go
import (
    "go/ast"
    "go/parser"
    "go/token"
)

type GoTestAnalyzer struct {
    fset *token.FileSet
}

func (gta *GoTestAnalyzer) AnalyzeTestFile(path string) (*TestFileAnalysis, error) {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }

    // Walk AST to find test patterns
    ast.Inspect(node, func(n ast.Node) bool {
        switch x := n.(type) {
        case *ast.FuncDecl:
            if isTestFunction(x) {
                analyzeTestFunction(x)
            }
        }
        return true
    })
}
```

**TR-010.2**: Pattern detection functions
```go
func isTableDrivenTest(fn *ast.FuncDecl) bool {
    // Look for []struct{} or []testCase patterns
}

func hasParallelCall(fn *ast.FuncDecl) bool {
    // Look for t.Parallel() call
}

func countSubtests(fn *ast.FuncDecl) int {
    // Count t.Run() calls
}

func hasHelperCall(fn *ast.FuncDecl) bool {
    // Look for t.Helper() call
}
```

**TR-010.3**: Findings generation
```go
// Generate findings for detected issues
if !hasParallelCall(testFunc) && isSuitableForParallel(testFunc) {
    findings = append(findings, &Finding{
        Type:     FindingTypeBestPractice,
        Severity: SeverityLow,
        Title:    "Test could use t.Parallel()",
        Location: getFunctionLocation(testFunc),
    })
}
```

#### Open Source Tools
- **go/ast**: Go standard library for AST parsing
- **go/parser**: Go standard library for parsing
- **go/token**: Go standard library for position tracking

#### Test Requirements
- Unit tests for each pattern detector
- Test with real Go projects (e.g., analyze golang.org/x/ repos)
- Validate against known good/bad test patterns
- Performance: Parse 100 test files in <3 seconds

---

### Story SS-011: Python Test Analysis (pytest/unittest)

**As a** Python developer
**I want** Ship Shape to analyze my Python tests for quality and best practices
**So that** I can maintain high-quality test suites

**Priority**: P0 (Critical)
**Story Points**: 13
**Dependencies**: SS-001, SS-004

#### Acceptance Criteria

**AC-011.1**: pytest Test Detection
- **Given** Python test files using pytest
- **When** Ship Shape analyzes test functions
- **Then** it identifies test_* functions
- **And** identifies pytest fixtures usage
- **And** detects parametrized tests (@pytest.mark.parametrize)
- **And** counts test cases per parametrized test

**AC-011.2**: unittest Test Detection
- **Given** Python tests using unittest
- **When** Ship Shape parses test classes
- **Then** it identifies TestCase subclasses
- **And** identifies test_* methods
- **And** detects setUp/tearDown methods
- **And** identifies setUpClass/tearDownClass

**AC-011.3**: Fixture Quality Analysis
- **Given** pytest fixtures
- **When** Ship Shape examines fixture definitions
- **Then** it validates fixture scope usage (function, class, module, session)
- **And** identifies fixture dependencies
- **And** flags fixtures with excessive scope

**AC-011.4**: Test Organization
- **Given** test files and structure
- **When** Ship Shape analyzes organization
- **Then** it validates test file naming (test_*.py, *_test.py)
- **And** checks for conftest.py presence and usage
- **And** validates test class organization

**AC-011.5**: Assertion Usage
- **Given** test functions
- **When** Ship Shape counts assertions
- **Then** it identifies assertion types (assert, assertEqual, etc.)
- **And** flags tests with zero assertions
- **And** counts average assertions per test

**AC-011.6**: Mock Usage Detection
- **Given** tests using mocks
- **When** Ship Shape analyzes imports and usage
- **Then** it identifies unittest.mock usage
- **And** identifies pytest-mock usage
- **And** flags excessive mocking (>5 mocks per test)

#### Technical Requirements

**TR-011.1**: Python AST parsing
```go
// Call Python script to parse AST and return JSON
func parsePythonAST(filePath string) (*PythonAST, error) {
    cmd := exec.Command("python3", "-c", astParserScript, filePath)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var ast PythonAST
    json.Unmarshal(output, &ast)
    return &ast, nil
}

// Alternative: Use tree-sitter for parsing
```

**TR-011.2**: Python helper script (embedded)
```python
# Embedded in Go binary as string
import ast
import json
import sys

def analyze_test_file(filepath):
    with open(filepath) as f:
        tree = ast.parse(f.read())

    analysis = {
        'functions': [],
        'classes': [],
        'fixtures': [],
        'imports': []
    }

    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            if node.name.startswith('test_'):
                analysis['functions'].append(analyze_function(node))
        # ... more analysis

    print(json.dumps(analysis))

if __name__ == '__main__':
    analyze_test_file(sys.argv[1])
```

**TR-011.3**: Alternative: tree-sitter
```go
import sitter "github.com/smacker/go-tree-sitter"
import "github.com/smacker/go-tree-sitter/python"

func parsePythonWithTreeSitter(content []byte) (*sitter.Node, error) {
    parser := sitter.NewParser()
    parser.SetLanguage(python.GetLanguage())
    tree := parser.Parse(nil, content)
    return tree.RootNode(), nil
}
```

#### Open Source Tools
- **Python AST**: Python standard library (for helper script)
- **tree-sitter**: Alternative parser (MIT)
- **smacker/go-tree-sitter**: Go bindings (MIT)

#### Test Requirements
- Test with pytest and unittest projects
- Test parametrized test detection
- Test fixture analysis accuracy
- Test with popular Python projects (requests, flask, django)
- Performance: Analyze 50 test files in <5 seconds

---

### Story SS-012: JavaScript/TypeScript Test Analysis (Jest/Vitest)

**As a** JavaScript/TypeScript developer
**I want** Ship Shape to analyze my Jest/Vitest tests
**So that** I can ensure high-quality test coverage and practices

**Priority**: P0 (Critical)
**Story Points**: 13
**Dependencies**: SS-001, SS-004

#### Acceptance Criteria

**AC-012.1**: Jest/Vitest Test Detection
- **Given** JavaScript/TypeScript test files
- **When** Ship Shape analyzes test structure
- **Then** it identifies describe() blocks
- **And** identifies test()/it() functions
- **And** counts nested describe blocks
- **And** validates test organization

**AC-012.2**: Test Hook Detection
- **Given** tests using setup/teardown hooks
- **When** Ship Shape examines test files
- **Then** it identifies beforeEach/afterEach
- **And** identifies beforeAll/afterAll
- **And** validates hook usage patterns
- **And** flags missing cleanup in afterEach

**AC-012.3**: Mock and Spy Analysis
- **Given** tests using mocks
- **When** Ship Shape analyzes mocking patterns
- **Then** it identifies jest.mock() usage
- **And** identifies vi.mock() for Vitest
- **And** counts jest.spyOn() calls
- **And** flags excessive mocking

**AC-012.4**: Snapshot Testing Detection
- **Given** tests using snapshots
- **When** Ship Shape finds expect().toMatchSnapshot()
- **Then** it counts snapshot tests
- **And** validates snapshot file existence
- **And** flags overuse of snapshots (>30% of tests)

**AC-012.5**: Async Test Handling
- **Given** async tests
- **When** Ship Shape analyzes test functions
- **Then** it identifies async/await usage
- **And** identifies Promise-based tests
- **And** flags missing await keywords
- **And** validates proper async error handling

**AC-012.6**: Test Organization and Naming
- **Given** test files and suites
- **When** Ship Shape checks naming
- **Then** it validates test file patterns (*.test.js, *.spec.ts)
- **And** checks for descriptive test names
- **And** validates describe block nesting depth (<4)

#### Technical Requirements

**TR-012.1**: TypeScript/JavaScript AST parsing
```go
// Use tree-sitter for JS/TS parsing
import (
    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/javascript"
    "github.com/smacker/go-tree-sitter/typescript/tsx"
)

func parseJavaScriptTest(filePath string) (*JSTestAnalysis, error) {
    content, _ := os.ReadFile(filePath)

    parser := sitter.NewParser()
    if strings.HasSuffix(filePath, ".ts") || strings.HasSuffix(filePath, ".tsx") {
        parser.SetLanguage(tsx.GetLanguage())
    } else {
        parser.SetLanguage(javascript.GetLanguage())
    }

    tree := parser.Parse(nil, content)
    root := tree.RootNode()

    return analyzeJSTestTree(root, content)
}
```

**TR-012.2**: Pattern detection
```go
func findDescribeBlocks(node *sitter.Node) []*DescribeBlock {
    // Query for call_expression where function is "describe"
}

func findTestFunctions(node *sitter.Node) []*TestFunction {
    // Query for call_expression where function is "test" or "it"
}

func findMockCalls(node *sitter.Node) []*MockCall {
    // Query for jest.mock, vi.mock, jest.spyOn
}
```

**TR-012.3**: Tree-sitter queries
```go
const testQuery = `
    (call_expression
      function: (identifier) @func_name (#match? @func_name "^(test|it|describe)$")
      arguments: (arguments) @args
    )
`
```

#### Open Source Tools
- **tree-sitter**: Universal parser (MIT)
- **go-tree-sitter**: Go bindings (MIT)
- **tree-sitter-javascript**: JavaScript grammar (MIT)
- **tree-sitter-typescript**: TypeScript grammar (MIT)

#### Test Requirements
- Test with Jest and Vitest projects
- Test TypeScript and JavaScript variants
- Test React Testing Library patterns
- Validate with popular projects (react, vue, svelte)
- Performance: Analyze 100 test files in <10 seconds

---

## Epic 3: Monorepo Support

**Epic Goal**: Provide comprehensive analysis for monorepo projects with multi-language support and aggregate scoring.

### Story SS-020: Monorepo Package-Level Analysis

**As a** developer in a monorepo
**I want** Ship Shape to analyze each package independently
**So that** I can see quality metrics for individual packages

**Priority**: P0 (Critical)
**Story Points**: 13
**Dependencies**: SS-003, SS-010, SS-011, SS-012

#### Acceptance Criteria

**AC-020.1**: Independent Package Analysis
- **Given** a monorepo with multiple packages
- **When** Ship Shape runs analysis
- **Then** it analyzes each package as a separate unit
- **And** generates individual scores for each package
- **And** produces separate findings per package
- **And** calculates package-specific metrics

**AC-020.2**: Parallel Package Processing
- **Given** a monorepo with N packages
- **When** Ship Shape analyzes the monorepo
- **Then** it processes packages in parallel (up to 4 concurrent)
- **And** respects package dependency order
- **And** fails gracefully if one package fails
- **And** continues analyzing remaining packages

**AC-020.3**: Package Context Isolation
- **Given** packages with different languages/frameworks
- **When** Ship Shape analyzes each package
- **Then** it detects language per package
- **And** applies package-specific analyzers
- **And** uses package-specific configuration if present

**AC-020.4**: Shared Infrastructure Detection
- **Given** a monorepo with root-level configs
- **When** Ship Shape analyzes shared infrastructure
- **Then** it identifies root package.json (workspace root)
- **And** identifies shared ESLint/Prettier configs
- **And** identifies shared GitHub Actions workflows
- **And** scores shared infrastructure quality

#### Technical Requirements

**TR-020.1**: Monorepo coordinator
```go
type MonorepoCoordinator struct {
    maxConcurrency int
    analyzer       *AnalysisEngine
}

func (mc *MonorepoCoordinator) AnalyzeMonorepo(
    ctx context.Context,
    repoCtx *RepositoryContext,
) (*MonorepoReport, error) {
    packages := repoCtx.Monorepo.Packages

    // Create semaphore for concurrency control
    sem := make(chan struct{}, mc.maxConcurrency)
    results := make(chan *PackageResult, len(packages))

    // Launch goroutines for each package
    for _, pkg := range packages {
        sem <- struct{}{} // Acquire
        go func(p *PackageInfo) {
            defer func() { <-sem }() // Release

            result := mc.analyzePackage(ctx, p)
            results <- result
        }(pkg)
    }

    // Collect results
    packageReports := make(map[string]*Report)
    for i := 0; i < len(packages); i++ {
        result := <-results
        if result.Error != nil {
            log.Warn("Package analysis failed", result.PackageName, result.Error)
            continue
        }
        packageReports[result.PackageName] = result.Report
    }

    return mc.generateMonorepoReport(packageReports), nil
}
```

**TR-020.2**: Package analysis with isolated context
```go
func (mc *MonorepoCoordinator) analyzePackage(
    ctx context.Context,
    pkg *PackageInfo,
) *PackageResult {
    // Create package-specific context
    pkgCtx := &RepositoryContext{
        RootPath:  pkg.Path,
        Type:      TypeSinglePackage,
        Languages: detectLanguages(pkg.Path),
        // ... isolated context
    }

    // Run full analysis pipeline for package
    results, err := mc.analyzer.Analyze(ctx, pkgCtx)
    return &PackageResult{
        PackageName: pkg.Name,
        Report:      generateReport(results),
        Error:       err,
    }
}
```

**TR-020.3**: Shared infrastructure analyzer
```go
func (mc *MonorepoCoordinator) analyzeSharedInfra(
    repoCtx *RepositoryContext,
) *SharedInfraScore {
    score := &SharedInfraScore{}

    // Check for root-level quality configs
    if hasRootESLint(repoCtx.RootPath) {
        score.ConfigQuality += 20
    }

    // Check for shared CI workflows
    if hasSharedCI(repoCtx.RootPath) {
        score.CIIntegration += 30
    }

    // Check for shared dependencies
    sharedDeps := analyzeSharedDependencies(repoCtx)
    score.SharedDeps = sharedDeps.Score

    score.Overall = (score.ConfigQuality + score.CIIntegration + score.SharedDeps) / 3
    return score
}
```

#### Open Source Tools
- **goroutines**: Go concurrency (standard library)
- **context**: Go context for cancellation (standard library)

#### Test Requirements
- Test with npm workspaces monorepo
- Test with Go multi-module workspace
- Test parallel processing with 10+ packages
- Test failure isolation (one package fails, others continue)
- Performance: Analyze 20 packages in <30 seconds

---

### Story SS-021: Monorepo Aggregate Scoring

**As a** monorepo maintainer
**I want** Ship Shape to provide an aggregate score for the entire monorepo
**So that** I can track overall repository quality

**Priority**: P0 (Critical)
**Story Points**: 5
**Dependencies**: SS-020

#### Acceptance Criteria

**AC-021.1**: Weighted Aggregate Score
- **Given** individual package scores
- **When** Ship Shape calculates aggregate
- **Then** it weights scores by package size (LOC)
- **And** larger packages have more influence
- **And** aggregate score is between 0-100

**AC-021.2**: Dimension-Level Aggregation
- **Given** dimension scores per package
- **When** Ship Shape aggregates dimensions
- **Then** it calculates weighted average per dimension
- **And** shows dimension score distribution across packages
- **And** identifies dimension outliers

**AC-021.3**: Consistency Analysis
- **Given** scores across multiple packages
- **When** Ship Shape analyzes consistency
- **Then** it calculates score variance
- **And** flags high variance as inconsistency risk
- **And** identifies tool adoption inconsistencies

**AC-021.4**: Aggregate Findings
- **Given** findings from all packages
- **When** Ship Shape generates monorepo report
- **Then** it groups findings by severity
- **And** shows finding count per package
- **And** identifies monorepo-wide patterns
- **And** prioritizes cross-package issues

#### Technical Requirements

**TR-021.1**: Weighted scoring algorithm
```go
func calculateAggregateScore(packageReports map[string]*Report) *ScoreCard {
    totalLOC := 0
    weightedScore := 0.0

    for _, report := range packageReports {
        pkgLOC := report.RepositoryContext.Structure.TotalLOC
        totalLOC += int(pkgLOC)
        weightedScore += report.Scores.Overall * float64(pkgLOC)
    }

    if totalLOC == 0 {
        return &ScoreCard{Overall: 0, Grade: "F"}
    }

    overallScore := weightedScore / float64(totalLOC)

    return &ScoreCard{
        Overall: overallScore,
        Grade:   calculateGrade(overallScore),
        Dimensions: aggregateDimensions(packageReports, totalLOC),
    }
}
```

**TR-021.2**: Consistency metrics
```go
func analyzeConsistency(packageReports map[string]*Report) *ConsistencyAnalysis {
    scores := extractScores(packageReports)

    return &ConsistencyAnalysis{
        ScoreVariance:     calculateVariance(scores),
        ToolConsistency:   analyzeToolUsage(packageReports),
        FrameworkSpread:   analyzeFrameworks(packageReports),
        QualityRange:      Range{Min: min(scores), Max: max(scores)},
    }
}
```

#### Open Source Tools
- **gonum/stat**: Statistical functions (BSD-3-Clause)

#### Test Requirements
- Test weighted scoring with various package sizes
- Test variance calculation accuracy
- Test with uniform vs. heterogeneous package quality
- Validate scoring matches manual calculations

---

## Epic 4: Test Quality Analysis

**Epic Goal**: Deep analysis of test quality, including test smells, organization, and best practices.

### Story SS-030: Test Smell Detection

**As a** developer
**I want** Ship Shape to detect test smells and anti-patterns
**So that** I can refactor poor-quality tests

**Priority**: P0 (Critical)
**Story Points**: 13
**Dependencies**: SS-010, SS-011, SS-012

#### Acceptance Criteria

**AC-030.1**: Mystery Guest Detection
- **Given** tests that rely on external resources
- **When** Ship Shape analyzes test dependencies
- **Then** it identifies file reads without explicit setup
- **And** identifies database calls in unit tests
- **And** identifies HTTP calls without mocks
- **And** flags as "Mystery Guest" smell

**AC-030.2**: Eager Test Detection
- **Given** tests testing multiple concerns
- **When** Ship Shape counts assertions and behaviors
- **Then** it identifies tests with >5 assertions
- **And** identifies tests testing unrelated functionality
- **And** flags as "Eager Test" smell
- **And** suggests splitting into multiple tests

**AC-030.3**: Lazy Test Detection
- **Given** test functions with multiple test cases
- **When** Ship Shape analyzes test structure
- **Then** it identifies single test function with many scenarios
- **And** suggests converting to table-driven tests (Go)
- **And** suggests parametrized tests (Python, Jest)

**AC-030.4**: Assertion Roulette Detection
- **Given** tests with multiple assertions
- **When** Ship Shape examines assertion messages
- **Then** it identifies assertions without messages
- **And** flags tests where assertion failure would be ambiguous
- **And** requires descriptive assertion messages

**AC-030.5**: Test Code Duplication
- **Given** multiple test files
- **When** Ship Shape analyzes test setup code
- **Then** it identifies duplicated setup logic
- **And** suggests extracting to test helpers/fixtures
- **And** calculates duplication percentage

**AC-030.6**: Sleepy Test Detection
- **Given** tests using sleep/delay
- **When** Ship Shape finds time.Sleep, setTimeout, etc.
- **Then** it flags as "Sleepy Test" anti-pattern
- **And** suggests using polling or explicit waits
- **And** marks as flaky test risk

#### Technical Requirements

**TR-030.1**: Test smell detector framework
```go
type TestSmellDetector interface {
    Name() string
    Detect(testFile *TestFileAST) []*Finding
    Severity() Severity
}

type MysteryGuestDetector struct{}

func (mgd *MysteryGuestDetector) Detect(testFile *TestFileAST) []*Finding {
    findings := []*Finding{}

    for _, testFunc := range testFile.Functions {
        // Check for file I/O without explicit setup
        if hasFileIO(testFunc) && !hasSetup(testFunc) {
            findings = append(findings, &Finding{
                Type:        FindingTypeTestSmell,
                Severity:    SeverityMedium,
                Title:       "Mystery Guest: Test depends on external file",
                Description: "Test reads file without explicit setup, making it unclear what data is needed",
                Location:    testFunc.Location,
            })
        }
    }

    return findings
}
```

**TR-030.2**: Pattern matching for each smell
```go
// Eager Test: count distinct concerns
func countTestConcerns(testFunc *TestFunction) int {
    // Heuristic: count number of different objects/functions tested
}

// Lazy Test: detect multiple test cases in one function
func isLazyTest(testFunc *TestFunction) bool {
    // Look for if/switch on test data without table-driven pattern
}

// Sleepy Test: find sleep calls
func hasSleepCall(testFunc *TestFunction) bool {
    // Search AST for time.Sleep, Thread.sleep, setTimeout, etc.
}
```

**TR-030.3**: Duplication detection
```go
// Simple token-based similarity
func calculateSimilarity(func1, func2 *TestFunction) float64 {
    tokens1 := tokenize(func1.Code)
    tokens2 := tokenize(func2.Code)

    // Jaccard similarity
    intersection := countCommon(tokens1, tokens2)
    union := len(tokens1) + len(tokens2) - intersection

    return float64(intersection) / float64(union)
}

// Flag if >70% similar
```

#### Open Source Tools
- **AST parsers**: go/ast, tree-sitter (from previous stories)
- **diff**: Go standard library for duplication detection

#### Test Requirements
- Test each smell detector independently
- Validate against known test smell examples
- Test false positive rate (<5%)
- Test with real-world codebases
- Performance: Analyze 100 tests in <5 seconds

---

## Epic 5: Coverage Analysis

**Epic Goal**: Comprehensive coverage analysis with branch coverage, mutation testing readiness, and coverage quality assessment.

### Story SS-040: Coverage Report Parsing

**As a** developer
**I want** Ship Shape to parse coverage reports from various tools
**So that** I can see unified coverage metrics across languages

**Priority**: P0 (Critical)
**Story Points**: 8
**Dependencies**: SS-001, SS-002

#### Acceptance Criteria

**AC-040.1**: Go Coverage Parsing (gocov, go test -cover)
- **Given** Go coverage output
- **When** Ship Shape parses coverage data
- **Then** it parses `go test -coverprofile` output
- **And** parses gocov JSON format
- **And** extracts line coverage percentage
- **And** extracts branch coverage if available

**AC-040.2**: Python Coverage Parsing (coverage.py)
- **Given** Python coverage.xml or .coverage file
- **When** Ship Shape parses coverage
- **Then** it parses XML coverage format
- **And** parses JSON format from `coverage json`
- **And** extracts line and branch coverage
- **And** identifies uncovered lines

**AC-040.3**: JavaScript Coverage Parsing (c8, Istanbul)
- **Given** JavaScript coverage reports
- **When** Ship Shape processes coverage
- **Then** it parses lcov format
- **And** parses JSON format from c8/nyc
- **And** parses Istanbul JSON
- **And** extracts statement, branch, function, line coverage

**AC-040.4**: Java Coverage Parsing (JaCoCo)
- **Given** JaCoCo XML reports
- **When** Ship Shape reads coverage
- **Then** it parses jacoco.xml format
- **And** extracts instruction, branch, line, method coverage
- **And** calculates overall coverage percentage

**AC-040.5**: Unified Coverage Model
- **Given** coverage from multiple languages
- **When** Ship Shape normalizes data
- **Then** it creates unified coverage model
- **And** calculates weighted average across languages
- **And** identifies coverage gaps per language

#### Technical Requirements

**TR-040.1**: Coverage parsers per format
```go
type CoverageParser interface {
    Parse(filePath string) (*CoverageData, error)
    SupportedFormats() []string
}

type GoCoverageParser struct{}

func (gcp *GoCoverageParser) Parse(filePath string) (*CoverageData, error) {
    // Parse coverage.out format
    // Example line: "github.com/user/repo/pkg/file.go:10.2,12.3 2 1"
    // Format: file:startLine.startCol,endLine.endCol numStatements count
}

type PythonCoverageParser struct{}

func (pcp *PythonCoverageParser) Parse(filePath string) (*CoverageData, error) {
    // Parse coverage.xml (Cobertura format)
    // Or parse coverage.json
}

type JavaScriptCoverageParser struct{}

func (jcp *JavaScriptCoverageParser) Parse(filePath string) (*CoverageData, error) {
    // Parse lcov.info or coverage-final.json
}
```

**TR-040.2**: Unified coverage model
```go
type CoverageData struct {
    TotalLines      int
    CoveredLines    int
    LinePercentage  float64

    TotalBranches   int
    CoveredBranches int
    BranchPercentage float64

    Files           []*FileCoverage
}

type FileCoverage struct {
    Path            string
    Lines           int
    CoveredLines    int
    Percentage      float64
    UncoveredRanges []*LineRange
}
```

**TR-040.3**: XML/JSON parsing
```go
// Use encoding/xml and encoding/json from standard library
import (
    "encoding/xml"
    "encoding/json"
)

// For complex XML (JaCoCo), define structs
type JaCoCoReport struct {
    XMLName  xml.Name `xml:"report"`
    Packages []struct {
        Name    string `xml:"name,attr"`
        Classes []struct {
            Name    string `xml:"name,attr"`
            Methods []struct {
                Name     string `xml:"name,attr"`
                Coverage struct {
                    Covered int `xml:"covered,attr"`
                    Missed  int `xml:"missed,attr"`
                } `xml:"counter"`
            } `xml:"method"`
        } `xml:"class"`
    } `xml:"package"`
}
```

#### Open Source Tools
- **encoding/xml**: Go standard library
- **encoding/json**: Go standard library
- **bufio**: For line-by-line parsing of coverage.out

#### Test Requirements
- Test with real coverage files from each language
- Test malformed coverage file handling
- Test partial coverage data
- Validate percentage calculations
- Performance: Parse 100 coverage files in <2 seconds

---

### Story SS-041: Coverage Quality Assessment

**As a** developer
**I want** Ship Shape to assess the quality of my test coverage
**So that** I know if my coverage is meaningful or superficial

**Priority**: P1 (High)
**Story Points**: 8
**Dependencies**: SS-040

#### Acceptance Criteria

**AC-041.1**: Branch Coverage Analysis
- **Given** coverage data with branch information
- **When** Ship Shape assesses coverage quality
- **Then** it calculates branch coverage percentage
- **And** compares to line coverage percentage
- **And** flags if branch coverage < line coverage - 15%
- **And** identifies files with low branch coverage

**AC-041.2**: Critical Path Coverage
- **Given** source code with identified critical paths
- **When** Ship Shape analyzes coverage
- **Then** it identifies uncovered error handling
- **And** identifies uncovered edge cases
- **And** prioritizes coverage gaps by criticality

**AC-041.3**: Coverage Trends Over Time
- **Given** historical coverage data
- **When** Ship Shape compares coverage
- **Then** it shows coverage trend (improving/declining)
- **And** flags significant coverage drops (>5%)
- **And** calculates coverage velocity

**AC-041.4**: Coverage Thresholds by Language
- **Given** language-specific coverage expectations
- **When** Ship Shape evaluates coverage
- **Then** it applies 80% threshold for Python
- **And** applies 80% threshold for Go
- **And** applies 70% threshold for JavaScript
- **And** applies 70% threshold for Java
- **And** flags coverage below thresholds

#### Technical Requirements

**TR-041.1**: Coverage quality scorer
```go
type CoverageAssessor struct {
    thresholds map[string]float64
}

func NewCoverageAssessor() *CoverageAssessor {
    return &CoverageAssessor{
        thresholds: map[string]float64{
            "Go":         80.0,
            "Python":     80.0,
            "JavaScript": 70.0,
            "TypeScript": 70.0,
            "Java":       70.0,
            "Rust":       75.0,
        },
    }
}

func (ca *CoverageAssessor) AssessCoverage(
    lang string,
    coverage *CoverageData,
) *CoverageScore {
    threshold := ca.thresholds[lang]

    score := &CoverageScore{
        LineCoverage:   coverage.LinePercentage,
        BranchCoverage: coverage.BranchPercentage,
        Threshold:      threshold,
    }

    // Calculate quality score (0-100)
    if coverage.LinePercentage >= threshold {
        score.Score = 100
    } else {
        score.Score = (coverage.LinePercentage / threshold) * 100
    }

    // Adjust for branch coverage
    if coverage.BranchPercentage > 0 {
        branchGap := coverage.LinePercentage - coverage.BranchPercentage
        if branchGap > 15 {
            score.Score -= 10 // Penalize for poor branch coverage
            score.Findings = append(score.Findings, &Finding{
                Title: "Branch coverage significantly lower than line coverage",
                Severity: SeverityMedium,
            })
        }
    }

    return score
}
```

**TR-041.2**: Critical path identifier (heuristic)
```go
func identifyCriticalPaths(coverage *CoverageData) []*CriticalPath {
    // Heuristics for critical paths:
    // 1. Error handling code (try/catch, if err != nil)
    // 2. Validation logic
    // 3. External API calls
    // 4. Database operations
    // 5. Authentication/authorization

    // Match uncovered lines to these patterns
}
```

#### Open Source Tools
- No external dependencies (pure Go)

#### Test Requirements
- Test threshold application per language
- Test branch coverage gap detection
- Test with projects at various coverage levels
- Validate scoring algorithm

---

## Epic 6: Tool Adoption & Recommendations

**Epic Goal**: Detect installed tools, assess tool adoption quality, and provide actionable recommendations.

### Story SS-050: Tool Database and Detection

**As a** developer
**I want** Ship Shape to detect which testing tools I'm using
**So that** it can provide relevant recommendations

**Priority**: P1 (High)
**Story Points**: 8
**Dependencies**: SS-002

#### Acceptance Criteria

**AC-050.1**: Tool Database Structure
- **Given** a YAML-based tool database
- **When** Ship Shape loads tool definitions
- **Then** it loads language-specific tool catalogs
- **And** loads tool metadata (name, category, priority, installation, benefits)
- **And** validates database schema on load

**AC-050.2**: Tool Detection by Dependencies
- **Given** a repository with dependency files
- **When** Ship Shape scans dependencies
- **Then** it matches dependencies to tool database
- **And** identifies exact versions installed
- **And** flags outdated tool versions

**AC-050.3**: Tool Detection by Config Files
- **Given** configuration files in repository
- **When** Ship Shape scans for configs
- **Then** it detects pytest.ini, jest.config.js, .eslintrc, etc.
- **And** matches configs to tools
- **And** validates tool is actually installed

**AC-050.4**: Tool Detection by Commands
- **Given** CI workflows or Makefiles
- **When** Ship Shape parses build commands
- **Then** it identifies tools from command usage
- **And** cross-references with dependencies

#### Technical Requirements

**TR-050.1**: Tool database schema (YAML)
```yaml
# data/tools/python.yml
language: Python
tools:
  - name: pytest
    category: testing
    priority: Essential
    description: Python testing framework with fixtures and plugins
    installation:
      methods:
        - pip: pip install pytest
        - pipx: pipx install pytest
    configuration:
      files:
        - pytest.ini
        - pyproject.toml
      example: |
        [tool.pytest.ini_options]
        addopts = "--strict-markers --cov=src"
        testpaths = ["tests"]
    detection:
      dependencies:
        - pytest
        - pytest-cov
        - pytest-xdist
      config_files:
        - pytest.ini
        - pyproject.toml
      imports:
        - pytest
    benefits:
      - Powerful fixture system
      - Extensive plugin ecosystem
      - Parametrized testing
    effort: Low
    references:
      - url: https://docs.pytest.org/
        title: pytest documentation
```

**TR-050.2**: Tool database loader
```go
type ToolDatabase struct {
    tools map[string][]*Tool // language -> tools
}

func LoadToolDatabase(dataDir string) (*ToolDatabase, error) {
    db := &ToolDatabase{
        tools: make(map[string][]*Tool),
    }

    // Load all YAML files from data/tools/
    files, _ := filepath.Glob(filepath.Join(dataDir, "tools", "*.yml"))
    for _, file := range files {
        langTools, err := parseToolFile(file)
        if err != nil {
            return nil, err
        }
        db.tools[langTools.Language] = langTools.Tools
    }

    return db, nil
}
```

**TR-050.3**: Tool detector
```go
type ToolDetector struct {
    database *ToolDatabase
}

func (td *ToolDetector) DetectTools(
    repoCtx *RepositoryContext,
) map[string]*DetectedTool {
    detected := make(map[string]*DetectedTool)

    for _, lang := range repoCtx.Languages {
        langTools := td.database.tools[lang.Name]

        for _, tool := range langTools {
            if td.isToolPresent(repoCtx, tool) {
                detected[tool.Name] = &DetectedTool{
                    Tool:    tool,
                    Version: td.detectVersion(repoCtx, tool),
                    ConfigFound: td.hasConfig(repoCtx, tool),
                }
            }
        }
    }

    return detected
}
```

#### Open Source Tools
- **gopkg.in/yaml.v3**: YAML parsing (Apache 2.0)

#### Test Requirements
- Test YAML parsing with all language files
- Test tool detection accuracy
- Test version detection
- Test config file detection
- Create comprehensive tool database for Python, Go, JS, Java

---

### Story SS-051: Tool Adoption Recommendations

**As a** developer
**I want** Ship Shape to recommend missing essential tools
**So that** I can improve my testing infrastructure

**Priority**: P1 (High)
**Story Points**: 5
**Dependencies**: SS-050

#### Acceptance Criteria

**AC-051.1**: Essential Tool Recommendations
- **Given** detected languages and missing tools
- **When** Ship Shape generates recommendations
- **Then** it recommends essential tools not yet adopted
- **And** prioritizes by tool priority (Essential > Recommended > Optional)
- **And** provides installation instructions
- **And** estimates adoption effort

**AC-051.2**: Tool Configuration Recommendations
- **Given** tools detected without configuration
- **When** Ship Shape analyzes tool setup
- **Then** it recommends adding configuration files
- **And** provides example configurations
- **And** explains configuration benefits

**AC-051.3**: Tool Upgrade Recommendations
- **Given** outdated tool versions
- **When** Ship Shape compares versions
- **Then** it recommends upgrading to latest stable
- **And** flags breaking changes if known
- **And** provides migration guides if available

**AC-051.4**: Tool Integration Recommendations
- **Given** tools with poor integration
- **When** Ship Shape analyzes tool usage
- **Then** it recommends CI integration if missing
- **And** recommends pre-commit hooks
- **And** recommends IDE integration

#### Technical Requirements

**TR-051.1**: Recommendation engine
```go
type RecommendationEngine struct {
    toolDB *ToolDatabase
}

func (re *RecommendationEngine) GenerateRecommendations(
    repoCtx *RepositoryContext,
    detectedTools map[string]*DetectedTool,
) []*Recommendation {
    recommendations := []*Recommendation{}

    for _, lang := range repoCtx.Languages {
        if !lang.IsPrimary {
            continue // Only recommend for primary languages
        }

        essentialTools := re.toolDB.GetEssentialTools(lang.Name)

        for _, tool := range essentialTools {
            if _, detected := detectedTools[tool.Name]; !detected {
                recommendations = append(recommendations, &Recommendation{
                    Type:        RecommendationTypeAdopt,
                    Priority:    tool.Priority,
                    Tool:        tool,
                    Title:       fmt.Sprintf("Adopt %s for %s testing", tool.Name, lang.Name),
                    Description: tool.Description,
                    Benefits:    tool.Benefits,
                    Effort:      tool.Effort,
                    Installation: tool.Installation,
                })
            }
        }
    }

    return recommendations
}
```

**TR-051.2**: Configuration checker
```go
func (re *RecommendationEngine) checkConfiguration(
    repoCtx *RepositoryContext,
    detectedTool *DetectedTool,
) *Recommendation {
    if !detectedTool.ConfigFound {
        return &Recommendation{
            Type:     RecommendationTypeConfigure,
            Priority: PriorityHigh,
            Tool:     detectedTool.Tool,
            Title:    fmt.Sprintf("Add configuration for %s", detectedTool.Tool.Name),
            Example:  detectedTool.Tool.Configuration.Example,
        }
    }
    return nil
}
```

#### Open Source Tools
- No external dependencies

#### Test Requirements
- Test recommendation generation for each language
- Test priority ordering
- Test with complete vs. incomplete tool setups
- Validate recommendation actionability

---

## Epic 9: CI/CD Integration (GitHub Actions)

**Epic Goal**: Seamless integration with GitHub Actions and CI/CD pipelines, including quality gates and automated reporting.

### Story SS-080: GitHub Actions Workflow Detection

**As a** developer using GitHub Actions
**I want** Ship Shape to analyze my CI workflows
**So that** it can assess CI/CD test quality

**Priority**: P0 (Critical)
**Story Points**: 8
**Dependencies**: SS-001

#### Acceptance Criteria

**AC-080.1**: Workflow File Discovery
- **Given** a GitHub repository
- **When** Ship Shape scans for CI configuration
- **Then** it finds all .github/workflows/*.yml files
- **And** parses each workflow file
- **And** extracts workflow metadata (name, triggers, jobs)

**AC-080.2**: Test Job Identification
- **Given** GitHub Actions workflows
- **When** Ship Shape analyzes jobs
- **Then** it identifies jobs that run tests (name contains "test", "check", "validate")
- **And** extracts test commands from job steps
- **And** identifies test frameworks used in CI

**AC-080.3**: Coverage Integration Detection
- **Given** test jobs in workflows
- **When** Ship Shape examines job steps
- **Then** it detects coverage report generation
- **And** detects coverage upload to codecov/coveralls
- **And** validates coverage thresholds in CI

**AC-080.4**: Parallelization Analysis
- **Given** workflow job configurations
- **When** Ship Shape analyzes execution strategy
- **Then** it detects matrix builds for parallel testing
- **And** identifies concurrent job execution
- **And** calculates potential parallelization gains

**AC-080.5**: Caching Analysis
- **Given** workflow job steps
- **When** Ship Shape looks for caching
- **Then** it detects actions/cache usage
- **And** validates cache keys
- **And** recommends caching if missing

#### Technical Requirements

**TR-080.1**: GitHub Actions workflow parser
```go
type GitHubWorkflow struct {
    Name string `yaml:"name"`
    On   interface{} `yaml:"on"` // Can be string or object
    Jobs map[string]*Job `yaml:"jobs"`
}

type Job struct {
    Name     string `yaml:"name"`
    RunsOn   interface{} `yaml:"runs-on"` // string or matrix
    Strategy *Strategy `yaml:"strategy"`
    Steps    []*Step `yaml:"steps"`
}

type Strategy struct {
    Matrix map[string]interface{} `yaml:"matrix"`
}

type Step struct {
    Name string `yaml:"name"`
    Uses string `yaml:"uses"`
    Run  string `yaml:"run"`
    With map[string]interface{} `yaml:"with"`
}

func parseGitHubWorkflow(filePath string) (*GitHubWorkflow, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var workflow GitHubWorkflow
    err = yaml.Unmarshal(data, &workflow)
    return &workflow, err
}
```

**TR-080.2**: Test job detector
```go
func isTestJob(job *Job) bool {
    testKeywords := []string{"test", "check", "validate", "verify", "spec"}

    jobName := strings.ToLower(job.Name)
    for _, keyword := range testKeywords {
        if strings.Contains(jobName, keyword) {
            return true
        }
    }

    // Check steps for test commands
    for _, step := range job.Steps {
        if containsTestCommand(step.Run) {
            return true
        }
    }

    return false
}

func containsTestCommand(command string) bool {
    testCommands := []string{"pytest", "go test", "npm test", "jest", "vitest"}
    cmdLower := strings.ToLower(command)

    for _, testCmd := range testCommands {
        if strings.Contains(cmdLower, testCmd) {
            return true
        }
    }

    return false
}
```

**TR-080.3**: CI analysis findings
```go
func (ga *GitHubActionsAnalyzer) Analyze(
    repoPath string,
) (*CIAnalysis, error) {
    analysis := &CIAnalysis{
        Platform:    "github-actions",
        Workflows:   []*WorkflowAnalysis{},
        Findings:    []*Finding{},
    }

    workflowFiles, _ := filepath.Glob(filepath.Join(repoPath, ".github/workflows/*.yml"))

    for _, file := range workflowFiles {
        workflow, _ := parseGitHubWorkflow(file)

        workflowAnalysis := ga.analyzeWorkflow(workflow)
        analysis.Workflows = append(analysis.Workflows, workflowAnalysis)

        // Generate findings
        if !workflowAnalysis.HasTestJob {
            analysis.Findings = append(analysis.Findings, &Finding{
                Type:     FindingTypeBestPractice,
                Severity: SeverityCritical,
                Title:    "No test job in GitHub Actions workflow",
                Location: &Location{FilePath: file},
            })
        }
    }

    return analysis, nil
}
```

#### Open Source Tools
- **gopkg.in/yaml.v3**: YAML parsing (Apache 2.0)

#### Test Requirements
- Test with various GitHub Actions workflows
- Test matrix build detection
- Test caching detection
- Validate with real open source projects (react, vue, go)

---

### Story SS-081: GitHub Actions Integration (Ship Shape Action)

**As a** developer
**I want** to run Ship Shape as a GitHub Action in my CI pipeline
**So that** I can automatically enforce quality standards

**Priority**: P0 (Critical)
**Story Points**: 13
**Dependencies**: SS-080, SS-070 (Quality Gates)

#### Acceptance Criteria

**AC-081.1**: GitHub Action Creation
- **Given** Ship Shape binary
- **When** packaged as GitHub Action
- **Then** it is available in GitHub Marketplace
- **And** has Docker-based action
- **And** has action.yml metadata file
- **And** supports all inputs (path, config, gates, output)

**AC-081.2**: Workflow Integration
- **Given** a repository wants to use Ship Shape
- **When** developer adds Ship Shape action to workflow
- **Then** action runs on every PR
- **And** action runs on push to main
- **And** action respects workflow triggers

**AC-081.3**: Quality Gate Enforcement
- **Given** Ship Shape action with quality gates
- **When** analysis completes
- **Then** action exits with code 0 if gates pass
- **And** action exits with code 1 if gates fail
- **And** action posts comment on PR with results

**AC-081.4**: Artifact Upload
- **Given** Ship Shape generates reports
- **When** action completes
- **Then** it uploads HTML report as artifact
- **And** it uploads JSON results as artifact
- **And** artifacts are available for 90 days

**AC-081.5**: Status Check Integration
- **Given** Ship Shape action in workflow
- **When** action runs on PR
- **Then** it creates GitHub status check
- **And** status shows pass/fail state
- **And** status links to detailed report
- **And** status can be required for merge

#### Technical Requirements

**TR-081.1**: action.yml definition
```yaml
name: 'Ship Shape - Testing Quality Analysis'
description: 'Analyze testing quality, coverage, and best practices across multi-language repositories'
branding:
  icon: 'anchor'
  color: 'blue'

inputs:
  path:
    description: 'Path to repository (default: current directory)'
    required: false
    default: '.'

  config:
    description: 'Path to Ship Shape config file'
    required: false
    default: '.shipshape.yml'

  gates:
    description: 'Enable quality gates (true/false)'
    required: false
    default: 'true'

  fail-on-gate:
    description: 'Fail workflow if quality gates not met'
    required: false
    default: 'true'

  upload-artifacts:
    description: 'Upload report artifacts'
    required: false
    default: 'true'

  comment-pr:
    description: 'Post results as PR comment'
    required: false
    default: 'true'

outputs:
  score:
    description: 'Overall quality score (0-100)'

  grade:
    description: 'Quality grade (A-F)'

  gates-passed:
    description: 'Whether quality gates passed (true/false)'

  report-url:
    description: 'URL to uploaded HTML report'

runs:
  using: 'docker'
  image: 'Dockerfile'
  args:
    - analyze
    - ${{ inputs.path }}
    - --config=${{ inputs.config }}
    - --output=json
    - --output=html
```

**TR-081.2**: Dockerfile for action
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o shipshape ./cmd/shipshape

# Runtime image with language tools
FROM ubuntu:22.04

# Install language runtimes for analysis
RUN apt-get update && apt-get install -y \
    python3 python3-pip \
    nodejs npm \
    default-jdk \
    git \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/shipshape /usr/local/bin/
COPY action-entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
```

**TR-081.3**: Action entrypoint script
```bash
#!/bin/bash
set -e

# Run Ship Shape analysis
shipshape analyze "$INPUT_PATH" \
  --config="$INPUT_CONFIG" \
  --output=json \
  --output=html \
  --output-dir=./shipshape-results

# Set outputs
SCORE=$(jq -r '.scores.overall' ./shipshape-results/report.json)
GRADE=$(jq -r '.scores.grade' ./shipshape-results/report.json)
GATES_PASSED=$(jq -r '.gates.passed' ./shipshape-results/report.json)

echo "score=$SCORE" >> $GITHUB_OUTPUT
echo "grade=$GRADE" >> $GITHUB_OUTPUT
echo "gates-passed=$GATES_PASSED" >> $GITHUB_OUTPUT

# Upload artifacts if enabled
if [ "$INPUT_UPLOAD_ARTIFACTS" = "true" ]; then
    echo "Uploading artifacts..."
fi

# Post PR comment if enabled
if [ "$INPUT_COMMENT_PR" = "true" ] && [ -n "$GITHUB_EVENT_PATH" ]; then
    # Post comment using GitHub API
    ./post-pr-comment.sh
fi

# Exit based on gates
if [ "$INPUT_FAIL_ON_GATE" = "true" ] && [ "$GATES_PASSED" != "true" ]; then
    echo "Quality gates failed"
    exit 1
fi

exit 0
```

**TR-081.4**: PR comment poster (using GitHub API)
```bash
#!/bin/bash

if [ -z "$GITHUB_TOKEN" ]; then
    echo "GITHUB_TOKEN not set, skipping PR comment"
    exit 0
fi

# Extract PR number from GITHUB_REF
PR_NUMBER=$(echo $GITHUB_REF | sed 's/refs\/pull\/\([0-9]*\)\/merge/\1/')

if [ -z "$PR_NUMBER" ]; then
    echo "Not a PR, skipping comment"
    exit 0
fi

# Generate markdown summary
SUMMARY=$(cat <<EOF
## Ship Shape Analysis Results

**Overall Score**: $(jq -r '.scores.overall' ./shipshape-results/report.json)/100
**Grade**: $(jq -r '.scores.grade' ./shipshape-results/report.json)

### Dimension Scores
$(jq -r '.scores.dimensions[] | "- **\(.dimension)**: \(.score)/100"' ./shipshape-results/report.json)

### Quality Gates
$(jq -r '.gates.results[] | "- \(.name): \(if .passed then "✅ PASS" else "❌ FAIL" end)"' ./shipshape-results/report.json)

[View Full Report](artifact-url)
EOF
)

# Post comment via GitHub API
curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"body\": $(echo "$SUMMARY" | jq -Rs .)}" \
  "https://api.github.com/repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments"
```

#### Open Source Tools
- **GitHub Actions**: GitHub's CI/CD platform (free for public repos)
- **jq**: JSON processor (MIT)
- **curl**: HTTP client (MIT)

#### Test Requirements
- Test action in real GitHub repository
- Test PR comment posting
- Test artifact upload
- Test quality gate failures
- Test with various workflow triggers
- Test with matrix builds

---

### Story SS-082: Pre-commit Hook Integration

**As a** developer
**I want** to run Ship Shape analysis on pre-commit
**So that** I catch quality issues before pushing

**Priority**: P2 (Medium)
**Story Points**: 5
**Dependencies**: SS-070 (Quality Gates)

#### Acceptance Criteria

**AC-082.1**: pre-commit Hook Definition
- **Given** Ship Shape binary
- **When** configured as pre-commit hook
- **Then** .pre-commit-hooks.yaml is available
- **And** hook runs on commit
- **And** hook validates only changed files

**AC-082.2**: Fast Incremental Analysis
- **Given** pre-commit hook execution
- **When** user commits files
- **Then** analysis completes in <10 seconds
- **And** only changed test files are analyzed
- **And** results are cached between runs

**AC-082.3**: Configurable Strictness
- **Given** pre-commit configuration
- **When** hook runs
- **Then** user can configure which gates apply
- **And** user can set warning vs. blocking mode
- **And** user can exclude specific checks

#### Technical Requirements

**TR-082.1**: .pre-commit-hooks.yaml
```yaml
- id: shipshape
  name: Ship Shape Test Quality Check
  entry: shipshape analyze
  language: golang
  files: '.*_test\.(go|py|js|ts|java|rs)$'
  pass_filenames: true
  args:
    - --fast
    - --changed-only
```

**TR-082.2**: Changed files analysis
```go
func analyzeChangedFiles(files []string) (*Report, error) {
    // Only analyze test files and their corresponding source files
    testFiles := filterTestFiles(files)

    // Run lightweight analysis
    results := &AnalysisResults{
        FastMode: true,
    }

    for _, testFile := range testFiles {
        // Quick analysis without full repository scan
        analyzeTestFile(testFile)
    }

    return generateQuickReport(results), nil
}
```

#### Open Source Tools
- **pre-commit**: Hook framework (MIT)

#### Test Requirements
- Test with various pre-commit scenarios
- Test performance with 10+ changed files
- Test incremental analysis accuracy

---

## Epic 10: Reporting & Dashboards

**Epic Goal**: Generate beautiful, actionable reports in multiple formats (HTML, JSON, CLI).

### Story SS-090: HTML Report Generation

**As a** developer
**I want** Ship Shape to generate a beautiful HTML report
**So that** I can easily share results with my team

**Priority**: P1 (High)
**Story Points**: 13
**Dependencies**: SS-070 (Scoring)

#### Acceptance Criteria

**AC-090.1**: HTML Report Structure
- **Given** analysis results
- **When** HTML report is generated
- **Then** report includes overall score prominently
- **And** includes dimension score breakdown
- **And** includes findings categorized by severity
- **And** includes coverage charts
- **And** includes trend graphs (if historical data exists)

**AC-090.2**: Interactive Elements
- **Given** HTML report
- **When** user opens report in browser
- **Then** findings are collapsible by category
- **And** clicking finding shows file location
- **And** charts are interactive (hover for details)
- **And** report is responsive (mobile-friendly)

**AC-090.3**: Styling and Branding
- **Given** HTML report generation
- **When** rendering report
- **Then** uses consistent color scheme
- **And** uses readable typography
- **And** includes Ship Shape branding
- **And** works without external CDN dependencies (self-contained)

**AC-090.4**: Report Sections
- **Given** complete analysis
- **When** report is generated
- **Then** includes executive summary
- **And** includes score breakdown with rationale
- **And** includes findings by category (test smells, coverage gaps, missing tools)
- **And** includes recommendations prioritized by impact
- **And** includes technical metrics appendix

#### Technical Requirements

**TR-090.1**: Use templ for type-safe HTML templates
```go
// Install: go install github.com/a-h/templ/cmd/templ@latest

// templates/report.templ
package templates

import "github.com/shipshape/pkg/core"

templ Report(report *core.Report) {
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <title>Ship Shape Report - { report.RepositoryContext.RootPath }</title>
        <style>
            @templ.Raw(embedCSS())
        </style>
    </head>
    <body>
        <header>
            <h1>Ship Shape Analysis Report</h1>
            <div class="score-badge" data-score={ fmt.Sprintf("%.1f", report.Scores.Overall) }>
                { fmt.Sprintf("%.1f", report.Scores.Overall) }
            </div>
            <div class="grade">Grade: { report.Scores.Grade }</div>
        </header>

        <main>
            @ExecutiveSummary(report)
            @DimensionScores(report.Scores.Dimensions)
            @FindingsSection(report.Findings)
            @RecommendationsSection(report.Recommendations)
            @TechnicalMetrics(report.Results)
        </main>

        <script>
            @templ.Raw(embedJS())
        </script>
    </body>
    </html>
}

templ ExecutiveSummary(report *core.Report) {
    <section id="summary">
        <h2>Executive Summary</h2>
        <div class="summary-grid">
            <div class="metric">
                <div class="label">Overall Score</div>
                <div class="value">{ fmt.Sprintf("%.1f/100", report.Scores.Overall) }</div>
            </div>
            <div class="metric">
                <div class="label">Lines of Code</div>
                <div class="value">{ fmt.Sprintf("%d", report.RepositoryContext.Structure.TotalLOC) }</div>
            </div>
            <div class="metric">
                <div class="label">Test Lines</div>
                <div class="value">{ fmt.Sprintf("%d", report.RepositoryContext.Structure.TestLOC) }</div>
            </div>
            <div class="metric">
                <div class="label">Coverage</div>
                <div class="value">{ fmt.Sprintf("%.1f%%", report.Results.Coverage.LinePercentage) }</div>
            </div>
        </div>
    </section>
}

templ DimensionScores(dimensions []*core.DimensionScore) {
    <section id="dimensions">
        <h2>Dimension Scores</h2>
        for _, dim := range dimensions {
            <div class="dimension">
                <div class="dimension-header">
                    <span class="dimension-name">{ dim.Dimension }</span>
                    <span class="dimension-score">{ fmt.Sprintf("%.1f/100", dim.Score) }</span>
                </div>
                <div class="progress-bar">
                    <div class="progress-fill" style={ fmt.Sprintf("width: %.1f%%", dim.Score) }></div>
                </div>
                <p class="dimension-rationale">{ dim.Rationale }</p>
            </div>
        }
    </section>
}
```

**TR-090.2**: Embedded CSS (no external dependencies)
```css
/* Embed in report */
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: #333;
    background: #f5f5f5;
}

header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 2rem;
    text-align: center;
}

.score-badge {
    font-size: 4rem;
    font-weight: bold;
    margin: 1rem 0;
}

.score-badge[data-score^="9"],
.score-badge[data-score^="10"] {
    color: #10b981;
}

/* ... more CSS */
```

**TR-090.3**: Interactive JavaScript (embedded)
```javascript
// Collapsible findings
document.querySelectorAll('.finding-header').forEach(header => {
    header.addEventListener('click', () => {
        header.parentElement.classList.toggle('expanded');
    });
});

// Chart rendering using embedded library
function renderCoverageChart(data) {
    // Simple SVG-based chart, no external dependencies
}
```

**TR-090.4**: Report generator
```go
func (rg *ReportGenerator) GenerateHTML(report *core.Report) ([]byte, error) {
    // Use templ to render
    component := templates.Report(report)

    var buf bytes.Buffer
    err := component.Render(context.Background(), &buf)
    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}
```

#### Open Source Tools
- **templ**: Type-safe HTML templating (MIT) - github.com/a-h/templ
- **htmx**: Optional enhancement for interactivity (BSD-2-Clause) - embedded

#### Test Requirements
- Test HTML generation with various report data
- Test rendering in multiple browsers
- Test mobile responsiveness
- Test with large reports (1000+ findings)
- Validate HTML structure

---

## Epic 11: Historical Tracking & Trends

**Epic Goal**: Track analysis results over time to show quality trends and improvements.

### Story SS-100: SQLite Storage for Historical Data

**As a** developer
**I want** Ship Shape to store historical analysis results
**So that** I can track quality trends over time

**Priority**: P1 (High)
**Story Points**: 8
**Dependencies**: SS-070 (Scoring)

#### Acceptance Criteria

**AC-100.1**: Database Initialization
- **Given** Ship Shape first run
- **When** storage is initialized
- **Then** SQLite database is created at ~/.shipshape/history.db
- **And** schema includes analysis_runs, dimension_scores, findings, metrics tables
- **And** indexes are created for query performance

**AC-100.2**: Analysis Result Storage
- **Given** completed analysis
- **When** storing results
- **Then** overall score is stored
- **And** dimension scores are stored
- **And** critical/high findings are stored
- **And** git commit hash is stored (if repo is git)
- **And** timestamp is recorded

**AC-100.3**: Historical Retrieval
- **Given** stored historical data
- **When** querying history
- **Then** recent N runs can be retrieved
- **And** runs within date range can be queried
- **And** trend data can be calculated

**AC-100.4**: Data Cleanup
- **Given** old historical data
- **When** cleanup runs
- **Then** runs older than configurable threshold are deleted
- **And** database is vacuumed to reclaim space

#### Technical Requirements

**TR-100.1**: SQLite schema (from architecture doc)
```sql
CREATE TABLE analysis_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_path TEXT NOT NULL,
    repository_hash TEXT,
    timestamp DATETIME NOT NULL,
    overall_score REAL,
    grade TEXT,
    duration_seconds INTEGER,
    analyzer_version TEXT
);

CREATE TABLE dimension_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    dimension TEXT NOT NULL,
    score REAL NOT NULL,
    weight REAL NOT NULL,
    rationale TEXT,
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE TABLE findings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    finding_id TEXT NOT NULL,
    check_id TEXT NOT NULL,
    type TEXT NOT NULL,
    severity TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    file_path TEXT,
    line_number INTEGER,
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE TABLE metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    metric_name TEXT NOT NULL,
    metric_value REAL NOT NULL,
    metric_category TEXT,
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE INDEX idx_runs_repo ON analysis_runs(repository_path, timestamp);
CREATE INDEX idx_findings_run ON findings(run_id);
CREATE INDEX idx_metrics_run ON metrics(run_id);
```

**TR-100.2**: Storage interface implementation
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
    db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    storage := &SQLiteStorage{db: db}
    if err := storage.initialize(); err != nil {
        return nil, err
    }

    return storage, nil
}

func (s *SQLiteStorage) StoreReport(report *core.Report) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert analysis run
    result, err := tx.Exec(`
        INSERT INTO analysis_runs
        (repository_path, repository_hash, timestamp, overall_score, grade, duration_seconds, analyzer_version)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `,
        report.RepositoryContext.RootPath,
        getGitHash(report.RepositoryContext.RootPath),
        report.GeneratedAt,
        report.Scores.Overall,
        report.Scores.Grade,
        report.Results.EndTime.Sub(report.Results.StartTime).Seconds(),
        version.Version,
    )

    runID, _ := result.LastInsertId()

    // Insert dimension scores
    for _, dim := range report.Scores.Dimensions {
        tx.Exec(`
            INSERT INTO dimension_scores (run_id, dimension, score, weight, rationale)
            VALUES (?, ?, ?, ?, ?)
        `, runID, dim.Dimension, dim.Score, dim.Weight, dim.Rationale)
    }

    return tx.Commit()
}

func (s *SQLiteStorage) GetRecentReports(repoPath string, limit int) ([]*core.Report, error) {
    rows, err := s.db.Query(`
        SELECT id, timestamp, overall_score, grade
        FROM analysis_runs
        WHERE repository_path = ?
        ORDER BY timestamp DESC
        LIMIT ?
    `, repoPath, limit)
    // ... parse and return reports
}
```

#### Open Source Tools
- **mattn/go-sqlite3**: SQLite driver for Go (MIT)

#### Test Requirements
- Test database initialization
- Test report storage and retrieval
- Test with multiple repositories
- Test data cleanup
- Performance: Store 1000 runs in <5 seconds

---

### Story SS-101: Trend Visualization and Analysis

**As a** developer
**I want** to see quality trends over time
**So that** I can track improvements and regressions

**Priority**: P1 (High)
**Story Points**: 8
**Dependencies**: SS-100

#### Acceptance Criteria

**AC-101.1**: Trend Data Calculation
- **Given** historical analysis data
- **When** generating trend report
- **Then** score trends over last 30 days are calculated
- **And** dimension score trends are calculated
- **And** finding count trends are calculated
- **And** coverage trends are calculated

**AC-101.2**: Trend Charts in HTML Report
- **Given** HTML report with historical data
- **When** report is rendered
- **Then** overall score trend chart is shown
- **And** dimension score trends are shown
- **And** charts show clear improvement/regression indicators

**AC-101.3**: Regression Detection
- **Given** new analysis with historical data
- **When** comparing to previous runs
- **Then** score drops >5 points are flagged
- **And** new critical findings are highlighted
- **And** coverage decreases are reported

**AC-101.4**: CLI Trend Summary
- **Given** trend data
- **When** running shipshape trends command
- **Then** ASCII chart shows score history
- **And** summary shows improvement/regression
- **And** highlights significant changes

#### Technical Requirements

**TR-101.1**: Trend calculator
```go
func (s *SQLiteStorage) GetMetricTrend(
    repoPath string,
    metricName string,
    days int,
) ([]MetricPoint, error) {
    cutoff := time.Now().AddDate(0, 0, -days)

    rows, err := s.db.Query(`
        SELECT timestamp, metric_value
        FROM metrics m
        JOIN analysis_runs r ON m.run_id = r.id
        WHERE r.repository_path = ?
          AND m.metric_name = ?
          AND r.timestamp >= ?
        ORDER BY r.timestamp ASC
    `, repoPath, metricName, cutoff)
    // ... parse and return
}
```

**TR-101.2**: Chart generation (SVG)
```go
func generateTrendChart(points []MetricPoint) string {
    // Generate simple SVG line chart
    svg := `<svg width="600" height="200" xmlns="http://www.w3.org/2000/svg">`

    // Calculate scales
    maxY := maxValue(points)
    scaleX := 600.0 / float64(len(points))
    scaleY := 200.0 / maxY

    // Draw line
    pathData := "M "
    for i, point := range points {
        x := float64(i) * scaleX
        y := 200 - (point.Value * scaleY)
        pathData += fmt.Sprintf("%.1f,%.1f ", x, y)
    }

    svg += fmt.Sprintf(`<path d="%s" fill="none" stroke="#667eea" stroke-width="2"/>`, pathData)
    svg += `</svg>`

    return svg
}
```

**TR-101.3**: CLI trend display (ASCII chart)
```go
func printTrendASCII(points []MetricPoint) {
    // Use github.com/gizak/termui or simple custom renderer

    fmt.Println("Score Trend (Last 30 Days)")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")

    // Normalize to terminal height (20 rows)
    maxY := maxValue(points)

    for row := 20; row >= 0; row-- {
        line := ""
        threshold := (float64(row) / 20.0) * maxY

        for _, point := range points {
            if point.Value >= threshold {
                line += "█"
            } else {
                line += " "
            }
        }

        fmt.Printf("%3.0f │%s\n", threshold, line)
    }

    fmt.Println("    └" + strings.Repeat("─", len(points)))
}
```

#### Open Source Tools
- No external charting library (pure SVG generation)
- Optional: **gizak/termui**: Terminal UI components (MIT)

#### Test Requirements
- Test trend calculation accuracy
- Test chart generation
- Test regression detection
- Test with various data patterns (improving, declining, stable)

---

## Epic 12: Organization-Wide Features

**Epic Goal**: Support organization-level analysis with cross-project dashboards and rankings.

### Story SS-110: Multi-Project Analysis

**As an** engineering manager
**I want** to analyze multiple projects in my organization
**So that** I can compare quality across teams

**Priority**: P2 (Medium)
**Story Points**: 13
**Dependencies**: SS-100 (Storage)

#### Acceptance Criteria

**AC-110.1**: Organization Configuration
- **Given** multiple repositories to analyze
- **When** configuring organization analysis
- **Then** user can provide list of repository paths
- **And** user can provide GitHub org name to auto-discover repos
- **And** configuration is saved to ~/.shipshape/orgs.yml

**AC-110.2**: Batch Analysis
- **Given** organization with N repositories
- **When** running organization analysis
- **Then** each repository is analyzed independently
- **And** analyses run in parallel (up to 4 concurrent)
- **And** failures in one repo don't stop others
- **And** progress is shown for each repository

**AC-110.3**: Organization Dashboard
- **Given** completed organization analysis
- **When** dashboard is generated
- **Then** all projects are ranked by score
- **And** score distribution is shown
- **And** dimension averages across org are shown
- **And** outliers (best and worst) are highlighted

**AC-110.4**: Cross-Project Insights
- **Given** organization analysis results
- **When** generating insights
- **Then** common issues across projects are identified
- **And** best practices from top projects are highlighted
- **And** tool adoption consistency is analyzed

#### Technical Requirements

**TR-110.1**: Organization config
```yaml
# ~/.shipshape/orgs.yml
organizations:
  - name: my-company
    repositories:
      - path: /path/to/repo1
        name: frontend
      - path: /path/to/repo2
        name: backend-api
      - path: /path/to/repo3
        name: data-pipeline

  - name: my-oss-projects
    github_org: my-username
    include_pattern: "^project-.*"
    exclude_pattern: "^archived-.*"
```

**TR-110.2**: Organization analyzer
```go
type OrganizationAnalyzer struct {
    analyzer      *AnalysisEngine
    storage       *SQLiteStorage
    maxConcurrent int
}

func (oa *OrganizationAnalyzer) AnalyzeOrganization(
    org *Organization,
) (*OrganizationReport, error) {
    repos := oa.discoverRepositories(org)

    sem := make(chan struct{}, oa.maxConcurrent)
    results := make(chan *RepoResult, len(repos))

    for _, repo := range repos {
        sem <- struct{}{}
        go func(r *Repository) {
            defer func() { <-sem }()

            report, err := oa.analyzer.Analyze(context.Background(), r.Path)
            results <- &RepoResult{
                Repository: r,
                Report:     report,
                Error:      err,
            }
        }(repo)
    }

    // Collect and aggregate results
    orgReport := &OrganizationReport{
        Organization: org.Name,
        Repositories: make(map[string]*core.Report),
    }

    for i := 0; i < len(repos); i++ {
        result := <-results
        if result.Error != nil {
            log.Warn("Failed to analyze", result.Repository.Name, result.Error)
            continue
        }
        orgReport.Repositories[result.Repository.Name] = result.Report
    }

    orgReport.Rankings = oa.calculateRankings(orgReport.Repositories)
    orgReport.Insights = oa.generateInsights(orgReport.Repositories)

    return orgReport, nil
}
```

**TR-110.3**: Ranking calculation
```go
func (oa *OrganizationAnalyzer) calculateRankings(
    repos map[string]*core.Report,
) []*Ranking {
    rankings := make([]*Ranking, 0, len(repos))

    for name, report := range repos {
        rankings = append(rankings, &Ranking{
            Repository: name,
            Score:      report.Scores.Overall,
            Grade:      report.Scores.Grade,
        })
    }

    // Sort by score descending
    sort.Slice(rankings, func(i, j int) bool {
        return rankings[i].Score > rankings[j].Score
    })

    // Assign ranks
    for i := range rankings {
        rankings[i].Rank = i + 1
    }

    return rankings
}
```

**TR-110.4**: Organization dashboard HTML
```go
// templates/org-dashboard.templ
templ OrganizationDashboard(orgReport *OrganizationReport) {
    <html>
    <head>
        <title>{ orgReport.Organization } - Ship Shape Dashboard</title>
    </head>
    <body>
        <h1>{ orgReport.Organization } Quality Dashboard</h1>

        <section id="overview">
            <div class="metric">
                <span>Average Score</span>
                <span>{ fmt.Sprintf("%.1f", orgReport.AverageScore) }</span>
            </div>
            <div class="metric">
                <span>Total Repositories</span>
                <span>{ fmt.Sprintf("%d", len(orgReport.Repositories)) }</span>
            </div>
        </section>

        <section id="rankings">
            <h2>Repository Rankings</h2>
            <table>
                <thead>
                    <tr>
                        <th>Rank</th>
                        <th>Repository</th>
                        <th>Score</th>
                        <th>Grade</th>
                        <th>Trend</th>
                    </tr>
                </thead>
                <tbody>
                    for _, ranking := range orgReport.Rankings {
                        <tr>
                            <td>{ fmt.Sprintf("%d", ranking.Rank) }</td>
                            <td>{ ranking.Repository }</td>
                            <td>{ fmt.Sprintf("%.1f", ranking.Score) }</td>
                            <td class={ "grade-" + ranking.Grade }>{ ranking.Grade }</td>
                            <td>{ renderTrend(ranking.Trend) }</td>
                        </tr>
                    }
                </tbody>
            </table>
        </section>
    </body>
    </html>
}
```

#### Open Source Tools
- **templ**: HTML templating (MIT)
- **goroutines**: Parallel processing

#### Test Requirements
- Test with 10+ repositories
- Test parallel processing
- Test ranking calculation
- Test dashboard generation
- Performance: Analyze 10 repos in <2 minutes

---

## Summary and Implementation Roadmap

### Release Plan

**v0.1.0 - Repository Discovery (Q1 2026)**
- SS-001: Language detection
- SS-002: Framework detection
- SS-003: Monorepo detection
- SS-004: Directory structure mapping

**v0.2.0 - Multi-Language Foundation (Q1 2026)**
- SS-010: Go test analysis
- SS-011: Python test analysis
- SS-012: JavaScript/TypeScript test analysis
- SS-020: Monorepo package analysis
- SS-021: Monorepo aggregate scoring

**v0.3.0 - Quality Analysis (Q2 2026)**
- SS-030: Test smell detection
- SS-040: Coverage parsing
- SS-041: Coverage quality assessment

**v0.4.0 - Tool Ecosystem (Q2 2026)**
- SS-050: Tool database and detection
- SS-051: Tool recommendations
- SS-070: Scoring and assessment (from Epic 7)

**v0.5.0 - CI/CD Integration (Q3 2026)**
- SS-080: GitHub Actions workflow detection
- SS-081: GitHub Actions integration
- SS-082: Pre-commit hooks
- SS-070: Quality gates (from Epic 8)

**v0.6.0 - Reporting (Q3 2026)**
- SS-090: HTML report generation
- SS-091: JSON output (from Epic 10)
- SS-092: CLI summary (from Epic 10)

**v0.7.0 - Historical Tracking (Q4 2026)**
- SS-100: SQLite storage
- SS-101: Trend visualization

**v0.8.0 - Organization Features (Q4 2026)**
- SS-110: Multi-project analysis

### Total Story Points: 180 points
**Estimated Development Time**: 36-45 weeks (9-11 months)

---

**Document Version**: 1.0.0
**Last Updated**: 2026-01-27
**Status**: Ready for Implementation Planning
