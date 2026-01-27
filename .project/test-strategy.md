# Ship Shape - Test Strategy

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Active
**Owner**: Engineering Team

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [The Meta-Testing Challenge](#the-meta-testing-challenge)
3. [Test Philosophy](#test-philosophy)
4. [Test Pyramid](#test-pyramid)
5. [Ground Truth Establishment](#ground-truth-establishment)
6. [Test Levels](#test-levels)
7. [Accuracy Validation](#accuracy-validation)
8. [Meta-Validation Approach](#meta-validation-approach)
9. [Continuous Validation Plan](#continuous-validation-plan)
10. [Dogfooding Strategy](#dogfooding-strategy)
11. [Test Data Management](#test-data-management)
12. [Quality Metrics](#quality-metrics)
13. [Tools and Infrastructure](#tools-and-infrastructure)

---

## Executive Summary

### The Challenge
Ship Shape is a **testing quality analysis tool** - we analyze tests, detect test smells, and assess code coverage. This creates a unique meta-problem: **How do we validate that our testing tool is correct when it's designed to validate testing?**

### Our Approach
We employ a **multi-layered validation strategy**:

1. **Traditional Testing**: Comprehensive unit, integration, and performance tests
2. **Ground Truth Validation**: Known-good and known-bad test examples with manual verification
3. **Academic Validation**: Comparison against peer-reviewed test smell research (tsDetect, etc.)
4. **Real-World Validation**: Testing on popular open-source projects with manual code review
5. **Dogfooding**: Ship Shape must analyze itself and achieve high scores
6. **Continuous Validation**: Ongoing accuracy monitoring and improvement

### Success Criteria
- **Test Coverage**: >90% unit, >80% integration
- **Accuracy**: >90% precision and recall on ground truth datasets
- **False Positive Rate**: <10% across all detectors
- **Self-Analysis**: Ship Shape scores ≥85/100 when analyzing itself
- **Real-World Validation**: Manual review confirms >85% of findings in OSS projects

---

## The Meta-Testing Challenge

### The Paradox
```
Ship Shape analyzes test quality
    ↓
Ship Shape itself has tests
    ↓
Who validates Ship Shape's tests?
    ↓
Ship Shape should validate its own tests
    ↓
But can we trust Ship Shape if it hasn't been validated?
```

### Breaking the Circular Dependency

**Phase 1: Bootstrap with External Validation**
- Use academic research as ground truth (tsDetect, SniffTest)
- Manual code review by testing experts
- Comparison with industry-standard tools
- Peer review of test patterns

**Phase 2: Establish Confidence**
- Build comprehensive ground truth datasets
- Validate accuracy metrics (precision >90%, recall >90%)
- Achieve low false positive rates (<10%)
- Demonstrate real-world effectiveness

**Phase 3: Self-Validation (Dogfooding)**
- Only after Phase 1 & 2, use Ship Shape on itself
- Continuously monitor self-analysis scores
- Any drop in self-score triggers investigation
- Treat Ship Shape's codebase as exemplar

### The Trust Chain

```
Academic Research (tsDetect papers, etc.)
    ↓ validates
Ground Truth Datasets (manually verified)
    ↓ trains/validates
Ship Shape Detectors (v1.0)
    ↓ analyzes
Ship Shape's Own Tests
    ↓ proves quality
Ship Shape is Trustworthy
    ↓ can analyze
User Codebases
```

---

## Test Philosophy

### Core Principles

1. **Test-Driven Development (TDD)**
   - Write tests before implementation
   - Red-Green-Refactor cycle
   - Tests as living documentation

2. **Exemplar Testing**
   - Ship Shape's tests must exemplify best practices
   - No test smells in our own tests
   - High coverage on critical paths
   - Clear, maintainable test code

3. **Evidence-Based Validation**
   - All thresholds backed by research
   - Manual verification of samples
   - Statistical significance in measurements
   - Reproducible validation processes

4. **Continuous Validation**
   - Tests run on every commit
   - Accuracy monitoring in CI/CD
   - Regular benchmark updates
   - Community feedback integration

5. **Defense in Depth**
   - Multiple validation layers
   - Independent verification methods
   - Cross-validation across datasets
   - Peer review of critical detectors

---

## Test Pyramid

### Ship Shape Test Pyramid

```
                 /\
                /  \
               /E2E \
              /------\
             /        \
            /Validation\
           /------------\
          /              \
         /  Integration   \
        /------------------\
       /                    \
      /    Unit Tests        \
     /------------------------\
```

**Distribution**:
- **Unit Tests**: 70% (fast, isolated, comprehensive)
- **Integration Tests**: 20% (workflow validation)
- **Validation Tests**: 5% (accuracy on ground truth)
- **End-to-End Tests**: 5% (full workflow on real repos)

### Inverted Pyramid Warning

We deliberately **avoid** the anti-pattern:
```
     ❌ Anti-Pattern
        /\
       /E2E\ (slow, brittle)
      /----\
     /Unit \
    /------\
```

---

## Ground Truth Establishment

### Ground Truth Dataset Structure

```
testdata/ground-truth/
├── test-smells/
│   ├── mystery-guest/
│   │   ├── positive/           # Known smells (should detect)
│   │   │   ├── go/
│   │   │   │   ├── file-access-no-setup.go
│   │   │   │   ├── network-call-no-mock.go
│   │   │   │   └── metadata.yml
│   │   │   ├── python/
│   │   │   └── javascript/
│   │   └── negative/           # Clean tests (should NOT detect)
│   │       ├── go/
│   │       ├── python/
│   │       └── javascript/
│   ├── assertion-roulette/
│   ├── eager-test/
│   └── ... (all 11 smell types)
├── coverage/
│   ├── formats/
│   │   ├── cobertura/
│   │   │   ├── sample-valid.xml
│   │   │   ├── sample-large.xml
│   │   │   └── expected-output.json
│   │   ├── lcov/
│   │   ├── coverage-py-json/
│   │   └── go-profile/
├── test-patterns/
│   ├── table-driven/
│   ├── parametrized/
│   ├── fixtures/
│   └── mocking/
└── real-world-samples/
    ├── kubernetes/           # Curated from OSS projects
    ├── django/
    ├── react/
    └── ...
```

### Ground Truth Curation Process

1. **Academic Baseline**
   - Use tsDetect's validated test smell examples
   - Reference SniffTest's benchmark suite
   - Incorporate examples from test smell papers

2. **Manual Verification**
   - Each ground truth example reviewed by 2+ engineers
   - Testing experts validate smell classifications
   - Document reasoning in metadata.yml

3. **Metadata Schema**
   ```yaml
   # testdata/ground-truth/test-smells/mystery-guest/positive/go/metadata.yml
   file: file-access-no-setup.go
   smell_type: mystery_guest
   classification: positive
   severity: high
   confidence: 0.95

   description: |
     Test opens a file from /etc without creating test fixture.
     This is a clear Mystery Guest smell.

   detection_rules:
     - pattern: "os.Open"
     - no_setup_detected: true
     - no_fixture: true

   verified_by:
     - name: "Jane Doe"
       date: "2026-01-15"
       role: "Senior Test Engineer"
     - name: "John Smith"
       date: "2026-01-16"
       role: "Testing Expert"

   references:
     - "tsDetect: https://testsmells.org/"
     - "Test Smells Paper: https://..."

   expected_detection:
     should_detect: true
     expected_confidence: ">= 0.85"
     expected_severity: "high"
   ```

4. **Continuous Curation**
   - Add new examples as edge cases discovered
   - Update based on community feedback
   - Quarterly review of ground truth accuracy

### Real-World Sampling

**Selection Criteria** for OSS projects:
- Popular projects (>5k stars)
- Diverse languages and frameworks
- Active maintenance (commits in last 3 months)
- Known for good/bad test practices
- Representative of different domains

**Sampling Strategy**:
1. Select 10 files per project (mix of good/bad)
2. Manual analysis by testing expert
3. Document all findings
4. Ship Shape analysis on same files
5. Compare results (precision/recall)

**Projects for Validation**:

| Language   | Project     | Stars | Test Quality | Purpose |
|------------|-------------|-------|--------------|---------|
| Go         | kubernetes  | 110k  | High         | Baseline for good tests |
| Go         | docker      | 68k   | High         | Production-grade |
| Python     | django      | 79k   | High         | Framework tests |
| Python     | requests    | 52k   | Medium       | Mixed quality |
| JavaScript | react       | 228k  | High         | Modern practices |
| JavaScript | vue         | 207k  | High         | Framework tests |
| TypeScript | vscode      | 163k  | High         | Large codebase |

---

## Test Levels

### 1. Unit Tests (70% of tests)

**Scope**: Individual functions, methods, and components

**What We Test**:
- ✅ Each analyzer independently (Go, Python, JS)
- ✅ Each smell detector individually
- ✅ Coverage parser for each format
- ✅ Utility functions and helpers
- ✅ Edge cases and error conditions

**Example Test Structure**:
```go
// internal/analyzer/goanalyzer/go_analyzer_test.go
package goanalyzer_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGoAnalyzer_DetectTableDrivenTest(t *testing.T) {
    tests := []struct {
        name           string
        inputFile      string
        wantDetected   bool
        wantConfidence float64
    }{
        {
            name:           "valid table-driven test",
            inputFile:      "testdata/table-driven-valid.go",
            wantDetected:   true,
            wantConfidence: 0.95,
        },
        {
            name:           "false positive - similar pattern",
            inputFile:      "testdata/table-driven-false-positive.go",
            wantDetected:   false,
            wantConfidence: 0.0,
        },
        {
            name:           "edge case - map-based table",
            inputFile:      "testdata/table-driven-map.go",
            wantDetected:   false, // Currently unsupported
            wantConfidence: 0.0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            analyzer := goanalyzer.NewGoTestAnalyzer()
            file := loadTestFile(t, tt.inputFile)

            result, err := analyzer.Analyze(context.Background(), &AnalysisRequest{
                File: file,
            })

            require.NoError(t, err)

            hasTableDriven := false
            for _, testFunc := range result.TestFunctions {
                if testFunc.IsTable {
                    hasTableDriven = true
                    break
                }
            }

            assert.Equal(t, tt.wantDetected, hasTableDriven)
        })
    }
}

// Benchmark performance
func BenchmarkGoAnalyzer_100Files(b *testing.B) {
    analyzer := goanalyzer.NewGoTestAnalyzer()
    files := loadBenchmarkFiles(b, 100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, file := range files {
            _, err := analyzer.Analyze(context.Background(), &AnalysisRequest{File: file})
            if err != nil {
                b.Fatalf("analysis failed: %v", err)
            }
        }
    }
}
```

**Coverage Requirements**:
- **Critical Paths**: 100% coverage
- **Public APIs**: 100% coverage
- **Overall**: >90% coverage
- **Mutation Testing**: >80% mutation score (future)

### 2. Integration Tests (20% of tests)

**Scope**: Component interactions and workflows

**What We Test**:
- ✅ Full analysis pipeline (discovery → analysis → scoring → report)
- ✅ Multi-language repository analysis
- ✅ Monorepo package coordination
- ✅ GitHub API integrations
- ✅ Coverage parsing → assessment workflow

**Example Integration Test**:
```go
// test/integration/full_analysis_test.go
package integration_test

func TestFullAnalysis_MultiLanguageRepo(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup: Create test repository
    repoPath := setupTestRepo(t, "testdata/repos/multi-language")
    defer cleanupTestRepo(t, repoPath)

    // Execute: Run full Ship Shape analysis
    cli := NewShipShapeCLI(t)
    output, err := cli.Run("analyze", repoPath, "--format", "json")
    require.NoError(t, err)

    // Verify: Parse and validate output
    var report Report
    err = json.Unmarshal(output, &report)
    require.NoError(t, err)

    // Assertions
    assert.Equal(t, 3, len(report.Languages)) // Go, Python, JavaScript
    assert.GreaterOrEqual(t, report.Score, 70)
    assert.Contains(t, report.DimensionScores, "test-quality")
    assert.Contains(t, report.DimensionScores, "coverage")

    // Verify Go analysis
    goResults := findLanguageResults(report, "go")
    assert.NotNil(t, goResults)
    assert.Greater(t, len(goResults.TestFunctions), 0)

    // Verify Python analysis
    pyResults := findLanguageResults(report, "python")
    assert.NotNil(t, pyResults)
    assert.Greater(t, len(pyResults.Fixtures), 0)
}

func TestMonorepoAnalysis_ParallelExecution(t *testing.T) {
    repoPath := setupTestRepo(t, "testdata/repos/monorepo-20-packages")
    defer cleanupTestRepo(t, repoPath)

    startTime := time.Now()

    cli := NewShipShapeCLI(t)
    output, err := cli.Run("analyze", repoPath, "--parallel", "4")
    require.NoError(t, err)

    duration := time.Since(startTime)

    var report Report
    err = json.Unmarshal(output, &report)
    require.NoError(t, err)

    // Verify all packages analyzed
    assert.Equal(t, 20, len(report.Packages))

    // Verify performance (should be <30s for 20 packages)
    assert.Less(t, duration, 30*time.Second, "Parallel analysis too slow")
}
```

### 3. Validation Tests (5% of tests)

**Scope**: Accuracy validation against ground truth

**What We Test**:
- ✅ Precision and recall on ground truth datasets
- ✅ False positive/negative rates
- ✅ Confidence score calibration
- ✅ Cross-validation across datasets

**Example Validation Test**:
```go
// test/validation/smell_detection_accuracy_test.go
package validation_test

func TestSmellDetection_MysteryGuest_Accuracy(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping validation test in short mode")
    }

    // Load ground truth dataset
    groundTruth := loadGroundTruth(t, "test-smells/mystery-guest")

    detector := smell.NewMysteryGuestDetector()

    var truePositives, falsePositives, trueNegatives, falseNegatives int

    // Test positive examples (should detect)
    for _, example := range groundTruth.Positive {
        smells, err := detector.Detect(example.TestFile)
        require.NoError(t, err)

        detected := len(smells) > 0
        if detected {
            truePositives++
        } else {
            falseNegatives++
            t.Logf("FALSE NEGATIVE: %s - %s", example.File, example.Description)
        }
    }

    // Test negative examples (should NOT detect)
    for _, example := range groundTruth.Negative {
        smells, err := detector.Detect(example.TestFile)
        require.NoError(t, err)

        detected := len(smells) > 0
        if detected {
            falsePositives++
            t.Logf("FALSE POSITIVE: %s - Detected: %v", example.File, smells)
        } else {
            trueNegatives++
        }
    }

    // Calculate metrics
    precision := float64(truePositives) / float64(truePositives+falsePositives)
    recall := float64(truePositives) / float64(truePositives+falseNegatives)
    f1Score := 2 * (precision * recall) / (precision + recall)
    falsePositiveRate := float64(falsePositives) / float64(falsePositives+trueNegatives)

    // Log results
    t.Logf("Mystery Guest Detection Accuracy:")
    t.Logf("  True Positives: %d", truePositives)
    t.Logf("  False Positives: %d", falsePositives)
    t.Logf("  True Negatives: %d", trueNegatives)
    t.Logf("  False Negatives: %d", falseNegatives)
    t.Logf("  Precision: %.2f%%", precision*100)
    t.Logf("  Recall: %.2f%%", recall*100)
    t.Logf("  F1 Score: %.2f", f1Score)
    t.Logf("  False Positive Rate: %.2f%%", falsePositiveRate*100)

    // Assert minimum thresholds
    assert.GreaterOrEqual(t, precision, 0.90, "Precision below threshold")
    assert.GreaterOrEqual(t, recall, 0.90, "Recall below threshold")
    assert.LessOrEqual(t, falsePositiveRate, 0.10, "False positive rate too high")
}

// Test all smell detectors
func TestAllSmellDetectors_AccuracyBaseline(t *testing.T) {
    smellTypes := []string{
        "mystery-guest",
        "assertion-roulette",
        "eager-test",
        "lazy-test",
        "conditional-logic",
        "general-fixture",
        "obscure-test",
        "sensitive-equality",
        "resource-optimism",
        "code-duplication",
        "flakiness",
    }

    results := make(map[string]*AccuracyMetrics)

    for _, smellType := range smellTypes {
        t.Run(smellType, func(t *testing.T) {
            metrics := runAccuracyTest(t, smellType)
            results[smellType] = metrics

            // Each smell type must meet minimum threshold
            assert.GreaterOrEqual(t, metrics.Precision, 0.85,
                "Smell %s: precision below minimum", smellType)
            assert.GreaterOrEqual(t, metrics.Recall, 0.85,
                "Smell %s: recall below minimum", smellType)
        })
    }

    // Generate accuracy report
    generateAccuracyReport(t, results)
}
```

### 4. End-to-End Tests (5% of tests)

**Scope**: Full workflows on real repositories

**What We Test**:
- ✅ CLI workflows
- ✅ GitHub Actions integration
- ✅ Report generation
- ✅ Real OSS project analysis

**Example E2E Test**:
```go
// test/e2e/github_actions_workflow_test.go
package e2e_test

func TestGitHubActionsWorkflow_RealRepo(t *testing.T) {
    if os.Getenv("E2E_TESTS") != "1" {
        t.Skip("Skipping E2E test (set E2E_TESTS=1 to run)")
    }

    // Clone a real repository
    repoURL := "https://github.com/testorg/sample-repo"
    repoPath := cloneRepo(t, repoURL)
    defer os.RemoveAll(repoPath)

    // Simulate GitHub Actions environment
    setGitHubEnv(t, map[string]string{
        "GITHUB_TOKEN":      os.Getenv("TEST_GITHUB_TOKEN"),
        "GITHUB_REPOSITORY": "testorg/sample-repo",
        "GITHUB_SHA":        "abc123",
        "GITHUB_REF":        "refs/pull/42/merge",
    })

    // Run Ship Shape as GitHub Action would
    action := NewGitHubAction(t)
    err := action.Run(map[string]string{
        "config":     ".shipshape.yml",
        "fail-on":    "high",
        "comment-pr": "true",
    })
    require.NoError(t, err)

    // Verify outputs
    assert.FileExists(t, "shipshape-report.html")

    // Verify GitHub API calls were made
    assert.True(t, action.CheckRunCreated())
    assert.True(t, action.PRCommentPosted())
}
```

---

## Accuracy Validation

### Metrics Tracking

**For Each Detector** (test smells, patterns, etc.):

```go
type AccuracyMetrics struct {
    TruePositives  int
    FalsePositives int
    TrueNegatives  int
    FalseNegatives int

    // Calculated
    Precision          float64  // TP / (TP + FP)
    Recall             float64  // TP / (TP + FN)
    F1Score            float64  // 2 * (P * R) / (P + R)
    Accuracy           float64  // (TP + TN) / (TP + FP + TN + FN)
    FalsePositiveRate  float64  // FP / (FP + TN)
    FalseNegativeRate  float64  // FN / (FN + TP)

    // Confidence calibration
    MeanConfidence     float64
    ConfidenceStdDev   float64
}
```

**Minimum Thresholds**:
- Precision: ≥90% (critical), ≥85% (acceptable)
- Recall: ≥90% (critical), ≥85% (acceptable)
- F1 Score: ≥0.88
- False Positive Rate: ≤10%

### Confusion Matrix Tracking

```
                 Predicted
                 Positive  Negative
Actual Positive    TP        FN
       Negative    FP        TN
```

**Example for Mystery Guest**:
```
Tested 100 examples:
- 50 positive (actual smells)
- 50 negative (clean code)

Results:
                 Detected  Not Detected
Actual Smell        47          3         (Recall: 94%)
Clean Code           4         46         (Precision: 92%)

Precision: 47/(47+4) = 92%
Recall: 47/(47+3) = 94%
F1: 2*(0.92*0.94)/(0.92+0.94) = 0.93
```

### Continuous Accuracy Monitoring

**CI/CD Integration**:
```yaml
# .github/workflows/accuracy-validation.yml
name: Accuracy Validation

on:
  pull_request:
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Validation Tests
        run: make test-validation

      - name: Generate Accuracy Report
        run: |
          go test -v ./test/validation/... -json > accuracy-results.json
          go run ./cmd/accuracy-report accuracy-results.json

      - name: Upload Accuracy Report
        uses: actions/upload-artifact@v3
        with:
          name: accuracy-report
          path: accuracy-report.html

      - name: Check Thresholds
        run: |
          go run ./cmd/accuracy-check accuracy-results.json \
            --min-precision 0.90 \
            --min-recall 0.90 \
            --max-false-positive-rate 0.10
```

---

## Meta-Validation Approach

### Self-Analysis (Dogfooding)

**Ship Shape Must Analyze Itself**

**Weekly Self-Analysis**:
```bash
# Run Ship Shape on its own codebase
shipshape analyze . --output shipshape-self-analysis.html

# CI/CD enforcement
shipshape analyze . --fail-on medium --min-score 85
```

**Self-Analysis Requirements**:
- **Overall Score**: ≥85/100
- **Test Quality**: ≥90/100
- **Coverage**: ≥90%
- **Zero High/Critical Test Smells**
- **All analyzers pass their own analysis**

**Self-Analysis CI Check**:
```yaml
# .github/workflows/dogfooding.yml
name: Self-Analysis (Dogfooding)

on:
  pull_request:
  push:
    branches: [main]

jobs:
  self-analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build Ship Shape
        run: make build

      - name: Analyze Ship Shape with Ship Shape
        run: |
          ./bin/shipshape analyze . \
            --config .shipshape.yml \
            --output shipshape-self-analysis.html \
            --fail-on high \
            --min-score 85

      - name: Upload Self-Analysis Report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: self-analysis-report
          path: shipshape-self-analysis.html

      - name: Comment on PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            // Post self-analysis results to PR
```

**What We Validate**:
1. **Our Own Tests**
   - No test smells detected in Ship Shape's tests
   - High coverage on critical paths
   - Best practices demonstrated

2. **Our Analyzers**
   - Go analyzer analyzes its own code correctly
   - Python analyzer (if any Python) works
   - Meta-patterns detected

3. **Our Scoring**
   - Ship Shape achieves high scores
   - Dimension scores balanced
   - No gaming the system

### Comparison with Academic Tools

**tsDetect Benchmark** (for Java, if we add Java support):
```go
func TestSmellDetection_CompareWithTsDetect(t *testing.T) {
    // Use tsDetect's validated test suite
    tsDetectDataset := loadTsDetectBenchmark(t)

    for _, example := range tsDetectDataset {
        // Run Ship Shape
        shipShapeResults := analyzeWithShipShape(t, example)

        // Compare with tsDetect findings
        agreement := compareFindings(shipShapeResults, example.TsDetectFindings)

        // We should agree on >90% of smells
        assert.GreaterOrEqual(t, agreement, 0.90,
            "Disagreement with tsDetect on %s", example.File)
    }
}
```

### Expert Review Process

**Monthly Expert Review**:
1. Select 20 random findings from production usage
2. Testing expert manually reviews each
3. Classify as True Positive or False Positive
4. Track false positive rate
5. Update patterns if FP rate >10%

**Expert Review Template**:
```markdown
## Finding Review

**File**: src/services/user_service_test.go:42
**Smell**: Mystery Guest
**Confidence**: 0.87
**Ship Shape Description**: "Test accesses file system without setup"

### Manual Analysis

**Is this a true positive?** [ ] Yes [ ] No

**Reasoning**:
[Expert explanation]

**Should we adjust detection?** [ ] Yes [ ] No

**Suggested Improvement**:
[If applicable]

**Reviewed By**: Jane Doe
**Date**: 2026-01-27
```

---

## Continuous Validation Plan

### Daily Validation

**Every Commit**:
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Test coverage ≥90%
- [ ] No linting errors
- [ ] Build succeeds

### Weekly Validation

**Every Sunday**:
- [ ] Full validation test suite
- [ ] Accuracy metrics calculated
- [ ] Self-analysis (dogfooding)
- [ ] Benchmark performance tests
- [ ] Ground truth validation

**Automated Report**:
```
Ship Shape Weekly Validation Report
Date: 2026-01-27

Test Coverage:
  Unit: 92.3% ✅
  Integration: 84.1% ✅
  Overall: 91.2% ✅

Accuracy Metrics:
  Mystery Guest:       Precision: 94% ✅  Recall: 92% ✅  FP Rate: 6% ✅
  Assertion Roulette:  Precision: 97% ✅  Recall: 95% ✅  FP Rate: 3% ✅
  Eager Test:          Precision: 89% ✅  Recall: 87% ⚠️   FP Rate: 11% ❌
  [... all smell types]

Self-Analysis (Dogfooding):
  Overall Score: 87/100 ✅
  Test Quality: 92/100 ✅
  Coverage: 91% ✅
  Critical Smells: 0 ✅

Performance Benchmarks:
  Go Analyzer (100 files): 1.8s ✅ (target: <2s)
  Python Analyzer (50 files): 4.2s ✅ (target: <5s)
  JS Analyzer (100 files): 8.9s ✅ (target: <10s)

⚠️ Action Items:
  - Eager Test: FP rate 11%, investigate pattern refinement
  - Recall for Eager Test below 90%, add more training examples
```

### Monthly Validation

**First Monday of Month**:
- [ ] Expert review of 20 random findings
- [ ] Real-world OSS project analysis
- [ ] Manual comparison with production findings
- [ ] Community feedback review
- [ ] Ground truth dataset update

### Quarterly Validation

**Every Quarter**:
- [ ] Comprehensive benchmark update
- [ ] Accuracy revalidation on all datasets
- [ ] Comparison with latest academic research
- [ ] Performance regression analysis
- [ ] Security audit

---

## Dogfooding Strategy

### Ship Shape Eats Its Own Dog Food

**Principle**: Ship Shape must analyze itself and achieve exemplary scores.

### Self-Analysis Enforcement

**Pre-Commit Hook**:
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running Ship Shape self-analysis..."

./bin/shipshape analyze . \
  --config .shipshape.yml \
  --format text \
  --fail-on high \
  --min-score 80

if [ $? -ne 0 ]; then
  echo "❌ Ship Shape self-analysis failed!"
  echo "Fix the issues before committing."
  exit 1
fi

echo "✅ Self-analysis passed"
```

**CI/CD Gate**:
```yaml
# Required check before merge
- name: Self-Analysis Quality Gate
  run: |
    ./bin/shipshape analyze . --min-score 85 || \
      (echo "Self-analysis failed. Ship Shape must achieve 85+ when analyzing itself." && exit 1)
```

### Self-Improvement Loop

```
Ship Shape Analyzes Itself
    ↓ finds issues
Fix Test Smells in Our Tests
    ↓ improves
Ship Shape's Tests Get Better
    ↓ validates
Detectors Become More Accurate
    ↓ analyses
Ship Shape Analyzes Itself Again
    ↓ scores higher
Continuous Improvement
```

### Example Self-Analysis Goals

**v1.0.0 Self-Analysis Targets**:
- Overall Score: ≥85/100
- Test Quality: ≥90/100
- Coverage Metrics: ≥90/100
- Test Smells: 0 critical, ≤5 medium
- Tool Adoption: ≥90/100
- CI/CD Analysis: ≥85/100

**Tracking**:
```
Version   Date        Score   Test Quality   Coverage   Smells (C/H/M)
-------   ----------  -----   ------------   --------   --------------
v0.1.0    2026-02-10   72          78           88          2/3/12
v0.2.0    2026-03-15   78          84           90          1/2/8
v0.5.0    2026-04-20   83          89           91          0/1/4
v1.0.0    2026-06-01   87          92           92          0/0/3
```

---

## Test Data Management

### Test Fixtures Organization

```
testdata/
├── ground-truth/              # Ground truth datasets
│   ├── test-smells/
│   ├── coverage/
│   └── test-patterns/
├── integration/               # Integration test repos
│   ├── single-language/
│   ├── multi-language/
│   └── monorepos/
├── validation/                # Validation datasets
│   ├── oss-samples/
│   └── synthetic/
├── benchmarks/                # Performance benchmarks
│   ├── small/  (10 files)
│   ├── medium/ (100 files)
│   └── large/  (1000 files)
└── regression/                # Regression test cases
    ├── bug-fixes/
    └── edge-cases/
```

### Fixture Management Best Practices

1. **Version Control**: All fixtures in Git
2. **Minimal Size**: Keep fixtures small and focused
3. **Documentation**: Each fixture has metadata.yml
4. **Realistic**: Based on real-world code
5. **Diverse**: Cover edge cases and variations

### Synthetic Test Generation

**For Coverage**:
```go
// Generate synthetic test files for benchmarking
func generateSyntheticTests(count int, language string) []TestFile {
    // Create realistic but synthetic test code
    // Used for performance benchmarking only
}
```

---

## Quality Metrics

### Test Quality Metrics (For Our Tests)

**Code Coverage**:
- Line Coverage: >90%
- Branch Coverage: >85%
- Function Coverage: >95%

**Test Quality**:
- No test smells in our tests
- All tests independent (can run in any order)
- Fast (<1s for 90% of unit tests)
- Deterministic (no flaky tests)

**Maintainability**:
- Clear test names (Given-When-Then or similar)
- Minimal test code duplication
- Good use of fixtures and helpers
- Well-documented edge cases

### Detector Quality Metrics

**For Each Smell Detector**:
- Precision: ≥90%
- Recall: ≥90%
- F1 Score: ≥0.88
- False Positive Rate: ≤10%
- Performance: <5s for 100 files

### System Quality Metrics

**Performance**:
- Small repos (<100 files): <10s
- Medium repos (100-1000 files): <30s
- Large repos (1000+ files): <60s

**Reliability**:
- Crash rate: <0.1%
- Error recovery: 100%
- Graceful degradation: Yes

---

## Tools and Infrastructure

### Testing Tools

**Go Testing**:
- `go test` - Built-in test runner
- `testify` - Assertions and mocking
- `golangci-lint` - Linting (including test smells)
- `go-mutesting` - Mutation testing (future)

**Coverage Tools**:
- `go tool cover` - Coverage reporting
- `gocov` - Coverage conversion
- `coverage` - Coverage badges

**Benchmarking**:
- `go test -bench` - Built-in benchmarking
- `benchstat` - Benchmark comparison
- Custom benchmark harness

**CI/CD**:
- GitHub Actions - Continuous testing
- Pre-commit hooks - Local validation
- Codecov - Coverage tracking

### Test Infrastructure

**Docker for E2E**:
```dockerfile
# test/e2e/Dockerfile
FROM golang:1.21

# Install testing dependencies
RUN apt-get update && apt-get install -y \
    git \
    python3 \
    nodejs \
    npm

# Setup test environment
COPY . /app
WORKDIR /app

CMD ["make", "test-e2e"]
```

**Test Database** (for historical tracking):
```sql
CREATE TABLE accuracy_metrics (
    id SERIAL PRIMARY KEY,
    detector_name VARCHAR(100),
    date DATE,
    precision DECIMAL(5,2),
    recall DECIMAL(5,2),
    f1_score DECIMAL(5,2),
    false_positive_rate DECIMAL(5,2),
    dataset_version VARCHAR(50)
);

CREATE TABLE self_analysis (
    id SERIAL PRIMARY KEY,
    version VARCHAR(20),
    date DATE,
    overall_score INT,
    test_quality_score INT,
    coverage_score INT,
    critical_smells INT,
    high_smells INT,
    medium_smells INT
);
```

---

## Appendix: Test Examples

### Example Unit Test (Table-Driven)

```go
func TestCoverageParser_Cobertura(t *testing.T) {
    tests := []struct {
        name           string
        inputFile      string
        wantLineRate   float64
        wantBranchRate float64
        wantError      bool
    }{
        {
            name:           "valid cobertura",
            inputFile:      "testdata/coverage/cobertura/valid.xml",
            wantLineRate:   0.85,
            wantBranchRate: 0.75,
            wantError:      false,
        },
        {
            name:           "malformed XML",
            inputFile:      "testdata/coverage/cobertura/malformed.xml",
            wantError:      true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewCoberturaParser()
            content, _ := os.ReadFile(tt.inputFile)

            report, err := parser.Parse(bytes.NewReader(content))

            if tt.wantError {
                assert.Error(t, err)
                return
            }

            require.NoError(t, err)
            assert.InDelta(t, tt.wantLineRate, report.OverallMetrics.LineRate, 0.01)
            assert.InDelta(t, tt.wantBranchRate, report.OverallMetrics.BranchRate, 0.01)
        })
    }
}
```

---

**Document Status**: Active - In Use
**Review Cycle**: Quarterly
**Next Review**: 2026-04-27
**Owner**: QA Lead + Engineering Team

---

**Last Updated**: 2026-01-27
**Version**: 1.0.0
**Author**: Senior Software Engineer
