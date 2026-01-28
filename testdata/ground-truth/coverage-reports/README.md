# Coverage Reports Ground Truth Examples

This directory contains sample coverage report files for validating Ship Shape's coverage parsing capabilities.

## Purpose

These coverage reports are used to:

1. **Validate Parsing**: Ensure Ship Shape correctly parses all supported coverage formats
2. **Test Edge Cases**: Handle malformed, empty, or unusual coverage reports
3. **Benchmark Performance**: Measure parsing speed for large reports
4. **Prevent Regressions**: Ensure updates don't break existing parsers

## Supported Formats

### Go Coverage (`go-cover`)

Go's native coverage format from `go test -coverprofile=coverage.out`.

**Format characteristics**:
- Text-based format
- Line-by-line coverage with block counts
- Function-level granularity

**Example**:
```
mode: set
github.com/example/pkg/file.go:10.15,12.2 1 1
github.com/example/pkg/file.go:14.20,16.2 1 0
```

**Tools**:
- `go test -coverprofile`
- `go tool cover`
- `gocov`

### Python Coverage (`coverage.py`)

Coverage.py generates XML and JSON formats.

**XML Format (Cobertura)**:
- Standard XML coverage format
- Package, class, and method granularity
- Line and branch coverage

**JSON Format**:
- Structured coverage data
- Per-file line coverage
- Summary statistics

**Tools**:
- `coverage.py`
- `pytest-cov`

### JavaScript Coverage (`istanbul`/`nyc`)

JavaScript coverage tools output LCOV and JSON formats.

**LCOV Format**:
- Text-based format
- Line and function coverage
- Industry standard

**Istanbul JSON Format**:
- Detailed per-file coverage
- Statement, branch, function, and line coverage
- Source map support

**Tools**:
- `nyc` (Istanbul CLI)
- `c8` (V8 coverage)
- `jest --coverage`

### Java Coverage (`jacoco`)

JaCoCo XML format for Java coverage.

**Format characteristics**:
- XML-based
- Package, class, method granularity
- Instruction and branch coverage

**Tools**:
- JaCoCo Maven plugin
- JaCoCo Gradle plugin

## Directory Structure

```
coverage-reports/
├── go/
│   ├── minimal.coverprofile           # Minimal valid coverage file
│   ├── multi-package.coverprofile     # Multiple packages
│   ├── zero-coverage.coverprofile     # 0% coverage
│   ├── full-coverage.coverprofile     # 100% coverage
│   └── metadata.yml
├── python/
│   ├── coverage.xml                   # Cobertura XML
│   ├── coverage.json                  # coverage.py JSON
│   ├── pytest-cov.xml                 # pytest-cov output
│   └── metadata.yml
├── javascript/
│   ├── coverage.lcov                  # LCOV format
│   ├── coverage-final.json            # Istanbul JSON
│   ├── c8-coverage.json               # c8 format
│   └── metadata.yml
└── java/
    ├── jacoco.xml                     # JaCoCo XML
    ├── jacoco-minimal.xml             # Minimal example
    └── metadata.yml
```

## Example Metadata

Each language directory includes `metadata.yml` describing the coverage reports:

```yaml
version: "1.0.0"
language: go
category: coverage-report
coverage_format: go-cover
description: "Sample Go coverage reports for parser validation"
verified_by:
  - "engineer1@example.com"
  - "engineer2@example.com"
verified_date: "2026-01-28"

files:
  - name: "minimal.coverprofile"
    description: "Minimal valid coverage file with single package"
    expected_coverage:
      line_coverage: 75.0
      files_covered: 1
      lines_total: 20
      lines_covered: 15

  - name: "multi-package.coverprofile"
    description: "Coverage across multiple packages"
    expected_coverage:
      line_coverage: 82.5
      files_covered: 8
      lines_total: 500
      lines_covered: 412

  - name: "zero-coverage.coverprofile"
    description: "Valid report with 0% coverage"
    expected_coverage:
      line_coverage: 0.0
      files_covered: 3
      lines_total: 150
      lines_covered: 0

test_requirements:
  min_go_version: "1.21"

notes: |
  Coverage files generated with:
  - go test -coverprofile=coverage.out ./...
  - go tool cover -func=coverage.out
```

## Validation Tests

Coverage parser tests verify:

1. **Correct Parsing**: Metrics match expected values
2. **Format Detection**: Auto-detect coverage format
3. **Error Handling**: Graceful handling of malformed files
4. **Performance**: Parse large files efficiently

```go
func TestParseCoverageReport(t *testing.T) {
    tests := []struct {
        file     string
        format   string
        expected CoverageMetrics
    }{
        {
            file:   "testdata/ground-truth/coverage-reports/go/minimal.coverprofile",
            format: "go-cover",
            expected: CoverageMetrics{
                LineCoverage:   75.0,
                FilesCovered:   1,
                LinesTotal:     20,
                LinesCovered:   15,
            },
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.file, func(t *testing.T) {
            metrics, err := ParseCoverageFile(tt.file, tt.format)
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, metrics)
        })
    }
}
```

## Adding New Coverage Reports

When adding coverage reports:

1. **Generate from real code**: Use actual test suites to generate reports
2. **Verify metrics**: Manually confirm coverage percentages
3. **Include edge cases**: 0%, 100%, partial coverage
4. **Document generation**: Note command used to create the report
5. **Update metadata.yml**: Add expected metrics

## Coverage Format References

- **Go**: https://go.dev/testing/coverage
- **Python**: https://coverage.readthedocs.io/
- **JavaScript (nyc)**: https://github.com/istanbuljs/nyc
- **Java (JaCoCo)**: https://www.jacoco.org/

---

**Last Updated**: 2026-01-28
**Formats**: 4 (go-cover, cobertura-xml, lcov, jacoco-xml)
**Status**: Structure only, coverage files to be added
