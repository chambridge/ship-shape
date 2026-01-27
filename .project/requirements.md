# Ship Shape - Detailed Requirements Specification

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Draft
**Author**: Senior Software Engineer

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Core Principles](#core-principles)
3. [Architecture Overview](#architecture-overview)
4. [Language & Monorepo Detection System](#language--monorepo-detection-system)
5. [Data Collection & Analysis Engine](#data-collection--analysis-engine)
6. [Scoring & Assessment System](#scoring--assessment-system)
7. [Dashboard & Visualization](#dashboard--visualization)
8. [Tool Ecosystem Analysis](#tool-ecosystem-analysis)
9. [Technology Stack](#technology-stack)
10. [Integration Points](#integration-points)
11. [Quality Gates & CI/CD](#quality-gates--cicd)
12. [Multi-Project Organization Features](#multi-project-organization-features)
13. [Non-Functional Requirements](#non-functional-requirements)

---

## Executive Summary

Ship Shape is a comprehensive testing quality and code analysis platform designed to provide deep, actionable insights into test health, coverage effectiveness, and testing best practices across multi-language repositories and monorepos. The tool automatically detects repository structure, programming languages, and existing infrastructure to provide context-aware, language-specific recommendations.

**Critical Design Principles:**
- **Repository-Aware**: Automatically detects languages, frameworks, and infrastructure present in the repository
- **Context-Specific**: Only recommends tools and practices applicable to detected languages/frameworks
- **Monorepo-Native**: First-class support for monorepo structures with multi-language projects
- **Evidence-Based**: All metrics and recommendations backed by academic research and industry standards
- **Actionable**: Provides specific, implementable recommendations with setup guides

---

## Core Principles

### 1. Repository Context Awareness

**REQ-CP-001**: The tool MUST analyze the repository structure before making any recommendations.

**REQ-CP-002**: The tool MUST detect all programming languages present in the repository through:
- File extension analysis
- Configuration file detection
- Dependency manifest parsing
- Build system identification

**REQ-CP-003**: The tool MUST only recommend tools, frameworks, and practices applicable to languages/technologies actually present in the repository.

**REQ-CP-004**: The tool MUST NOT penalize repositories for lacking tools/practices specific to languages they don't use.

### 2. Monorepo Support

**REQ-CP-005**: The tool MUST detect and support common monorepo structures:
- Workspace-based monorepos (npm workspaces, yarn workspaces, pnpm workspaces)
- Go modules with multi-module workspaces
- Python monorepos with multiple packages
- Gradle/Maven multi-module projects
- Lerna/Nx/Turborepo managed monorepos
- Custom monorepo structures with multiple projects

**REQ-CP-006**: The tool MUST provide both aggregate scores for the entire monorepo AND individual scores for each sub-project/package.

**REQ-CP-007**: The tool MUST identify shared testing infrastructure and evaluate its quality across the monorepo.

### 3. Multi-Language Intelligence

**REQ-CP-008**: The tool MUST maintain language-specific best practices databases for:
- Python (pytest, unittest, nose2, doctest)
- JavaScript/TypeScript (Jest, Mocha, Vitest, Cypress, Playwright)
- Go (testing, testify, ginkgo)
- Java (JUnit 5, TestNG, Mockito)
- Rust (cargo test, proptest)
- Ruby (RSpec, Minitest)
- C# (NUnit, xUnit.net, MSTest)
- C/C++ (Google Test, Catch2)

**REQ-CP-009**: Analysis depth MUST be proportional to language presence (primary vs. secondary languages).

### 4. Evidence-Based Analysis

**REQ-CP-010**: All metrics MUST reference academic research or industry standards with citations.

**REQ-CP-011**: Recommendations MUST include:
- Justification based on evidence
- Expected benefit (quantified where possible)
- Implementation effort estimate
- Priority level (Critical, High, Medium, Low)

---

## Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Ship Shape CLI/API                       │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  Repository  │   │  Execution   │   │   Report     │
│  Discovery   │   │   Engine     │   │  Generator   │
│   Engine     │   │              │   │              │
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
                │                       │
        ┌───────┴───────┬───────────────┴───────┬────────────┐
        │               │                       │            │
        ▼               ▼                       ▼            ▼
┌──────────────┐ ┌─────────────┐ ┌──────────────┐ ┌─────────────┐
│   Language   │ │  Framework  │ │   Coverage   │ │ Tool Adoption│
│   Detectors  │ │  Analyzers  │ │  Analyzers   │ │  Analyzers  │
└──────────────┘ └─────────────┘ └──────────────┘ └─────────────┘
```

### Component Responsibilities

**REQ-ARCH-001**: Repository Discovery Engine MUST:
- Detect repository type (single-project, monorepo, multi-language)
- Identify all programming languages with percentage distribution
- Locate all project roots and module boundaries
- Parse workspace configurations
- Identify build systems and dependency managers
- Map project structure (source, test, config directories)

**REQ-ARCH-002**: Execution Engine MUST:
- Orchestrate analyzer execution with dependency management
- Support parallel execution where dependencies allow
- Maintain execution context (language, framework, project scope)
- Handle partial failures gracefully
- Provide progress reporting

**REQ-ARCH-003**: Analyzer Registry MUST:
- Register language-specific analyzers dynamically
- Match analyzers to detected languages/frameworks
- Support analyzer versioning and updates
- Allow custom analyzer plugins

**REQ-ARCH-004**: Assessor Registry MUST:
- Register scoring assessors by category
- Apply only relevant assessors based on repository context
- Support weighted scoring configurations
- Enable custom assessment rules

---

## Language & Monorepo Detection System

### Language Detection

**REQ-LANG-001**: Language detection MUST analyze:
- File extensions with count and percentage distribution
- Primary vs. secondary languages (>10% presence = primary)
- Configuration files (package.json, requirements.txt, go.mod, pom.xml, Cargo.toml, etc.)
- Build files (Makefile, CMakeLists.txt, build.gradle, etc.)
- Language-specific directories (.venv, node_modules, target, etc.)

**REQ-LANG-002**: Detection output MUST include:
```json
{
  "languages": [
    {
      "name": "Go",
      "percentage": 65.3,
      "file_count": 247,
      "primary": true,
      "test_frameworks_detected": ["testing", "testify"],
      "coverage_tools_detected": ["gocov"],
      "build_system": "go modules"
    },
    {
      "name": "Python",
      "percentage": 28.2,
      "file_count": 89,
      "primary": true,
      "test_frameworks_detected": ["pytest"],
      "coverage_tools_detected": ["coverage.py"],
      "build_system": "pip/setuptools"
    },
    {
      "name": "JavaScript",
      "percentage": 6.5,
      "file_count": 21,
      "primary": false,
      "test_frameworks_detected": [],
      "coverage_tools_detected": [],
      "build_system": "npm"
    }
  ],
  "repository_type": "multi-language",
  "monorepo_detected": false
}
```

### Monorepo Detection

**REQ-MONO-001**: Monorepo detection MUST identify:
- Workspace configuration files (lerna.json, nx.json, pnpm-workspace.yaml, etc.)
- Multiple package.json files with workspace references
- Go workspace files (go.work)
- Multiple Maven/Gradle modules
- Multiple Python packages with separate setup.py/pyproject.toml
- Custom monorepo structures with multiple root-level projects

**REQ-MONO-002**: Monorepo analysis MUST produce:
```json
{
  "monorepo": true,
  "monorepo_type": "npm-workspaces",
  "packages": [
    {
      "name": "@myorg/api",
      "path": "packages/api",
      "language": "TypeScript",
      "test_framework": "Jest",
      "dependencies": ["@myorg/shared"]
    },
    {
      "name": "@myorg/web",
      "path": "packages/web",
      "language": "TypeScript",
      "test_framework": "Vitest",
      "dependencies": ["@myorg/shared"]
    },
    {
      "name": "@myorg/shared",
      "path": "packages/shared",
      "language": "TypeScript",
      "test_framework": "Jest",
      "dependencies": []
    }
  ],
  "shared_infrastructure": {
    "root_jest_config": true,
    "shared_tsconfig": true,
    "common_ci_workflow": true
  }
}
```

**REQ-MONO-003**: Scoring for monorepos MUST provide:
- Overall monorepo score (aggregate)
- Individual package scores
- Shared infrastructure quality score
- Cross-package consistency metrics
- Dependency health between packages

### Framework Detection

**REQ-FRAME-001**: For each detected language, the tool MUST detect:
- Test frameworks in use (via imports, config files, dependencies)
- Coverage tools configured
- Linting/formatting tools
- Build systems and task runners
- CI/CD integration patterns

**REQ-FRAME-002**: Framework detection MUST be version-aware:
- Detect framework versions from lockfiles
- Flag outdated or deprecated versions
- Recommend version upgrades with breaking change warnings

---

## Data Collection & Analysis Engine

### Analysis Layers

The analysis engine MUST execute in layered phases with clear dependencies:

#### Layer 1: Repository Discovery (<1s)

**REQ-ANLZ-001**: Discovery phase MUST complete:
- Language detection with distribution analysis
- Monorepo structure identification
- Test file location mapping
- Configuration file discovery
- Dependency manifest parsing

**Output**: Repository context object used by all subsequent layers

#### Layer 2: Test Structure Analysis (5-30s)

**REQ-ANLZ-002**: Structure analysis MUST:
- Categorize tests by type (unit, integration, e2e, performance, etc.)
- Calculate test pyramid/trophy distribution
- Identify test file organization patterns
- Detect test suite structure (nested suites, test groups)
- Map tests to source code (test-to-code traceability)

**Language-Specific Patterns**:
- **Python**: Detect pytest fixtures, unittest classes, doctest usage
- **JavaScript**: Identify describe/it blocks, test file naming conventions
- **Go**: Recognize table-driven tests, TestMain usage, subtests
- **Java**: Find JUnit test classes, test lifecycle methods

#### Layer 3: Test Quality Analysis (30s-5m)

**REQ-ANLZ-003**: Quality analysis MUST detect:
- Test smells (11 types) with language-specific patterns
- Assertion quality and specificity
- Test naming conventions compliance
- Test size appropriateness
- Test independence violations
- Mocking pattern quality

**Test Smell Detection by Language**:

| Test Smell | Python | JavaScript | Go | Java |
|------------|--------|------------|-----|------|
| Mystery Guest | ✓ | ✓ | ✓ | ✓ |
| Eager Test | ✓ | ✓ | ✓ | ✓ |
| Lazy Test | ✓ | ✓ | ✓ | ✓ |
| Obscure Test | ✓ | ✓ | ✓ | ✓ |
| Conditional Logic | ✓ | ✓ | ✓ | ✓ |
| General Fixture | ✓ (pytest) | ✓ (beforeEach) | ✓ (setup funcs) | ✓ (@Before) |
| Code Duplication | ✓ | ✓ | ✓ | ✓ |
| Assertion Roulette | ✓ | ✓ | ✓ | ✓ |
| Sensitive Equality | ✓ | ✓ | ✓ | ✓ |
| Resource Optimism | ✓ | ✓ | ✓ | ✓ |
| Flakiness | ✓ | ✓ | ✓ | ✓ |

#### Layer 4: Coverage Analysis (1-10m)

**REQ-ANLZ-004**: Coverage analysis MUST:
- Integrate with language-specific coverage tools
- Parse coverage reports (XML, JSON, LCOV formats)
- Calculate line, branch, function coverage
- Identify critical paths and verify coverage
- Generate uncovered code reports
- Support incremental coverage analysis

**Coverage Tool Integration**:
- **Python**: coverage.py, pytest-cov
- **JavaScript/TypeScript**: Istanbul/nyc, c8
- **Go**: go test -cover, gocov
- **Java**: JaCoCo, Cobertura
- **Rust**: cargo-tarpaulin, cargo-llvm-cov
- **C#**: coverlet, dotCover

**REQ-ANLZ-005**: Coverage thresholds MUST be language-appropriate:
```yaml
coverage_thresholds:
  line_coverage:
    critical: 90%  # High-reliability systems
    standard: 80%  # Most applications
    baseline: 70%  # Minimum acceptable
  branch_coverage:
    critical: 85%
    standard: 70%
    baseline: 60%
  function_coverage:
    critical: 95%
    standard: 85%
    baseline: 75%
```

#### Layer 5: Test Execution Analysis (variable timing)

**REQ-ANLZ-006**: Execution analysis MUST (when enabled):
- Run test suites with timing instrumentation
- Detect flaky tests through multiple runs
- Profile test performance
- Identify slow tests (language-specific thresholds)
- Analyze parallelization opportunities
- Measure test isolation

**Performance Thresholds**:
- Unit tests: <100ms (Python/JS), <50ms (Go), <200ms (Java)
- Integration tests: <1s (most languages)
- E2E tests: <30s (browser-based), <5s (API)

**REQ-ANLZ-007**: Flakiness detection MUST:
- Run tests multiple times (configurable: 3-10 runs)
- Track pass/fail patterns
- Identify environmental dependencies
- Detect timing-related issues
- Classify flakiness severity (occasional vs. frequent)

#### Layer 6: Mutation Testing (10m-2h, optional)

**REQ-ANLZ-008**: Mutation testing MUST:
- Integrate with language-specific mutation tools
- Calculate mutation score
- Identify surviving mutants
- Focus on critical code paths
- Provide killed/survived/timeout/error counts

**Mutation Tool Integration**:
- **JavaScript/TypeScript**: Stryker
- **Java**: PITest
- **Python**: mutmut, cosmic-ray
- **C#**: Stryker.NET

**REQ-ANLZ-009**: Mutation analysis SHOULD be configurable for scope:
- Critical code only (fast feedback)
- Changed code only (CI/CD mode)
- Full codebase (comprehensive analysis)

#### Layer 7: CI/CD Integration Analysis (<30s)

**REQ-ANLZ-010**: CI/CD analysis MUST examine:
- Test workflow configuration
- Parallelization strategy
- Caching mechanisms
- Test result reporting
- Quality gate enforcement
- Dependency caching
- Test failure handling

**Supported CI/CD Platforms**:
- GitHub Actions (.github/workflows/)
- GitLab CI (.gitlab-ci.yml)
- Jenkins (Jenkinsfile)
- CircleCI (.circleci/config.yml)
- Travis CI (.travis.yml)
- Buildkite (.buildkite/)

### Static Analysis Techniques

**REQ-STATIC-001**: AST (Abstract Syntax Tree) parsing MUST support:
- **Python**: ast module, parso
- **JavaScript/TypeScript**: Babel parser, TypeScript compiler API
- **Go**: go/ast, go/parser
- **Java**: JavaParser, Eclipse JDT
- **Rust**: syn crate

**REQ-STATIC-002**: Pattern matching MUST detect:
- Test method signatures
- Assertion patterns
- Setup/teardown methods
- Mock/stub usage
- Fixture definitions
- Test decorators/annotations

### Dynamic Analysis Techniques

**REQ-DYNAMIC-001**: When executing tests, the tool MUST:
- Respect existing test configurations
- Use project's dependency versions
- Execute in isolated environments
- Capture stdout/stderr
- Measure resource usage
- Track exit codes

**REQ-DYNAMIC-002**: Execution isolation MUST:
- Use virtual environments (Python: venv/virtualenv)
- Use containerization when configured
- Restore original state after execution
- Handle cleanup of temporary resources

---

## Scoring & Assessment System

### Multi-Dimensional Scoring

**REQ-SCORE-001**: Overall quality score MUST be calculated from weighted dimensions:

```yaml
dimensions:
  test_coverage:
    weight: 30%
    components:
      - line_coverage (40%)
      - branch_coverage (30%)
      - function_coverage (15%)
      - critical_path_coverage (15%)

  test_quality:
    weight: 25%
    components:
      - test_smell_count (30%)
      - test_independence (25%)
      - assertion_quality (20%)
      - naming_quality (15%)
      - test_organization (10%)

  tool_adoption:
    weight: 10%
    components:
      - essential_tools (50%)
      - recommended_tools (30%)
      - configuration_maturity (20%)

  test_performance:
    weight: 10%
    components:
      - execution_time (40%)
      - flakiness (40%)
      - parallelization (20%)

  maintainability:
    weight: 15%
    components:
      - test_organization (35%)
      - code_duplication (25%)
      - complexity (25%)
      - documentation (15%)

  code_quality:
    weight: 10%
    components:
      - linting_compliance (40%)
      - formatting_consistency (25%)
      - type_safety (20%)
      - best_practices (15%)
```

**REQ-SCORE-002**: Dimension weights MUST be adjustable per language/project type.

**REQ-SCORE-003**: For monorepos, scoring MUST provide:
- Aggregate score (weighted average of all packages)
- Per-package scores
- Shared infrastructure score
- Cross-package consistency score

### Language-Specific Scoring Adjustments

**REQ-SCORE-004**: Scoring MUST adapt to language characteristics:

**Go-specific adjustments**:
- Emphasize table-driven test patterns (+5% quality if >70% usage)
- Reward t.Parallel() usage (+3% performance)
- Expect higher coverage thresholds (90% line coverage standard)

**Python-specific adjustments**:
- Reward pytest fixture usage (+5% quality if well-organized)
- Emphasize parametrization (+3% quality)
- Consider doctest as supplementary testing

**JavaScript/TypeScript-specific adjustments**:
- Reward proper TypeScript typing in tests (+5% quality)
- Consider snapshot testing appropriately (warning if >20% of tests)
- Emphasize async test handling quality

**Java-specific adjustments**:
- Reward proper test lifecycle usage
- Emphasize assertion library quality (AssertJ > plain JUnit)
- Consider integration test separation (Maven/Gradle phases)

### Evidence-Based Thresholds

**REQ-SCORE-005**: All thresholds MUST reference supporting research:

```yaml
coverage_thresholds:
  line_coverage_excellent: 90%
    # Reference: "Minimum Acceptable Code Coverage" (Antinyan et al., 2018)
    # DOI: 10.1145/3194164.3194175

  branch_coverage_good: 70%
    # Reference: "An Empirical Study of Branch Coverage" (Gopinath et al., 2014)
    # DOI: 10.1145/2635868.2635879

  mutation_score_excellent: 80%
    # Reference: "Mutation Testing: An Empirical Evaluation" (Andrews et al., 2005)
    # DOI: 10.1109/TSE.2005.93
```

### Grade Assignment

**REQ-SCORE-006**: Letter grades MUST map to scores:

| Grade | Score Range | Interpretation |
|-------|-------------|----------------|
| A+ | 95-100 | Exceptional testing quality |
| A | 90-94 | Excellent testing quality |
| B | 80-89 | Good testing quality |
| C | 70-79 | Acceptable testing quality |
| D | 60-69 | Below standard, needs improvement |
| F | 0-59 | Failing, significant gaps |

**REQ-SCORE-007**: Grade thresholds MAY be configured per organization.

---

## Dashboard & Visualization

### Report Formats

**REQ-DASH-001**: The tool MUST generate reports in:
- HTML (interactive, standalone)
- JSON (machine-readable)
- Plain text (CLI-friendly)
- Markdown (documentation-friendly)
- CSV (trend analysis)

### HTML Dashboard Components

**REQ-DASH-002**: HTML report MUST include:

1. **Executive Summary Section**
   - Overall quality score with letter grade
   - Repository type and language distribution
   - Key findings summary (top 5 issues)
   - Trending indicators (if historical data available)

2. **Language Distribution View**
   - Pie chart of language percentages
   - Table with file counts and frameworks detected
   - Test framework adoption per language

3. **Monorepo Overview** (if applicable)
   - Package/module listing with individual scores
   - Dependency graph visualization
   - Shared infrastructure summary

4. **Score Breakdown Dashboard**
   - Radar chart showing all dimensions
   - Individual dimension scores with progress bars
   - Comparison to thresholds (if configured)

5. **Test Coverage Dashboard**
   - Coverage percentage with visual indicators
   - Line, branch, function coverage metrics
   - Coverage by module/package
   - Uncovered critical code highlighting
   - Coverage trend chart (if historical data)
   - Interactive coverage heatmap (file-level)

6. **Test Quality Dashboard**
   - Test smell counts by type
   - Test pyramid/trophy visualization
   - Test independence metrics
   - Flaky test list with details
   - Test size distribution chart

7. **Test Performance Dashboard**
   - Execution time statistics
   - Slowest tests table
   - Flakiness detection results
   - Parallelization opportunity analysis
   - Performance trend chart

8. **Tool Adoption Dashboard**
   - Category-wise tool adoption (Testing, Quality, Security, Automation)
   - Detected tools with versions
   - Missing essential tools (prioritized)
   - Configuration maturity indicators
   - Tool recommendation cards with:
     * Tool name and description
     * Benefits (quantified)
     * Effort estimate
     * Setup guide link
     * Priority level

9. **Code Quality Dashboard**
   - Linting/formatting compliance
   - Type safety metrics (TypeScript, typed Python, etc.)
   - Best practice adherence
   - Technical debt indicators

10. **CI/CD Integration Dashboard**
    - Workflow configuration analysis
    - Parallelization assessment
    - Caching strategy evaluation
    - Test result reporting quality
    - Optimization recommendations

11. **Findings & Recommendations**
    - Filterable table of all findings
    - Severity-based grouping
    - Evidence and references
    - Remediation guidance
    - Auto-fix availability indicators

**REQ-DASH-003**: All charts and visualizations MUST:
- Be responsive (mobile-friendly)
- Support dark/light themes
- Include accessible color schemes
- Provide data export (CSV/JSON)
- Include interactive tooltips

### CLI Output

**REQ-DASH-004**: CLI output MUST provide:

```
╔══════════════════════════════════════════════════════════╗
║           Ship Shape Analysis Results                    ║
╚══════════════════════════════════════════════════════════╝

Repository: /path/to/repo
Type: Monorepo (npm workspaces)
Languages: Go (65%), Python (28%), JavaScript (7%)

Overall Score: 82/100 (Grade: B)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Dimension Breakdown:
  ✓ Test Coverage        88/100 (30%) ████████░░
  ✓ Test Quality         78/100 (25%) ███████░░░
  ⚠ Tool Adoption        65/100 (10%) ██████░░░░
  ✓ Test Performance     85/100 (10%) ████████░░
  ✓ Maintainability      80/100 (15%) ████████░░
  ✓ Code Quality         75/100 (10%) ███████░░░

Language-Specific Results:
┌────────────┬────────┬──────────────┬───────┐
│ Language   │ Files  │ Test Framework│ Score │
├────────────┼────────┼──────────────┼───────┤
│ Go         │ 247    │ testing      │ 85/100│
│ Python     │ 89     │ pytest       │ 80/100│
│ JavaScript │ 21     │ (none)       │ 60/100│
└────────────┴────────┴──────────────┴───────┘

Top Issues:
  🔴 CRITICAL: JavaScript package has no test framework (packages/scripts)
  🟡 HIGH: 12 flaky tests detected in Go codebase
  🟡 HIGH: Branch coverage below threshold in 3 Python modules
  🟡 MEDIUM: 8 test smells detected (Eager Test pattern)
  🟡 MEDIUM: Missing essential tool: golangci-lint

Recommendations: Run 'shipshape report' for detailed HTML report
                 Run 'shipshape tools --setup golangci-lint' for setup guide
```

**REQ-DASH-005**: CLI MUST support:
- Color output (with --no-color flag)
- Compact mode (summary only)
- Verbose mode (detailed findings)
- JSON output mode (for scripting)
- Progress indicators during analysis

### Interactive Features

**REQ-DASH-006**: HTML dashboard MUST support:
- Filtering findings by severity/type/language
- Searching findings by keyword
- Sorting tables by any column
- Expanding/collapsing detail sections
- Drill-down from summary to details
- Direct links to source files (file:// URIs)
- Copy-to-clipboard for code examples

---

## Tool Ecosystem Analysis

### Tool Detection

**REQ-TOOL-001**: Tool detection MUST identify tools through:
- Dependency manifests (package.json, requirements.txt, go.mod, pom.xml, etc.)
- Configuration files (.eslintrc, pytest.ini, .golangci.yml, etc.)
- CI/CD workflow files
- Pre-commit hook configurations
- Scripts and Makefiles
- Lock files for version pinning

**REQ-TOOL-002**: Tool categorization MUST include:

```yaml
tool_categories:
  testing:
    - Test frameworks
    - Test runners
    - Coverage tools
    - Mutation testing tools
    - E2E testing frameworks
    - Performance testing tools
    - Contract testing tools

  quality:
    - Linters
    - Formatters
    - Static analyzers
    - Type checkers
    - Complexity analyzers
    - Code smell detectors

  security:
    - Security scanners
    - Dependency vulnerability scanners
    - Secret detection tools
    - SAST tools
    - License compliance tools

  automation:
    - Pre-commit hooks
    - Git hooks
    - CI/CD tools
    - Build automation
    - Release automation
```

**REQ-TOOL-003**: Tool registry MUST maintain for each tool:
- Tool name and aliases
- Applicable languages
- Category and subcategory
- Priority level (Essential, Recommended, Advanced, Optional)
- Current stable version
- Installation methods
- Configuration requirements
- Integration patterns
- Benefits (quantified where possible)
- Setup effort estimate
- Official documentation URL
- Sample configurations

### Language-Specific Tool Recommendations

**REQ-TOOL-004**: Tool recommendations MUST be language-aware:

**Python Tools**:
```yaml
essential:
  - pytest (testing framework)
  - coverage.py (coverage analysis)
  - black (code formatting)
  - mypy (type checking)

recommended:
  - pytest-cov (coverage integration)
  - pytest-xdist (parallel testing)
  - ruff (fast linting)
  - pre-commit (git hooks)

advanced:
  - mutmut (mutation testing)
  - hypothesis (property-based testing)
  - bandit (security linting)
```

**Go Tools**:
```yaml
essential:
  - testing (built-in, verify usage)
  - gofmt (built-in, verify usage)
  - go vet (built-in, verify usage)

recommended:
  - golangci-lint (comprehensive linting)
  - testify (assertion library)
  - gocov (coverage reporting)
  - gotestsum (better test output)

advanced:
  - ginkgo (BDD framework)
  - gomock (mocking framework)
  - go-fuzz (fuzzing)
```

**JavaScript/TypeScript Tools**:
```yaml
essential:
  - Jest or Vitest (testing framework)
  - ESLint (linting)
  - Prettier (formatting)
  - TypeScript (type safety)

recommended:
  - @testing-library/* (testing utilities)
  - c8 or nyc (coverage)
  - husky (git hooks)
  - lint-staged (pre-commit)

advanced:
  - Stryker (mutation testing)
  - Cypress or Playwright (E2E)
  - MSW (API mocking)
```

**REQ-TOOL-005**: Tool adoption scoring MUST:
- Calculate category-wise adoption percentages
- Weight essential tools higher than optional tools
- Bonus for proper configuration (not just installed)
- Bonus for CI/CD integration
- Normalize scores to 0-100 scale per category
- Aggregate category scores for overall tool adoption score

### Tool Setup Guides

**REQ-TOOL-006**: For each missing essential tool, provide:

```markdown
# Setup Guide: pytest-cov

**Category**: Testing > Coverage
**Priority**: Essential
**Effort**: 15 minutes
**Expected Benefit**: Integrated coverage reporting, 30% faster workflow

## Installation

```bash
pip install pytest-cov
```

## Configuration

Add to `pytest.ini` or `pyproject.toml`:

```ini
[tool:pytest]
addopts = --cov=src --cov-report=html --cov-report=term
```

## CI/CD Integration

Add to `.github/workflows/test.yml`:

```yaml
- name: Run tests with coverage
  run: pytest --cov --cov-report=xml

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
```

## Verification

```bash
pytest --cov
```

Expected output: Coverage report with percentages

## References
- Official docs: https://pytest-cov.readthedocs.io/
- Coverage.py docs: https://coverage.readthedocs.io/
```

**REQ-TOOL-007**: Setup guides MUST be:
- Language/framework-specific
- Include copy-paste ready commands
- Provide configuration examples
- Show CI/CD integration
- Include verification steps
- Reference official documentation

### Tool Comparison & Alternatives

**REQ-TOOL-008**: When multiple tool options exist, provide comparison:

```yaml
test_frameworks_python:
  - name: pytest
    adoption: 75%
    pros: [fixtures, plugins, parametrization]
    cons: [learning curve for fixtures]
    best_for: Most projects

  - name: unittest
    adoption: 20%
    pros: [built-in, familiar to Java devs]
    cons: [verbose, limited features]
    best_for: Simple projects, stdlib preference

  - name: nose2
    adoption: 3%
    pros: [unittest extension]
    cons: [maintenance status unclear]
    best_for: Legacy projects

recommendation: pytest (industry standard, rich ecosystem)
```

---

## Technology Stack

### Core Implementation

**REQ-TECH-001**: Primary implementation language MUST be Go for:
- Type safety and performance
- Strong concurrency support (analysis parallelization)
- Excellent CLI tooling
- Cross-platform compilation
- Fast startup time
- Single binary distribution

**REQ-TECH-002**: Core libraries and frameworks:
- **CLI**: cobra (command-line interface)
- **Configuration**: viper (config management)
- **Output**: tablewriter, color (terminal output)
- **HTML**: html/template (report generation)
- **JSON**: encoding/json (standard library)
- **Concurrency**: goroutines, channels (parallel execution)

### Language-Specific Analysis Components

**REQ-TECH-003**: Analysis components MAY use language-native tools:
- **Python analysis**: Python scripts/libraries for AST parsing
- **JavaScript analysis**: Node.js scripts for TS/JS AST parsing
- **Java analysis**: Java-based AST parsers (invoked via CLI)

**REQ-TECH-004**: Component communication MUST use:
- Standard input/output for data exchange
- JSON for structured data
- Exit codes for status signaling
- File-based intermediate results for large datasets

### External Tool Integration

**REQ-TECH-005**: Coverage tool integration:
- Execute coverage tools via shell commands
- Parse output files (XML, JSON, LCOV formats)
- Support custom coverage tool configurations
- Handle missing/failed coverage gracefully

**REQ-TECH-006**: Supported coverage formats:
- Cobertura XML (multi-language)
- JaCoCo XML (Java)
- LCOV (JavaScript, C++)
- Coverage.py JSON/XML (Python)
- Go coverage profiles (Go)

### Data Storage

**REQ-TECH-007**: Configuration storage:
- YAML for user configuration (.shipshape.yml)
- TOML support for alternative configuration
- JSON for CI/CD integration configs
- Environment variables for CI/CD overrides

**REQ-TECH-008**: Results storage:
- JSON for machine-readable results
- HTML for human-readable reports
- SQLite for historical data (optional)
- CSV for trend exports

### Build System

**REQ-TECH-009**: Build requirements:
- Go 1.21+ for core application
- Make for build automation
- Docker for containerized execution (optional)
- Cross-compilation for Windows, macOS, Linux

**REQ-TECH-010**: Distribution:
- Single binary per platform
- Homebrew formula (macOS/Linux)
- Chocolatey package (Windows)
- Docker image
- GitHub releases with artifacts

---

## Integration Points

### CI/CD Integration

**REQ-INT-001**: GitHub Actions integration:

```yaml
name: Ship Shape Analysis

on: [pull_request]

jobs:
  test-quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Ship Shape
        uses: shipshape/action@v1
        with:
          config: .shipshape.yml
          fail-on: critical,high

      - name: Upload HTML Report
        uses: actions/upload-artifact@v3
        with:
          name: shipshape-report
          path: shipshape-report.html

      - name: Comment PR
        uses: shipshape/pr-comment@v1
        with:
          report: shipshape-results.json
```

**REQ-INT-002**: GitLab CI integration:

```yaml
shipshape:
  stage: test
  image: shipshape/cli:latest
  script:
    - shipshape analyze --config .shipshape.yml --format json > results.json
    - shipshape gate --config .shipshape.yml
  artifacts:
    reports:
      shipshape: results.json
    paths:
      - shipshape-report.html
```

**REQ-INT-003**: Quality gate enforcement:
- Exit code 0 = all gates passed
- Exit code 1 = blocking gates failed
- Exit code 2 = warning gates triggered
- JSON output includes gate results

### API Integration (Future)

**REQ-INT-004**: REST API MUST support (planned Q4 2026):
- POST /analyze (trigger analysis)
- GET /results/:id (retrieve results)
- GET /trends/:repo (historical trends)
- POST /gates/evaluate (evaluate quality gates)
- GET /recommendations/:repo (get recommendations)

### Coverage Tool Integration

**REQ-INT-005**: Coverage integration MUST:
- Detect existing coverage configuration
- Execute coverage tools with project settings
- Parse coverage output in multiple formats
- Generate unified coverage report
- Support custom coverage commands

**Example Python coverage integration**:
```bash
# Detect: pytest.ini with --cov option or .coveragerc file
# Execute: pytest --cov=src --cov-report=xml --cov-report=json
# Parse: coverage.xml + coverage.json
# Output: Unified coverage metrics
```

### IDE Integration (Future)

**REQ-INT-006**: IDE extensions SHOULD support (future):
- VS Code extension with inline findings
- IntelliJ plugin for real-time analysis
- Language server protocol implementation
- Real-time test quality feedback

---

## Quality Gates & CI/CD

### Gate Configuration

**REQ-GATE-001**: Quality gates MUST be configurable via YAML:

```yaml
# .shipshape.yml
gates:
  blocking:
    - metric: test_coverage.line_coverage
      threshold: 80
      operator: ">="

    - metric: test_coverage.branch_coverage
      threshold: 70
      operator: ">="

    - metric: test_coverage.uncovered_critical
      threshold: 0
      operator: "=="

    - metric: test_quality.flaky_tests
      threshold: 0
      operator: "=="

    - metric: test_quality.score
      threshold: 70
      operator: ">="

  warning:
    - metric: test_quality.test_smells
      threshold: 10
      operator: "<="

    - metric: test_pyramid.compliance
      threshold: 0.7
      operator: ">="

    - metric: tool_adoption.essential_missing
      threshold: 0
      operator: "=="

  trend:
    - metric: test_coverage.line_coverage
      delta: 0
      operator: ">="
      message: "Coverage cannot decrease"

    - metric: test_quality.flaky_tests
      delta: 0
      operator: "<="
      message: "Cannot add new flaky tests"
```

**REQ-GATE-002**: Gate evaluation MUST:
- Evaluate all configured gates
- Support arithmetic operators: >=, <=, ==, !=, >, <
- Support trend-based gates (comparing to baseline)
- Provide detailed failure messages
- Allow per-language gate configuration
- Support per-package gates in monorepos

### Trend Analysis

**REQ-GATE-003**: Trend tracking MUST:
- Store historical results (optional SQLite database)
- Compare current results to previous runs
- Calculate deltas for key metrics
- Detect regressions
- Support baseline override for major refactors

**REQ-GATE-004**: Trend visualization MUST show:
- Coverage trends over time
- Test quality score trends
- Tool adoption progress
- Flaky test count trends
- Performance metrics trends

### CI/CD Optimization

**REQ-GATE-005**: CI/CD analysis MUST detect:
- Missing test parallelization opportunities
- Inefficient caching strategies
- Unnecessary test reruns
- Suboptimal workflow structure
- Missing failure notifications

**REQ-GATE-006**: CI/CD recommendations MUST include:
- Parallelization configuration examples
- Caching strategy improvements
- Workflow optimization suggestions
- Test sharding patterns
- Failure handling improvements

---

## Multi-Project Organization Features

### Organization Dashboard (Planned Q4 2026)

**REQ-ORG-001**: Organization dashboard MUST support:
- Multi-repository view
- Aggregated organization-wide metrics
- Cross-repo comparisons
- Ranked project listing
- Team-based filtering

**REQ-ORG-002**: Repository ranking MUST provide:
- Sortable table by any dimension
- Filterable by language, team, or status
- Comparison to organization baseline
- Identification of outliers (best and worst)
- Trend indicators (improving/declining)

**Example organization dashboard**:
```
Organization: MyOrg (42 repositories analyzed)

Overall Health: B (82/100)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Top Performers:
  1. api-gateway        A+ (96/100) ↑ +3
  2. auth-service       A  (92/100) →
  3. payment-processor  A  (91/100) ↑ +5

Needs Attention:
  40. legacy-monolith   D  (65/100) ↓ -2
  41. prototype-v2      D  (62/100) →
  42. old-scripts       F  (45/100) ↓ -5

By Language:
  Go:         85/100 (15 repos)
  Python:     78/100 (18 repos)
  JavaScript: 72/100 (9 repos)

By Team:
  Platform:   88/100 (12 repos)
  Backend:    82/100 (18 repos)
  Frontend:   75/100 (12 repos)
```

**REQ-ORG-003**: Cross-repo analysis MUST identify:
- Common missing tools across repos
- Inconsistent practices across teams
- Opportunities for shared testing infrastructure
- Best practices to spread across organization
- Technical debt hotspots

### Benchmarking

**REQ-ORG-004**: Benchmarking MUST support:
- Internal benchmarks (organization average)
- Industry benchmarks (by language/domain)
- Custom benchmark definitions
- Percentile rankings (top 10%, top 25%, etc.)

**REQ-ORG-005**: Benchmark visualization:
```
Your Repository vs. Organization

Test Coverage:     Your repo: 88%  ██████████░
                   Org avg:   82%  █████████░░
                   Top 10%:   95%  ███████████

Test Quality:      Your repo: 78%  ████████░░░
                   Org avg:   75%  ████████░░░
                   Top 10%:   90%  ██████████░

Tool Adoption:     Your repo: 65%  ███████░░░░
                   Org avg:   70%  ████████░░░
                   Top 10%:   95%  ███████████
```

### Rollout Tracking

**REQ-ORG-006**: Tool adoption rollout tracking MUST:
- Track tool adoption across all repositories
- Identify repositories missing specific tools
- Calculate organization-wide adoption percentages
- Visualize rollout progress over time
- Generate rollout status reports

**Example rollout tracking**:
```
Tool Rollout Status: golangci-lint

Organization Adoption: 73% (11/15 Go repositories)

Adopted:
  ✓ api-gateway
  ✓ auth-service
  ✓ payment-processor
  ... (8 more)

Not Yet Adopted:
  ✗ legacy-service (team: Backend)
  ✗ prototype-api (team: Platform)
  ✗ scripts-repo (team: DevOps)
  ✗ experimental-v2 (team: R&D)

Rollout Plan:
  Week 1-2: Setup guides and training sessions
  Week 3-4: Migrate Backend team repos (2 repos)
  Week 5-6: Migrate Platform team repos (1 repo)
  Week 7+:  Evaluate R&D repos for applicability
```

---

## Non-Functional Requirements

### Performance

**REQ-PERF-001**: Analysis performance targets:
- Repository discovery: <1 second
- Test structure analysis: <30 seconds for 1000 test files
- Static analysis: <5 minutes for 10,000 LOC
- Full analysis (without execution): <10 minutes for typical repo
- Report generation: <5 seconds

**REQ-PERF-002**: Memory usage:
- Maximum heap: <2GB for typical repository
- Streaming analysis for large files
- Incremental processing for monorepos

**REQ-PERF-003**: Parallelization:
- Analyze independent modules in parallel
- Support configurable worker count
- Respect system resources (CPU, memory limits)

### Scalability

**REQ-SCALE-001**: Repository size support:
- Single repos: up to 1M LOC
- Monorepos: up to 100 packages
- Test files: up to 10,000 test files
- Coverage data: up to 100MB coverage reports

**REQ-SCALE-002**: Organization support (future):
- Up to 1000 repositories
- Up to 100 teams
- Historical data retention: 2 years
- Concurrent analyses: 10 repositories

### Reliability

**REQ-REL-001**: Error handling:
- Graceful degradation on tool failures
- Partial results on analysis failures
- Clear error messages with remediation guidance
- Crash recovery and resume capability

**REQ-REL-002**: Validation:
- Configuration file validation with helpful errors
- Input sanitization for security
- Output format validation
- Coverage data integrity checks

### Usability

**REQ-USE-001**: Configuration:
- Zero-config for standard repositories
- Sensible defaults based on detected languages
- Override capability for all settings
- Configuration validation with suggestions

**REQ-USE-002**: Documentation:
- Comprehensive user guide
- API documentation (future)
- Setup guides for all recommended tools
- Troubleshooting guide
- Best practices documentation

**REQ-USE-003**: Feedback:
- Progress indicators for long operations
- Informative log output (with verbosity levels)
- Clear success/failure messages
- Actionable error messages

### Extensibility

**REQ-EXT-001**: Plugin system (future):
- Custom analyzer plugins
- Custom assessor plugins
- Custom report formats
- Custom quality gates

**REQ-EXT-002**: Configuration extensions:
- Custom tool definitions
- Custom test smell patterns
- Custom threshold configurations
- Custom weight adjustments

### Security

**REQ-SEC-001**: Execution security:
- No arbitrary code execution without explicit config
- Sandboxed test execution (optional)
- Credential handling (read-only, no storage)
- Dependency vulnerability scanning for tool recommendations

**REQ-SEC-002**: Data privacy:
- No external data transmission by default
- Local-only analysis
- Optional cloud integration with explicit consent
- No source code included in telemetry (if enabled)

### Compatibility

**REQ-COMPAT-001**: Platform support:
- Linux (Ubuntu 20.04+, RHEL 8+)
- macOS (11.0+)
- Windows (10+, WSL2 recommended)

**REQ-COMPAT-002**: Runtime dependencies:
- Go runtime: Not required (static binary)
- Language runtimes: Required for dynamic analysis (Python, Node.js, Java, etc.)
- Coverage tools: Optional but recommended

**REQ-COMPAT-003**: Version compatibility:
- Semantic versioning (semver)
- Backward compatible configuration
- Migration guides for breaking changes
- Deprecation warnings (1 version advance notice)

---

## Implementation Priorities

### Phase 1: MVP (Q1 2026)
- Core Go implementation
- Repository discovery and language detection
- Basic monorepo support
- Test structure analysis
- Coverage integration (Python, Go, JavaScript)
- Test quality analysis (test smells)
- HTML and JSON report generation
- CLI with essential commands

### Phase 2: Core Features (Q2 2026)
- All language support (Java, Rust, C#, Ruby, C++)
- Framework-specific analyzers
- Tool adoption analysis
- Quality gates and CI/CD integration
- Trend analysis (historical tracking)
- Performance analysis
- Extended test quality metrics

### Phase 3: Advanced Features (Q3 2026)
- Mutation testing integration
- Advanced CI/CD optimization
- Auto-fix capabilities
- IDE integration (VS Code)
- Enhanced monorepo support
- Custom analyzer plugins

### Phase 4: Enterprise Features (Q4 2026)
- Organization dashboard
- Multi-repository ranking
- Web-based dashboard
- REST API
- Team and RBAC support
- Advanced benchmarking
- Tool rollout tracking

---

## Success Criteria

**REQ-SUCCESS-001**: The tool is successful when:
- 90%+ of analyzed repositories receive actionable recommendations
- Coverage gaps are identified with <5% false positive rate
- Tool recommendations are language-appropriate 100% of the time
- Monorepo analysis correctly identifies all packages/modules 95%+ of the time
- Analysis completes in <10 minutes for 90% of repositories
- Quality gates prevent regression in 95%+ of cases
- Users report improved test quality after implementing recommendations

**REQ-SUCCESS-002**: Quality metrics:
- Code coverage: >90% unit test coverage for core analyzers
- Test quality: Ship Shape's own score: A grade (90+)
- Documentation completeness: 100% of public APIs documented
- Build time: <2 minutes for full build and test
- Release frequency: Monthly releases in Q1-Q2, bi-weekly in Q3-Q4

---

## Appendices

### Appendix A: Supported Test Frameworks

See detailed framework support matrix in research document.

### Appendix B: Evidence Base

All thresholds and recommendations are based on:
- 50+ academic research papers (ICSE, FSE, ASE, ISSTA conferences)
- Industry standards (ISO 29119, ISTQB)
- Framework vendor research
- Open source best practices

### Appendix C: Configuration Examples

See detailed configuration examples for:
- Single-language repositories
- Multi-language repositories
- Monorepos (npm, Go workspaces, Maven multi-module)
- CI/CD integration (GitHub Actions, GitLab CI, Jenkins)
- Quality gate configurations
- Custom analyzer configurations

### Appendix D: Glossary

- **Test Pyramid**: Testing strategy with many unit tests, fewer integration tests, and minimal E2E tests
- **Test Trophy**: Modern variant emphasizing integration tests over pure unit tests
- **Test Smell**: Anti-pattern in test code that indicates potential quality issues
- **Mutation Testing**: Technique that modifies source code to verify test effectiveness
- **Flaky Test**: Test that exhibits non-deterministic behavior (passes/fails inconsistently)
- **Critical Path**: Code paths essential for core functionality (security, data integrity, business logic)
- **Test Independence**: Tests that don't depend on execution order or shared state

---

**Document Status**: Draft v1.0
**Next Review**: After technical spike completion
**Approval Required**: Architecture review, security review, UX review
