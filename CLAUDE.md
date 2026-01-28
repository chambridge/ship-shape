# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ship Shape is a comprehensive testing quality and code analysis platform for multi-language repositories and monorepos. Currently in **Phase 0: Foundation & Setup** (v0.0.0), targeting v1.0.0 release after 24 weeks of development.

**Core Philosophy**:
- Context-First Analysis: Discover repository structure before analysis
- Language Isolation: Each language analyzer is independent and composable
- Type-Safe Core: Go-based for performance and reliability
- Plugin Architecture: Extensible through well-defined interfaces
- Graceful Degradation: Partial results on failures, not complete failure

## Development Commands

### Building and Running
```bash
# Build the binary (creates ./bin/shipshape)
make build

# Run with arguments
make run ARGS="version"
make run ARGS="analyze ."

# Install to $GOPATH/bin
make install

# Clean build artifacts
make clean
```

### Testing and Quality Checks
```bash
# Run all tests with coverage
make test

# Generate HTML coverage report
make coverage

# Run all quality checks (recommended before commit)
make check  # Runs: fmt, vet, lint, test

# Individual quality checks
make fmt         # Format code with gofmt
make vet         # Run go vet static analysis
make lint        # Run golangci-lint (30+ linters)
make security    # Run gosec security scanner
make actionlint  # Validate GitHub Actions workflows
```

### Dependencies
```bash
# Download and verify dependencies
make deps

# Tidy module dependencies
make tidy
```

### Logging

Ship Shape uses structured logging with `log/slog` (Go 1.21+) through the `internal/logger` package. Logging is controlled by CLI flags and provides consistent, leveled output.

**Usage in Code**:
```go
import "github.com/chambridge/ship-shape/internal/logger"

// Package-level convenience functions
logger.Debug("detailed diagnostic information", "key", value)
logger.Info("general informational messages", "key", value)
logger.Warn("warning messages for recoverable issues", "key", value)
logger.Error("error messages", "error", err, "context", "additional info")

// Context-aware logging
logger.InfoContext(ctx, "operation completed", "duration", time.Since(start))

// Logger with attributes
log := logger.With("component", "analyzer", "language", "go")
log.Info("starting analysis")  // All logs will include component and language fields

// Logger with groups
log := logger.WithGroup("request")
log.Info("API call", "method", "GET", "path", "/api/v1")  // Groups method and path under "request"
```

**CLI Flags**:
```bash
# Default (INFO level)
./bin/shipshape version

# Verbose mode (DEBUG level)
./bin/shipshape version --verbose
./bin/shipshape version -v

# Quiet mode (ERROR level only)
./bin/shipshape version --quiet
./bin/shipshape version -q

# Disable colored output
./bin/shipshape version --no-color
```

**Configuration**:
- Logs go to `stderr` by default
- Text format for human readability
- JSON format available for machine parsing (future enhancement)
- Structured key-value pairs for rich context
- Thread-safe and race-condition free

**Important Development Guidelines**:
- **ALWAYS use make targets** instead of invoking tools directly (e.g., `make test` not `go test`, `make lint` not `golangci-lint`)
- Keep make targets up-to-date when adding new development workflows
- Use logger functions instead of `fmt.Printf` for all diagnostic output
- Use `logger.Debug` for development/troubleshooting messages
- Use `logger.Info` for normal operation messages
- Use `logger.Warn` for recoverable issues
- Use `logger.Error` for errors (but still return errors to callers)
- Always include relevant context as key-value pairs

**Quality Standards (Enforced by CI)**:
- Go 1.21+ required (tested on 1.21, 1.22, 1.23)
- >90% unit test coverage target
- Zero linting errors (golangci-lint with 30+ linters)
- All commits must be signed with DCO (`git commit -s`)
- Conventional commit format required

## Architecture

### High-Level System Design

```
CLI Layer (cmd/shipshape)
    ↓
Core Engine (Coordinator)
    ↓
Repository Context (Shared State)
    ├─ Languages, Frameworks, Structure
    └─ Monorepo Configuration
    ↓
┌─────────────┬─────────────┬──────────────┐
│  Discovery  │  Execution  │   Report     │
│   Engine    │   Engine    │  Generator   │
└─────────────┴─────────────┴──────────────┘
    ↓
┌──────────────┬──────────────┬──────────────┐
│  Analyzer    │  Assessor    │    Tool      │
│  Registry    │  Registry    │  Detector    │
└──────────────┴──────────────┴──────────────┘
    ↓
Language-Specific Plugins (Go, Python, JS/TS, Java, Rust)
```

### 7-Layer Analysis Model

Ship Shape processes repositories through seven distinct layers:

1. **Discovery** - Language, framework, and structure detection
2. **Structure** - Directory mapping and LOC analysis
3. **Quality** - Test smell detection and best practices
4. **Coverage** - Coverage report parsing and analysis
5. **Execution** - Test performance and flakiness (optional)
6. **Mutation** - Mutation testing analysis (optional)
7. **CI/CD** - CI pipeline analysis and optimization

Each layer can operate independently and fail gracefully without blocking subsequent layers.

### Directory Structure and Conventions

```
ship-shape/
├── cmd/shipshape/          # CLI entry point and command handlers
│   ├── main.go             # Entry point, calls rootCmd.Execute()
│   ├── root.go             # Root command with global flags
│   └── version.go          # Version command
│
├── internal/               # Private application code (not importable)
│   ├── analyzer/           # Language-specific analyzers (registry pattern)
│   │   ├── go/             # Go AST analysis
│   │   ├── python/         # Python analysis (tree-sitter)
│   │   └── javascript/     # JS/TS analysis (tree-sitter)
│   ├── coverage/           # Coverage report parsers
│   ├── detector/           # Framework and tool detection
│   └── report/             # Report generation (HTML, JSON, etc.)
│
├── pkg/types/              # Public shared types and interfaces
│   ├── repository.go       # Repository context and metadata
│   ├── analyzer.go         # Analyzer interface and types
│   ├── assessor.go         # Assessment and scoring types
│   └── results.go          # Analysis results and findings
│
├── testdata/               # Test fixtures and validation data
│   ├── ground-truth/       # Curated test cases with known results
│   ├── integration/        # Integration test repositories
│   └── validation/         # Real-world OSS projects for validation
│
├── docs/v1.0.0/            # Version-specific documentation
│   ├── cicd.md             # CI/CD pipeline documentation
│   └── [other docs]        # To be added during implementation
│
└── .project/               # Comprehensive planning documents
    ├── requirements.md     # Detailed requirements (2,300+ lines)
    ├── architecture.md     # Technical architecture (3,400+ lines)
    ├── user-stories.md     # 27 user stories (180 story points)
    └── plan.md             # 24-week implementation roadmap
```

**Important Conventions**:
- `internal/` - Private code that cannot be imported by external projects
- `pkg/` - Public interfaces and types that may be used by plugins or extensions
- `cmd/` - Only CLI logic, delegate to `internal/` packages
- `testdata/` - Non-Go test fixtures (actual test files live beside code)

### Plugin Architecture Pattern

Ship Shape uses a **registry-based plugin system** for analyzers and assessors:

```go
// Analyzer Registry Pattern (to be implemented)
type Analyzer interface {
    Name() string
    Languages() []Language
    Analyze(ctx context.Context, repo *Repository) (*AnalysisResult, error)
}

// Registry manages all analyzers
type AnalyzerRegistry struct {
    analyzers map[Language][]Analyzer
}

// Usage pattern
registry := NewAnalyzerRegistry()
registry.Register(NewGoAnalyzer())
registry.Register(NewPythonAnalyzer())

// Execute all analyzers for detected languages
results := registry.AnalyzeRepository(ctx, repo)
```

**Key Design Principles**:
- Each analyzer is self-contained and stateless
- Analyzers declare their supported languages
- Registry handles parallel execution and error aggregation
- Context-aware: analyzers receive full repository context
- Fail gracefully: one analyzer failure doesn't stop others

### Multi-Language Support Strategy

Each language has:
1. **Detector** - Identifies if language is present (file extensions, manifests)
2. **Analyzer** - Performs language-specific analysis (AST parsing, test detection)
3. **Assessor** - Evaluates quality using language-specific best practices

**AST Parsing Approach**:
- **Go**: Use `go/ast` standard library (native, fast, accurate)
- **Python**: tree-sitter or embedded Python AST parser
- **JavaScript/TypeScript**: tree-sitter with language grammars
- **Java**: tree-sitter
- **Rust**: tree-sitter

**Framework Detection**:
- Scan for config files (package.json, pyproject.toml, go.mod, Cargo.toml)
- Parse dependency lists for known testing frameworks
- Detect test file patterns (e.g., `*_test.go`, `test_*.py`, `*.spec.js`)

### Monorepo Handling

Ship Shape is **monorepo-native** with first-class support:

1. **Workspace Detection**:
   - npm/yarn/pnpm workspaces (package.json)
   - Go workspaces (go.work)
   - Maven multi-module (pom.xml)
   - Lerna/Nx monorepos

2. **Analysis Strategy**:
   - Detect workspace structure first
   - Analyze each package independently
   - Generate per-package scores
   - Aggregate to workspace-level score
   - Report shared infrastructure quality

3. **Scoring Model**:
   - Individual package scores (0-100)
   - Weighted average for workspace (by LOC or test count)
   - Shared infrastructure score separate from packages

## Test Strategy and Validation

### Testing Requirements (Non-Negotiable)

- **TDD Required**: Write tests before implementing functionality
- **Coverage Targets**:
  - Unit tests: >90% coverage
  - Integration tests: >80% coverage
- **Test Types**:
  - Unit tests (beside implementation files)
  - Integration tests (testdata/integration/)
  - Ground truth validation (testdata/ground-truth/)

### Ground Truth Datasets

Critical for validating analyzer accuracy:

**Structure** (`testdata/ground-truth/`):
```
ground-truth/
├── go/
│   ├── test-smells/
│   │   ├── eager-test/      # Known eager test examples
│   │   ├── assertion-roulette/
│   │   └── ...
│   └── coverage/            # Known coverage scenarios
├── python/
└── javascript/
```

**Requirements**:
- Each example must have expected results (JSON manifests)
- ≥90% precision/recall target for smell detection
- Continuous validation in CI (daily automated tests)

### Meta-Validation ("Dogfooding")

Ship Shape must analyze **itself** and achieve:
- ≥85/100 overall quality score
- >90% test coverage
- Zero high-severity test smells
- All best practices followed

This is validated continuously in CI.

## Configuration System

Ship Shape uses a hierarchical configuration approach:

**Config File Locations** (searched in order):
1. `--config` flag value (if provided)
2. `.shipshape.yml` (current directory)
3. Repository root (searches for .git)
4. `$HOME/.config/shipshape/.shipshape.yml`
5. `/etc/shipshape/.shipshape.yml`

**Environment Variables**:
- Prefix: `SHIPSHAPE_`
- Example: `SHIPSHAPE_VERBOSE=true`
- Overrides config file values

**Example Configuration** (see `.shipshape.example.yml`):
- Quality thresholds and gates
- Language-specific settings
- Exclusion patterns
- Output formatting

## Important Implementation Patterns

### Error Handling

Use structured errors with context:
```go
// Good
return fmt.Errorf("failed to parse coverage file %s: %w", filename, err)

// Bad
return err  // loses context
```

### Logging

Use leveled logging (to be implemented):
- `--verbose` flag: DEBUG level
- `--quiet` flag: ERROR only
- Default: INFO level

### Concurrent Analysis

Analyzers run in parallel when possible:
- Use `errgroup` for coordinated goroutines
- Context-aware cancellation
- Aggregate results safely

### Configuration Binding

Global flags bound to Viper configuration:
```go
// Explicitly ignore binding errors (BindPFlag always succeeds for valid flags)
_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
```

## CI/CD Pipeline

Comprehensive GitHub Actions pipeline (see `docs/v1.0.0/cicd.md`):

**Jobs** (run in parallel):
- Build (Go 1.21, 1.22, 1.23)
- Test (with Codecov upload)
- Lint (golangci-lint with 30+ linters)
- Format Check (gofmt compliance)
- Vet (go vet static analysis)
- Security Scan (gosec with SARIF upload)
- Cross-Compile (5 platforms)
- Quality Gate (coverage threshold enforcement)

**Local Validation Before Push**:
```bash
make actionlint  # Validate workflow files
make check       # Run all quality checks
make security    # Run security scan (optional)
```

## Common Development Workflows

### Adding a New Language Analyzer

1. Create package in `internal/analyzer/<language>/`
2. Implement `Analyzer` interface (from `pkg/types/`)
3. Add language detector in `internal/detector/`
4. Create ground truth test cases in `testdata/ground-truth/<language>/`
5. Register analyzer in registry (when implemented)
6. Add integration tests
7. Update documentation

### Adding a Test Smell Detector

1. Research smell definition and detection rules
2. Add to appropriate language analyzer
3. Create ground truth examples (true positives and false positives)
4. Implement AST-based detection logic
5. Add unit tests with >90% coverage
6. Validate against ground truth (≥90% precision/recall)
7. Document detection algorithm

### Working with the Planning Documents

Critical references in `.project/`:

- **architecture.md**: Technical design, interfaces, Go code examples
- **user-stories.md**: 27 stories with acceptance criteria (180 SP total)
- **plan.md**: 24-week roadmap with dependencies and milestones
- **requirements.md**: Detailed functional and non-functional requirements

When implementing a feature:
1. Find corresponding user story (SS-XXX) in user-stories.md
2. Review acceptance criteria
3. Check architecture.md for design patterns
4. Verify story point estimate is reasonable
5. Update plan.md with actual progress

## Development Standards

### Commit Requirements

All commits must:
- Use conventional commit format: `type(scope): description`
- Include DCO sign-off: `git commit -s`
- Pass all CI checks (build, test, lint)

**Commit Types**:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Test additions/updates
- `refactor:` - Code refactoring
- `chore:` - Maintenance, dependencies
- `ci:` - CI/CD changes

### Code Quality Standards

**Complexity Limits** (enforced by golangci-lint):
- Cyclomatic complexity: ≤15
- Cognitive complexity: ≤20
- Nested if depth: ≤4

**Required Checks**:
- errcheck: All error returns checked or explicitly ignored with `_`
- gosec: No security vulnerabilities
- gocyclo/gocognit: Within complexity limits
- revive: Package comments, exported names documented
- wsl: Whitespace rules (no cuddling of different statement types)

### Deprecated Linters

- `exportloopref` - Removed (Go 1.22+ has loopvar)
- Use `copyloopvar` instead

## Release Process (Planned)

Version tagging triggers release workflow:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

Release workflow will:
- Run full CI pipeline
- Build for 5 platforms (Linux/macOS/Windows on amd64/arm64)
- Generate changelog
- Create GitHub release
- Upload binaries as release assets

## Reference Documentation

- **CI/CD Details**: `docs/v1.0.0/cicd.md`
- **Architecture Deep Dive**: `.project/architecture.md` (3,400+ lines)
- **User Stories**: `.project/user-stories.md` (27 stories, 180 SP)
- **Implementation Roadmap**: `.project/plan.md` (24-week plan)
- **Requirements Spec**: `.project/requirements.md` (2,300+ lines)
