# Ship Shape: Testing Quality & Code Analysis Tool

## Design Proposal for a Production-Ready Testing Excellence Platform

**Author**: Senior Software Engineer Analysis
**Date**: 2026-01-27
**Based on**: AgentReady architecture analysis
**Focus**: Testing quality, test coverage, test harnesses, code quality, best practices

---

## Executive Summary

Ship Shape is a proposed comprehensive testing quality and code analysis platform that extends beyond AI-readiness to provide deep, actionable insights into test health, coverage effectiveness, and testing best practices. Unlike AgentReady's focus on AI agent compatibility, Ship Shape targets test quality excellence, comprehensive coverage analysis, and testing framework optimization.

**Key Differentiators:**
- Multi-dimensional testing quality scoring (unit, integration, e2e, performance)
- Test harness detection and optimization for language-specific frameworks
- **Tool adoption analysis** - Compare tools in use vs. available free/opensource tools
- Real-time test quality feedback with auto-fix capabilities
- Test trend analysis and quality gates for CI/CD integration
- Framework-specific best practices (pytest, Jest, JUnit, Cypress, etc.)
- Test performance and flakiness detection
- **Actionable recommendations** with setup guides for missing essential tools

---

## 1. Research Gathering Strategy

### 1.1 Evidence-Based Testing Metrics Sources

Drawing from AgentReady's research-driven approach, Ship Shape will base all metrics on peer-reviewed research and industry testing standards:

#### **A. Academic Research Sources**

| Source Type | Examples | What We Extract |
|-------------|----------|-----------------|
| **Software Testing** | ICST, ISSTA, ASE proceedings | Test effectiveness metrics, coverage adequacy |
| **Empirical Studies** | MSR (Mining Software Repositories) | Test-defect correlations, coverage thresholds |
| **Testing Research** | IEEE Software, ACM TOSEM | Mutation testing effectiveness, test smell patterns |
| **Quality Metrics** | Software Quality Journal | Test quality indicators, flakiness predictors |

**Key Papers to Reference:**
1. "Coverage Is Not Strongly Correlated with Test Suite Effectiveness" (Inozemtseva & Holmes, ICSE 2014)
2. "An Empirical Study of Flaky Tests" (Luo et al., FSE 2014)
3. "Mutation Testing: An Empirical Evaluation" (Andrews et al., IEEE TSE 2005)
4. "Test Smells - The Bad Smell of Test Code" (Van Deursen et al., XP 2001)
5. "How Developers Test Machine Learning Applications" (Humbatova et al., FSE 2020)
6. "The Impact of Test Ownership and Team Structure on the Reliability of Production Software" (Bird et al., Microsoft Research 2011)

#### **B. Industry Standards & Testing Frameworks**

| Standard | Focus Area | Integration Approach |
|----------|------------|---------------------|
| **ISO 29119** | Software testing standards | Map test categories to international standards |
| **ISTQB** | International testing certification | Test level categorization (unit, integration, system) |
| **Test Pyramid (Fowler)** | Test distribution strategy | Validate test suite balance |
| **Testing Trophy (Dodds)** | Modern testing philosophy | Integration test emphasis validation |
| **Google Testing Blog** | Large-scale testing practices | Flakiness mitigation, test sizing |

#### **C. Testing Tool Vendor Research**

Learn from established testing tools and frameworks:

| Tool/Framework | Strengths to Adopt | Data Source |
|----------------|-------------------|-------------|
| **pytest** | Fixture patterns, parametrization | Plugin ecosystem, best practices docs |
| **Jest** | Snapshot testing, mocking patterns | Configuration patterns, matchers |
| **JUnit** | Test lifecycle hooks, assertions | Jupiter API, vintage patterns |
| **Cypress** | E2E best practices, selector strategies | Best practices guide, plugin patterns |
| **TestNG** | Data providers, parallel execution | XML configuration patterns |
| **RSpec** | Behavior-driven patterns | Shared examples, let declarations |
| **Go testing** | Table-driven tests, subtests | Standard library patterns |
| **Coverage.py** | Branch coverage, source analysis | Coverage metrics methodology |
| **Istanbul/nyc** | JavaScript coverage standards | Threshold configuration |
| **Stryker** | Mutation testing patterns | Mutator configurations |

#### **D. Language Ecosystem Testing Best Practices**

| Language | Testing Frameworks | Metrics to Extract |
|----------|-------------------|-------------------|
| **Python** | pytest, unittest, nose2, doctest | Fixture usage, parametrization patterns, coverage targets |
| **Go** | testing, testify, ginkgo | Table-driven tests, subtests, benchmark patterns |
| **JavaScript/TypeScript** | Jest, Mocha, Vitest, Cypress, Playwright | Async testing, mocking strategies, snapshot usage |
| **Java** | JUnit 5, TestNG, Mockito, AssertJ | Annotation patterns, assertion styles, test containers |
| **Rust** | cargo test, proptest | Property-based testing, doc tests, integration tests |
| **C/C++** | Google Test, Catch2, CppUnit | Fixture patterns, assertion macros, test discovery |
| **Ruby** | RSpec, Minitest | BDD patterns, shared contexts, factory patterns |
| **C#** | NUnit, xUnit.net, MSTest | Theory/fact patterns, fixtures, async testing |

#### **E. Test Harness & CI/CD Platform Research**

| Platform | Purpose | Integration Patterns |
|----------|---------|---------------------|
| **GitHub Actions** | CI/CD workflows | Test matrix strategies, caching patterns |
| **GitLab CI** | Pipeline testing | Test stage optimization, artifact management |
| **Jenkins** | Classic CI/CD | Test result aggregation, parallel execution |
| **CircleCI** | Cloud testing | Test splitting, timing-based parallelism |
| **Travis CI** | Open source testing | Build matrix configuration |
| **Buildkite** | Distributed testing | Agent-based parallelism |

### 1.2 Research Collection Methodology

**Automated Research Pipeline:**

```python
class TestingResearchCollector:
    """Continuously update evidence-based testing rule justifications"""

    def collect_testing_citations(self, metric: str) -> list[Citation]:
        """
        Search ICST, ISSTA, ASE for testing metric justification
        Returns: Citations with DOI, abstract, key findings
        """

    def validate_coverage_thresholds(self, language: str) -> CoverageThresholds:
        """
        Cross-reference coverage thresholds against industry data
        Example: Is 80% line coverage still valid for Python in 2026?
        """

    def extract_framework_patterns(self, framework: str) -> TestPatterns:
        """
        Parse testing framework best practices
        Example: Extract pytest fixture patterns, Jest mocking strategies
        """

    def analyze_test_trends(self) -> list[TestingTrend]:
        """
        Identify emerging testing patterns from open source
        Example: Rise of contract testing, shift to integration focus
        """
```

**Research Documentation Requirements:**
- Every testing check must cite at least 2 sources (academic OR framework documentation)
- Coverage thresholds must be justified with empirical data
- Test smell detection based on published anti-patterns
- Annual review process to update based on new testing research

---

## 2. Repository Assessment Mechanisms

### 2.1 Multi-Layer Testing Analysis Architecture

Ship Shape employs a comprehensive testing-focused analysis pipeline:

```
Layer 1: Test Discovery (fast, <1s)
         ├─ Test file detection (naming patterns, decorators, locations)
         ├─ Test framework identification (pytest, Jest, JUnit, etc.)
         ├─ Test configuration parsing (pytest.ini, jest.config.js, etc.)
         └─ Test-to-code ratio calculation

Layer 2: Test Structure Analysis (fast, 5-30s)
         ├─ Test categorization (unit, integration, e2e, performance)
         ├─ Test pyramid/trophy validation
         ├─ Test file organization analysis
         └─ Naming convention compliance

Layer 3: Test Quality Analysis (medium, 30s-5m)
         ├─ Test smell detection (long tests, unclear assertions, etc.)
         ├─ Assertion quality analysis
         ├─ Test independence verification
         ├─ Setup/teardown pattern analysis
         └─ Mock/stub usage patterns

Layer 4: Coverage Analysis (medium, 1-10m depending on size)
         ├─ Line coverage measurement
         ├─ Branch coverage analysis
         ├─ Function/method coverage
         ├─ Path coverage (advanced)
         └─ Uncovered critical path detection

Layer 5: Test Execution Analysis (variable, requires running tests)
         ├─ Test execution time measurement
         ├─ Flaky test detection
         ├─ Test failure pattern analysis
         ├─ Parallel execution compatibility
         └─ Resource usage profiling

Layer 6: Mutation Testing (intensive, 10m-2h)
         ├─ Mutation score calculation
         ├─ Surviving mutant analysis
         ├─ Test effectiveness measurement
         └─ Critical code mutation coverage

Layer 7: Integration & CI/CD Analysis (fast, <30s)
         ├─ CI/CD test configuration validation
         ├─ Test parallelization assessment
         ├─ Test caching strategy analysis
         └─ Test result reporting quality
```

### 2.2 Analyzer Architecture

**Core Abstraction (improving on AgentReady's BaseAssessor):**

```go
// Using Go for performance and type safety
package analyzers

import (
    "context"
    "time"
)

// Severity levels for findings
type Severity int

const (
    SeverityInfo Severity = iota
    SeverityLow
    SeverityMedium
    SeverityHigh
    SeverityCritical
)

// FindingType categorizes the issue
type FindingType int

const (
    TypeTestCoverage FindingType = iota
    TypeTestQuality
    TypeTestPerformance
    TypeTestMaintainability
    TypeCodeQuality
    TypeDocumentation
    TypeBestPractices
)

// TestCategory classifies test types
type TestCategory int

const (
    TestCategoryUnit TestCategory = iota
    TestCategoryIntegration
    TestCategoryEndToEnd
    TestCategoryPerformance
    TestCategoryContract
    TestCategorySmoke
    TestCategoryRegression
)

// Location in source code
type Location struct {
    FilePath  string `json:"file_path"`
    StartLine int    `json:"start_line"`
    EndLine   int    `json:"end_line"`
    StartCol  int    `json:"start_col,omitempty"`
    EndCol    int    `json:"end_col,omitempty"`
    Snippet   string `json:"snippet,omitempty"`
}

// Finding represents a single detected issue
type Finding struct {
    ID          string       `json:"id"`
    CheckID     string       `json:"check_id"`
    Type        FindingType  `json:"type"`
    Severity    Severity     `json:"severity"`
    Title       string       `json:"title"`
    Description string       `json:"description"`
    Location    Location     `json:"location"`
    Evidence    []string     `json:"evidence"`
    Remediation Remediation  `json:"remediation"`
    References  []Reference  `json:"references"`
    Confidence  float64      `json:"confidence"`
}

// Remediation guidance
type Remediation struct {
    Description string        `json:"description"`
    Effort      string        `json:"effort"`
    AutoFix     *AutoFix      `json:"auto_fix,omitempty"`
    Examples    []CodeExample `json:"examples"`
}

// Reference to supporting research
type Reference struct {
    Title   string `json:"title"`
    URL     string `json:"url"`
    DOI     string `json:"doi,omitempty"`
    Summary string `json:"summary"`
}

// TestFramework information
type TestFramework struct {
    Name       string   `json:"name"`      // pytest, Jest, JUnit, etc.
    Version    string   `json:"version"`
    ConfigFile string   `json:"config_file"`
    Plugins    []string `json:"plugins"`
}

// TestMetrics collected during analysis
type TestMetrics struct {
    TotalTests      int               `json:"total_tests"`
    TestsByCategory map[TestCategory]int `json:"tests_by_category"`
    LineCoverage    float64           `json:"line_coverage"`
    BranchCoverage  float64           `json:"branch_coverage"`
    MutationScore   float64           `json:"mutation_score,omitempty"`
    AvgTestTime     time.Duration     `json:"avg_test_time"`
    FlakyTests      int               `json:"flaky_tests"`
    TestSmells      int               `json:"test_smells"`
}
```

### 2.3 Specific Analyzer Categories

#### **A. Test Coverage Analyzers**

| Analyzer | Metrics | Tools | Thresholds |
|----------|---------|-------|------------|
| **LineCoverageAnalyzer** | Line coverage percentage | coverage.py, Istanbul, gocov | >80% good, >90% excellent |
| **BranchCoverageAnalyzer** | Branch coverage percentage | coverage.py (branch), Istanbul | >70% good, >85% excellent |
| **FunctionCoverageAnalyzer** | Function/method coverage | Language-specific tools | >85% good |
| **CriticalPathCoverageAnalyzer** | Coverage of critical code paths | Custom analysis | 100% for critical paths |
| **UncoveredCodeAnalyzer** | Identifies uncovered dangerous code | Pattern matching | Zero uncovered error handling |

**Example: Line Coverage Analyzer (Go)**

```go
type LineCoverageAnalyzer struct {
    baseAnalyzer
}

func (a *LineCoverageAnalyzer) Analyze(ctx *AnalysisContext) (*AnalysisResult, error) {
    result := &AnalysisResult{
        AnalyzerID: a.ID(),
        Status:     "success",
    }

    // Detect testing framework
    framework := detectTestFramework(ctx.Repository)

    // Run coverage tool based on language
    var coverage CoverageData
    var err error

    switch ctx.Repository.PrimaryLanguage {
    case "Python":
        coverage, err = runCoveragePy(ctx.Repository.Path)
    case "JavaScript", "TypeScript":
        coverage, err = runIstanbul(ctx.Repository.Path)
    case "Go":
        coverage, err = runGoCoverage(ctx.Repository.Path)
    case "Java":
        coverage, err = runJacoco(ctx.Repository.Path)
    default:
        return result.withSkipped("Language not supported"), nil
    }

    if err != nil {
        return result.withError(err), nil
    }

    // Analyze coverage quality
    if coverage.LineCoverage < 80.0 {
        result.Findings = append(result.Findings, Finding{
            CheckID:     "coverage-low",
            Type:        TypeTestCoverage,
            Severity:    SeverityHigh,
            Title:       "Low line coverage",
            Description: fmt.Sprintf("Line coverage is %.1f%%, below recommended 80%%", coverage.LineCoverage),
            Remediation: Remediation{
                Description: "Increase test coverage by adding tests for uncovered code paths",
                Effort:      "medium",
                Examples: []CodeExample{
                    {
                        Explanation: "Focus on uncovered critical paths first",
                        Code:        generateCoverageExample(coverage.UncoveredCriticalPaths),
                    },
                },
            },
            References: []Reference{
                {
                    Title: "Coverage Is Not Strongly Correlated with Test Suite Effectiveness",
                    DOI:   "10.1145/2568225.2568271",
                    Summary: "While coverage isn't everything, 80%+ coverage correlates with lower defect rates",
                },
            },
            Confidence: 1.0,
        })
    }

    // Identify uncovered critical code
    for _, path := range coverage.UncoveredCriticalPaths {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "coverage-critical-missing",
            Type:     TypeTestCoverage,
            Severity: SeverityCritical,
            Title:    "Critical code path not tested",
            Location: Location{
                FilePath:  path.FilePath,
                StartLine: path.StartLine,
                EndLine:   path.EndLine,
            },
            Confidence: 0.9,
        })
    }

    result.Metrics = Metrics{
        CustomMetrics: map[string]any{
            "line_coverage":     coverage.LineCoverage,
            "branch_coverage":   coverage.BranchCoverage,
            "uncovered_lines":   coverage.UncoveredLines,
            "critical_uncovered": len(coverage.UncoveredCriticalPaths),
        },
    }

    return result, nil
}
```

#### **B. Test Quality Analyzers**

| Analyzer | Detections | Quality Indicators |
|----------|-----------|-------------------|
| **TestSmellAnalyzer** | Mystery guest, eager test, lazy test, etc. | Test independence, clarity |
| **AssertionQualityAnalyzer** | Assertion count, specificity, clarity | >2 assertions avg, specific matchers |
| **TestNamingAnalyzer** | Naming conventions, clarity | Descriptive names, pattern compliance |
| **TestSizeAnalyzer** | Test length, complexity | <50 LOC per test, CC <5 |
| **TestIndependenceAnalyzer** | Shared state, test order dependencies | Zero dependencies |
| **MockingPatternAnalyzer** | Mock overuse, mock leakage | Appropriate mocking level |

**Common Test Smells Detected:**

```go
type TestSmellAnalyzer struct {
    baseAnalyzer
    smellDetectors []TestSmellDetector
}

type TestSmell int

const (
    SmellMysteryGuest TestSmell = iota  // Unclear test dependencies
    SmellEagerTest                       // Tests too much
    SmellLazyTest                        // Asserts too little
    SmellObscureTest                     // Unclear what's being tested
    SmellConditionalLogic                // If/switch in test
    SmellGeneralFixture                  // Setup does too much
    SmellTestCodeDuplication             // Repeated test code
    SmellAssertionRoulette               // Multiple asserts without messages
    SmellSensitiveEquality               // Fragile equality checks
    SmellResourceOptimism                // Assumes external resource available
    SmellFlakiness                       // Non-deterministic behavior
)

func (a *TestSmellAnalyzer) Analyze(ctx *AnalysisContext) (*AnalysisResult, error) {
    result := &AnalysisResult{AnalyzerID: a.ID()}

    testFiles := findTestFiles(ctx.Repository)

    for _, testFile := range testFiles {
        ast := parseTestFile(testFile)

        for _, testFunc := range extractTestFunctions(ast) {
            // Detect test smells
            smells := a.detectSmells(testFunc)

            for _, smell := range smells {
                severity := a.smellSeverity(smell.Type)

                result.Findings = append(result.Findings, Finding{
                    CheckID:  fmt.Sprintf("test-smell-%s", smell.Type),
                    Type:     TypeTestQuality,
                    Severity: severity,
                    Title:    smell.Title,
                    Location: smell.Location,
                    Description: smell.Description,
                    Remediation: Remediation{
                        Description: smell.Remediation,
                        Effort:      smell.EffortEstimate,
                        Examples:    smell.RefactoredExamples,
                    },
                    References: []Reference{
                        {
                            Title: "Test Smells - The Bad Smell of Test Code",
                            URL:   "https://testsmells.org/",
                            Summary: "Catalog of test code smells and refactorings",
                        },
                    },
                })
            }
        }
    }

    return result, nil
}
```

#### **C. Test Harness & Framework Analyzers**

| Analyzer | Focus | Frameworks Supported |
|----------|-------|---------------------|
| **TestFrameworkDetector** | Identifies test frameworks in use | pytest, Jest, JUnit, Cypress, Go testing, RSpec, etc. |
| **FrameworkConfigAnalyzer** | Validates configuration best practices | Framework-specific configs |
| **FixturePatternAnalyzer** | Analyzes fixture usage patterns | Pytest fixtures, JUnit @Before/@After, etc. |
| **TestRunnerOptimizer** | Suggests parallelization, caching | All major frameworks |
| **CICDTestIntegration** | CI/CD test workflow analysis | GitHub Actions, GitLab CI, Jenkins |

**Framework-Specific Best Practices:**

```go
// Python pytest analyzer
type PytestAnalyzer struct {
    baseAnalyzer
}

func (a *PytestAnalyzer) Analyze(ctx *AnalysisContext) (*AnalysisResult, error) {
    result := &AnalysisResult{AnalyzerID: a.ID()}

    // Check for pytest.ini or pyproject.toml
    config := findPytestConfig(ctx.Repository)

    if config == nil {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "pytest-no-config",
            Type:     TypeBestPractices,
            Severity: SeverityMedium,
            Title:    "Missing pytest configuration",
            Description: "No pytest.ini or pyproject.toml [tool.pytest.ini_options] found",
            Remediation: Remediation{
                Description: "Create pytest.ini with recommended settings",
                Examples: []CodeExample{{
                    Code: `[pytest]
testpaths = tests
python_files = test_*.py
python_classes = Test*
python_functions = test_*
addopts =
    --strict-markers
    --cov=src
    --cov-report=term-missing
    --cov-report=html
    --cov-fail-under=80
markers =
    unit: Unit tests
    integration: Integration tests
    slow: Slow-running tests`,
                }},
            },
        })
    }

    // Analyze fixture usage
    fixtures := analyzePytestFixtures(ctx.Repository)

    // Check for fixture anti-patterns
    if fixtures.HasFixtureLeakage {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "pytest-fixture-leakage",
            Type:     TypeTestQuality,
            Severity: SeverityHigh,
            Title:    "Fixture state leakage detected",
            Description: "Fixtures are leaking state between tests",
        })
    }

    // Check for parametrization usage
    if fixtures.ParametrizationUsage < 0.1 && fixtures.SimilarTests > 5 {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "pytest-missing-parametrization",
            Type:     TypeTestMaintainability,
            Severity: SeverityMedium,
            Title:    "Consider using @pytest.mark.parametrize",
            Description: "Found repeated similar tests that could be parametrized",
            Remediation: Remediation{
                Description: "Use parametrization to reduce test duplication",
                Examples: []CodeExample{{
                    Before: `def test_add_positive():
    assert add(1, 2) == 3

def test_add_negative():
    assert add(-1, -2) == -3

def test_add_zero():
    assert add(0, 0) == 0`,
                    After: `@pytest.mark.parametrize("a,b,expected", [
    (1, 2, 3),
    (-1, -2, -3),
    (0, 0, 0),
])
def test_add(a, b, expected):
    assert add(a, b) == expected`,
                }},
            },
        })
    }

    return result, nil
}

// JavaScript/TypeScript Jest analyzer
type JestAnalyzer struct {
    baseAnalyzer
}

func (a *JestAnalyzer) Analyze(ctx *AnalysisContext) (*AnalysisResult, error) {
    result := &AnalysisResult{AnalyzerID: a.ID()}

    // Find jest.config.js/ts or package.json jest config
    config := findJestConfig(ctx.Repository)

    // Check for recommended configuration
    recommendations := []struct {
        check       bool
        checkID     string
        title       string
        description string
        example     string
    }{
        {
            check:       !config.HasCoverageThresholds,
            checkID:     "jest-no-coverage-thresholds",
            title:       "Missing coverage thresholds",
            description: "Jest config lacks coverage thresholds",
            example: `module.exports = {
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  }
};`,
        },
        {
            check:       !config.HasTestMatch,
            checkID:     "jest-no-test-match",
            title:       "Missing testMatch pattern",
            description: "Explicit test file patterns recommended",
            example: `testMatch: [
  "**/__tests__/**/*.[jt]s?(x)",
  "**/?(*.)+(spec|test).[jt]s?(x)"
]`,
        },
    }

    for _, rec := range recommendations {
        if rec.check {
            result.Findings = append(result.Findings, Finding{
                CheckID:     rec.checkID,
                Type:        TypeBestPractices,
                Severity:    SeverityLow,
                Title:       rec.title,
                Description: rec.description,
                Remediation: Remediation{
                    Examples: []CodeExample{{Code: rec.example}},
                },
            })
        }
    }

    // Analyze snapshot usage
    snapshots := analyzeJestSnapshots(ctx.Repository)
    if snapshots.TooManySnapshots {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "jest-snapshot-overuse",
            Type:     TypeTestQuality,
            Severity: SeverityMedium,
            Title:    "Excessive snapshot test usage",
            Description: fmt.Sprintf("%d snapshot tests found - consider explicit assertions", snapshots.Count),
            Remediation: Remediation{
                Description: "Snapshots should supplement, not replace, explicit assertions",
            },
        })
    }

    return result, nil
}

// Go testing analyzer
type GoTestingAnalyzer struct {
    baseAnalyzer
}

func (a *GoTestingAnalyzer) Analyze(ctx *AnalysisContext) (*AnalysisResult, error) {
    result := &AnalysisResult{AnalyzerID: a.ID()}

    // Analyze table-driven test usage
    tests := analyzeGoTests(ctx.Repository)

    // Check for table-driven test pattern
    if tests.TableDrivenPercentage < 0.7 && tests.TotalTests > 10 {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "go-missing-table-driven",
            Type:     TypeBestPractices,
            Severity: SeverityMedium,
            Title:    "Low table-driven test usage",
            Description: "Only " + fmt.Sprintf("%.0f%%", tests.TableDrivenPercentage*100) + " of tests use table-driven pattern",
            Remediation: Remediation{
                Description: "Use table-driven tests for better test coverage and maintainability",
                Examples: []CodeExample{{
                    Code: `func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}`,
                }},
            },
            References: []Reference{
                {
                    Title: "Table Driven Tests - Go Wiki",
                    URL:   "https://github.com/golang/go/wiki/TableDrivenTests",
                },
            },
        })
    }

    // Check for t.Parallel() usage
    if !tests.UsesParallel && tests.TotalTests > 20 {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "go-no-parallel-tests",
            Type:     TypeTestPerformance,
            Severity: SeverityLow,
            Title:    "Tests not marked for parallel execution",
            Description: "Consider using t.Parallel() to speed up test execution",
        })
    }

    return result, nil
}

// Cypress E2E analyzer
type CypressAnalyzer struct {
    baseAnalyzer
}

func (a *CypressAnalyzer) Analyze(ctx *AnalysisContext) (*AnalysisResult, error) {
    result := &AnalysisResult{AnalyzerID: a.ID()}

    // Find cypress.config.js/ts
    config := findCypressConfig(ctx.Repository)

    if config == nil {
        return result.withSkipped("Cypress not detected"), nil
    }

    // Analyze test organization
    tests := analyzeCypressTests(ctx.Repository)

    // Check for selector best practices
    if tests.HasDataCySelectors < 0.8 {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "cypress-poor-selectors",
            Type:     TypeBestPractices,
            Severity: SeverityMedium,
            Title:    "Not using data-cy selectors",
            Description: "Cypress tests should use data-cy attributes for stability",
            Remediation: Remediation{
                Description: "Add data-cy attributes and use cy.get('[data-cy=...]')",
                Examples: []CodeExample{{
                    Before: `cy.get('.submit-button').click()`,
                    After:  `cy.get('[data-cy=submit-button]').click()`,
                }},
            },
        })
    }

    // Check for proper waiting patterns
    if tests.HasImproperWaits {
        result.Findings = append(result.Findings, Finding{
            CheckID:  "cypress-improper-waits",
            Type:     TypeTestQuality,
            Severity: SeverityHigh,
            Title:    "Using cy.wait() with arbitrary timeouts",
            Description: "Replace cy.wait(ms) with proper assertions",
            Remediation: Remediation{
                Examples: []CodeExample{{
                    Before: `cy.get('button').click()
cy.wait(1000)
cy.get('.result').should('be.visible')`,
                    After: `cy.get('button').click()
cy.get('.result').should('be.visible')`,
                }},
            },
        })
    }

    return result, nil
}
```

#### **D. Test Performance Analyzers**

| Analyzer | Metrics | Quality Indicators |
|----------|---------|-------------------|
| **TestExecutionTimeAnalyzer** | Individual test timing | <100ms unit, <1s integration |
| **FlakinesDetector** | Non-deterministic test detection | Zero flaky tests |
| **TestParallelizationAnalyzer** | Parallel execution opportunities | Safe parallelization identification |
| **TestResourceAnalyzer** | Memory, CPU usage during tests | Resource leak detection |

#### **E. Mutation Testing Analyzers**

| Analyzer | Focus | Tools |
|----------|-------|-------|
| **MutationScoreAnalyzer** | Overall mutation score | Stryker, PITest, mutmut |
| **SurvivingMutantAnalyzer** | Identifies weak tests | Mutation testing tools |
| **CriticalCodeMutationAnalyzer** | Mutation coverage of critical code | Custom analysis |

---

### 2.4 Parallel Execution Engine

**Performance Optimization (learning from AgentReady's sequential issues):**

```go
type TestAnalysisEngine struct {
    analyzers  []Analyzer
    maxWorkers int
    cache      Cache
}

func (e *TestAnalysisEngine) Run(ctx *AnalysisContext) (*AssessmentReport, error) {
    // Phase 1: Quick test discovery
    testDiscovery := e.runTestDiscovery(ctx)

    // Phase 2: Group analyzers by dependencies
    layers := e.buildDependencyLayers()

    report := &AssessmentReport{
        StartTime: time.Now(),
        TestMetrics: testDiscovery.Metrics,
    }

    // Phase 3: Execute layers in parallel
    for _, layer := range layers {
        layerResults := make(chan *AnalysisResult, len(layer))

        for _, analyzer := range layer {
            go func(a Analyzer) {
                result, err := a.Analyze(ctx)
                if err != nil {
                    result = &AnalysisResult{
                        AnalyzerID: a.ID(),
                        Status:     "failed",
                        Error:      err.Error(),
                    }
                }
                layerResults <- result
            }(analyzer)
        }

        for i := 0; i < len(layer); i++ {
            result := <-layerResults
            report.Results = append(report.Results, result)
        }
    }

    // Phase 4: Aggregate and score
    report.Score = e.calculateScore(report.Results)
    report.EndTime = time.Now()

    return report, nil
}
```

---

## 3. Scoring Methodology

### 3.1 Multi-Dimensional Testing Score

**Testing-Focused Scoring Dashboard:**

```go
type QualityScore struct {
    Overall         float64           `json:"overall"`
    Dimensions      DimensionScores   `json:"dimensions"`
    Trends          TrendAnalysis     `json:"trends"`
    Gates           []QualityGate     `json:"gates"`
    TestPyramid     PyramidAnalysis   `json:"test_pyramid"`
}

type DimensionScores struct {
    TestCoverage    CoverageScore        `json:"test_coverage"`       // 30%
    TestQuality     TestQualityScore     `json:"test_quality"`        // 25%
    TestPerformance PerformanceScore     `json:"test_performance"`    // 15%
    Maintainability MaintainabilityScore `json:"maintainability"`     // 15%
    CodeQuality     CodeQualityScore     `json:"code_quality"`        // 10%
    Documentation   DocumentationScore   `json:"documentation"`       // 5%
}

type CoverageScore struct {
    Score            float64         `json:"score"`           // 0-100
    Grade            string          `json:"grade"`           // A+, A, B, C, D, F
    LineCoverage     float64         `json:"line_coverage"`   // Percentage
    BranchCoverage   float64         `json:"branch_coverage"` // Percentage
    FunctionCoverage float64         `json:"function_coverage"`
    MutationScore    float64         `json:"mutation_score,omitempty"`
    UncoveredCritical int            `json:"uncovered_critical"` // Critical paths not tested
    CoverageByModule map[string]float64 `json:"coverage_by_module"`
}

type TestQualityScore struct {
    Score               float64            `json:"score"`
    Grade               string             `json:"grade"`
    TestSmells          int                `json:"test_smells"`
    TestIndependence    float64            `json:"test_independence"` // 0-1
    AssertionQuality    float64            `json:"assertion_quality"`
    NamingQuality       float64            `json:"naming_quality"`
    TestsByCategory     map[TestCategory]int `json:"tests_by_category"`
    PyramidCompliance   float64            `json:"pyramid_compliance"` // How well balanced
}

type PerformanceScore struct {
    Score             float64       `json:"score"`
    Grade             string        `json:"grade"`
    AvgUnitTestTime   time.Duration `json:"avg_unit_test_time"`
    AvgIntegTestTime  time.Duration `json:"avg_integ_test_time"`
    TotalTestTime     time.Duration `json:"total_test_time"`
    FlakyTests        int           `json:"flaky_tests"`
    SlowTests         int           `json:"slow_tests"`
    ParallelizationScore float64    `json:"parallelization_score"`
}

type PyramidAnalysis struct {
    UnitTests        int     `json:"unit_tests"`
    IntegrationTests int     `json:"integration_tests"`
    E2ETests         int     `json:"e2e_tests"`
    Ratio            string  `json:"ratio"` // e.g., "70:20:10"
    Compliance       float64 `json:"compliance"` // How close to ideal pyramid
    Recommendation   string  `json:"recommendation"`
}
```

### 3.2 Scoring Algorithm Design

#### **A. Test Coverage Scoring (Evidence-Based)**

```go
func (s *CoverageScorer) Calculate(metrics map[string]float64) CoverageScore {
    lineCov := metrics["line_coverage"]
    branchCov := metrics["branch_coverage"]
    funcCov := metrics["function_coverage"]
    mutationScore := metrics["mutation_score"]
    uncoveredCritical := int(metrics["uncovered_critical"])

    // Weighted sub-scores
    lineScore := s.scoreLineCoverage(lineCov)
    branchScore := s.scoreBranchCoverage(branchCov)
    funcScore := s.scoreFunctionCoverage(funcCov)
    mutationScoreVal := s.scoreMutationTesting(mutationScore)
    criticalPenalty := float64(uncoveredCritical) * 10.0 // -10 per uncovered critical path

    // Weighted average with critical path penalty
    score := (lineScore*0.3 + branchScore*0.3 + funcScore*0.2 + mutationScoreVal*0.2) - criticalPenalty

    if score < 0 {
        score = 0
    }

    return CoverageScore{
        Score:             score,
        Grade:             scoreToGrade(score),
        LineCoverage:      lineCov,
        BranchCoverage:    branchCov,
        FunctionCoverage:  funcCov,
        MutationScore:     mutationScore,
        UncoveredCritical: uncoveredCritical,
    }
}

// Research-based thresholds (Google Testing Blog, Microsoft Research)
func (s *CoverageScorer) scoreLineCoverage(coverage float64) float64 {
    switch {
    case coverage >= 90:
        return 100.0 // Excellent
    case coverage >= 80:
        return 85.0  // Good (industry standard)
    case coverage >= 70:
        return 70.0  // Acceptable
    case coverage >= 60:
        return 50.0  // Poor
    default:
        return 25.0  // Critical
    }
}

func (s *CoverageScorer) scoreBranchCoverage(coverage float64) float64 {
    // Branch coverage typically 10-15% lower than line coverage
    switch {
    case coverage >= 85:
        return 100.0
    case coverage >= 70:
        return 85.0
    case coverage >= 60:
        return 70.0
    default:
        return 40.0
    }
}

func (s *CoverageScorer) scoreMutationTesting(mutationScore float64) float64 {
    // Mutation score is strong indicator of test effectiveness
    switch {
    case mutationScore >= 80:
        return 100.0 // Excellent test quality
    case mutationScore >= 70:
        return 80.0  // Good
    case mutationScore >= 60:
        return 60.0  // Adequate
    default:
        return 30.0  // Weak tests
    }
}
```

#### **B. Test Quality Scoring**

```go
func (s *TestQualityScorer) Calculate(metrics TestMetrics, findings []Finding) TestQualityScore {
    score := 100.0

    // Deduct for test smells
    testSmellCount := countFindingsByType(findings, TypeTestQuality)
    score -= float64(testSmellCount) * 2.0 // -2 points per smell

    // Deduct for flaky tests
    score -= float64(metrics.FlakyTests) * 5.0 // -5 points per flaky test (severe)

    // Reward for good test independence
    if metrics.TestIndependence >= 0.95 {
        score += 10.0
    }

    // Analyze test pyramid compliance
    pyramidScore := s.analyzeTestPyramid(metrics.TestsByCategory)
    score = score * pyramidScore // Multiply by pyramid compliance (0-1)

    if score < 0 {
        score = 0
    }
    if score > 100 {
        score = 100
    }

    return TestQualityScore{
        Score:            score,
        Grade:            scoreToGrade(score),
        TestSmells:       testSmellCount,
        FlakyTests:       metrics.FlakyTests,
        PyramidCompliance: pyramidScore,
    }
}

func (s *TestQualityScorer) analyzeTestPyramid(testsByCategory map[TestCategory]int) float64 {
    unit := testsByCategory[TestCategoryUnit]
    integration := testsByCategory[TestCategoryIntegration]
    e2e := testsByCategory[TestCategoryEndToEnd]

    total := float64(unit + integration + e2e)
    if total == 0 {
        return 0.0
    }

    unitPct := float64(unit) / total
    integPct := float64(integration) / total
    e2ePct := float64(e2e) / total

    // Ideal pyramid: 70% unit, 20% integration, 10% e2e
    // Trophy model: 45% unit, 40% integration, 15% e2e (modern preference)
    idealPyramid := []float64{0.70, 0.20, 0.10}
    idealTrophy := []float64{0.45, 0.40, 0.15}

    actual := []float64{unitPct, integPct, e2ePct}

    // Calculate distance from ideal (use closer of pyramid or trophy)
    pyramidDistance := euclideanDistance(actual, idealPyramid)
    trophyDistance := euclideanDistance(actual, idealTrophy)

    distance := math.Min(pyramidDistance, trophyDistance)

    // Convert distance to score (0 distance = 1.0, large distance = 0.0)
    compliance := math.Max(0, 1.0 - distance)

    return compliance
}
```

#### **C. Overall Score Composition**

**Testing-First Weighting:**

```go
func CalculateOverallScore(dimensions DimensionScores) float64 {
    // Testing-focused weighting
    weights := map[string]float64{
        "test_coverage":    0.30,  // 30% - Highest priority
        "test_quality":     0.25,  // 25% - Test effectiveness
        "test_performance": 0.15,  // 15% - Test speed and reliability
        "maintainability":  0.15,  // 15% - Long-term health
        "code_quality":     0.10,  // 10% - Code health
        "documentation":    0.05,  // 5%  - Knowledge transfer
    }

    overall := (dimensions.TestCoverage.Score * weights["test_coverage"]) +
               (dimensions.TestQuality.Score * weights["test_quality"]) +
               (dimensions.TestPerformance.Score * weights["test_performance"]) +
               (dimensions.Maintainability.Score * weights["maintainability"]) +
               (dimensions.CodeQuality.Score * weights["code_quality"]) +
               (dimensions.Documentation.Score * weights["documentation"])

    return overall
}
```

### 3.3 Quality Gates (CI/CD Integration)

**Testing-Focused Quality Gates:**

```yaml
# .shipshape/quality-gates.yml
quality_gates:
  # Block merge if any condition fails
  blocking:
    - test_coverage.line_coverage >= 80
    - test_coverage.branch_coverage >= 70
    - test_coverage.uncovered_critical == 0
    - test_quality.flaky_tests == 0
    - test_quality.score >= 70
    - test_performance.avg_unit_test_time < "100ms"

  # Warn but allow merge
  warning:
    - test_coverage.line_coverage >= 75
    - test_quality.test_smells <= 10
    - test_pyramid.compliance >= 0.7
    - test_performance.slow_tests <= 5

  # Trend-based gates
  trend_gates:
    - test_coverage.line_coverage.delta >= 0     # Cannot decrease coverage
    - test_quality.flaky_tests.delta <= 0        # Cannot add flaky tests
    - test_performance.total_time.delta <= "10s" # Tests can't get >10s slower
```

### 3.4 Test Pyramid/Trophy Analysis

**Visual Test Distribution Analysis:**

```go
type PyramidAnalyzer struct{}

func (a *PyramidAnalyzer) Analyze(tests TestMetrics) PyramidAnalysis {
    total := tests.UnitTests + tests.IntegrationTests + tests.E2ETests

    if total == 0 {
        return PyramidAnalysis{
            Recommendation: "No tests found - start with unit tests",
        }
    }

    unitPct := float64(tests.UnitTests) / float64(total) * 100
    integPct := float64(tests.IntegrationTests) / float64(total) * 100
    e2ePct := float64(tests.E2ETests) / float64(total) * 100

    ratio := fmt.Sprintf("%.0f:%.0f:%.0f", unitPct, integPct, e2ePct)

    // Determine compliance
    compliance := 1.0
    recommendation := "Test distribution looks good"

    if e2ePct > 20 {
        compliance -= 0.2
        recommendation = "Too many E2E tests - consider moving some to integration tests"
    }
    if unitPct < 50 {
        compliance -= 0.3
        recommendation = "Not enough unit tests - aim for 70% unit tests"
    }
    if integPct > 40 {
        compliance -= 0.1
        recommendation = "High integration test ratio - verify they're not slow"
    }

    return PyramidAnalysis{
        UnitTests:        tests.UnitTests,
        IntegrationTests: tests.IntegrationTests,
        E2ETests:         tests.E2ETests,
        Ratio:            ratio,
        Compliance:       math.Max(0, compliance),
        Recommendation:   recommendation,
    }
}
```

---

## 4. Test Harness Deep Dive

### 4.1 Framework-Specific Optimizations

#### **Python - pytest**

```yaml
# Ship Shape pytest recommendations
pytest_best_practices:
  configuration:
    - Use pytest.ini or pyproject.toml
    - Set testpaths, python_files, python_classes, python_functions
    - Configure markers for test categorization
    - Set coverage thresholds with --cov-fail-under

  fixtures:
    - Prefer function-scoped fixtures for test independence
    - Use conftest.py for shared fixtures
    - Avoid fixture leakage with proper cleanup
    - Use @pytest.fixture(autouse=True) sparingly

  parametrization:
    - Use @pytest.mark.parametrize for similar test cases
    - Prefer parametrize over loops in tests
    - Use indirect parametrization for complex setup

  plugins:
    - pytest-cov: Coverage measurement
    - pytest-xdist: Parallel execution
    - pytest-mock: Mocking support
    - pytest-asyncio: Async test support
    - pytest-timeout: Timeout protection

  performance:
    - Run with -n auto (pytest-xdist) for parallelization
    - Use --lf (last-failed) during development
    - Cache expensive fixtures with scope="session"
    - Profile slow tests with --durations=10
```

#### **JavaScript/TypeScript - Jest**

```yaml
jest_best_practices:
  configuration:
    - Define testMatch patterns
    - Set coverageThreshold
    - Configure transform for TypeScript
    - Use projects for monorepos

  test_organization:
    - Use describe blocks for grouping
    - Follow AAA pattern (Arrange, Act, Assert)
    - One assertion concept per test
    - Clear test names: "should ... when ..."

  mocking:
    - Use jest.mock() for module mocks
    - Prefer jest.spyOn() for partial mocks
    - Clean up with mockClear/mockReset
    - Avoid manual mocks when possible

  performance:
    - Run with --maxWorkers for parallelization
    - Use --onlyChanged during development
    - Avoid large snapshots
    - Set testTimeout for slow tests

  typescript:
    - Use @types/jest
    - Configure ts-jest transformer
    - Use type-safe matchers
```

#### **Go - testing package**

```yaml
go_testing_best_practices:
  table_driven_tests:
    - Use table-driven pattern for multiple cases
    - Use t.Run() for subtests
    - Name test cases descriptively

  parallelization:
    - Call t.Parallel() in independent tests
    - Avoid t.Parallel() with shared state
    - Use testdata directory for fixtures

  benchmarks:
    - Use b.N loop in benchmarks
    - Reset timer after setup: b.ResetTimer()
    - Run with -benchmem for allocations

  helpers:
    - Call t.Helper() in test helper functions
    - Return errors instead of calling t.Fatal in helpers

  coverage:
    - Run with -cover flag
    - Use -coverprofile for detailed analysis
    - Aim for >80% coverage

  performance:
    - Run tests with -race detector
    - Use -short flag for quick feedback
    - Parallelize with -parallel flag
```

#### **E2E - Cypress**

```yaml
cypress_best_practices:
  selectors:
    - Use data-cy attributes for stability
    - Avoid CSS class selectors
    - Avoid text content selectors

  waiting:
    - Never use cy.wait() with arbitrary timeouts
    - Use assertions for automatic retry
    - Use cy.intercept() for network waits

  organization:
    - Use custom commands for reuse
    - Page object pattern for complex pages
    - Separate test data from tests

  performance:
    - Minimize visit() calls
    - Use cy.session() for authentication
    - Run in parallel with Cypress Cloud
    - Use --browser for different browsers

  best_practices:
    - Test user journeys, not implementation
    - Keep tests independent
    - Clean up test data after each test
    - Use beforeEach for common setup
```

### 4.2 CI/CD Test Integration Patterns

#### **GitHub Actions**

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: ["3.10", "3.11", "3.12"]

    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: ${{ matrix.python-version }}

      - name: Cache dependencies
        uses: actions/cache@v3
        with:
          path: ~/.cache/pip
          key: ${{ runner.os }}-pip-${{ hashFiles('requirements.txt') }}

      - name: Install dependencies
        run: |
          pip install -r requirements.txt
          pip install pytest pytest-cov pytest-xdist

      - name: Run tests with coverage
        run: |
          pytest -n auto --cov=src --cov-report=xml --cov-report=term-missing

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.xml
          fail_ci_if_error: true

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: testpass
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - name: Run integration tests
        run: pytest tests/integration/ -m integration

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Cypress run
        uses: cypress-io/github-action@v6
        with:
          start: npm start
          wait-on: 'http://localhost:3000'
          wait-on-timeout: 120
          record: true
        env:
          CYPRESS_RECORD_KEY: ${{ secrets.CYPRESS_RECORD_KEY }}

  test-quality-gate:
    needs: [unit-tests, integration-tests, e2e-tests]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Ship Shape
        run: |
          shipshape analyze --format=json > results.json
          shipshape gate --config=.shipshape/gates.yml
```

### 4.3 Test Performance Optimization

```go
type TestPerformanceOptimizer struct{}

func (o *TestPerformanceOptimizer) AnalyzeAndOptimize(repo *Repository) []Recommendation {
    recommendations := []Recommendation{}

    // Detect parallelization opportunities
    if canParallelize(repo) {
        recommendations = append(recommendations, Recommendation{
            Type: "parallelization",
            Description: "Tests can run in parallel",
            Implementation: map[string]string{
                "pytest":  "pytest -n auto",
                "go":      "go test -parallel 4",
                "jest":    "jest --maxWorkers=50%",
            },
        })
    }

    // Detect caching opportunities
    if hasExpensiveSetup(repo) {
        recommendations = append(recommendations, Recommendation{
            Type: "caching",
            Description: "Cache expensive setup",
            Implementation: map[string]string{
                "pytest":  "Use session-scoped fixtures",
                "jest":    "Use setupFilesAfterEnv for one-time setup",
                "cypress": "Use cy.session() for auth state",
            },
        })
    }

    // Detect slow tests
    slowTests := findSlowTests(repo)
    if len(slowTests) > 0 {
        recommendations = append(recommendations, Recommendation{
            Type: "slow-tests",
            Description: fmt.Sprintf("%d slow tests detected", len(slowTests)),
            SlowTests: slowTests,
            Suggestions: []string{
                "Profile tests to find bottlenecks",
                "Use mocks instead of real dependencies",
                "Move slow tests to integration category",
                "Optimize database queries in tests",
            },
        })
    }

    return recommendations
}
```

---

## 5. Tool Adoption Analysis & Scoring

### 5.1 Overview

One of Ship Shape's key differentiators is evaluating repositories based on **tool adoption maturity** - comparing what free/opensource tools are currently in use versus what's available for the language ecosystem. This provides actionable insights into missed opportunities for quality improvement at zero cost.

**Core Philosophy:**
- **Detection**: Automatically identify tools already integrated
- **Gap Analysis**: Compare against comprehensive tool catalog
- **Scoring**: Reward tool adoption, penalize missing essential tools
- **Recommendations**: Suggest specific tools with setup instructions
- **ROI Calculation**: Show value of missing tools vs. implementation effort

### 5.2 Tool Categories & Detection

#### **A. Testing Tools**

| Category | Purpose | Detection Method |
|----------|---------|------------------|
| **Test Frameworks** | Unit/integration testing | Config files, imports, dependencies |
| **Coverage Tools** | Code coverage measurement | Config files, CI/CD scripts, dependencies |
| **Mutation Testing** | Test effectiveness validation | Config files, dependencies |
| **E2E Frameworks** | End-to-end testing | Config files, test directories |
| **Performance Testing** | Load/stress testing | Dependencies, test files |
| **Test Runners** | Test execution orchestration | CI/CD configs, package.json scripts |

#### **B. Code Quality Tools**

| Category | Purpose | Detection Method |
|----------|---------|------------------|
| **Linters** | Static code analysis | Config files (.eslintrc, .pylintrc, etc.) |
| **Formatters** | Code formatting | Config files (.prettierrc, pyproject.toml) |
| **Type Checkers** | Static type checking | Config files (tsconfig.json, mypy.ini) |
| **Complexity Analyzers** | Code complexity metrics | CI/CD integrations |
| **Duplication Detectors** | Code clone detection | CI/CD integrations |

#### **C. Security Tools**

| Category | Purpose | Detection Method |
|----------|---------|------------------|
| **Dependency Scanners** | Vulnerability detection | GitHub Dependabot, Snyk, etc. |
| **Secret Scanners** | Credential leak detection | git-secrets, trufflehog configs |
| **SAST Tools** | Static security analysis | Semgrep, Bandit, gosec configs |
| **License Scanners** | License compliance | Config files, CI/CD |

#### **D. CI/CD & Automation Tools**

| Category | Purpose | Detection Method |
|----------|---------|------------------|
| **CI/CD Platforms** | Continuous integration | .github/workflows, .gitlab-ci.yml, etc. |
| **Pre-commit Hooks** | Local quality gates | .pre-commit-config.yaml, husky config |
| **Build Tools** | Build automation | Makefile, build scripts |
| **Container Tools** | Containerization | Dockerfile, docker-compose.yml |

### 5.3 Language-Specific Tool Catalog

#### **Python Ecosystem**

```go
var PythonToolCatalog = ToolCatalog{
    Testing: []Tool{
        {Name: "pytest", Category: "framework", Priority: "essential", Free: true},
        {Name: "unittest", Category: "framework", Priority: "standard", Free: true},
        {Name: "coverage.py", Category: "coverage", Priority: "essential", Free: true},
        {Name: "pytest-cov", Category: "coverage", Priority: "essential", Free: true},
        {Name: "pytest-xdist", Category: "parallelization", Priority: "recommended", Free: true},
        {Name: "pytest-mock", Category: "mocking", Priority: "recommended", Free: true},
        {Name: "mutmut", Category: "mutation", Priority: "advanced", Free: true},
        {Name: "hypothesis", Category: "property-based", Priority: "advanced", Free: true},
        {Name: "tox", Category: "test-runner", Priority: "recommended", Free: true},
    },
    Quality: []Tool{
        {Name: "black", Category: "formatter", Priority: "essential", Free: true},
        {Name: "pylint", Category: "linter", Priority: "recommended", Free: true},
        {Name: "flake8", Category: "linter", Priority: "recommended", Free: true},
        {Name: "ruff", Category: "linter", Priority: "recommended", Free: true},
        {Name: "mypy", Category: "type-checker", Priority: "essential", Free: true},
        {Name: "isort", Category: "import-formatter", Priority: "recommended", Free: true},
        {Name: "radon", Category: "complexity", Priority: "optional", Free: true},
        {Name: "bandit", Category: "security", Priority: "essential", Free: true},
    },
    Security: []Tool{
        {Name: "bandit", Category: "sast", Priority: "essential", Free: true},
        {Name: "safety", Category: "dependency-scan", Priority: "essential", Free: true},
        {Name: "pip-audit", Category: "dependency-scan", Priority: "essential", Free: true},
        {Name: "semgrep", Category: "sast", Priority: "recommended", Free: true},
    },
    Automation: []Tool{
        {Name: "pre-commit", Category: "hooks", Priority: "essential", Free: true},
        {Name: "tox", Category: "automation", Priority: "recommended", Free: true},
        {Name: "nox", Category: "automation", Priority: "optional", Free: true},
    },
}
```

#### **JavaScript/TypeScript Ecosystem**

```go
var JavaScriptToolCatalog = ToolCatalog{
    Testing: []Tool{
        {Name: "jest", Category: "framework", Priority: "essential", Free: true},
        {Name: "vitest", Category: "framework", Priority: "recommended", Free: true},
        {Name: "mocha", Category: "framework", Priority: "alternative", Free: true},
        {Name: "istanbul", Category: "coverage", Priority: "essential", Free: true},
        {Name: "nyc", Category: "coverage", Priority: "essential", Free: true},
        {Name: "cypress", Category: "e2e", Priority: "essential", Free: true},
        {Name: "playwright", Category: "e2e", Priority: "recommended", Free: true},
        {Name: "stryker", Category: "mutation", Priority: "advanced", Free: true},
        {Name: "testing-library", Category: "framework", Priority: "recommended", Free: true},
    },
    Quality: []Tool{
        {Name: "eslint", Category: "linter", Priority: "essential", Free: true},
        {Name: "prettier", Category: "formatter", Priority: "essential", Free: true},
        {Name: "typescript", Category: "type-checker", Priority: "essential", Free: true},
        {Name: "tsc", Category: "type-checker", Priority: "essential", Free: true},
        {Name: "jscpd", Category: "duplication", Priority: "optional", Free: true},
    },
    Security: []Tool{
        {Name: "npm audit", Category: "dependency-scan", Priority: "essential", Free: true},
        {Name: "snyk", Category: "dependency-scan", Priority: "recommended", Free: true},
        {Name: "semgrep", Category: "sast", Priority: "recommended", Free: true},
    },
    Automation: []Tool{
        {Name: "husky", Category: "hooks", Priority: "essential", Free: true},
        {Name: "lint-staged", Category: "hooks", Priority: "recommended", Free: true},
        {Name: "commitlint", Category: "commit-validation", Priority: "recommended", Free: true},
    },
}
```

#### **Go Ecosystem**

```go
var GoToolCatalog = ToolCatalog{
    Testing: []Tool{
        {Name: "testing", Category: "framework", Priority: "standard", Free: true},
        {Name: "testify", Category: "framework", Priority: "recommended", Free: true},
        {Name: "ginkgo", Category: "framework", Priority: "alternative", Free: true},
        {Name: "go test -cover", Category: "coverage", Priority: "essential", Free: true},
        {Name: "go test -race", Category: "race-detection", Priority: "essential", Free: true},
        {Name: "gomock", Category: "mocking", Priority: "recommended", Free: true},
        {Name: "go-mutesting", Category: "mutation", Priority: "advanced", Free: true},
    },
    Quality: []Tool{
        {Name: "gofmt", Category: "formatter", Priority: "essential", Free: true},
        {Name: "goimports", Category: "import-formatter", Priority: "essential", Free: true},
        {Name: "golangci-lint", Category: "linter", Priority: "essential", Free: true},
        {Name: "staticcheck", Category: "linter", Priority: "essential", Free: true},
        {Name: "go vet", Category: "linter", Priority: "essential", Free: true},
        {Name: "gocyclo", Category: "complexity", Priority: "optional", Free: true},
    },
    Security: []Tool{
        {Name: "gosec", Category: "sast", Priority: "essential", Free: true},
        {Name: "govulncheck", Category: "dependency-scan", Priority: "essential", Free: true},
        {Name: "semgrep", Category: "sast", Priority: "recommended", Free: true},
    },
    Automation: []Tool{
        {Name: "pre-commit", Category: "hooks", Priority: "recommended", Free: true},
        {Name: "golangci-lint", Category: "automation", Priority: "essential", Free: true},
    },
}
```

#### **Java Ecosystem**

```go
var JavaToolCatalog = ToolCatalog{
    Testing: []Tool{
        {Name: "junit5", Category: "framework", Priority: "essential", Free: true},
        {Name: "testng", Category: "framework", Priority: "alternative", Free: true},
        {Name: "mockito", Category: "mocking", Priority: "essential", Free: true},
        {Name: "jacoco", Category: "coverage", Priority: "essential", Free: true},
        {Name: "pitest", Category: "mutation", Priority: "advanced", Free: true},
        {Name: "assertj", Category: "assertions", Priority: "recommended", Free: true},
        {Name: "testcontainers", Category: "integration", Priority: "recommended", Free: true},
    },
    Quality: []Tool{
        {Name: "checkstyle", Category: "linter", Priority: "essential", Free: true},
        {Name: "pmd", Category: "linter", Priority: "recommended", Free: true},
        {Name: "spotbugs", Category: "bug-detection", Priority: "essential", Free: true},
        {Name: "google-java-format", Category: "formatter", Priority: "recommended", Free: true},
        {Name: "errorprone", Category: "compiler-plugin", Priority: "recommended", Free: true},
    },
    Security: []Tool{
        {Name: "spotbugs-security", Category: "sast", Priority: "essential", Free: true},
        {Name: "dependency-check", Category: "dependency-scan", Priority: "essential", Free: true},
        {Name: "semgrep", Category: "sast", Priority: "recommended", Free: true},
    },
    Automation: []Tool{
        {Name: "maven", Category: "build", Priority: "standard", Free: true},
        {Name: "gradle", Category: "build", Priority: "alternative", Free: true},
    },
}
```

### 5.4 Tool Detection Implementation

```go
type ToolDetector struct {
    catalog map[string]ToolCatalog
}

type DetectedTool struct {
    Tool          Tool
    Detected      bool
    DetectionPath string  // Where found (e.g., "requirements.txt", ".github/workflows/test.yml")
    Version       string
    ConfigFile    string
    Configured    bool    // Has dedicated config file
}

type ToolAdoptionAnalysis struct {
    Language          string
    AvailableTools    []Tool
    DetectedTools     []DetectedTool
    MissingTools      []Tool
    AdoptionScore     float64         // 0-100
    CategoryScores    map[string]float64
    Recommendations   []ToolRecommendation
}

func (d *ToolDetector) AnalyzeRepository(repo *Repository) ToolAdoptionAnalysis {
    analysis := ToolAdoptionAnalysis{
        Language:      repo.PrimaryLanguage,
        CategoryScores: make(map[string]float64),
    }

    catalog := d.catalog[repo.PrimaryLanguage]
    analysis.AvailableTools = catalog.AllTools()

    // Detect tools by category
    testingTools := d.detectTestingTools(repo, catalog.Testing)
    qualityTools := d.detectQualityTools(repo, catalog.Quality)
    securityTools := d.detectSecurityTools(repo, catalog.Security)
    automationTools := d.detectAutomationTools(repo, catalog.Automation)

    analysis.DetectedTools = append(analysis.DetectedTools, testingTools...)
    analysis.DetectedTools = append(analysis.DetectedTools, qualityTools...)
    analysis.DetectedTools = append(analysis.DetectedTools, securityTools...)
    analysis.DetectedTools = append(analysis.DetectedTools, automationTools...)

    // Calculate adoption score
    analysis.AdoptionScore = d.calculateAdoptionScore(analysis)
    analysis.CategoryScores = d.calculateCategoryScores(analysis)

    // Generate recommendations
    analysis.MissingTools = d.findMissingEssentialTools(catalog, analysis.DetectedTools)
    analysis.Recommendations = d.generateRecommendations(analysis)

    return analysis
}

func (d *ToolDetector) detectTestingTools(repo *Repository, tools []Tool) []DetectedTool {
    detected := []DetectedTool{}

    for _, tool := range tools {
        result := DetectedTool{Tool: tool}

        switch repo.PrimaryLanguage {
        case "Python":
            result = d.detectPythonTool(repo, tool)
        case "JavaScript", "TypeScript":
            result = d.detectJavaScriptTool(repo, tool)
        case "Go":
            result = d.detectGoTool(repo, tool)
        case "Java":
            result = d.detectJavaTool(repo, tool)
        }

        if result.Detected {
            detected = append(detected, result)
        }
    }

    return detected
}

func (d *ToolDetector) detectPythonTool(repo *Repository, tool Tool) DetectedTool {
    result := DetectedTool{Tool: tool}

    // Check requirements.txt
    if fileExists(filepath.Join(repo.Path, "requirements.txt")) {
        content := readFile(filepath.Join(repo.Path, "requirements.txt"))
        if strings.Contains(content, tool.Name) {
            result.Detected = true
            result.DetectionPath = "requirements.txt"
        }
    }

    // Check pyproject.toml
    if fileExists(filepath.Join(repo.Path, "pyproject.toml")) {
        content := readFile(filepath.Join(repo.Path, "pyproject.toml"))
        if strings.Contains(content, tool.Name) {
            result.Detected = true
            result.DetectionPath = "pyproject.toml"
        }
    }

    // Check for tool-specific config files
    configFiles := map[string]string{
        "pytest":     "pytest.ini",
        "black":      "pyproject.toml",
        "mypy":       "mypy.ini",
        "pylint":     ".pylintrc",
        "bandit":     ".bandit",
        "pre-commit": ".pre-commit-config.yaml",
    }

    if configFile, exists := configFiles[tool.Name]; exists {
        if fileExists(filepath.Join(repo.Path, configFile)) {
            result.Configured = true
            result.ConfigFile = configFile
        }
    }

    return result
}

func (d *ToolDetector) detectJavaScriptTool(repo *Repository, tool Tool) DetectedTool {
    result := DetectedTool{Tool: tool}

    // Check package.json
    packageJSON := filepath.Join(repo.Path, "package.json")
    if fileExists(packageJSON) {
        var pkg struct {
            Dependencies    map[string]string `json:"dependencies"`
            DevDependencies map[string]string `json:"devDependencies"`
        }

        if readJSON(packageJSON, &pkg) == nil {
            if _, ok := pkg.Dependencies[tool.Name]; ok {
                result.Detected = true
                result.DetectionPath = "package.json (dependencies)"
            }
            if _, ok := pkg.DevDependencies[tool.Name]; ok {
                result.Detected = true
                result.DetectionPath = "package.json (devDependencies)"
            }
        }
    }

    // Check for tool-specific config files
    configFiles := map[string][]string{
        "eslint":    {".eslintrc.js", ".eslintrc.json", ".eslintrc.yml"},
        "prettier":  {".prettierrc", ".prettierrc.json", ".prettierrc.yml"},
        "jest":      {"jest.config.js", "jest.config.ts"},
        "cypress":   {"cypress.config.js", "cypress.config.ts"},
        "playwright": {"playwright.config.js", "playwright.config.ts"},
        "husky":     {".husky/"},
    }

    if configs, exists := configFiles[tool.Name]; exists {
        for _, configFile := range configs {
            if fileExists(filepath.Join(repo.Path, configFile)) {
                result.Configured = true
                result.ConfigFile = configFile
                break
            }
        }
    }

    return result
}

func (d *ToolDetector) detectGoTool(repo *Repository, tool Tool) DetectedTool {
    result := DetectedTool{Tool: tool}

    // Check go.mod
    goMod := filepath.Join(repo.Path, "go.mod")
    if fileExists(goMod) {
        content := readFile(goMod)
        if strings.Contains(content, tool.Name) {
            result.Detected = true
            result.DetectionPath = "go.mod"
        }
    }

    // Check for tool usage in code
    if tool.Name == "testing" {
        // Standard library - check for *_test.go files
        testFiles := findFiles(repo.Path, "*_test.go")
        if len(testFiles) > 0 {
            result.Detected = true
            result.DetectionPath = "test files"
        }
    }

    // Check CI/CD configs for tool usage
    if d.checkCICD(repo, tool.Name) {
        result.Detected = true
        result.DetectionPath = "CI/CD configuration"
    }

    return result
}
```

### 5.5 Adoption Scoring Algorithm

```go
func (d *ToolDetector) calculateAdoptionScore(analysis ToolAdoptionAnalysis) float64 {
    score := 0.0
    maxScore := 0.0

    // Weight tools by priority
    weights := map[string]float64{
        "essential":   10.0,  // Must-have tools
        "recommended": 5.0,   // Should-have tools
        "advanced":    3.0,   // Nice-to-have tools
        "optional":    1.0,   // Optional tools
        "alternative": 5.0,   // Alternative to standard
        "standard":    10.0,  // Language standard tools
    }

    for _, tool := range analysis.AvailableTools {
        weight := weights[tool.Priority]
        maxScore += weight

        // Check if tool is detected
        for _, detected := range analysis.DetectedTools {
            if detected.Tool.Name == tool.Name {
                points := weight

                // Bonus for proper configuration
                if detected.Configured {
                    points += weight * 0.2 // 20% bonus
                }

                // Bonus for CI/CD integration
                if d.isInCICD(detected) {
                    points += weight * 0.1 // 10% bonus
                }

                score += points
                break
            }
        }
    }

    // Normalize to 0-100
    if maxScore > 0 {
        return (score / maxScore) * 100
    }
    return 0
}

func (d *ToolDetector) calculateCategoryScores(analysis ToolAdoptionAnalysis) map[string]float64 {
    scores := make(map[string]float64)

    categories := []string{"Testing", "Quality", "Security", "Automation"}

    for _, category := range categories {
        categoryTools := d.getToolsByCategory(analysis.AvailableTools, category)
        detectedInCategory := d.getDetectedByCategory(analysis.DetectedTools, category)

        essential := 0
        essentialDetected := 0

        for _, tool := range categoryTools {
            if tool.Priority == "essential" || tool.Priority == "standard" {
                essential++

                for _, detected := range detectedInCategory {
                    if detected.Tool.Name == tool.Name {
                        essentialDetected++
                        break
                    }
                }
            }
        }

        if essential > 0 {
            scores[category] = (float64(essentialDetected) / float64(essential)) * 100
        } else {
            scores[category] = 100.0 // No essential tools in category
        }
    }

    return scores
}
```

### 5.6 Tool Recommendations

```go
type ToolRecommendation struct {
    Tool           Tool
    Priority       string
    Reason         string
    SetupGuide     string
    EstimatedEffort string
    ExpectedBenefit string
    ConfigExample  string
}

func (d *ToolDetector) generateRecommendations(analysis ToolAdoptionAnalysis) []ToolRecommendation {
    recommendations := []ToolRecommendation{}

    // Prioritize essential missing tools
    for _, tool := range analysis.MissingTools {
        if tool.Priority == "essential" || tool.Priority == "standard" {
            rec := ToolRecommendation{
                Tool:            tool,
                Priority:        "HIGH",
                Reason:          d.getReasonForTool(tool),
                SetupGuide:      d.getSetupGuide(tool, analysis.Language),
                EstimatedEffort: d.estimateEffort(tool),
                ExpectedBenefit: d.estimateBenefit(tool),
                ConfigExample:   d.getConfigExample(tool, analysis.Language),
            }
            recommendations = append(recommendations, rec)
        }
    }

    // Add recommended tools
    for _, tool := range analysis.MissingTools {
        if tool.Priority == "recommended" {
            rec := ToolRecommendation{
                Tool:            tool,
                Priority:        "MEDIUM",
                Reason:          d.getReasonForTool(tool),
                SetupGuide:      d.getSetupGuide(tool, analysis.Language),
                EstimatedEffort: d.estimateEffort(tool),
                ExpectedBenefit: d.estimateBenefit(tool),
            }
            recommendations = append(recommendations, rec)
        }
    }

    return recommendations
}

func (d *ToolDetector) getSetupGuide(tool Tool, language string) string {
    guides := map[string]map[string]string{
        "Python": {
            "pytest": `1. Install: pip install pytest pytest-cov
2. Create pytest.ini:
   [pytest]
   testpaths = tests
   python_files = test_*.py
   addopts = --cov=src --cov-report=term-missing
3. Run: pytest`,

            "black": `1. Install: pip install black
2. Create pyproject.toml:
   [tool.black]
   line-length = 100
   target-version = ['py311']
3. Run: black .`,

            "mypy": `1. Install: pip install mypy
2. Create mypy.ini:
   [mypy]
   python_version = 3.11
   strict = true
3. Run: mypy src/`,

            "pre-commit": `1. Install: pip install pre-commit
2. Create .pre-commit-config.yaml:
   repos:
     - repo: https://github.com/psf/black
       rev: 23.9.1
       hooks:
         - id: black
3. Run: pre-commit install`,
        },

        "JavaScript": {
            "jest": `1. Install: npm install --save-dev jest @types/jest
2. Create jest.config.js:
   module.exports = {
     coverageThreshold: {
       global: { lines: 80, branches: 80 }
     }
   };
3. Add to package.json: "test": "jest"`,

            "eslint": `1. Install: npm install --save-dev eslint
2. Initialize: npx eslint --init
3. Create .eslintrc.json with your rules
4. Add to package.json: "lint": "eslint ."`,

            "prettier": `1. Install: npm install --save-dev prettier
2. Create .prettierrc.json:
   {
     "semi": true,
     "singleQuote": true,
     "tabWidth": 2
   }
3. Add to package.json: "format": "prettier --write ."`,

            "husky": `1. Install: npm install --save-dev husky
2. Initialize: npx husky-init
3. Edit .husky/pre-commit to run tests/linters`,
        },

        "Go": {
            "golangci-lint": `1. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
2. Create .golangci.yml with your configuration
3. Run: golangci-lint run`,

            "gosec": `1. Install: go install github.com/securego/gosec/v2/cmd/gosec@latest
2. Run: gosec ./...
3. Add to CI/CD pipeline`,
        },
    }

    if langGuides, ok := guides[language]; ok {
        if guide, ok := langGuides[tool.Name]; ok {
            return guide
        }
    }

    return fmt.Sprintf("Install %s following official documentation", tool.Name)
}
```

### 5.7 Integration into Overall Scoring

**Tool Adoption Score contributes to overall quality score:**

```go
type DimensionScores struct {
    TestCoverage    CoverageScore        `json:"test_coverage"`       // 30%
    TestQuality     TestQualityScore     `json:"test_quality"`        // 25%
    ToolAdoption    ToolAdoptionScore    `json:"tool_adoption"`       // NEW: 10%
    TestPerformance PerformanceScore     `json:"test_performance"`    // 10%
    Maintainability MaintainabilityScore `json:"maintainability"`     // 15%
    CodeQuality     CodeQualityScore     `json:"code_quality"`        // 10%
}

type ToolAdoptionScore struct {
    Score              float64            `json:"score"`           // 0-100
    Grade              string             `json:"grade"`
    ToolsDetected      int                `json:"tools_detected"`
    ToolsAvailable     int                `json:"tools_available"`
    EssentialMissing   int                `json:"essential_missing"`
    CategoryBreakdown  map[string]float64 `json:"category_breakdown"`
    TopRecommendations []ToolRecommendation `json:"top_recommendations"`
}

func CalculateOverallScore(dimensions DimensionScores) float64 {
    // Updated weighting to include tool adoption
    weights := map[string]float64{
        "test_coverage":    0.30,  // 30%
        "test_quality":     0.25,  // 25%
        "tool_adoption":    0.10,  // 10% - NEW
        "test_performance": 0.10,  // 10%
        "maintainability":  0.15,  // 15%
        "code_quality":     0.10,  // 10%
    }

    overall := (dimensions.TestCoverage.Score * weights["test_coverage"]) +
               (dimensions.TestQuality.Score * weights["test_quality"]) +
               (dimensions.ToolAdoption.Score * weights["tool_adoption"]) +
               (dimensions.TestPerformance.Score * weights["test_performance"]) +
               (dimensions.Maintainability.Score * weights["maintainability"]) +
               (dimensions.CodeQuality.Score * weights["code_quality"])

    return overall
}
```

### 5.8 Report Example

**Tool Adoption Section in HTML Report:**

```html
<section class="tool-adoption">
  <h2>Tool Adoption Analysis</h2>

  <div class="score-card">
    <div class="score">78/100</div>
    <div class="grade">B</div>
  </div>

  <h3>Category Breakdown</h3>
  <div class="category-scores">
    <div class="category">
      <span class="name">Testing Tools</span>
      <div class="progress-bar">
        <div class="progress" style="width: 85%">85%</div>
      </div>
      <span class="tools">5/6 essential tools detected</span>
    </div>

    <div class="category">
      <span class="name">Quality Tools</span>
      <div class="progress-bar">
        <div class="progress" style="width: 70%">70%</div>
      </div>
      <span class="tools">4/6 essential tools detected</span>
    </div>

    <div class="category warning">
      <span class="name">Security Tools</span>
      <div class="progress-bar">
        <div class="progress" style="width: 40%">40%</div>
      </div>
      <span class="tools">2/5 essential tools detected ⚠️</span>
    </div>

    <div class="category">
      <span class="name">Automation Tools</span>
      <div class="progress-bar">
        <div class="progress" style="width: 90%">90%</div>
      </div>
      <span class="tools">3/3 essential tools detected</span>
    </div>
  </div>

  <h3>Top Recommendations (Free Tools)</h3>
  <div class="recommendations">
    <div class="recommendation priority-high">
      <h4>🔴 High Priority: Add Snyk or pip-audit</h4>
      <p><strong>Category:</strong> Security - Dependency Scanning</p>
      <p><strong>Why:</strong> No dependency vulnerability scanning detected</p>
      <p><strong>Effort:</strong> 15 minutes</p>
      <p><strong>Benefit:</strong> Detect vulnerable dependencies automatically</p>
      <details>
        <summary>Setup Guide</summary>
        <pre><code>1. Install: pip install pip-audit
2. Run: pip-audit
3. Add to CI/CD:
   - name: Security scan
     run: pip-audit --require-hashes
        </code></pre>
      </details>
    </div>

    <div class="recommendation priority-high">
      <h4>🔴 High Priority: Add pre-commit hooks</h4>
      <p><strong>Category:</strong> Automation</p>
      <p><strong>Why:</strong> No pre-commit hooks detected</p>
      <p><strong>Effort:</strong> 30 minutes</p>
      <p><strong>Benefit:</strong> Catch issues before commit, faster feedback</p>
      <details>
        <summary>Setup Guide</summary>
        <pre><code>1. Install: pip install pre-commit
2. Create .pre-commit-config.yaml
3. Run: pre-commit install
        </code></pre>
      </details>
    </div>

    <div class="recommendation priority-medium">
      <h4>🟡 Medium Priority: Add mutation testing with mutmut</h4>
      <p><strong>Category:</strong> Testing - Effectiveness</p>
      <p><strong>Why:</strong> No mutation testing detected</p>
      <p><strong>Effort:</strong> 1 hour</p>
      <p><strong>Benefit:</strong> Verify test quality, find weak tests</p>
    </div>
  </div>
</section>
```

### 5.9 CLI Output

```bash
$ shipshape analyze --verbose

Tool Adoption Analysis
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Overall Tool Adoption Score: 78/100 (Grade: B)

Category Breakdown:
  ✅ Testing Tools      85% (5/6 essential tools)
  ✅ Quality Tools      70% (4/6 essential tools)
  ⚠️  Security Tools    40% (2/5 essential tools)
  ✅ Automation Tools   90% (3/3 essential tools)

Detected Tools:
  ✓ pytest (configured: pytest.ini)
  ✓ coverage.py (configured: .coveragerc)
  ✓ black (configured: pyproject.toml)
  ✓ mypy (no config file)
  ✓ GitHub Actions (configured)
  ✓ pre-commit (configured: .pre-commit-config.yaml)

Missing Essential Tools:
  ✗ bandit - Python security linter
  ✗ safety - Dependency vulnerability scanner
  ✗ pylint - Code quality checker

Top Recommendations:
  🔴 HIGH: Add bandit for security scanning
     Setup: pip install bandit && bandit -r src/
     Effort: 10 minutes | Benefit: Detect security issues

  🔴 HIGH: Add safety for dependency scanning
     Setup: pip install safety && safety check
     Effort: 5 minutes | Benefit: Find vulnerable packages

  🟡 MEDIUM: Add pylint for code quality
     Setup: pip install pylint && pylint src/
     Effort: 30 minutes | Benefit: Catch code quality issues

Run 'shipshape tools --setup <tool-name>' for detailed setup instructions.
```

---

## 6. Implementation Roadmap

### Q1 2026: MVP
- Core test analysis engine (Go)
- Test discovery for Python, JavaScript, Go
- Coverage analysis integration
- **Tool detection engine (Python, JavaScript, Go)**
- **Language-specific tool catalogs**
- CLI with analyze/report commands
- HTML/JSON reporting
- GitHub Actions integration

### Q2 2026: Framework Deep Dive
- pytest deep analysis
- Jest/Vitest analysis
- Go testing patterns
- Test quality metrics
- Test smell detection
- Cypress/Playwright E2E analysis
- **Tool adoption scoring and recommendations**
- **Setup guides for missing tools**

### Q3 2026: Advanced Testing
- Mutation testing integration
- Flakiness detection
- Test performance profiling
- Quality gates for CI/CD
- Historical test tracking
- Test harness optimization
- **Expand tool catalog to Java, Rust, C#**
- **CI/CD integration detection**

### Q4 2026: Enterprise Features
- Multi-framework project support
- Test ROI analysis
- AI-powered test generation suggestions
- Web dashboard for test metrics
- Test analytics and insights
- **Tool adoption trend tracking**
- **Automated tool setup (one-click install)**

---

## Conclusion

Ship Shape represents a comprehensive testing quality platform that prioritizes test excellence over security scanning. By combining evidence-based testing metrics, framework-specific best practices, and deep test harness analysis, it addresses the critical need for testing quality assurance in modern software development.

**Key Innovations:**
1. Multi-dimensional testing quality scoring (coverage, quality, performance)
2. Framework-specific deep analysis (pytest, Jest, JUnit, Cypress, etc.)
3. Test pyramid/trophy compliance validation
4. Comprehensive test smell detection
5. Test harness optimization recommendations
6. CI/CD test workflow analysis and optimization
7. **Tool adoption analysis** - Gap analysis comparing tools in use vs. available free/opensource tools
8. **Actionable recommendations** - Setup guides with effort estimates for missing essential tools

The design leverages modern software engineering practices while focusing on what truly matters: effective, maintainable, and performant tests that provide confidence in code quality.

**Next Steps:**
1. Implement test discovery engine
2. Build framework-specific analyzers
3. Integrate coverage tools
4. Develop test quality metrics
5. Build tool catalog and detection engine
6. Create tool recommendation system with setup guides
7. Test against real-world repositories
8. Iterate based on developer feedback

---

**Document Version**: 3.0
**Last Updated**: 2026-01-27
**Status**: Design Proposal - Testing Focused with Tool Adoption Analysis
**Feedback**: Open for review and iteration

**Changelog:**
- v3.0: Added comprehensive Tool Adoption Analysis & Scoring section
- v2.0: Shifted focus from security to testing quality
- v1.0: Initial design with security focus
