# SPIKE-003: Test Smell Detection Framework

## Overview
This spike validates the technical approach for detecting test smells and anti-patterns across multiple languages, leveraging existing research and open-source tools while building a unified detection framework for Ship Shape.

**Associated User Stories**: SS-030
**Risk Level**: HIGH
**Priority**: P0 (Critical)
**Target Completion**: Week 3-4 of implementation

## Problem Statement
Test smells are design flaws that affect test code quality, maintainability, and effectiveness. Ship Shape needs to detect 11+ test smell types across Python, JavaScript/TypeScript, Go, and Java with high accuracy (>90%) and provide actionable remediation guidance.

**Key Challenges**:
1. Multi-language detection (different patterns per language)
2. High accuracy requirements (low false positives)
3. Performance constraints (<5s for 100 test files)
4. Actionable remediation guidance
5. Extensibility for new smell types

## Existing Tools and Research

### Academic Research
- **tsDetect** ([GitHub](https://github.com/TestSmells/TSDetect), [Paper](https://testsmells.org/assets/publications/FSE2020_TechnicalPaper.pdf)): Open-source test smell detection tool for Java
  - Detects 19 test smell types
  - 96% precision, 97% recall
  - Available as IntelliJ plugin and CLI
  - Java-only limitation

- **SniffTest** ([GitHub](https://github.com/MaierFlorian/SniffTest)): JUnit test smell detector
  - Detects 5 test smell types
  - 87-100% accuracy
  - Java/JUnit specific

### Key Learnings from Research
- Pattern-based detection using AST works well (tsDetect proven)
- Detection rules can be formalized and reused
- Most test smells are language-agnostic conceptually
- Implementation patterns vary by language/framework

## Spike Objectives
- [ ] Evaluate tsDetect detection rules for reuse
- [ ] Design language-agnostic smell detection framework
- [ ] Prototype 11 core test smell detectors
- [ ] Validate detection accuracy (>90% precision/recall)
- [ ] Benchmark performance across languages
- [ ] Create remediation guidance database
- [ ] Test with real-world codebases

## Test Smell Catalog

### Priority 1: Critical Smells (High Impact)
1. **Mystery Guest**: Test depends on external resources without clear setup
2. **Resource Optimism**: Test assumes resources (files, network) are always available
3. **Flaky Test**: Non-deterministic test behavior
4. **Assertion Roulette**: Multiple assertions without messages
5. **Conditional Logic**: If/else or loops in test code

### Priority 2: High Impact Smells
6. **Eager Test**: Single test verifies multiple methods/behaviors
7. **Lazy Test**: Multiple tests for same production method
8. **General Fixture**: Shared fixture used by few tests
9. **Obscure Test**: Unclear test intent or poor naming
10. **Sensitive Equality**: Using toString() or == for object comparison

### Priority 3: Maintainability Smells
11. **Code Duplication**: Repeated test code without helpers
12. **Long Test**: Excessive lines of code or assertions
13. **Magic Number**: Hard-coded values without explanation

## Technical Investigation Areas

### 1. Language-Agnostic Detection Framework

**Design Pattern**: Strategy pattern with language-specific implementations

```go
// Core abstraction
type TestSmellDetector interface {
    // Detect smells in a test file
    Detect(testFile *TestFile) ([]*TestSmell, error)

    // Get smell type this detector handles
    SmellType() TestSmellType

    // Get supported languages
    SupportedLanguages() []string

    // Get detection confidence level
    Confidence() float64
}

type TestSmell struct {
    Type         TestSmellType
    Severity     Severity
    Location     SourceLocation
    Description  string
    Remediation  string
    Confidence   float64
    Evidence     *Evidence
}

type TestSmellType string
const (
    MysteryGuest      TestSmellType = "mystery_guest"
    ResourceOptimism  TestSmellType = "resource_optimism"
    AssertionRoulette TestSmellType = "assertion_roulette"
    EagerTest         TestSmellType = "eager_test"
    LazyTest          TestSmellType = "lazy_test"
    // ... more types
)

// Registry for all detectors
type SmellDetectorRegistry struct {
    detectors map[TestSmellType][]TestSmellDetector
}

func (sdr *SmellDetectorRegistry) Register(detector TestSmellDetector)
func (sdr *SmellDetectorRegistry) DetectAll(testFile *TestFile) ([]*TestSmell, error)
```

### 2. Detection Rule Patterns (Inspired by tsDetect)

**Mystery Guest Detection**:
```go
type MysteryGuestDetector struct {
    fileSystemPatterns []string
    networkPatterns    []string
    databasePatterns   []string
}

func (mgd *MysteryGuestDetector) Detect(testFile *TestFile) ([]*TestSmell, error) {
    smells := []*TestSmell{}

    for _, testFunc := range testFile.TestFunctions {
        // Check for file operations without explicit setup
        if hasFileOperations(testFunc) && !hasSetupForFiles(testFunc) {
            smells = append(smells, &TestSmell{
                Type:        MysteryGuest,
                Severity:    SeverityHigh,
                Location:    testFunc.Location,
                Description: "Test accesses file system without setup",
                Remediation: "Add explicit file setup in test or fixture",
                Confidence:  0.85,
            })
        }

        // Check for network calls
        if hasNetworkCalls(testFunc) && !hasMocking(testFunc) {
            smells = append(smells, &TestSmell{
                Type:        MysteryGuest,
                Severity:    SeverityHigh,
                Location:    testFunc.Location,
                Description: "Test makes network calls without mocking",
                Remediation: "Mock network calls or use test fixtures",
                Confidence:  0.90,
            })
        }
    }

    return smells, nil
}

// Language-specific detection patterns
func hasFileOperations(testFunc *TestFunction) bool {
    switch testFunc.Language {
    case "go":
        return containsPatterns(testFunc.Body, []string{
            "os.Open", "os.ReadFile", "ioutil.ReadFile",
        })
    case "python":
        return containsPatterns(testFunc.Body, []string{
            "open(", "with open", "pathlib.Path",
        })
    case "javascript":
        return containsPatterns(testFunc.Body, []string{
            "fs.readFile", "fs.readFileSync", "readFile(",
        })
    }
    return false
}
```

**Assertion Roulette Detection**:
```go
type AssertionRouletteDetector struct{}

func (ard *AssertionRouletteDetector) Detect(testFile *TestFile) ([]*TestSmell, error) {
    smells := []*TestSmell{}

    for _, testFunc := range testFile.TestFunctions {
        assertions := testFunc.Assertions

        if len(assertions) > 1 {
            assertionsWithoutMessages := 0
            for _, assertion := range assertions {
                if !hasAssertionMessage(assertion) {
                    assertionsWithoutMessages++
                }
            }

            if assertionsWithoutMessages > 1 {
                smells = append(smells, &TestSmell{
                    Type:        AssertionRoulette,
                    Severity:    SeverityMedium,
                    Location:    testFunc.Location,
                    Description: fmt.Sprintf("%d assertions without messages", assertionsWithoutMessages),
                    Remediation: "Add descriptive messages to assertions",
                    Confidence:  0.95,
                })
            }
        }
    }

    return smells, nil
}
```

**Eager Test Detection**:
```go
type EagerTestDetector struct{}

func (etd *EagerTestDetector) Detect(testFile *TestFile) ([]*TestSmell, error) {
    smells := []*TestSmell{}

    for _, testFunc := range testFile.TestFunctions {
        // Multiple production methods called in single test
        productionMethods := extractProductionMethodCalls(testFunc)

        if len(productionMethods) > 1 {
            smells = append(smells, &TestSmell{
                Type:        EagerTest,
                Severity:    SeverityMedium,
                Location:    testFunc.Location,
                Description: fmt.Sprintf("Test calls %d production methods", len(productionMethods)),
                Remediation: "Split into separate tests, one per method",
                Confidence:  0.80,
                Evidence: &Evidence{
                    Methods: productionMethods,
                },
            })
        }
    }

    return smells, nil
}
```

### 3. Language-Specific Pattern Libraries

**Go Patterns**:
```go
var goSmellPatterns = map[TestSmellType]*PatternSet{
    MysteryGuest: {
        FilePatterns:    []string{"os.Open", "os.ReadFile", "ioutil.ReadFile"},
        NetworkPatterns: []string{"http.Get", "http.Post", "net.Dial"},
        DBPatterns:      []string{"sql.Open", "db.Query"},
    },
    ConditionalLogic: {
        Patterns: []string{"if ", "for ", "switch ", "select {"},
    },
}
```

**Python Patterns**:
```go
var pythonSmellPatterns = map[TestSmellType]*PatternSet{
    MysteryGuest: {
        FilePatterns:    []string{"open(", "with open", "pathlib.Path"},
        NetworkPatterns: []string{"requests.", "urllib.", "httpx."},
        DBPatterns:      []string{"sqlite3.", "psycopg2.", "pymongo."},
    },
    ConditionalLogic: {
        Patterns: []string{"if ", "for ", "while "},
    },
}
```

**JavaScript Patterns**:
```go
var jsSmellPatterns = map[TestSmellType]*PatternSet{
    MysteryGuest: {
        FilePatterns:    []string{"fs.readFile", "fs.writeFile", "readFileSync"},
        NetworkPatterns: []string{"fetch(", "axios.", "http.get"},
        DBPatterns:      []string{"mongoose.", "sequelize.", "prisma."},
    },
}
```

## Prototype Requirements

### Deliverable 1: Core Detection Framework
**Files**: `internal/smell/detector.go`, `internal/smell/registry.go`
- Smell detector interface
- Detector registry
- Evidence collection system
- Confidence scoring

### Deliverable 2: Priority 1 Smell Detectors
**Files**: `internal/smell/detectors/*.go`
- Mystery Guest detector
- Resource Optimism detector
- Assertion Roulette detector
- Conditional Logic detector
- Flaky Test detector (pattern-based heuristics)

### Deliverable 3: Language Pattern Libraries
**Files**: `internal/smell/patterns/*.go`
- Go pattern definitions
- Python pattern definitions
- JavaScript pattern definitions
- Java pattern definitions

### Deliverable 4: Remediation Guidance Database
**Files**: `internal/smell/remediation.go`, `data/remediations.yaml`
```yaml
remediations:
  mystery_guest:
    description: "Test depends on external resources without explicit setup"
    impact: "Tests become flaky and environment-dependent"
    fix:
      - "Add explicit resource setup in test setup/fixture"
      - "Use mocking for external dependencies"
      - "Create test data within test or fixture"
    examples:
      go: |
        // Before
        func TestReadConfig(t *testing.T) {
          config := LoadConfig("/etc/app/config.yml") // Mystery guest!
        }

        // After
        func TestReadConfig(t *testing.T) {
          tmpfile, _ := os.CreateTemp("", "config.yml")
          defer os.Remove(tmpfile.Name())
          // Create test config...
          config := LoadConfig(tmpfile.Name())
        }
      python: |
        # Before
        def test_process_file():
          result = process_file("data.csv")  # Mystery guest!

        # After
        def test_process_file(tmp_path):
          test_file = tmp_path / "data.csv"
          test_file.write_text("test,data\n1,2")
          result = process_file(str(test_file))
```

### Deliverable 5: Validation Test Suite
**Files**: `internal/smell/detectors/*_test.go`
- Known smell test cases per type
- False positive test cases
- Multi-language test coverage
- Performance benchmarks

## Performance Benchmarks

| Language   | Test Files | Target Time | Smells Detected | Accuracy Target |
|------------|-----------|-------------|-----------------|-----------------|
| Go         | 100       | <3s         | All 11 types    | >90%           |
| Python     | 100       | <5s         | All 11 types    | >90%           |
| JavaScript | 100       | <5s         | All 11 types    | >90%           |

## Validation Strategy

### 1. Accuracy Validation
Create benchmark test suite with known smells:
```
testdata/
├── go/
│   ├── mystery_guest_positive.go      # Should detect
│   ├── mystery_guest_negative.go      # Should NOT detect
│   ├── assertion_roulette_positive.go
│   └── ...
├── python/
│   ├── mystery_guest_positive.py
│   └── ...
└── javascript/
    ├── mystery_guest_positive.test.js
    └── ...
```

Calculate precision and recall:
```
Precision = True Positives / (True Positives + False Positives)
Recall = True Positives / (True Positives + False Negatives)
F1 Score = 2 * (Precision * Recall) / (Precision + Recall)
```

### 2. Real-World Validation
Test on popular open-source projects:
- **Go**: kubernetes, docker, prometheus
- **Python**: django, flask, pytest
- **JavaScript**: react, vue, express

Compare results with manual code review (sample 50 files).

## Risk Mitigation

### Risk 1: High false positive rate
**Mitigation**:
- Confidence scoring for each detection
- Configurable thresholds per smell type
- User feedback mechanism
- Continuous refinement of patterns

### Risk 2: Language-specific pattern complexity
**Mitigation**:
- Start with common, well-defined patterns
- Incremental pattern addition
- Community contributions for patterns
- Regular pattern validation

### Risk 3: Performance degradation
**Mitigation**:
- Pattern matching optimization
- Parallel detection (independent detectors)
- Caching of AST analysis results
- Incremental detection (changed files only)

## Go/No-Go Decision Criteria

### GO if:
- ✅ Priority 1 smells detected with >90% accuracy
- ✅ Performance benchmarks met
- ✅ False positive rate <10%
- ✅ Framework extensible for new smell types
- ✅ Remediation guidance comprehensive

### NO-GO if:
- ❌ Accuracy <80% on any Priority 1 smell
- ❌ Performance >2x slower than targets
- ❌ False positive rate >20%
- ❌ Framework too complex to extend

### Alternative Approach:
If pattern-based detection insufficient:
- Machine learning-based detection (trained on tsDetect dataset)
- Hybrid approach (patterns + ML)
- Integration with existing tools (tsDetect for Java)

## Spike Deliverables

1. **Detection Framework Implementation**
   - Core interfaces and abstractions
   - Detector registry
   - Evidence collection system

2. **Priority 1 Smell Detectors**
   - 5 detectors with high accuracy
   - Multi-language support
   - Comprehensive tests

3. **Remediation Guidance Database**
   - YAML-based remediation data
   - Code examples per language
   - Integration with detectors

4. **Accuracy Validation Report**
   - Precision/recall metrics per smell
   - False positive analysis
   - Real-world validation results

5. **Performance Analysis**
   - Benchmark results
   - Scalability testing
   - Optimization recommendations

## Integration Guidelines

```go
// Usage in analysis pipeline
registry := smell.NewRegistry()
registry.RegisterAll(smell.DefaultDetectors())

analyzer := smell.NewAnalyzer(registry)
results, err := analyzer.AnalyzeTestFile(testFile)

for _, smell := range results.Smells {
    fmt.Printf("%s: %s (confidence: %.2f)\n",
        smell.Type, smell.Description, smell.Confidence)
    fmt.Printf("Remediation: %s\n", smell.Remediation)
}
```

## Success Metrics
- [ ] 11 test smell types implemented
- [ ] >90% accuracy on validation suite
- [ ] <10% false positive rate
- [ ] Performance targets met
- [ ] Remediation guidance complete
- [ ] Real-world validation successful

## Timeline
- **Week 1**: Core framework and 2 detectors
- **Week 2**: 3 more Priority 1 detectors
- **Week 3**: Remediation database and validation
- **Week 4**: Real-world testing and optimization

## Sources and References
- [tsDetect Tool (GitHub)](https://github.com/TestSmells/TSDetect)
- [tsDetect Paper (PDF)](https://testsmells.org/assets/publications/FSE2020_TechnicalPaper.pdf)
- [SniffTest (GitHub)](https://github.com/MaierFlorian/SniffTest)
- [Test Smells 20 years later (Research)](https://link.springer.com/article/10.1007/s10664-022-10207-5)
- [Best Code Smell Detection Tools 2026](https://www.getpanto.ai/blog/best-code-smell-detection-tools-to-optimize-code-quality)
