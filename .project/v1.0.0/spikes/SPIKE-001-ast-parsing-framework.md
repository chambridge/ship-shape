# SPIKE-001: Multi-Language AST Parsing Framework

## Overview
This spike validates the technical approach for parsing and analyzing test code across multiple programming languages (Go, Python, JavaScript/TypeScript, Java, Rust, etc.) to extract test patterns, structure, and quality metrics.

**Associated User Stories**: SS-010, SS-011, SS-012
**Risk Level**: HIGH
**Priority**: P0 (Critical)
**Target Completion**: Week 1-2 of implementation

## Problem Statement
Ship Shape needs to perform deep static analysis of test code across multiple languages to detect:
- Test structure and organization patterns
- Test framework usage (pytest, Jest, testing, JUnit, etc.)
- Test smells and anti-patterns
- Assertion quality and completeness
- Mock/stub usage patterns
- Table-driven tests, parametrization, fixtures

**Key Challenges**:
1. Each language requires different parsing approach (AST libraries vary)
2. Performance constraints (<3s for 100 test files)
3. Accuracy requirements (>95% pattern detection)
4. Maintainability (plugin architecture for new languages)
5. Cross-platform compatibility

## Spike Objectives
- [ ] Validate go/ast for Go test analysis (built-in advantage)
- [ ] Evaluate Python AST approaches (Python stdlib vs tree-sitter)
- [ ] Validate tree-sitter for JavaScript/TypeScript
- [ ] Design unified AST analyzer interface
- [ ] Prototype pattern detection for 3+ languages
- [ ] Benchmark parsing performance
- [ ] Assess maintenance burden for each approach

## Technical Investigation Areas

### 1. Go AST Parsing (go/ast package)
**Approach**: Use Go standard library `go/ast` and `go/parser`

**Validation Questions**:
- Can we reliably detect table-driven test patterns?
- Can we identify t.Parallel() usage and subtests?
- Can we detect test helpers and fixture patterns?
- What's the performance for parsing 100+ test files?

**Prototype Requirements**:
```go
// Detect table-driven tests
func detectTableDrivenTest(fn *ast.FuncDecl) bool
// Detect t.Parallel() calls
func hasParallelCall(fn *ast.FuncDecl) bool
// Count subtests (t.Run calls)
func countSubtests(fn *ast.FuncDecl) int
// Detect assertion libraries (testify, etc.)
func detectAssertionLibrary(file *ast.File) string
```

**Success Criteria**:
- Parse 100 Go test files in <2 seconds
- 100% accuracy on table-driven test detection
- 95%+ accuracy on test pattern detection

### 2. Python AST Parsing
**Option A: Python stdlib `ast` module via subprocess**
- ✅ Pros: Native Python parsing, accurate, comprehensive
- ❌ Cons: Requires Python runtime, subprocess overhead, cross-platform complexity

**Option B: tree-sitter with Python grammar**
- ✅ Pros: Pure Go integration, fast, no runtime dependency
- ❌ Cons: Query language learning curve, potential accuracy issues

**Validation Questions**:
- Can tree-sitter accurately detect pytest fixtures?
- Can we identify parametrized tests (@pytest.mark.parametrize)?
- What's the performance difference between approaches?
- Can we handle both pytest and unittest patterns?

**Prototype Requirements**:
```go
// Detect test functions
func detectPytestTests(content []byte) ([]*TestFunction, error)
// Detect fixtures and their scopes
func detectPytestFixtures(content []byte) ([]*Fixture, error)
// Detect parametrized tests
func detectParametrizedTests(content []byte) ([]*ParametrizedTest, error)
// Detect assertion usage
func detectAssertions(content []byte) ([]*Assertion, error)
```

**Success Criteria**:
- Parse 50 Python test files in <5 seconds
- 90%+ accuracy on pytest pattern detection
- Handle both pytest and unittest
- Works without Python runtime (tree-sitter approach)

### 3. JavaScript/TypeScript AST Parsing
**Approach**: tree-sitter with JavaScript/TypeScript grammars

**Validation Questions**:
- Can we accurately parse JSX/TSX test files?
- Can we detect describe/it blocks with proper nesting?
- Can we identify beforeEach/afterEach hooks?
- Can we detect Jest vs Vitest differences?
- What's the performance for large test suites?

**Prototype Requirements**:
```go
// Detect test blocks (describe, it, test)
func detectTestBlocks(content []byte, lang string) ([]*TestBlock, error)
// Detect hooks (beforeEach, afterEach, beforeAll, afterAll)
func detectTestHooks(content []byte) ([]*TestHook, error)
// Detect mocking patterns (jest.mock, vi.mock, jest.spyOn)
func detectMockUsage(content []byte) ([]*MockPattern, error)
// Detect snapshot tests
func detectSnapshotTests(content []byte) ([]*SnapshotTest, error)
```

**Success Criteria**:
- Parse 100 JS/TS test files in <10 seconds
- Handle both .js, .ts, .jsx, .tsx files
- 95%+ accuracy on test structure detection
- Differentiate Jest and Vitest patterns

## Unified Analyzer Interface Design

```go
// Language-agnostic analyzer interface
type TestFileAnalyzer interface {
    // Analyze a test file and return structured results
    Analyze(filePath string) (*TestFileAnalysis, error)

    // Get supported file extensions
    SupportedExtensions() []string

    // Get analyzer capabilities
    Capabilities() AnalyzerCapabilities
}

type TestFileAnalysis struct {
    FilePath      string
    Language      string
    Framework     string
    TestFunctions []*TestFunction
    Fixtures      []*Fixture
    Hooks         []*TestHook
    Assertions    []*Assertion
    Smells        []*TestSmell
    Metrics       *TestMetrics
}

type TestFunction struct {
    Name          string
    Location      SourceLocation
    Assertions    int
    Mocks         int
    IsTable       bool
    IsParallel    bool
    IsAsync       bool
    Subtests      []*TestFunction
    Complexity    int
}

// Analyzer registry for plugin-style architecture
type AnalyzerRegistry struct {
    analyzers map[string]TestFileAnalyzer
}

func (ar *AnalyzerRegistry) Register(lang string, analyzer TestFileAnalyzer)
func (ar *AnalyzerRegistry) GetAnalyzer(lang string) (TestFileAnalyzer, error)
```

## Prototype Requirements

### Deliverable 1: Go Test Analyzer Prototype
**Files**: `internal/analyzer/go_test_analyzer.go`
- Implement Go AST parsing for test patterns
- Detect table-driven tests, t.Parallel(), subtests
- Benchmark parsing performance
- Unit tests with real Go test files

### Deliverable 2: Python Test Analyzer Prototype (tree-sitter)
**Files**: `internal/analyzer/python_test_analyzer.go`
- Implement tree-sitter Python parsing
- Detect pytest and unittest patterns
- Compare accuracy vs subprocess Python AST approach
- Performance benchmarks

### Deliverable 3: JavaScript/TypeScript Test Analyzer Prototype
**Files**: `internal/analyzer/js_test_analyzer.go`
- Implement tree-sitter JS/TS parsing
- Detect Jest/Vitest test patterns
- Handle JSX/TSX files
- Performance benchmarks

### Deliverable 4: Unified Interface and Registry
**Files**: `internal/analyzer/interface.go`, `internal/analyzer/registry.go`
- Define common interfaces
- Implement analyzer registry
- Plugin registration mechanism
- Integration tests

## Performance Benchmarks

| Language   | Files | Target Time | Approach         |
|------------|-------|-------------|------------------|
| Go         | 100   | <2s         | go/ast           |
| Python     | 50    | <5s         | tree-sitter      |
| JavaScript | 100   | <10s        | tree-sitter      |
| TypeScript | 100   | <10s        | tree-sitter      |

## Risk Mitigation

### Risk 1: tree-sitter accuracy for complex patterns
**Mitigation**:
- Create comprehensive test suite with known patterns
- Compare tree-sitter results with native parsers
- Maintain fallback to subprocess approach if needed
- Document known limitations

### Risk 2: Performance degradation at scale
**Mitigation**:
- Implement streaming/incremental parsing
- Use worker pools for parallel processing
- Cache parsed ASTs for unchanged files
- Optimize tree-sitter queries

### Risk 3: Maintenance burden for multiple parsers
**Mitigation**:
- Design clean abstraction layer
- Document parser-specific patterns
- Create comprehensive test suites per language
- Consider using tree-sitter for all languages (if viable)

## Go/No-Go Decision Criteria

### GO if:
- ✅ All three language prototypes achieve >90% accuracy
- ✅ Performance benchmarks meet targets
- ✅ Unified interface design is clean and extensible
- ✅ tree-sitter proves viable for Python (no runtime dependency)
- ✅ Maintenance burden is acceptable

### NO-GO if:
- ❌ tree-sitter accuracy <85% for Python patterns
- ❌ Performance >2x slower than targets
- ❌ Interface design too complex or brittle
- ❌ Cross-platform issues with tree-sitter

### Alternative Approach:
If tree-sitter proves insufficient for Python:
- Use subprocess approach with embedded Python script
- Ship minimal Python analyzer as separate binary
- Use Language Server Protocol (LSP) for parsing

## Spike Deliverables

1. **Working Prototypes** (3 analyzers)
   - Go test analyzer with AST parsing
   - Python test analyzer (tree-sitter preferred)
   - JavaScript/TypeScript test analyzer

2. **Performance Benchmarks Report**
   - Parsing time measurements
   - Memory usage analysis
   - Scalability testing results

3. **Accuracy Assessment Report**
   - Pattern detection accuracy per language
   - False positive/negative rates
   - Known limitations documented

4. **Unified Interface Design Document**
   - Interface definitions
   - Plugin architecture design
   - Integration patterns

5. **Go/No-Go Recommendation**
   - Clear decision with supporting evidence
   - Risk assessment
   - Alternative approaches if needed

## Integration Guidelines

Upon successful spike completion:

1. **Code Structure**:
```
internal/analyzer/
├── interface.go          # Unified interfaces
├── registry.go           # Analyzer registry
├── go_analyzer.go        # Go implementation
├── python_analyzer.go    # Python implementation
├── js_analyzer.go        # JS/TS implementation
└── testdata/             # Test fixtures
```

2. **Configuration**:
```yaml
analyzers:
  go:
    enabled: true
    parallel: true
  python:
    enabled: true
    approach: tree-sitter  # or subprocess
  javascript:
    enabled: true
    handle_jsx: true
```

3. **Usage Pattern**:
```go
registry := analyzer.NewRegistry()
registry.Register("go", analyzer.NewGoAnalyzer())
registry.Register("python", analyzer.NewPythonAnalyzer())
registry.Register("javascript", analyzer.NewJSAnalyzer())

analyzer, _ := registry.GetAnalyzer("go")
results, _ := analyzer.Analyze("path/to/test.go")
```

## Success Metrics
- [ ] All prototypes implemented and tested
- [ ] Performance benchmarks meet or exceed targets
- [ ] Accuracy >90% for all languages
- [ ] Unified interface proven extensible
- [ ] Documentation complete
- [ ] Go/No-Go decision made with evidence

## Timeline
- **Week 1**: Go and Python prototype implementation
- **Week 2**: JavaScript/TypeScript prototype and interface design
- **Week 3**: Performance optimization and accuracy testing
- **Week 4**: Documentation and Go/No-Go decision

## References
- Go AST: https://pkg.go.dev/go/ast
- tree-sitter: https://tree-sitter.github.io/tree-sitter/
- tree-sitter Go bindings: https://github.com/smacker/go-tree-sitter
- Python AST: https://docs.python.org/3/library/ast.html
