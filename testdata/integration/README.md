# Integration Test Repositories

This directory contains complete mini-repositories for end-to-end integration testing of Ship Shape.

## Purpose

Integration test repositories are used to:

1. **End-to-End Validation**: Test Ship Shape's complete analysis workflow
2. **Real-World Scenarios**: Simulate actual project structures and configurations
3. **Multi-Language Testing**: Verify behavior with multiple languages in one repo
4. **Monorepo Support**: Test monorepo detection and analysis
5. **Edge Cases**: Handle unusual or complex repository structures

## Repository Types

### Simple Single-Language Projects

Minimal projects with one language and standard structure.

**Examples**:
- `simple-go/` - Basic Go module with tests
- `simple-python/` - Python package with pytest
- `simple-javascript/` - Node.js project with Jest
- `simple-java/` - Maven project with JUnit

**Characteristics**:
- Single language
- Standard project structure
- One testing framework
- Clear test organization

### Multi-Language Projects

Projects containing multiple programming languages.

**Examples**:
- `web-app/` - Frontend (JS) + Backend (Go)
- `data-pipeline/` - Python scripts + Go services
- `microservices/` - Multiple languages in separate services

**Characteristics**:
- 2+ languages in one repository
- Separate test suites per language
- Different testing frameworks
- Shared configuration

### Monorepo Projects

Repositories with multiple packages or workspaces.

**Examples**:
- `npm-workspace/` - npm/yarn/pnpm workspaces
- `go-workspace/` - Go workspace (go.work)
- `lerna-monorepo/` - Lerna-managed packages

**Characteristics**:
- Package/workspace configuration
- Shared dependencies
- Per-package test suites
- Workspace-level tooling

### Edge Cases

Unusual or problematic repository structures.

**Examples**:
- `no-tests/` - Project with production code but no tests
- `mixed-patterns/` - Multiple test frameworks in same language
- `nested-modules/` - Deeply nested project structure
- `unconventional/` - Non-standard directory layout

## Directory Structure

```
integration/
├── simple-go/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── main_test.go
│   ├── internal/
│   │   └── helper/
│   │       ├── helper.go
│   │       └── helper_test.go
│   ├── .shipshape.yml           # Ship Shape configuration
│   └── expected-results.yml     # Expected analysis results
│
├── simple-python/
│   ├── pyproject.toml
│   ├── src/
│   │   └── mypackage/
│   │       ├── __init__.py
│   │       └── module.py
│   ├── tests/
│   │   └── test_module.py
│   └── expected-results.yml
│
├── web-app/
│   ├── frontend/                # JavaScript/React
│   │   ├── package.json
│   │   └── src/
│   ├── backend/                 # Go
│   │   ├── go.mod
│   │   └── cmd/
│   └── expected-results.yml
│
└── monorepo-npm/
    ├── package.json             # Workspace root
    ├── packages/
    │   ├── package-a/
    │   └── package-b/
    └── expected-results.yml
```

## Expected Results Format

Each integration test includes `expected-results.yml` defining expected Ship Shape output:

```yaml
version: "1.0.0"
repository: simple-go
description: "Basic Go module with standard test structure"

expected_analysis:
  languages_detected:
    - language: go
      file_count: 4
      percentage: 100.0

  frameworks_detected:
    - language: go
      framework: testing
      version: "builtin"

  test_files:
    total: 2
    by_language:
      go: 2

  test_coverage:
    overall: 85.0
    by_package:
      - package: "."
        coverage: 80.0
      - package: "./internal/helper"
        coverage: 95.0

  test_smells:
    total: 1
    by_type:
      eager-test: 1
    findings:
      - file: "main_test.go"
        function: "TestMainFunction"
        smell: "eager-test"
        severity: "medium"

  quality_score:
    overall: 82.0
    test_coverage: 85.0
    test_quality: 75.0
    tool_adoption: 90.0

validation_tolerance:
  coverage_delta: 2.0      # ±2% coverage is acceptable
  score_delta: 5.0         # ±5 points on scores
  smell_count_exact: true  # Smell count must match exactly
```

## Integration Tests

Integration tests verify Ship Shape's end-to-end behavior:

```go
func TestIntegrationSimpleGo(t *testing.T) {
    // Load expected results
    expected := loadExpectedResults(t, "testdata/integration/simple-go/expected-results.yml")

    // Run Ship Shape analysis
    result := runShipShape(t, "testdata/integration/simple-go")

    // Verify languages detected
    assert.Equal(t, expected.LanguagesDetected, result.Languages)

    // Verify test coverage within tolerance
    assertWithinTolerance(t, expected.TestCoverage.Overall, result.Coverage.Overall, 2.0)

    // Verify test smells detected
    assert.Len(t, result.Smells, expected.TestSmells.Total)

    // Verify quality score
    assertWithinTolerance(t, expected.QualityScore.Overall, result.Score, 5.0)
}
```

## Creating New Integration Tests

When adding a new integration test repository:

1. **Choose Scenario**: Identify what specific behavior to test
2. **Create Minimal Repo**: Keep it as simple as possible while demonstrating the scenario
3. **Add Real Tests**: Include actual working tests (not stubs)
4. **Document Structure**: Add README.md explaining the purpose
5. **Define Expected Results**: Create `expected-results.yml` with precise expectations
6. **Set Tolerance Levels**: Define acceptable variance for non-deterministic metrics
7. **Verify Locally**: Run Ship Shape on the test repo before committing

## Maintenance

Integration test repositories should be:

- **Self-Contained**: No external dependencies that can change
- **Version-Locked**: Pin dependency versions
- **Documented**: Clear README explaining purpose and structure
- **Minimal**: Only include what's necessary to test the scenario
- **Updated**: Keep in sync with Ship Shape changes

## Common Scenarios

### Testing Monorepo Detection

```
monorepo-npm/
├── package.json                 # workspaces: ["packages/*"]
└── packages/
    ├── web/
    │   ├── package.json
    │   └── src/
    └── api/
        ├── package.json
        └── src/

Expected: Ship Shape detects npm workspace structure
```

### Testing Multi-Framework Detection

```
mixed-frameworks/
├── go.mod
├── legacy_test.go              # Uses testing package
└── modern_test.go              # Uses testify

Expected: Ship Shape detects both testing patterns
```

### Testing No Tests Scenario

```
no-tests/
├── go.mod
├── main.go                     # Production code only
└── internal/
    └── pkg/

Expected: Ship Shape reports 0% coverage, recommends adding tests
```

## Validation Frequency

Integration tests run:

- **Every PR**: Prevent regressions
- **Nightly**: Comprehensive validation
- **Before Release**: Full test suite execution

## Performance Benchmarks

Integration tests also track performance:

```yaml
performance_benchmarks:
  max_analysis_time_seconds: 5.0
  max_memory_mb: 256
  max_file_opens: 1000
```

---

**Last Updated**: 2026-01-28
**Repositories**: 0 (structure only, repositories to be added)
**Status**: Ready for integration test repositories
