# Ship Shape - Configuration Schema

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Complete
**Author**: Senior Software Engineer

---

## Table of Contents

1. [Overview](#overview)
2. [Schema Version](#schema-version)
3. [Analysis Configuration](#analysis-configuration)
4. [Scoring Configuration](#scoring-configuration)
5. [Coverage Configuration](#coverage-configuration)
6. [Quality Configuration](#quality-configuration)
7. [Quality Gates](#quality-gates)
8. [Output Configuration](#output-configuration)
9. [CI/CD Integration](#cicd-integration)
10. [Historical Tracking](#historical-tracking)
11. [Advanced Configuration](#advanced-configuration)
12. [Validation Rules](#validation-rules)
13. [Complete Examples](#complete-examples)

---

## Overview

Ship Shape uses YAML configuration files (`.shipshape.yml`) to customize analysis behavior, scoring weights, quality gates, and reporting options.

### File Naming Conventions

**Supported Names** (in order of precedence):
1. `.shipshape.yml` (recommended)
2. `.shipshape.yaml`
3. `shipshape.yml`
4. `shipshape.yaml`
5. `.shipshape.toml` (alternate format)
6. `.shipshape.json` (alternate format)

### Configuration Discovery

Ship Shape searches for configuration in:
1. Current directory (`.`)
2. Repository root (finds `.git` directory)
3. Home directory (`~/.config/shipshape/`)
4. System config (`/etc/shipshape/`)

### Quick Start

```yaml
# Minimal configuration
version: 1

gates:
  min-score: 70
```

---

## Schema Version

### `version`

Configuration schema version (required).

**Type**: Integer
**Required**: Yes
**Valid Values**: `1`
**Default**: None

```yaml
version: 1
```

**Purpose**: Ensures configuration compatibility with Ship Shape version.

---

## Analysis Configuration

### `analysis`

Controls analysis behavior, parallelization, and language detection.

#### `analysis.parallel`

Enable parallel analysis of multiple files/packages.

**Type**: Boolean
**Default**: `true`

```yaml
analysis:
  parallel: true
```

**Performance Impact**:
- `true`: Analyze files concurrently (faster)
- `false`: Sequential analysis (slower, lower memory)

#### `analysis.max-workers`

Maximum number of parallel workers.

**Type**: Integer
**Range**: `0-16`
**Default**: `0` (auto-detect from CPU cores)

```yaml
analysis:
  max-workers: 8
```

**Recommendations**:
- `0`: Auto-detect (recommended)
- `1-4`: Small projects or CI environments
- `8-16`: Large projects with available CPU

#### `analysis.timeout`

Global analysis timeout.

**Type**: Duration String
**Format**: `30s`, `5m`, `1h`
**Default**: `10m`

```yaml
analysis:
  timeout: 10m
```

**Common Values**:
- `5m`: Small projects (<1000 files)
- `10m`: Medium projects (1000-5000 files)
- `30m`: Large monorepos (>5000 files)

#### `analysis.analyzer-timeout`

Timeout per individual analyzer.

**Type**: Duration String
**Default**: `5m`

```yaml
analysis:
  analyzer-timeout: 5m
```

**Purpose**: Prevents single analyzer from blocking entire analysis.

#### `analysis.cache-ast`

Enable AST caching for faster re-analysis.

**Type**: Boolean
**Default**: `true`

```yaml
analysis:
  cache-ast: true
```

**Benefits**:
- Faster re-runs (unchanged files use cached AST)
- Reduced CPU usage
- Trade-off: Increased disk usage

#### `analysis.languages`

Languages to analyze.

**Type**: Array of Strings
**Valid Values**: `go`, `python`, `javascript`, `typescript`, `java`, `rust`
**Default**: `[]` (auto-detect)

```yaml
analysis:
  languages:
    - go
    - python
    - javascript
```

**Auto-Detection**: If empty, Ship Shape detects languages from file extensions.

#### `analysis.exclude`

Path patterns to exclude (glob format).

**Type**: Array of Strings
**Default**: Common exclusions (node_modules, vendor, etc.)

```yaml
analysis:
  exclude:
    - "**/node_modules/**"
    - "**/vendor/**"
    - "**/.venv/**"
    - "**/venv/**"
    - "**/__pycache__/**"
    - "**/target/**"
    - "**/build/**"
    - "**/dist/**"
    - "**/.next/**"
    - "**/.git/**"
```

**Pattern Syntax**:
- `**` matches zero or more directories
- `*` matches any characters within directory
- `?` matches single character

#### `analysis.include`

Path patterns to include (overrides exclude).

**Type**: Array of Strings
**Default**: None

```yaml
analysis:
  include:
    - "src/**"
    - "pkg/**"
    - "internal/**"
```

**Use Case**: Explicitly include paths that might be excluded by default.

### Complete Example

```yaml
analysis:
  parallel: true
  max-workers: 8
  timeout: 10m
  analyzer-timeout: 5m
  cache-ast: true
  languages:
    - go
    - python
  exclude:
    - "**/vendor/**"
    - "**/node_modules/**"
  include:
    - "src/**"
    - "pkg/**"
```

---

## Scoring Configuration

### `scoring`

Controls dimension weights and tier-based scoring.

#### `scoring.weights`

Dimension weights (must sum to 1.0).

**Type**: Object (Float values)
**Required**: No
**Validation**: Must sum to exactly `1.0`

```yaml
scoring:
  weights:
    coverage: 0.30        # 30% weight
    quality: 0.25         # 25% weight
    performance: 0.15     # 15% weight
    tools: 0.15           # 15% weight
    maintainability: 0.10 # 10% weight
    code-quality: 0.05    # 5% weight
```

**Dimensions**:

| Dimension | Description | Default Weight |
|-----------|-------------|----------------|
| `coverage` | Test coverage metrics (line, branch, function) | 0.30 |
| `quality` | Test quality (smells, patterns, organization) | 0.25 |
| `performance` | Test execution speed, flakiness | 0.15 |
| `tools` | Testing tool adoption and best practices | 0.15 |
| `maintainability` | Test-to-code ratio, duplication | 0.10 |
| `code-quality` | Integration with static analysis | 0.05 |

**Customization**:
```yaml
# Prioritize coverage
scoring:
  weights:
    coverage: 0.40
    quality: 0.30
    performance: 0.10
    tools: 0.10
    maintainability: 0.05
    code-quality: 0.05
```

#### `scoring.tier-weights`

Tier-based attribute weights (from AgentReady pattern).

**Type**: Object (Float values)
**Required**: No
**Validation**: Must sum to exactly `1.0`

```yaml
scoring:
  tier-weights:
    essential: 0.50   # Tier 1: Must-have fundamentals
    critical: 0.30    # Tier 2: Quality & maintainability
    important: 0.15   # Tier 3: Excellence & optimization
    advanced: 0.05    # Tier 4: Cutting-edge practices
```

**Tier System**:

| Tier | Name | Weight | Focus |
|------|------|--------|-------|
| 1 | Essential | 50% | Must-have fundamentals |
| 2 | Critical | 30% | Quality & maintainability |
| 3 | Important | 15% | Excellence & optimization |
| 4 | Advanced | 5% | Cutting-edge practices |

See [attribute-tier-system.md](./attribute-tier-system.md) for complete tier definitions.

#### `scoring.min-thresholds`

Minimum thresholds for each dimension (0-100).

**Type**: Object (Integer values)
**Range**: `0-100`
**Default**: See below

```yaml
scoring:
  min-thresholds:
    coverage: 50      # At least 50% coverage
    quality: 60       # At least 60% quality score
    performance: 0    # No minimum (varies by project)
    tools: 40         # Some basic tooling required
    maintainability: 50
    code-quality: 0
```

**Behavior**: If any dimension falls below threshold, overall grade is capped.

**Example**:
- Coverage: 40% (below 50% threshold)
- Quality: 90%
- Result: Overall grade capped at 'D' or 'F'

### Complete Example

```yaml
scoring:
  weights:
    coverage: 0.30
    quality: 0.25
    performance: 0.15
    tools: 0.15
    maintainability: 0.10
    code-quality: 0.05

  tier-weights:
    essential: 0.50
    critical: 0.30
    important: 0.15
    advanced: 0.05

  min-thresholds:
    coverage: 50
    quality: 60
    performance: 0
    tools: 40
    maintainability: 50
    code-quality: 0
```

---

## Coverage Configuration

### `coverage`

Controls coverage analysis and thresholds.

#### `coverage.paths`

Coverage file paths to analyze.

**Type**: Array of Strings
**Default**: `[]` (auto-detect)

```yaml
coverage:
  paths:
    - "coverage.out"
    - "coverage.xml"
    - "coverage/lcov.info"
```

**Auto-Detection** (if empty):
- **Go**: `coverage.out`, `coverage.txt`
- **Python**: `coverage.xml`, `.coverage`, `htmlcov/`
- **JavaScript**: `coverage/lcov.info`, `coverage-final.json`
- **Java**: `target/site/jacoco/jacoco.xml`

#### `coverage.thresholds`

Coverage thresholds (percentages).

**Type**: Object (Integer values)
**Range**: `0-100`
**Default**: See below

```yaml
coverage:
  thresholds:
    line: 80          # Target line coverage: ≥80%
    branch: 70        # Target branch coverage: ≥70%
    function: 90      # Target function coverage: ≥90%
```

**Recommendations**:

| Project Type | Line | Branch | Function |
|--------------|------|--------|----------|
| **New Projects** | 90% | 80% | 95% |
| **Established** | 80% | 70% | 90% |
| **Legacy** | 60% | 50% | 70% |

#### `coverage.critical-paths-required`

Report uncovered critical paths as high severity.

**Type**: Boolean
**Default**: `true`

```yaml
coverage:
  critical-paths-required: true
```

**Critical Paths**:
- Error handling code
- Security-sensitive functions
- Business logic core

### Complete Example

```yaml
coverage:
  paths:
    - "coverage.out"
    - "coverage/lcov.info"
  thresholds:
    line: 80
    branch: 70
    function: 90
  critical-paths-required: true
```

---

## Quality Configuration

### `quality`

Controls test quality analysis (smells, patterns).

#### `quality.smells`

Test smell detection configuration.

##### `quality.smells.enabled`

Enable test smell detection.

**Type**: Boolean
**Default**: `true`

```yaml
quality:
  smells:
    enabled: true
```

##### `quality.smells.max-per-100-tests`

Maximum acceptable smells per 100 tests.

**Type**: Integer
**Range**: `0-100`
**Default**: `5`

```yaml
quality:
  smells:
    max-per-100-tests: 5
```

**Interpretation**:
- `0`: Zero tolerance (very strict)
- `5`: Production quality (recommended)
- `10`: Moderate quality
- `20+`: Technical debt accumulating

##### `quality.smells.severity-overrides`

Override default severity for specific smells.

**Type**: Object (String values)
**Valid Severities**: `critical`, `high`, `medium`, `low`, `info`

```yaml
quality:
  smells:
    severity-overrides:
      mystery-guest: critical
      assertion-roulette: high
      eager-test: medium
```

##### `quality.smells.detect`

Enable/disable specific smell detectors.

**Type**: Object (Boolean values)
**Default**: All enabled

```yaml
quality:
  smells:
    detect:
      mystery-guest: true
      eager-test: true
      lazy-test: true
      assertion-roulette: true
      conditional-logic: true
      general-fixture: true
      obscure-test: true
      sensitive-equality: true
      resource-optimism: true
      code-duplication: true
      flakiness: true
```

**Test Smell Types**:

| Smell | Description | Default Severity |
|-------|-------------|------------------|
| `mystery-guest` | External dependencies not visible | High |
| `eager-test` | Testing multiple behaviors | Medium |
| `lazy-test` | Multiple test methods for one behavior | Low |
| `assertion-roulette` | Multiple asserts without context | Medium |
| `conditional-logic` | If/else in tests | Medium |
| `general-fixture` | Overly complex setup | Medium |
| `obscure-test` | Unclear test purpose | High |
| `sensitive-equality` | Brittle equality checks | Low |
| `resource-optimism` | Assuming resources available | High |
| `code-duplication` | Duplicated test code | Low |
| `flakiness` | Non-deterministic tests | Critical |

#### `quality.patterns`

Framework-specific best practices detection.

##### Go Patterns

```yaml
quality:
  patterns:
    go:
      min-table-driven-percentage: 50  # At least 50% table-driven tests
      min-parallel-percentage: 30      # At least 30% t.Parallel() usage
      require-testify: false           # Don't require testify library
```

**Go Best Practices**:
- **Table-driven tests**: Reduces duplication, improves maintainability
- **Parallel execution**: Faster test runs with `t.Parallel()`
- **Testify library**: Optional assertion library

##### Python Patterns

```yaml
quality:
  patterns:
    python:
      min-parametrize-percentage: 50   # At least 50% parametrized tests
      min-fixture-usage: 40            # At least 40% fixture usage
      prefer-pytest: true              # Recommend pytest over unittest
```

**Python Best Practices**:
- **Parametrized tests**: `@pytest.mark.parametrize` for data-driven tests
- **Fixture usage**: pytest fixtures for setup/teardown
- **pytest framework**: Modern testing framework (vs unittest)

##### JavaScript/TypeScript Patterns

```yaml
quality:
  patterns:
    javascript:
      min-describe-it-percentage: 80   # At least 80% describe/it structure
      min-async-await-percentage: 60   # Modern async patterns
      prefer-jest: false               # Don't require specific framework
```

**JavaScript Best Practices**:
- **describe/it structure**: Organized test suites
- **async/await**: Modern async pattern (vs callbacks/promises)
- **Framework flexibility**: Jest, Mocha, Vitest, etc.

### Complete Example

```yaml
quality:
  smells:
    enabled: true
    max-per-100-tests: 5
    severity-overrides:
      mystery-guest: critical
      flakiness: critical
    detect:
      mystery-guest: true
      eager-test: true
      lazy-test: true
      assertion-roulette: true
      conditional-logic: true
      general-fixture: true
      obscure-test: true
      sensitive-equality: true
      resource-optimism: true
      code-duplication: true
      flakiness: true

  patterns:
    go:
      min-table-driven-percentage: 50
      min-parallel-percentage: 30
      require-testify: false
    python:
      min-parametrize-percentage: 50
      min-fixture-usage: 40
      prefer-pytest: true
    javascript:
      min-describe-it-percentage: 80
      min-async-await-percentage: 60
      prefer-jest: false
```

---

## Quality Gates

### `gates`

Controls quality gate behavior for CI/CD integration.

#### `gates.fail-on`

Fail analysis if findings of specified severity are found.

**Type**: String (Enum)
**Options**: `critical`, `high`, `medium`, `low`, `info`, `none`
**Default**: `none`

```yaml
gates:
  fail-on: high
```

**Behavior**:
- `critical`: Fail only on critical findings
- `high`: Fail on high or critical findings
- `medium`: Fail on medium, high, or critical findings
- `low`: Fail on low, medium, high, or critical findings
- `info`: Fail on any findings
- `none`: Never fail based on findings

#### `gates.min-score`

Fail if overall score falls below threshold.

**Type**: Integer
**Range**: `0-100`
**Default**: `0` (disabled)

```yaml
gates:
  min-score: 70
```

**Grading Scale**:
- `90-100`: A (Excellent)
- `80-89`: B (Good)
- `70-79`: C (Satisfactory)
- `60-69`: D (Needs Improvement)
- `0-59`: F (Failing)

#### `gates.fail-on-regression`

Fail on score regression (requires baseline).

**Type**: Object
**Default**: Disabled

```yaml
gates:
  fail-on-regression:
    enabled: true
    baseline-file: shipshape-baseline.json
    max-score-drop: 5  # Fail if score drops >5 points
```

**Use Case**: Prevent quality decay in CI/CD.

#### `gates.fail-on-coverage-drop`

Fail on coverage regression.

**Type**: Object
**Default**: Disabled

```yaml
gates:
  fail-on-coverage-drop:
    enabled: true
    max-drop-percentage: 2  # Fail if coverage drops >2%
```

#### `gates.custom-gates`

Define custom gate rules.

**Type**: Array of Objects
**Default**: None

```yaml
gates:
  custom-gates:
    - name: "no-critical-smells"
      type: smell-count
      severity: critical
      max-count: 0

    - name: "min-test-count"
      type: test-count
      min-count: 10

    - name: "max-test-duration"
      type: test-duration
      max-duration: 10s
```

**Custom Gate Types**:

| Type | Parameters | Description |
|------|------------|-------------|
| `smell-count` | `severity`, `max-count` | Limit smells by severity |
| `test-count` | `min-count` | Minimum number of tests |
| `test-duration` | `max-duration` | Maximum test execution time |
| `coverage-metric` | `metric`, `min-value` | Custom coverage metric |

### Complete Example

```yaml
gates:
  fail-on: high
  min-score: 70

  fail-on-regression:
    enabled: true
    baseline-file: shipshape-baseline.json
    max-score-drop: 5

  fail-on-coverage-drop:
    enabled: true
    max-drop-percentage: 2

  custom-gates:
    - name: "no-critical-smells"
      type: smell-count
      severity: critical
      max-count: 0
    - name: "min-test-count"
      type: test-count
      min-count: 10
```

---

## Output Configuration

### `output`

Controls report generation and formatting.

#### `output.format`

Default output format.

**Type**: String (Enum)
**Options**: `html`, `json`, `markdown`, `text`
**Default**: `html`

```yaml
output:
  format: html
```

#### `output.path`

Default output path.

**Type**: String
**Default**: `shipshape-report.{format}`

```yaml
output:
  path: shipshape-report.html
```

#### `output.formats`

Generate multiple formats simultaneously.

**Type**: Array of Objects

```yaml
output:
  formats:
    - type: html
      path: shipshape-report.html
    - type: json
      path: shipshape-results.json
    - type: markdown
      path: SHIPSHAPE.md
```

#### `output.report`

Report customization options.

```yaml
output:
  report:
    title: "Ship Shape Analysis Report"
    include-findings: true
    include-code-snippets: true
    max-snippet-lines: 10
    theme: auto  # light, dark, auto
    open-browser: false
```

**Report Options**:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `title` | String | "Ship Shape Analysis Report" | Report title |
| `include-findings` | Boolean | `true` | Include detailed findings |
| `include-code-snippets` | Boolean | `true` | Include code snippets in findings |
| `max-snippet-lines` | Integer | `10` | Maximum lines per snippet |
| `theme` | String | `auto` | HTML theme: light, dark, auto |
| `open-browser` | Boolean | `false` | Open report after generation |

### Complete Example

```yaml
output:
  format: html
  path: shipshape-report.html

  formats:
    - type: html
      path: shipshape-report.html
    - type: json
      path: shipshape-results.json
    - type: markdown
      path: SHIPSHAPE.md

  report:
    title: "Ship Shape Analysis Report"
    include-findings: true
    include-code-snippets: true
    max-snippet-lines: 10
    theme: auto
    open-browser: false
```

---

## CI/CD Integration

### `ci`

CI/CD integration features.

#### `ci.enabled`

Enable CI/CD integration features.

**Type**: Boolean
**Default**: `true` (auto-detected)

```yaml
ci:
  enabled: true
```

**Auto-Detection**: Detects `CI=true` or `GITHUB_ACTIONS=true` environment variables.

#### `ci.upload-results`

Upload results to CI system.

**Type**: Boolean
**Default**: `true`

```yaml
ci:
  upload-results: true
```

#### `ci.github-actions`

GitHub Actions specific settings.

```yaml
ci:
  github-actions:
    create-check: true
    comment-pr: true
    update-comment: true
    create-annotations: true
    upload-artifact: true
    artifact-name: shipshape-report
```

**GitHub Actions Options**:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `create-check` | Boolean | `true` | Create GitHub check run |
| `comment-pr` | Boolean | `true` | Comment on pull requests |
| `update-comment` | Boolean | `true` | Update existing comment (don't spam) |
| `create-annotations` | Boolean | `true` | Create code annotations for findings |
| `upload-artifact` | Boolean | `true` | Upload report as artifact |
| `artifact-name` | String | `shipshape-report` | Artifact name |

#### `ci.artifacts`

Generic CI artifact settings.

```yaml
ci:
  artifacts:
    directory: .shipshape-artifacts
```

### Complete Example

```yaml
ci:
  enabled: true
  upload-results: true

  github-actions:
    create-check: true
    comment-pr: true
    update-comment: true
    create-annotations: true
    upload-artifact: true
    artifact-name: shipshape-report

  artifacts:
    directory: .shipshape-artifacts
```

---

## Historical Tracking

### `history`

Track analysis results over time.

#### `history.enabled`

Enable historical tracking.

**Type**: Boolean
**Default**: `false`

```yaml
history:
  enabled: true
```

#### `history.database`

Database file for historical data (SQLite).

**Type**: String
**Default**: `.shipshape/history.db`

```yaml
history:
  database: .shipshape/history.db
```

#### `history.retention`

Retention policy for historical data.

```yaml
history:
  retention:
    days: 365              # Keep data for 365 days
    max-analyses: 1000     # Maximum analyses to keep
```

### Complete Example

```yaml
history:
  enabled: true
  database: .shipshape/history.db
  retention:
    days: 365
    max-analyses: 1000
```

---

## Advanced Configuration

### `advanced`

Advanced settings for debugging and performance tuning.

#### `advanced.debug`

Enable debug logging.

**Type**: Boolean
**Default**: `false`

```yaml
advanced:
  debug: true
```

#### `advanced.log-file`

Log file path.

**Type**: String
**Default**: `stderr`

```yaml
advanced:
  log-file: shipshape-debug.log
```

#### `advanced.profile`

Enable profiling (for performance debugging).

```yaml
advanced:
  profile:
    enabled: true
    cpu: shipshape-cpu.prof
    memory: shipshape-mem.prof
```

#### `advanced.analyzers`

Custom analyzer configuration.

```yaml
advanced:
  analyzers:
    go-test-analyzer:
      timeout: 2m
      options:
        detect-benchmarks: true
    python-pytest-analyzer:
      timeout: 5m
      options:
        detect-markers: true
```

### Complete Example

```yaml
advanced:
  debug: false
  log-file: shipshape-debug.log

  profile:
    enabled: false
    cpu: shipshape-cpu.prof
    memory: shipshape-mem.prof

  analyzers:
    go-test-analyzer:
      timeout: 2m
      options:
        detect-benchmarks: true
```

---

## Validation Rules

### Weight Validation

**Rule**: Dimension weights must sum to exactly `1.0`

```yaml
# ✓ Valid
scoring:
  weights:
    coverage: 0.30
    quality: 0.25
    performance: 0.15
    tools: 0.15
    maintainability: 0.10
    code-quality: 0.05
  # Sum: 1.0

# ✗ Invalid
scoring:
  weights:
    coverage: 0.50
    quality: 0.50
  # Sum: 1.0 (missing dimensions)
```

**Rule**: Tier weights must sum to exactly `1.0`

```yaml
# ✓ Valid
scoring:
  tier-weights:
    essential: 0.50
    critical: 0.30
    important: 0.15
    advanced: 0.05
  # Sum: 1.0

# ✗ Invalid
scoring:
  tier-weights:
    essential: 0.60
    critical: 0.30
    important: 0.15
    advanced: 0.05
  # Sum: 1.10
```

### Threshold Validation

**Rule**: All thresholds must be in range `0-100`

```yaml
# ✓ Valid
coverage:
  thresholds:
    line: 80
    branch: 70
    function: 90

# ✗ Invalid
coverage:
  thresholds:
    line: 150  # Out of range
```

### Duration Validation

**Rule**: Durations must use valid format

```yaml
# ✓ Valid
analysis:
  timeout: 10m   # 10 minutes
  timeout: 30s   # 30 seconds
  timeout: 1h    # 1 hour

# ✗ Invalid
analysis:
  timeout: 10    # Missing unit
  timeout: 10min # Invalid unit (use 'm')
```

---

## Complete Examples

### Minimal Configuration

```yaml
# Absolute minimum
version: 1

gates:
  min-score: 70
```

### Standard Configuration

```yaml
# Recommended for most projects
version: 1

analysis:
  parallel: true
  timeout: 10m
  languages:
    - go
    - python

scoring:
  weights:
    coverage: 0.30
    quality: 0.25
    performance: 0.15
    tools: 0.15
    maintainability: 0.10
    code-quality: 0.05

coverage:
  thresholds:
    line: 80
    branch: 70
    function: 90

quality:
  smells:
    enabled: true
    max-per-100-tests: 5

gates:
  fail-on: high
  min-score: 70

output:
  format: html
  path: shipshape-report.html

ci:
  enabled: true
  github-actions:
    comment-pr: true
```

### Strict Configuration (Production)

```yaml
# High standards for production systems
version: 1

analysis:
  parallel: true
  timeout: 15m
  cache-ast: true

scoring:
  weights:
    coverage: 0.35
    quality: 0.30
    performance: 0.15
    tools: 0.10
    maintainability: 0.05
    code-quality: 0.05

  min-thresholds:
    coverage: 80
    quality: 80
    tools: 60

coverage:
  thresholds:
    line: 90
    branch: 80
    function: 95
  critical-paths-required: true

quality:
  smells:
    enabled: true
    max-per-100-tests: 2
    severity-overrides:
      mystery-guest: critical
      flakiness: critical

gates:
  fail-on: high
  min-score: 85

  fail-on-regression:
    enabled: true
    baseline-file: shipshape-baseline.json
    max-score-drop: 3

  fail-on-coverage-drop:
    enabled: true
    max-drop-percentage: 1

  custom-gates:
    - name: "no-critical-smells"
      type: smell-count
      severity: critical
      max-count: 0

output:
  formats:
    - type: html
      path: shipshape-report.html
    - type: json
      path: shipshape-results.json

ci:
  enabled: true
  github-actions:
    create-check: true
    comment-pr: true
    create-annotations: true

history:
  enabled: true
  database: .shipshape/history.db
```

### Comprehensive Configuration (All Options)

See [.shipshape.example.yml](../../.shipshape.example.yml) for complete example with all options documented.

---

## Migration Guide

### From Version 0.x to 1.0

```yaml
# Before (0.x)
coverage: 80
quality_gates:
  - type: smell
    max: 5

# After (1.0)
version: 1

coverage:
  thresholds:
    line: 80

gates:
  fail-on: high

quality:
  smells:
    max-per-100-tests: 5
```

---

## Best Practices

### 1. Start Minimal, Grow Gradually

```yaml
# Start with minimal config
version: 1
gates:
  min-score: 60

# Gradually increase standards
version: 1
gates:
  min-score: 70  # After improvements
```

### 2. Use Environment-Specific Configs

```bash
# Development
.shipshape.dev.yml (lenient)

# CI/CD
.shipshape.ci.yml (moderate)

# Production
.shipshape.production.yml (strict)
```

### 3. Document Custom Settings

```yaml
# Custom configuration for legacy codebase
version: 1

# Lower thresholds due to technical debt
# Target: Increase by 5% per quarter
coverage:
  thresholds:
    line: 60  # Goal: 80% by Q4 2026
    branch: 50  # Goal: 70% by Q4 2026
```

### 4. Version Control Configuration

```bash
# Commit configuration to git
git add .shipshape.yml
git commit -m "chore: add Ship Shape configuration"
```

---

## Troubleshooting

### Configuration Not Found

```bash
# Verify configuration is in expected location
ls -la .shipshape.yml

# Or specify path explicitly
shipshape analyze -c /path/to/config.yml
```

### Validation Errors

```bash
# Validate configuration
shipshape validate .shipshape.yml

# Strict validation
shipshape validate --strict .shipshape.yml
```

### Weight Sum Errors

```yaml
# Error: weights sum to 0.95 (not 1.0)
scoring:
  weights:
    coverage: 0.30
    quality: 0.25
    performance: 0.15
    tools: 0.15
    maintainability: 0.10
    # Missing code-quality: 0.05
```

---

## Additional Resources

- **Example Configuration**: [.shipshape.example.yml](../../.shipshape.example.yml)
- **CLI Reference**: [cli-reference.md](./cli-reference.md)
- **Attribute Tier System**: [attribute-tier-system.md](./attribute-tier-system.md)
- **Architecture**: [../architecture.md](../architecture.md)

---

**Document Status**: Complete - Ready for Implementation
**Last Updated**: 2026-01-27
**Version**: 1.0.0
**Author**: Senior Software Engineer
