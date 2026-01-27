# Ship Shape - CLI Reference

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Complete
**Author**: Senior Software Engineer

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Global Flags](#global-flags)
4. [Commands](#commands)
   - [analyze](#shipshape-analyze)
   - [report](#shipshape-report)
   - [gate](#shipshape-gate)
   - [compare](#shipshape-compare)
   - [init](#shipshape-init)
   - [validate](#shipshape-validate)
   - [version](#shipshape-version)
5. [Environment Variables](#environment-variables)
6. [Exit Codes](#exit-codes)
7. [Configuration Files](#configuration-files)
8. [Shell Completion](#shell-completion)
9. [Examples](#examples)

---

## Overview

Ship Shape is a comprehensive test quality analysis tool that evaluates your codebase's testing practices, coverage metrics, and overall test health.

**Binary Name**: `shipshape`
**Module Path**: `github.com/chambridge/ship-shape`
**Supported Languages**: Go, Python, JavaScript/TypeScript, Java, Rust

### Quick Start

```bash
# Analyze current directory
shipshape analyze

# Analyze with quality gates
shipshape analyze --fail-on high --min-score 80

# Generate report from previous analysis
shipshape report -i results.json -o report.html

# Initialize configuration
shipshape init
```

---

## Installation

### Go Install (Recommended)

```bash
go install github.com/chambridge/ship-shape/cmd/shipshape@latest
```

### Binary Download

```bash
# Linux
curl -L https://github.com/chambridge/ship-shape/releases/latest/download/shipshape-linux-amd64 -o shipshape
chmod +x shipshape
sudo mv shipshape /usr/local/bin/

# macOS
curl -L https://github.com/chambridge/ship-shape/releases/latest/download/shipshape-darwin-arm64 -o shipshape
chmod +x shipshape
sudo mv shipshape /usr/local/bin/

# Windows
# Download shipshape-windows-amd64.exe from releases
```

### From Source

```bash
git clone https://github.com/chambridge/ship-shape.git
cd ship-shape
make build
sudo make install
```

### Verify Installation

```bash
shipshape version
```

---

## Global Flags

Global flags are available for all commands.

### `--config, -c <file>`

Configuration file path.

**Default**: Auto-discovered (`.shipshape.yml`, `shipshape.yml`, etc.)
**Type**: String
**Example**: `shipshape analyze -c .shipshape.production.yml`

**Search Order**:
1. Current directory (`.`)
2. Repository root (find `.git` directory)
3. Home directory (`~/.config/shipshape/`)
4. System config (`/etc/shipshape/`)

### `--verbose, -v`

Enable verbose output with detailed logging.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape analyze -v`

**Output Example**:
```
Using config file: /path/to/.shipshape.yml
Discovered 42 test files
Running go-test-analyzer...
Running coverage-parser...
Analysis complete in 3.2s
```

### `--quiet, -q`

Quiet mode - show errors only.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape analyze -q`

**Note**: Mutually exclusive with `--verbose`

### `--no-color`

Disable colored terminal output.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape analyze --no-color`

**Use Cases**:
- CI/CD environments
- Log file generation
- Terminal compatibility issues

---

## Commands

### `shipshape analyze`

Analyze a repository and generate quality assessment.

#### Synopsis

```bash
shipshape analyze [path] [flags]
```

#### Arguments

- `[path]` - Repository path to analyze (default: current directory `.`)

#### Flags

##### `--output, -o <file>`

Output file path.

**Default**: `shipshape-report.html`
**Type**: String
**Example**: `shipshape analyze -o results.html`

##### `--format, -f <format>`

Output format.

**Options**: `html`, `json`, `markdown`, `text`
**Default**: `html`
**Type**: String
**Example**: `shipshape analyze -f json -o results.json`

**Format Details**:
- `html` - Interactive HTML report with charts
- `json` - Machine-readable JSON for CI/CD integration
- `markdown` - Human-readable Markdown for documentation
- `text` - Plain text for terminal output

##### `--fail-on <severity>`

Fail if findings of specified severity are found.

**Options**: `critical`, `high`, `medium`, `low`, `info`, `none`
**Default**: `none`
**Type**: String
**Example**: `shipshape analyze --fail-on high`

**Exit Behavior**:
- Exits with code `1` if findings ≥ specified severity
- Useful for CI/CD quality gates

##### `--min-score <score>`

Fail if overall score falls below threshold.

**Range**: `0-100`
**Default**: `0` (disabled)
**Type**: Integer
**Example**: `shipshape analyze --min-score 80`

##### `--parallel`

Enable parallel analysis of files/packages.

**Default**: `true`
**Type**: Boolean
**Example**: `shipshape analyze --parallel=false`

##### `--max-workers <n>`

Maximum number of parallel workers.

**Range**: `1-16`
**Default**: `0` (auto-detect from CPU cores)
**Type**: Integer
**Example**: `shipshape analyze --max-workers 4`

##### `--timeout <duration>`

Global analysis timeout.

**Format**: `30s`, `5m`, `1h`
**Default**: `10m`
**Type**: Duration
**Example**: `shipshape analyze --timeout 5m`

##### `--languages <langs>`

Comma-separated languages to analyze.

**Options**: `go`, `python`, `javascript`, `typescript`, `java`, `rust`
**Default**: Auto-detect
**Type**: String
**Example**: `shipshape analyze --languages go,python`

##### `--exclude <patterns>`

Exclude paths matching glob patterns (can be specified multiple times).

**Format**: Glob pattern
**Type**: String (repeatable)
**Example**:
```bash
shipshape analyze --exclude "**/vendor/**" --exclude "**/node_modules/**"
```

##### `--cache`

Enable AST caching for faster re-analysis.

**Default**: `true`
**Type**: Boolean
**Example**: `shipshape analyze --cache=false`

##### `--no-upload`

Don't upload results to CI system (GitHub Actions, etc.).

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape analyze --no-upload`

#### Examples

```bash
# Analyze current directory, output to HTML
shipshape analyze

# Analyze specific directory with JSON output
shipshape analyze /path/to/repo -f json -o results.json

# Analyze and fail on high severity findings
shipshape analyze --fail-on high --min-score 80

# Analyze only Go and Python files
shipshape analyze --languages go,python

# Analyze with custom config
shipshape analyze -c .shipshape.yml

# Analyze with exclusions
shipshape analyze --exclude "**/vendor/**" --exclude "**/test_data/**"

# Fast analysis with 8 workers, 2-minute timeout
shipshape analyze --max-workers 8 --timeout 2m
```

#### Exit Codes

- `0` - Success, no quality gate failures
- `1` - Quality gate failure (based on `--fail-on` or `--min-score`)
- `2` - Analysis error (parsing failure, timeout, etc.)
- `3` - Configuration error

---

### `shipshape report`

Generate reports from previous analysis results.

#### Synopsis

```bash
shipshape report [flags]
```

#### Flags

##### `--input, -i <file>`

Input analysis results file (JSON format).

**Required**: Yes
**Type**: String
**Example**: `shipshape report -i results.json`

##### `--output, -o <file>`

Output report file.

**Default**: `shipshape-report.<format>`
**Type**: String
**Example**: `shipshape report -i results.json -o report.html`

##### `--format, -f <format>`

Output format.

**Options**: `html`, `markdown`, `text`
**Default**: `html`
**Type**: String
**Example**: `shipshape report -i results.json -f markdown`

##### `--template <file>`

Custom template file (for HTML/Markdown).

**Type**: String
**Example**: `shipshape report -i results.json --template my-template.html`

**Template Format**:
- HTML templates use Go `html/template` syntax
- Markdown templates use Go `text/template` syntax

##### `--open`

Open report in browser after generation (HTML only).

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape report -i results.json --open`

#### Examples

```bash
# Generate HTML report from JSON results
shipshape report -i results.json -o report.html

# Generate Markdown report and open
shipshape report -i results.json -f markdown --open

# Use custom template
shipshape report -i results.json --template my-template.html

# Generate multiple formats
shipshape report -i results.json -o report.html
shipshape report -i results.json -f markdown -o REPORT.md
shipshape report -i results.json -f text -o report.txt
```

#### Exit Codes

- `0` - Success
- `2` - Report generation error
- `3` - Invalid input file

---

### `shipshape gate`

Evaluate quality gates without full analysis (fast check).

#### Synopsis

```bash
shipshape gate [flags]
```

#### Flags

##### `--input, -i <file>`

Input analysis results file (JSON format).

**Required**: Yes
**Type**: String
**Example**: `shipshape gate -i results.json`

##### `--fail-on <severity>`

Fail if findings of specified severity are found.

**Options**: `critical`, `high`, `medium`, `low`, `info`
**Type**: String
**Example**: `shipshape gate -i results.json --fail-on high`

##### `--min-score <score>`

Fail if overall score below threshold.

**Range**: `0-100`
**Type**: Integer
**Example**: `shipshape gate -i results.json --min-score 80`

##### `--config, -c <file>`

Configuration file with gate definitions.

**Type**: String
**Example**: `shipshape gate -i results.json -c .shipshape.yml`

#### Examples

```bash
# Check if previous analysis passes quality gate
shipshape gate -i results.json --min-score 80

# Use gates defined in config
shipshape gate -i results.json -c .shipshape.yml

# Fail on high severity findings
shipshape gate -i results.json --fail-on high

# Combined gates
shipshape gate -i results.json --fail-on high --min-score 75
```

#### Exit Codes

- `0` - Quality gates passed
- `1` - Quality gates failed
- `2` - Error evaluating gates

---

### `shipshape compare`

Compare two analysis results (e.g., baseline vs current).

#### Synopsis

```bash
shipshape compare <baseline> <current> [flags]
```

#### Arguments

- `<baseline>` - Baseline analysis results (JSON)
- `<current>` - Current analysis results (JSON)

#### Flags

##### `--output, -o <file>`

Output comparison report.

**Default**: Standard output
**Type**: String
**Example**: `shipshape compare baseline.json current.json -o comparison.html`

##### `--format, -f <format>`

Output format.

**Options**: `html`, `markdown`, `text`
**Default**: `text`
**Type**: String
**Example**: `shipshape compare baseline.json current.json -f html`

##### `--show-improvements`

Show improvements, not just regressions.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape compare baseline.json current.json --show-improvements`

##### `--fail-on-regression`

Fail if any metric regressed.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape compare baseline.json current.json --fail-on-regression`

#### Examples

```bash
# Compare baseline to current (text output)
shipshape compare baseline.json current.json

# Generate HTML comparison report
shipshape compare baseline.json current.json -f html -o comparison.html

# Show all changes (improvements and regressions)
shipshape compare baseline.json current.json --show-improvements

# Fail CI if any regression
shipshape compare baseline.json current.json --fail-on-regression

# Combined options
shipshape compare baseline.json current.json \
  -f html -o comparison.html \
  --show-improvements \
  --fail-on-regression
```

#### Exit Codes

- `0` - Success, no regressions (or not checking)
- `1` - Regression detected (if `--fail-on-regression`)
- `2` - Comparison error

---

### `shipshape init`

Initialize `.shipshape.yml` configuration file.

#### Synopsis

```bash
shipshape init [flags]
```

#### Flags

##### `--interactive, -i`

Interactive configuration wizard.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape init -i`

**Wizard Steps**:
1. Select languages to analyze
2. Set coverage thresholds
3. Choose quality gate settings
4. Configure output preferences
5. Review and save

##### `--template <name>`

Use predefined template.

**Options**: `minimal`, `standard`, `strict`, `comprehensive`
**Default**: `standard`
**Type**: String
**Example**: `shipshape init --template minimal`

**Template Descriptions**:
- `minimal` - Bare minimum configuration
- `standard` - Recommended defaults for most projects
- `strict` - High standards for production systems
- `comprehensive` - All available options with documentation

##### `--overwrite`

Overwrite existing `.shipshape.yml`.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape init --overwrite`

#### Examples

```bash
# Interactive wizard
shipshape init -i

# Create minimal config
shipshape init --template minimal

# Create comprehensive config (all options)
shipshape init --template comprehensive

# Overwrite existing config
shipshape init --template standard --overwrite
```

#### Output

Creates `.shipshape.yml` in current directory.

#### Exit Codes

- `0` - Success
- `2` - File already exists (without `--overwrite`)
- `3` - Invalid template

---

### `shipshape validate`

Validate `.shipshape.yml` configuration file.

#### Synopsis

```bash
shipshape validate [file] [flags]
```

#### Arguments

- `[file]` - Configuration file path (default: auto-discover)

#### Flags

##### `--strict`

Enable strict validation (warn on unknown keys).

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape validate --strict`

**Strict Mode**:
- Warns on unknown configuration keys
- Validates threshold ranges
- Checks weight sums (must equal 1.0)
- Validates regex patterns

#### Examples

```bash
# Validate default config
shipshape validate

# Validate specific config file
shipshape validate .shipshape.production.yml

# Strict validation
shipshape validate --strict

# Validate with verbose output
shipshape validate -v .shipshape.yml
```

#### Output

```bash
# Success
✓ Configuration is valid
  Version: 1
  Languages: go, python
  Quality gates: fail-on=high, min-score=70

# Error
✗ Configuration is invalid
  Error: scoring.weights must sum to 1.0 (currently: 0.95)
  Error: coverage.thresholds.line must be 0-100 (got: 150)
```

#### Exit Codes

- `0` - Valid configuration
- `1` - Invalid configuration
- `2` - Configuration file not found

---

### `shipshape version`

Show version information.

#### Synopsis

```bash
shipshape version [flags]
```

#### Flags

##### `--short`

Show version number only.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape version --short`

**Output**: `1.0.0`

##### `--json`

Output as JSON.

**Default**: `false`
**Type**: Boolean
**Example**: `shipshape version --json`

**Output**:
```json
{
  "version": "1.0.0",
  "build_time": "2026-01-27_15:30:00",
  "git_commit": "abc1234",
  "go_version": "go1.21.5",
  "platform": "darwin/arm64"
}
```

#### Examples

```bash
# Show full version info
shipshape version

# Show version number only
shipshape version --short

# JSON output
shipshape version --json

# Use in scripts
VERSION=$(shipshape version --short)
echo "Using Ship Shape v$VERSION"
```

#### Output Example

```
Ship Shape v1.0.0
Build Date: 2026-01-27_15:30:00
Git Commit: abc1234
Go Version: go1.21.5
Platform: darwin/arm64
```

---

## Environment Variables

Ship Shape respects environment variables for configuration.

### Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `SHIPSHAPE_CONFIG` | Configuration file path | Auto-discover | `/path/to/config.yml` |

### Analysis

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `SHIPSHAPE_PARALLEL` | Enable parallel analysis | `true` | `true`, `false` |
| `SHIPSHAPE_TIMEOUT` | Analysis timeout | `10m` | `5m`, `30m`, `1h` |
| `SHIPSHAPE_MAX_WORKERS` | Maximum parallel workers | `0` (auto) | `4`, `8`, `16` |

### Output

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `SHIPSHAPE_OUTPUT_FORMAT` | Output format | `html` | `json`, `markdown` |
| `SHIPSHAPE_OUTPUT_PATH` | Output file path | `shipshape-report.html` | `./results.json` |
| `SHIPSHAPE_NO_COLOR` | Disable colored output | `false` | `true`, `false` |

### Quality Gates

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `SHIPSHAPE_FAIL_ON` | Fail on severity | `none` | `critical`, `high` |
| `SHIPSHAPE_MIN_SCORE` | Minimum score threshold | `0` | `70`, `80`, `90` |

### CI/CD Integration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `CI` | CI environment detected | Auto-detect | `true`, `false` |
| `GITHUB_ACTIONS` | GitHub Actions detected | Auto-detect | `true`, `false` |
| `SHIPSHAPE_NO_UPLOAD` | Disable CI uploads | `false` | `true`, `false` |

### Precedence

Configuration precedence (highest to lowest):

1. **CLI flags** (e.g., `--fail-on high`)
2. **Environment variables** (e.g., `SHIPSHAPE_FAIL_ON=high`)
3. **Configuration file** (`.shipshape.yml`)
4. **Built-in defaults**

### Example Usage

```bash
# Set configuration via environment
export SHIPSHAPE_CONFIG=/path/to/.shipshape.yml
export SHIPSHAPE_FAIL_ON=high
export SHIPSHAPE_MIN_SCORE=80

# Run analysis (uses environment variables)
shipshape analyze

# Override with flags (flags take precedence)
shipshape analyze --fail-on critical --min-score 90
```

---

## Exit Codes

Ship Shape uses standard exit codes for scripting and CI/CD integration.

| Code | Name | Description | Commands |
|------|------|-------------|----------|
| `0` | Success | Operation completed successfully | All |
| `1` | Quality Gate Failure | Quality gates failed (findings or score) | analyze, gate, compare |
| `2` | Analysis/Operation Error | Analysis failed, report generation failed, etc. | analyze, report, gate, compare |
| `3` | Configuration Error | Invalid or missing configuration | analyze, init, validate |

### Usage in Scripts

```bash
#!/bin/bash

# Run analysis and capture exit code
shipshape analyze --fail-on high --min-score 80
EXIT_CODE=$?

case $EXIT_CODE in
  0)
    echo "✓ Analysis passed all quality gates"
    ;;
  1)
    echo "✗ Quality gate failure"
    exit 1
    ;;
  2)
    echo "✗ Analysis error"
    exit 2
    ;;
  3)
    echo "✗ Configuration error"
    exit 3
    ;;
esac
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Run Ship Shape Analysis
  run: shipshape analyze --fail-on high --min-score 80
  continue-on-error: false  # Fail build on non-zero exit

- name: Upload Report
  if: always()  # Upload even on failure
  uses: actions/upload-artifact@v3
  with:
    name: shipshape-report
    path: shipshape-report.html
```

---

## Configuration Files

### Discovery Order

Ship Shape searches for configuration files in this order:

1. **CLI flag**: `--config /path/to/config.yml`
2. **Current directory**: `.shipshape.yml`, `.shipshape.yaml`
3. **Repository root** (finds `.git`): `.shipshape.yml`, `shipshape.yml`
4. **Home directory**: `~/.config/shipshape/.shipshape.yml`
5. **System config**: `/etc/shipshape/.shipshape.yml`

### Supported Formats

| Format | Extensions | Priority |
|--------|-----------|----------|
| YAML | `.yml`, `.yaml` | Primary (recommended) |
| TOML | `.toml` | Alternate |
| JSON | `.json` | Alternate |

### Example Structure

```yaml
# .shipshape.yml
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

gates:
  fail-on: high
  min-score: 70

output:
  format: html
  path: shipshape-report.html
```

See [configuration-schema.md](./configuration-schema.md) for complete reference.

---

## Shell Completion

Ship Shape supports shell completion for all major shells.

### Bash

```bash
# Generate completion script
shipshape completion bash > /etc/bash_completion.d/shipshape

# Or add to .bashrc
echo "source <(shipshape completion bash)" >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```bash
# Generate completion script
shipshape completion zsh > "${fpath[1]}/_shipshape"

# Or add to .zshrc
echo "source <(shipshape completion zsh)" >> ~/.zshrc
source ~/.zshrc
```

### Fish

```bash
# Generate completion script
shipshape completion fish > ~/.config/fish/completions/shipshape.fish

# Or source directly
shipshape completion fish | source
```

### PowerShell

```powershell
# Generate completion script
shipshape completion powershell > shipshape.ps1

# Add to profile
shipshape completion powershell >> $PROFILE
```

### Testing Completion

```bash
# After installation, test by typing and pressing TAB
shipshape <TAB>
# Should show: analyze, report, gate, compare, init, validate, version

shipshape analyze --<TAB>
# Should show: --output, --format, --fail-on, --min-score, etc.
```

---

## Examples

### Basic Workflows

#### Analyze Current Project

```bash
# Simple analysis
shipshape analyze

# View report
open shipshape-report.html
```

#### CI/CD Integration

```bash
# Analyze with strict gates
shipshape analyze \
  --fail-on high \
  --min-score 80 \
  -f json \
  -o results.json

# Check gates
shipshape gate -i results.json --min-score 80

# Generate HTML report for artifact
shipshape report -i results.json -o report.html
```

#### Progressive Analysis

```bash
# Baseline (first run)
shipshape analyze -f json -o baseline.json

# Make improvements...

# Current analysis
shipshape analyze -f json -o current.json

# Compare results
shipshape compare baseline.json current.json --show-improvements
```

### Advanced Workflows

#### Multi-Format Output

```bash
# Generate all formats
shipshape analyze -f json -o results.json
shipshape report -i results.json -f html -o report.html
shipshape report -i results.json -f markdown -o REPORT.md
shipshape report -i results.json -f text -o report.txt
```

#### Language-Specific Analysis

```bash
# Go only
shipshape analyze --languages go

# Python only
shipshape analyze --languages python

# Multiple languages
shipshape analyze --languages go,python,javascript
```

#### Custom Configuration

```bash
# Development (lenient)
shipshape analyze -c .shipshape.dev.yml

# Production (strict)
shipshape analyze -c .shipshape.production.yml --fail-on high
```

### Pre-Commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

shipshape analyze \
  --fail-on critical \
  --timeout 30s \
  --quiet \
  -f text

exit $?
```

### GitHub Actions

```yaml
name: Ship Shape Analysis

on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install Ship Shape
        run: go install github.com/chambridge/ship-shape/cmd/shipshape@latest

      - name: Run Analysis
        run: |
          shipshape analyze \
            --fail-on high \
            --min-score 80 \
            -f json \
            -o results.json

      - name: Generate Report
        if: always()
        run: shipshape report -i results.json -o report.html

      - name: Upload Report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: shipshape-report
          path: |
            results.json
            report.html
```

---

## Troubleshooting

### Common Issues

#### Configuration Not Found

```
Error: config file not found
```

**Solution**:
- Ensure `.shipshape.yml` exists in current directory or repository root
- Use `-c` flag to specify path explicitly
- Run `shipshape init` to create configuration

#### Analysis Timeout

```
Error: analysis timed out after 10m
```

**Solution**:
- Increase timeout: `shipshape analyze --timeout 30m`
- Reduce parallelism: `shipshape analyze --max-workers 2`
- Exclude large directories: `shipshape analyze --exclude "**/vendor/**"`

#### Quality Gate Failure

```
Exit code 1: quality gate failure
```

**Solution**:
- Review findings in report
- Adjust thresholds in `.shipshape.yml`
- Use `--fail-on` flag to control severity level

### Debug Mode

```bash
# Enable verbose output
shipshape analyze -v

# Show configuration being used
shipshape validate -v
```

---

## Additional Resources

- **Project Documentation**: [.project/v1.0.0/](../)
- **Configuration Schema**: [configuration-schema.md](./configuration-schema.md)
- **Architecture**: [../architecture.md](../architecture.md)
- **User Stories**: [../user-stories.md](../user-stories.md)
- **GitHub Repository**: https://github.com/chambridge/ship-shape

---

**Document Status**: Complete - Ready for Implementation
**Last Updated**: 2026-01-27
**Version**: 1.0.0
**Author**: Senior Software Engineer
