# Ship Shape - Project Basics

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Design Complete
**Author**: Senior Software Engineer

---

## Table of Contents

1. [Go Module Path](#go-module-path)
2. [CLI Framework Selection](#cli-framework-selection)
3. [Configuration Format](#configuration-format)
4. [CLI Interface Design](#cli-interface-design)
5. [Example Configuration](#example-configuration)

---

## Go Module Path

### Decision: `github.com/chambridge/ship-shape`

**Rationale**:
- Standard Go module path format
- GitHub is the primary hosting platform
- Clear, memorable repository name
- No hyphens in import path (Go convention)

**Module Declaration**:
```go
// go.mod
module github.com/chambridge/ship-shape

go 1.21
```

**Import Example**:
```go
import (
    "github.com/chambridge/ship-shape/pkg/core"
    "github.com/chambridge/ship-shape/pkg/analyzers"
)
```

**Alternative Considered**:
- ❌ `github.com/chambridge/shipshape` - Too similar to legacy Google tool
- ❌ `github.com/chambridge/ship-shape-analyzer` - Too verbose
- ✅ `github.com/chambridge/ship-shape` - **SELECTED**

---

## CLI Framework Selection

### Decision: **Cobra** + **Viper** (for configuration)

**Why Cobra**:
- ✅ Industry standard (used by kubectl, docker, gh, hugo)
- ✅ 38k+ GitHub stars, actively maintained
- ✅ Excellent documentation and examples
- ✅ Built-in help generation
- ✅ Subcommand support (analyze, report, gate, etc.)
- ✅ Flag parsing with persistent flags
- ✅ Shell completion (bash, zsh, fish, PowerShell)
- ✅ Easy testing with `cobra.Command.Execute()`

**Why Viper** (companion to Cobra):
- ✅ Configuration file support (YAML, TOML, JSON)
- ✅ Environment variable binding
- ✅ Flag binding
- ✅ Configuration precedence (flags > env > config file > defaults)
- ✅ Watch config files for changes

**Dependencies**:
```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.2
)
```

**Alternatives Considered**:
- ❌ `urfave/cli` - Less feature-rich, simpler but less powerful
- ❌ `flag` (stdlib) - Too basic for complex CLI
- ❌ `kong` - Good but less community adoption
- ✅ **Cobra + Viper** - **SELECTED**

---

## Configuration Format

### Decision: **YAML** (primary), with TOML/JSON support

**Why YAML**:
- ✅ Most common in Go CLI tools (Kubernetes, Docker Compose, GitHub Actions)
- ✅ Human-readable and writable
- ✅ Supports comments (critical for config files)
- ✅ Better for hierarchical data than TOML
- ✅ Viper supports all formats (YAML, TOML, JSON)

**Configuration File Names** (in order of precedence):
1. `.shipshape.yml` (primary)
2. `.shipshape.yaml`
3. `shipshape.yml`
4. `shipshape.yaml`
5. `.shipshape.toml` (alternate)
6. `.shipshape.json` (alternate)

**Configuration Search Paths**:
1. Current directory (`.`)
2. Repository root (find `.git` and look there)
3. Home directory (`~/.config/shipshape/`)
4. System config (`/etc/shipshape/`)

**Example**:
```yaml
# .shipshape.yml
version: 1

# Analysis configuration
analysis:
  parallel: true
  timeout: 5m
  languages:
    - go
    - python
    - javascript

# Scoring configuration
scoring:
  weights:
    coverage: 0.30
    quality: 0.25
    performance: 0.15
    tools: 0.15
    maintainability: 0.10
    code-quality: 0.05

# Quality gates
gates:
  fail-on: high
  min-score: 70

# Reporting
output:
  format: html
  path: shipshape-report.html
```

**Alternatives Considered**:
- ❌ TOML only - Less familiar to most developers
- ❌ JSON only - No comments, harder to read
- ✅ **YAML primary, TOML/JSON supported** - **SELECTED**

---

## CLI Interface Design

### Command Structure

```
shipshape
├── analyze      # Analyze repository (main command)
├── report       # Generate reports from previous analysis
├── gate         # Evaluate quality gates
├── compare      # Compare two analyses
├── init         # Initialize .shipshape.yml
├── validate     # Validate .shipshape.yml
└── version      # Show version
```

### Global Flags

Available for all commands:

```
--config, -c <file>    Configuration file path (default: auto-discover)
--verbose, -v          Verbose output
--quiet, -q            Quiet mode (errors only)
--no-color             Disable colored output
--help, -h             Show help
```

---

### Command: `shipshape analyze`

**Purpose**: Analyze a repository and generate quality assessment

**Usage**:
```bash
shipshape analyze [path] [flags]
```

**Arguments**:
- `[path]` - Repository path (default: current directory)

**Flags**:
```
--output, -o <file>           Output file path (default: shipshape-report.html)
--format, -f <format>         Output format: html, json, markdown, text (default: html)
--fail-on <severity>          Fail if findings of severity found: critical, high, medium, low, info (default: none)
--min-score <score>           Fail if overall score below threshold (0-100)
--parallel                    Enable parallel analysis (default: true)
--max-workers <n>             Maximum parallel workers (default: auto)
--timeout <duration>          Analysis timeout (default: 10m)
--languages <langs>           Comma-separated languages to analyze (default: auto-detect)
--exclude <patterns>          Exclude paths matching patterns (glob)
--cache                       Enable AST caching (default: true)
--no-upload                   Don't upload results to CI system
```

**Examples**:
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
```

**Exit Codes**:
- `0` - Success, no quality gate failures
- `1` - Quality gate failure (based on `--fail-on` or `--min-score`)
- `2` - Analysis error
- `3` - Configuration error

---

### Command: `shipshape report`

**Purpose**: Generate reports from previous analysis results

**Usage**:
```bash
shipshape report [flags]
```

**Flags**:
```
--input, -i <file>      Input analysis results (JSON format)
--output, -o <file>     Output report file
--format, -f <format>   Output format: html, markdown, text (default: html)
--template <file>       Custom template file (for HTML/Markdown)
--open                  Open report in browser after generation
```

**Examples**:
```bash
# Generate HTML report from JSON results
shipshape report -i results.json -o report.html

# Generate Markdown report and open
shipshape report -i results.json -f markdown --open

# Use custom template
shipshape report -i results.json --template my-template.html
```

---

### Command: `shipshape gate`

**Purpose**: Evaluate quality gates without full analysis (fast check)

**Usage**:
```bash
shipshape gate [flags]
```

**Flags**:
```
--input, -i <file>      Input analysis results (JSON format)
--fail-on <severity>    Fail if findings of severity found
--min-score <score>     Fail if overall score below threshold
--config, -c <file>     Configuration file with gate definitions
```

**Examples**:
```bash
# Check if previous analysis passes quality gate
shipshape gate -i results.json --min-score 80

# Use gates defined in config
shipshape gate -i results.json -c .shipshape.yml
```

**Exit Codes**:
- `0` - Quality gates passed
- `1` - Quality gates failed
- `2` - Error evaluating gates

---

### Command: `shipshape compare`

**Purpose**: Compare two analysis results (e.g., baseline vs current)

**Usage**:
```bash
shipshape compare <baseline> <current> [flags]
```

**Arguments**:
- `<baseline>` - Baseline analysis results (JSON)
- `<current>` - Current analysis results (JSON)

**Flags**:
```
--output, -o <file>     Output comparison report
--format, -f <format>   Output format: html, markdown, text (default: text)
--show-improvements     Show improvements, not just regressions
--fail-on-regression    Fail if any metric regressed
```

**Examples**:
```bash
# Compare baseline to current
shipshape compare baseline.json current.json

# Generate HTML comparison report
shipshape compare baseline.json current.json -f html -o comparison.html

# Fail CI if any regression
shipshape compare baseline.json current.json --fail-on-regression
```

---

### Command: `shipshape init`

**Purpose**: Initialize `.shipshape.yml` configuration file

**Usage**:
```bash
shipshape init [flags]
```

**Flags**:
```
--interactive, -i       Interactive configuration wizard
--template <name>       Use predefined template: minimal, standard, comprehensive
--overwrite             Overwrite existing .shipshape.yml
```

**Examples**:
```bash
# Interactive wizard
shipshape init -i

# Create minimal config
shipshape init --template minimal

# Create comprehensive config (all options)
shipshape init --template comprehensive
```

**Output**: Creates `.shipshape.yml` in current directory

---

### Command: `shipshape validate`

**Purpose**: Validate `.shipshape.yml` configuration file

**Usage**:
```bash
shipshape validate [file] [flags]
```

**Arguments**:
- `[file]` - Configuration file path (default: auto-discover)

**Flags**:
```
--strict        Enable strict validation (warn on unknown keys)
```

**Examples**:
```bash
# Validate default config
shipshape validate

# Validate specific config file
shipshape validate .shipshape.production.yml

# Strict validation
shipshape validate --strict
```

**Exit Codes**:
- `0` - Valid configuration
- `1` - Invalid configuration
- `2` - Configuration file not found

---

### Command: `shipshape version`

**Purpose**: Show version information

**Usage**:
```bash
shipshape version [flags]
```

**Flags**:
```
--short         Show version number only
--json          Output as JSON
```

**Examples**:
```bash
# Show full version info
shipshape version

# Show version number only
shipshape version --short

# JSON output
shipshape version --json
```

**Output Example**:
```
Ship Shape v1.0.0
Build Date: 2026-06-01
Go Version: go1.21.5
Platform: darwin/arm64
Commit: abc1234
```

---

## Environment Variables

Ship Shape respects the following environment variables:

```bash
# Configuration
SHIPSHAPE_CONFIG=/path/to/config.yml

# Analysis
SHIPSHAPE_PARALLEL=true
SHIPSHAPE_TIMEOUT=10m
SHIPSHAPE_MAX_WORKERS=8

# Output
SHIPSHAPE_OUTPUT_FORMAT=json
SHIPSHAPE_OUTPUT_PATH=./results.json
SHIPSHAPE_NO_COLOR=true

# Quality Gates
SHIPSHAPE_FAIL_ON=high
SHIPSHAPE_MIN_SCORE=80

# CI/CD Integration
CI=true                    # Detected automatically
GITHUB_ACTIONS=true        # Detected automatically
SHIPSHAPE_NO_UPLOAD=false
```

**Precedence** (highest to lowest):
1. CLI flags (e.g., `--fail-on high`)
2. Environment variables (e.g., `SHIPSHAPE_FAIL_ON=high`)
3. Configuration file (`.shipshape.yml`)
4. Built-in defaults

---

## Shell Completion

Ship Shape supports shell completion for all major shells:

```bash
# Generate completion script
shipshape completion bash > /etc/bash_completion.d/shipshape
shipshape completion zsh > ~/.zsh/completion/_shipshape
shipshape completion fish > ~/.config/fish/completions/shipshape.fish
shipshape completion powershell > shipshape.ps1
```

**Installation**:
```bash
# Bash
echo "source <(shipshape completion bash)" >> ~/.bashrc

# Zsh
echo "source <(shipshape completion zsh)" >> ~/.zshrc

# Fish
shipshape completion fish | source
```

---

## CI/CD Integration Examples

### GitHub Actions

```yaml
# .github/workflows/shipshape.yml
name: Ship Shape Analysis

on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Ship Shape
        run: |
          curl -L shipshape.dev/install.sh | sh
          shipshape analyze --fail-on high --min-score 80 -f json -o results.json

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: shipshape-report
          path: results.json
```

### Pre-commit Hook

```bash
# .git/hooks/pre-commit
#!/bin/bash
shipshape analyze --fail-on critical --timeout 30s --quiet
exit $?
```

---

## Summary

### Final Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Module Path** | `github.com/chambridge/ship-shape` | Standard Go convention |
| **CLI Framework** | Cobra + Viper | Industry standard, feature-rich |
| **Config Format** | YAML (primary) | Human-readable, supports comments |
| **Binary Name** | `shipshape` | No hyphen, easy to type |
| **Main Command** | `analyze` | Primary use case |
| **Config File** | `.shipshape.yml` | Hidden file, YAML extension |

### Dependencies Required

```go
require (
    github.com/spf13/cobra v1.8.0       // CLI framework
    github.com/spf13/viper v1.18.2      // Configuration management
    gopkg.in/yaml.v3 v3.0.1             // YAML parsing
)
```

---

**Document Status**: Complete - Ready for Implementation
**Next Steps**: Create `.shipshape.example.yml` and initialize Go module
**Owner**: Architecture Team

---

**Last Updated**: 2026-01-27
**Version**: 1.0.0
**Author**: Senior Software Engineer
