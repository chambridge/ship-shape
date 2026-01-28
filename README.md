# Ship Shape

**A comprehensive testing quality and code analysis platform for multi-language repositories**

[![CI Pipeline](https://github.com/chambridge/ship-shape/actions/workflows/ci.yml/badge.svg)](https://github.com/chambridge/ship-shape/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chambridge/ship-shape)](https://goreportcard.com/report/github.com/chambridge/ship-shape)
[![codecov](https://codecov.io/gh/chambridge/ship-shape/branch/main/graph/badge.svg)](https://codecov.io/gh/chambridge/ship-shape)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/chambridge/ship-shape)](https://github.com/chambridge/ship-shape/releases)
[![Status](https://img.shields.io/badge/status-in%20development-yellow.svg)](https://github.com/chambridge/ship-shape)

---

## Overview

Ship Shape is an intelligent testing quality analysis tool that provides deep, actionable insights into test health, coverage effectiveness, and testing best practices across multi-language repositories and monorepos.

Unlike traditional coverage tools that only measure line coverage, Ship Shape:
- 🔍 **Analyzes test quality** - Detects test smells, anti-patterns, and best practice violations
- 🎯 **Context-aware** - Automatically detects languages, frameworks, and infrastructure
- 🏢 **Monorepo-native** - First-class support for monorepo structures with aggregate scoring
- 📊 **Multi-dimensional scoring** - Evaluates coverage, quality, performance, tools, and maintainability
- 🔬 **Evidence-based** - All metrics backed by academic research and industry standards
- 🚀 **CI/CD integrated** - Native GitHub Actions support with quality gates
- 📈 **Historical tracking** - Trend analysis and regression detection over time

---

## Key Features

### Repository-Aware Analysis
- Automatic language detection (Python, Go, JavaScript/TypeScript, Java, Rust, Ruby, C#, C/C++)
- Framework and tool detection from configuration files
- Intelligent recommendations based on detected infrastructure
- Only suggests tools applicable to your tech stack

### Comprehensive Test Quality Analysis
- **Test Smell Detection**: Mystery Guest, Eager Test, Lazy Test, Assertion Roulette, and more
- **Coverage Analysis**: Line, branch, and mutation coverage support
- **Best Practices**: Framework-specific patterns (table-driven tests, fixtures, mocking strategies)
- **AST-based Analysis**: Deep code understanding through syntax tree parsing

### Monorepo Support
- Workspace detection (npm, yarn, pnpm, Go workspaces, Maven multi-module)
- Independent package analysis with aggregate scoring
- Shared infrastructure quality assessment
- Cross-package consistency analysis

### CI/CD Integration
- **GitHub Actions**: Native action with PR comments, status checks, and artifact uploads
- **Quality Gates**: Blocking, warning, and trend-based gates
- **Pre-commit Hooks**: Fast incremental analysis for local development

### Multi-Dimensional Scoring
1. **Coverage** (25%) - Line, branch, and mutation coverage
2. **Test Quality** (25%) - Test organization, patterns, and smells
3. **Performance** (15%) - Test execution speed and flakiness
4. **Tool Adoption** (15%) - Testing infrastructure quality
5. **Maintainability** (10%) - Test-to-code ratio and duplication
6. **Code Quality** (10%) - Static analysis integration

### Historical Tracking
- SQLite-based storage for trend analysis
- Score trends over time with regression detection
- Coverage evolution tracking
- Finding trend analysis

---

## Architecture

Ship Shape is built with Go for performance, type safety, and cross-platform support:

```
┌─────────────────────────────────────────────────────────────┐
│                      Ship Shape CLI                          │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  Repository  │   │  Execution   │   │   Report     │
│  Discovery   │   │   Engine     │   │  Generator   │
└──────────────┘   └──────────────┘   └──────────────┘
        │                   │                   │
        └───────────────────┼───────────────────┘
                            │
                ┌───────────┴───────────┐
                │                       │
        ┌───────▼──────┐       ┌───────▼──────┐
        │   Analyzer   │       │   Assessor   │
        │   Registry   │       │   Registry   │
        └──────────────┘       └──────────────┘
```

### 7-Layer Analysis Model
1. **Discovery** - Language, framework, and structure detection
2. **Structure** - Directory mapping and LOC analysis
3. **Quality** - Test smell detection and best practices
4. **Coverage** - Coverage report parsing and analysis
5. **Execution** - Test performance and flakiness (optional)
6. **Mutation** - Mutation testing analysis (optional)
7. **CI/CD** - CI pipeline analysis and optimization

---

## Project Documentation

This repository contains comprehensive planning and design documents:

### Planning Documents
- **[.project/requirements.md](.project/requirements.md)** - Detailed requirements specification (2,300+ lines)
- **[.project/architecture.md](.project/architecture.md)** - Technical architecture with Go code examples (3,400+ lines)
- **[.project/user-stories.md](.project/user-stories.md)** - User stories with acceptance criteria (27 stories, 180 story points)

### Research
- **[.research/shipshape.md](.research/shipshape.md)** - Research and design analysis (98KB)

### Project Status
- **Current Version**: v0.0.0 (Planning Phase)
- **Target v0.1.0**: Q1 2026 - Repository Discovery
- **Target v1.0.0**: Q4 2026 - Full Feature Set

---

## Technology Stack

### Core
- **Language**: Go 1.21+
- **Storage**: SQLite (historical data)
- **Templating**: templ (type-safe HTML)

### Language Analysis
- **Go**: go/ast (standard library)
- **Python**: tree-sitter or embedded Python AST parser
- **JavaScript/TypeScript**: tree-sitter
- **Language Detection**: go-enry

### Dependencies (Minimal, Open Source)
- `go-enry/go-enry` - Language detection (Apache 2.0)
- `smacker/go-tree-sitter` - Universal parser (MIT)
- `gopkg.in/yaml.v3` - YAML parsing (Apache 2.0)
- `pelletier/go-toml` - TOML parsing (MIT)
- `mattn/go-sqlite3` - SQLite driver (MIT)
- `a-h/templ` - HTML templating (MIT)

---

## Getting Started

> **Note**: Ship Shape is currently in the planning and design phase. Implementation is scheduled to begin Q1 2026.

### Planned Installation (Future)

```bash
# Install via go install
go install github.com/shipshape/cmd/shipshape@latest

# Or download pre-built binary
curl -L https://github.com/shipshape/releases/latest/download/shipshape-$(uname -s)-$(uname -m) -o shipshape
chmod +x shipshape
sudo mv shipshape /usr/local/bin/
```

### Planned Usage (Future)

```bash
# Analyze current repository
shipshape analyze

# Analyze with configuration
shipshape analyze --config .shipshape.yml

# Run with quality gates
shipshape analyze --gates

# Generate HTML report
shipshape analyze --output html

# View trends
shipshape trends --days 30

# Organization analysis
shipshape org analyze --config ~/.shipshape/orgs.yml
```

### GitHub Actions Integration (Planned)

```yaml
name: Test Quality Check
on: [pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Ship Shape Analysis
        uses: shipshape/action@v1
        with:
          gates: true
          fail-on-gate: true
          comment-pr: true
          upload-artifacts: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Development Roadmap

### Phase 0: Foundation & Setup
- [x] Requirements specification
- [x] Technical architecture design
- [x] User stories and acceptance criteria
- [x] Go module initialization
- [x] CLI framework (Cobra)
- [x] Configuration system (Viper)
- [x] CI/CD pipeline (GitHub Actions)
- [x] Development documentation
- [ ] Logging framework
- [ ] Test infrastructure
- [ ] Repository discovery engine
- [ ] Ground truth datasets

### Phase 1: Technical Spike Validation
- [ ] AST parsing framework
- [ ] Monorepo analysis coordination
- [ ] Test smell detection
- [ ] Coverage parsing
- [ ] GitHub Actions integration
- [ ] HTML report generation

### Phase 2: Multi-Language Support
- [ ] Go test analysis with AST parsing
- [ ] Python test analysis (pytest/unittest)
- [ ] JavaScript/TypeScript analysis (Jest/Vitest)
- [ ] Monorepo detection and handling
- [ ] Package-level analysis

### Phase 3: Quality Analysis
- [ ] Test smell detection framework
- [ ] Coverage report parsing (all languages)
- [ ] Coverage quality assessment
- [ ] Tool database and detection

### Phase 4: Integration & Reporting
- [ ] GitHub Actions integration
- [ ] Quality gates system
- [ ] HTML report generation
- [ ] Pre-commit hooks

### Phase 5: Advanced Features
- [ ] Historical tracking and trends
- [ ] SQLite storage implementation
- [ ] Organization-wide analysis
- [ ] REST API (optional)

---

## Contributing

Contributions are welcome! Ship Shape is currently in active development.

### Development Setup

```bash
# Clone repository
git clone https://github.com/chambridge/ship-shape.git
cd ship-shape

# Install dependencies
go mod download

# Build
make build

# Run locally
./bin/shipshape version

# Run all quality checks
make check

# See all available commands
make help
```

### Development Standards
- **Test-Driven Development**: >90% unit test coverage required
- **Conventional Commits**: All commits follow conventional commit format
- **DCO Sign-off**: All commits must be signed off (`git commit -s`)
- **Code Quality**: Must pass `make check` (runs fmt, vet, lint, test)
- **Workflow Validation**: Run `make actionlint` before pushing workflow changes

See [docs/v1.0.0/cicd.md](docs/v1.0.0/cicd.md) for detailed CI/CD documentation.

---

## Research and References

Ship Shape's design is informed by extensive research in software testing quality:

- **Test Smells**: Based on research from Bavota et al., Fowler, and Meszaros
- **Coverage Metrics**: Google Testing Blog, Microsoft Research studies
- **Mutation Testing**: IEEE studies on mutation testing effectiveness
- **Best Practices**: Industry standards from Google, Microsoft, ThoughtWorks

See [.research/shipshape.md](.research/shipshape.md) for complete research analysis.

---

## License

> License to be determined

---

## Support

- **Issues**: [GitHub Issues](https://github.com/chambridge/ship-shape/issues)
- **Documentation**: [docs/v1.0.0/](docs/v1.0.0/)
- **Planning**: [.project/](.project/) - Requirements, architecture, user stories, roadmap

---

## Project Status

**Current Phase**: Phase 0 - Foundation & Setup 🚧

**Completed**:
- ✅ Requirements specification (2,300+ lines)
- ✅ Technical architecture (3,400+ lines with Go examples)
- ✅ User stories with acceptance criteria (27 stories, 180 story points)
- ✅ Research analysis (98KB)
- ✅ 24-week implementation roadmap
- ✅ Go module and project structure
- ✅ CLI framework with Cobra
- ✅ Configuration system with Viper
- ✅ Comprehensive CI/CD pipeline
- ✅ Quality tooling (golangci-lint, actionlint, gosec)
- ✅ Development documentation (CLAUDE.md, CI/CD docs)

**In Progress**:
- 🚧 Logging framework
- 🚧 Test infrastructure
- 🚧 Ground truth dataset structure

**Next Steps**:
1. Complete Phase 0 foundation tasks
2. Repository discovery engine implementation
3. Language detection framework
4. Begin technical spike validation

---

**Built with ⚓ by the Ship Shape team**

*Making testing quality measurable, actionable, and maintainable.*
