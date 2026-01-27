# Ship Shape Analyzer Plugin System Design

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Design
**Author**: Senior Software Engineer

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Design Philosophy](#design-philosophy)
3. [Core Interface Definitions](#core-interface-definitions)
4. [Registration Mechanisms](#registration-mechanisms)
5. [Language-Specific Analyzer Implementations](#language-specific-analyzer-implementations)
6. [Plugin Discovery and Lifecycle](#plugin-discovery-and-lifecycle)
7. [Error Handling and Graceful Degradation](#error-handling-and-graceful-degradation)
8. [Performance and Concurrency](#performance-and-concurrency)
9. [Configuration and Extensibility](#configuration-and-extensibility)
10. [Testing Strategy](#testing-strategy)
11. [Implementation Roadmap](#implementation-roadmap)
12. [Examples and Usage Patterns](#examples-and-usage-patterns)

---

## Executive Summary

The Ship Shape Analyzer Plugin System provides a **flexible, extensible architecture** for analyzing test code across multiple programming languages. Drawing from proven patterns in AgentReady and informed by technical spike validation, this design enables:

- **Language-agnostic core** with language-specific plugin implementations
- **Dynamic plugin registration** for easy extensibility
- **Graceful degradation** when analyzers fail or are unavailable
- **Parallel execution** with proper isolation and resource management
- **Type-safe interfaces** leveraging Go's strong typing
- **Zero-configuration defaults** with override capabilities

### Key Design Principles

1. **Strategy Pattern**: Each analyzer is an independent, stateless strategy
2. **Registry Pattern**: Centralized plugin registration and discovery
3. **Factory Pattern**: Analyzer creation abstracted through factories
4. **Graceful Degradation**: Partial failures don't stop analysis
5. **Composition over Inheritance**: Shallow hierarchies, rich composition

---

## Design Philosophy

### Lessons from AgentReady

The AgentReady codebase demonstrates excellent plugin architecture through its `BaseAssessor` pattern:

```python
class BaseAssessor(ABC):
    @abstractmethod
    def assess(self, repository) -> Finding:
        """Assess a single attribute"""

    def is_applicable(self, repository) -> bool:
        """Check if assessor applies to this repository"""

    def calculate_proportional_score(self, value, min_val, max_val) -> float:
        """Standard scoring calculation"""
```

**Key Takeaways Applied to Ship Shape**:
- Clean separation of concerns (assess vs applicability)
- Standardized scoring mechanisms
- Language/framework awareness
- Isolated error handling per analyzer

### Ship Shape Specific Requirements

Based on requirements and spike findings:

1. **Multi-Language AST Parsing** (SPIKE-001)
   - Go: Native `go/ast` package
   - Python: tree-sitter or subprocess approach
   - JavaScript/TypeScript: tree-sitter
   - Unified interface abstracts parsing differences

2. **Monorepo Support** (SPIKE-002)
   - Package-level isolation
   - Parallel execution (4-8 concurrent analyzers optimal)
   - Context-aware analysis per package

3. **Performance Requirements**
   - Go: 100 files in <2s
   - Python: 50 files in <5s
   - JavaScript: 100 files in <10s
   - Total analysis: <10 minutes for typical repository

4. **Extensibility**
   - Plugin developers can add new languages without modifying core
   - Custom analyzers via external plugins
   - Override built-in analyzers

---

## Core Interface Definitions

### 1. Base Analyzer Interface

```go
package analyzer

import (
    "context"
    "time"
)

// Analyzer is the primary interface for all test file analyzers
type Analyzer interface {
    // Analyze a single test file and return structured results
    Analyze(ctx context.Context, req *AnalysisRequest) (*AnalysisResult, error)

    // GetMetadata returns analyzer capabilities and configuration
    GetMetadata() *AnalyzerMetadata

    // IsApplicable checks if this analyzer can process the given file
    IsApplicable(file *FileInfo) bool

    // Initialize prepares the analyzer (called once before first use)
    Initialize(config *AnalyzerConfig) error

    // Cleanup releases any resources held by the analyzer
    Cleanup() error
}

// AnalysisRequest encapsulates a file analysis request
type AnalysisRequest struct {
    // File to analyze
    File *FileInfo

    // Repository context
    Repository *RepositoryContext

    // Package context (for monorepos)
    Package *PackageContext

    // Analysis options
    Options *AnalysisOptions
}

// AnalysisResult contains structured analysis output
type AnalysisResult struct {
    // File analyzed
    FilePath string

    // Language and framework detected
    Language  string
    Framework string

    // Test functions discovered
    TestFunctions []*TestFunction

    // Test fixtures/helpers
    Fixtures []*Fixture

    // Test hooks (setup/teardown)
    Hooks []*TestHook

    // Assertions found
    Assertions []*Assertion

    // Test smells detected
    Smells []*TestSmell

    // Performance metrics
    Metrics *TestMetrics

    // Warnings or non-critical issues
    Warnings []string

    // Analysis metadata
    Metadata *ResultMetadata
}

// AnalyzerMetadata describes analyzer capabilities
type AnalyzerMetadata struct {
    // Unique analyzer identifier
    ID string

    // Human-readable name
    Name string

    // Supported language(s)
    Languages []string

    // Supported file extensions
    Extensions []string

    // Supported test frameworks
    Frameworks []string

    // Analyzer version
    Version string

    // Capabilities this analyzer provides
    Capabilities []Capability

    // Dependencies (external tools required)
    Dependencies []*Dependency

    // Performance characteristics
    Performance *PerformanceInfo
}
```

### 2. Language-Specific Extensions

```go
// GoAnalyzer extends Analyzer with Go-specific capabilities
type GoAnalyzer interface {
    Analyzer

    // Detect table-driven test patterns
    DetectTableDrivenTests(file *FileInfo) ([]*TableDrivenTest, error)

    // Detect t.Parallel() usage
    HasParallelTests(file *FileInfo) (bool, []string, error)

    // Count subtests (t.Run calls)
    CountSubtests(file *FileInfo) (int, error)

    // Detect testify assertion library usage
    DetectAssertionLibrary(file *FileInfo) (string, error)
}

// PythonAnalyzer extends Analyzer with Python-specific capabilities
type PythonAnalyzer interface {
    Analyzer

    // Detect pytest fixtures
    DetectPytestFixtures(file *FileInfo) ([]*PytestFixture, error)

    // Detect parametrized tests
    DetectParametrizedTests(file *FileInfo) ([]*ParametrizedTest, error)

    // Differentiate pytest vs unittest
    DetectTestFramework(file *FileInfo) (PythonFramework, error)

    // Detect assertion styles
    DetectAssertionStyle(file *FileInfo) (AssertionStyle, error)
}

// JSAnalyzer extends Analyzer with JavaScript/TypeScript capabilities
type JSAnalyzer interface {
    Analyzer

    // Detect test blocks (describe, it, test)
    DetectTestBlocks(file *FileInfo) ([]*TestBlock, error)

    // Detect hooks (beforeEach, afterEach)
    DetectTestHooks(file *FileInfo) ([]*JSTestHook, error)

    // Detect mocking patterns
    DetectMockUsage(file *FileInfo) ([]*MockPattern, error)

    // Detect snapshot tests
    DetectSnapshotTests(file *FileInfo) ([]*SnapshotTest, error)

    // Differentiate Jest vs Vitest
    DetectJSFramework(file *FileInfo) (JSFramework, error)
}
```

### 3. Supporting Types

```go
// FileInfo contains information about a file to analyze
type FileInfo struct {
    Path      string
    Name      string
    Extension string
    Content   []byte
    Size      int64
    Language  string
}

// RepositoryContext provides repository-level context
type RepositoryContext struct {
    RootPath     string
    Languages    []*LanguageInfo
    Frameworks   []*FrameworkInfo
    IsMonorepo   bool
    MonorepoType string
}

// PackageContext provides package-level context (monorepo)
type PackageContext struct {
    Name         string
    Path         string
    Language     string
    Framework    string
    Dependencies []string
}

// AnalysisOptions configures analysis behavior
type AnalysisOptions struct {
    // Enable deep analysis (more expensive)
    DeepAnalysis bool

    // Maximum analysis time per file
    Timeout time.Duration

    // Enable parallel processing
    Parallel bool

    // Cache parsed ASTs
    CacheAST bool

    // Language-specific options
    LanguageOptions map[string]interface{}
}

// Capability represents an analyzer capability
type Capability string

const (
    CapabilityTestDetection      Capability = "test_detection"
    CapabilityFixtureDetection   Capability = "fixture_detection"
    CapabilitySmellDetection     Capability = "smell_detection"
    CapabilityAssertionAnalysis  Capability = "assertion_analysis"
    CapabilityComplexityAnalysis Capability = "complexity_analysis"
    CapabilityMockDetection      Capability = "mock_detection"
    CapabilityPerformanceMetrics Capability = "performance_metrics"
)

// Dependency represents an external tool dependency
type Dependency struct {
    Name     string
    Version  string
    Optional bool
    Purpose  string
}

// PerformanceInfo describes analyzer performance characteristics
type PerformanceInfo struct {
    // Typical files per second
    FilesPerSecond float64

    // Memory usage per file (MB)
    MemoryPerFile float64

    // Supports parallel execution
    SupportsParallel bool
}
```

### 4. Test Analysis Data Structures

```go
// TestFunction represents a single test function
type TestFunction struct {
    Name          string
    Location      SourceLocation
    Type          TestType
    Assertions    int
    Mocks         int
    IsTable       bool
    IsParallel    bool
    IsAsync       bool
    Subtests      []*TestFunction
    Complexity    int
    Duration      time.Duration
    Tags          []string
    Description   string
}

// SourceLocation identifies code location
type SourceLocation struct {
    File      string
    Line      int
    Column    int
    EndLine   int
    EndColumn int
}

// TestType categorizes test types
type TestType string

const (
    TestTypeUnit        TestType = "unit"
    TestTypeIntegration TestType = "integration"
    TestTypeE2E         TestType = "e2e"
    TestTypePerformance TestType = "performance"
    TestTypeUnknown     TestType = "unknown"
)

// Fixture represents test fixtures/helpers
type Fixture struct {
    Name      string
    Location  SourceLocation
    Type      FixtureType
    Scope     FixtureScope
    UsedBy    []string
}

// FixtureType categorizes fixtures
type FixtureType string

const (
    FixtureTypeSetup    FixtureType = "setup"
    FixtureTypeTeardown FixtureType = "teardown"
    FixtureTypeFactory  FixtureType = "factory"
    FixtureTypeBuilder  FixtureType = "builder"
)

// FixtureScope defines fixture lifecycle
type FixtureScope string

const (
    FixtureScopeFunction FixtureScope = "function"
    FixtureScopeModule   FixtureScope = "module"
    FixtureScopeSession  FixtureScope = "session"
)

// TestHook represents setup/teardown hooks
type TestHook struct {
    Name     string
    Location SourceLocation
    Type     HookType
    Scope    HookScope
}

// HookType categorizes hook types
type HookType string

const (
    HookTypeBeforeEach HookType = "before_each"
    HookTypeAfterEach  HookType = "after_each"
    HookTypeBeforeAll  HookType = "before_all"
    HookTypeAfterAll   HookType = "after_all"
)

// Assertion represents a test assertion
type Assertion struct {
    Location      SourceLocation
    Type          AssertionType
    IsSpecific    bool
    HasMessage    bool
    Operator      string
}

// AssertionType categorizes assertion styles
type AssertionType string

const (
    AssertionTypeEquality    AssertionType = "equality"
    AssertionTypeTruthiness  AssertionType = "truthiness"
    AssertionTypeComparison  AssertionType = "comparison"
    AssertionTypeException   AssertionType = "exception"
    AssertionTypeContainment AssertionType = "containment"
)

// TestSmell represents detected test anti-patterns
type TestSmell struct {
    Type        SmellType
    Severity    Severity
    Location    SourceLocation
    Description string
    Remediation string
    Confidence  float64
}

// SmellType categorizes test smells
type SmellType string

const (
    SmellMysteryGuest        SmellType = "mystery_guest"
    SmellEagerTest           SmellType = "eager_test"
    SmellLazyTest            SmellType = "lazy_test"
    SmellObscureTest         SmellType = "obscure_test"
    SmellConditionalLogic    SmellType = "conditional_logic"
    SmellGeneralFixture      SmellType = "general_fixture"
    SmellCodeDuplication     SmellType = "code_duplication"
    SmellAssertionRoulette   SmellType = "assertion_roulette"
    SmellSensitiveEquality   SmellType = "sensitive_equality"
    SmellResourceOptimism    SmellType = "resource_optimism"
    SmellFlakiness           SmellType = "flakiness"
)

// Severity levels for findings
type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
    SeverityInfo     Severity = "info"
)

// TestMetrics contains quantitative test metrics
type TestMetrics struct {
    TotalTests           int
    TestsByType          map[TestType]int
    AverageComplexity    float64
    AverageAssertions    float64
    TotalAssertions      int
    TotalMocks           int
    TestCoverage         float64
    LinesOfTestCode      int
    TestToCodeRatio      float64
}

// ResultMetadata contains analysis metadata
type ResultMetadata struct {
    AnalyzerID      string
    AnalyzerVersion string
    AnalysisTime    time.Duration
    Timestamp       time.Time
    Success         bool
    ErrorMessage    string
}
```

---

## Registration Mechanisms

### 1. Analyzer Registry

The registry is the central hub for analyzer management, implementing a thread-safe plugin registry pattern.

```go
package analyzer

import (
    "fmt"
    "sync"
)

// Registry manages analyzer plugins
type Registry struct {
    mu        sync.RWMutex
    analyzers map[string]Analyzer
    factories map[string]AnalyzerFactory
    metadata  map[string]*AnalyzerMetadata
}

// AnalyzerFactory creates analyzer instances
type AnalyzerFactory func(config *AnalyzerConfig) (Analyzer, error)

// NewRegistry creates a new analyzer registry
func NewRegistry() *Registry {
    return &Registry{
        analyzers: make(map[string]Analyzer),
        factories: make(map[string]AnalyzerFactory),
        metadata:  make(map[string]*AnalyzerMetadata),
    }
}

// Register adds an analyzer to the registry
func (r *Registry) Register(id string, analyzer Analyzer) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.analyzers[id]; exists {
        return fmt.Errorf("analyzer %s already registered", id)
    }

    r.analyzers[id] = analyzer
    r.metadata[id] = analyzer.GetMetadata()

    return nil
}

// RegisterFactory adds an analyzer factory to the registry
func (r *Registry) RegisterFactory(id string, factory AnalyzerFactory) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.factories[id]; exists {
        return fmt.Errorf("factory for %s already registered", id)
    }

    r.factories[id] = factory

    return nil
}

// Get retrieves an analyzer by ID
func (r *Registry) Get(id string) (Analyzer, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    analyzer, exists := r.analyzers[id]
    if !exists {
        return nil, fmt.Errorf("analyzer %s not found", id)
    }

    return analyzer, nil
}

// GetByLanguage retrieves all analyzers for a language
func (r *Registry) GetByLanguage(language string) []Analyzer {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var result []Analyzer
    for _, analyzer := range r.analyzers {
        metadata := analyzer.GetMetadata()
        for _, lang := range metadata.Languages {
            if lang == language {
                result = append(result, analyzer)
                break
            }
        }
    }

    return result
}

// GetByExtension retrieves analyzers for a file extension
func (r *Registry) GetByExtension(extension string) []Analyzer {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var result []Analyzer
    for _, analyzer := range r.analyzers {
        metadata := analyzer.GetMetadata()
        for _, ext := range metadata.Extensions {
            if ext == extension {
                result = append(result, analyzer)
                break
            }
        }
    }

    return result
}

// GetApplicable finds analyzers applicable to a file
func (r *Registry) GetApplicable(file *FileInfo) []Analyzer {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var result []Analyzer
    for _, analyzer := range r.analyzers {
        if analyzer.IsApplicable(file) {
            result = append(result, analyzer)
        }
    }

    return result
}

// List returns all registered analyzer IDs
func (r *Registry) List() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    ids := make([]string, 0, len(r.analyzers))
    for id := range r.analyzers {
        ids = append(ids, id)
    }

    return ids
}

// GetMetadata retrieves metadata for an analyzer
func (r *Registry) GetMetadata(id string) (*AnalyzerMetadata, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    metadata, exists := r.metadata[id]
    if !exists {
        return nil, fmt.Errorf("metadata for %s not found", id)
    }

    return metadata, nil
}

// ListByCapability returns analyzers with a specific capability
func (r *Registry) ListByCapability(capability Capability) []Analyzer {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var result []Analyzer
    for _, analyzer := range r.analyzers {
        metadata := analyzer.GetMetadata()
        for _, cap := range metadata.Capabilities {
            if cap == capability {
                result = append(result, analyzer)
                break
            }
        }
    }

    return result
}
```

### 2. Plugin Discovery

```go
package analyzer

import (
    "fmt"
    "plugin"
    "path/filepath"
)

// PluginDiscovery handles dynamic plugin loading
type PluginDiscovery struct {
    registry   *Registry
    pluginDirs []string
}

// NewPluginDiscovery creates a plugin discovery manager
func NewPluginDiscovery(registry *Registry, pluginDirs []string) *PluginDiscovery {
    return &PluginDiscovery{
        registry:   registry,
        pluginDirs: pluginDirs,
    }
}

// DiscoverPlugins finds and loads analyzer plugins
func (pd *PluginDiscovery) DiscoverPlugins() error {
    for _, dir := range pd.pluginDirs {
        plugins, err := filepath.Glob(filepath.Join(dir, "*.so"))
        if err != nil {
            return fmt.Errorf("failed to discover plugins in %s: %w", dir, err)
        }

        for _, pluginPath := range plugins {
            if err := pd.loadPlugin(pluginPath); err != nil {
                // Log error but continue loading other plugins
                fmt.Printf("Warning: failed to load plugin %s: %v\n", pluginPath, err)
            }
        }
    }

    return nil
}

// loadPlugin loads a single plugin
func (pd *PluginDiscovery) loadPlugin(path string) error {
    p, err := plugin.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open plugin: %w", err)
    }

    // Look for NewAnalyzer function
    symbol, err := p.Lookup("NewAnalyzer")
    if err != nil {
        return fmt.Errorf("plugin missing NewAnalyzer function: %w", err)
    }

    // Cast to factory function
    factory, ok := symbol.(func(*AnalyzerConfig) (Analyzer, error))
    if !ok {
        return fmt.Errorf("NewAnalyzer has incorrect signature")
    }

    // Create analyzer instance
    analyzer, err := factory(nil)
    if err != nil {
        return fmt.Errorf("failed to create analyzer: %w", err)
    }

    // Register the analyzer
    metadata := analyzer.GetMetadata()
    return pd.registry.Register(metadata.ID, analyzer)
}
```

### 3. Initialization System

```go
package analyzer

import (
    "context"
    "fmt"
)

// InitializationManager handles analyzer lifecycle
type InitializationManager struct {
    registry *Registry
    configs  map[string]*AnalyzerConfig
}

// NewInitializationManager creates an initialization manager
func NewInitializationManager(registry *Registry) *InitializationManager {
    return &InitializationManager{
        registry: registry,
        configs:  make(map[string]*AnalyzerConfig),
    }
}

// SetConfig sets configuration for an analyzer
func (im *InitializationManager) SetConfig(analyzerID string, config *AnalyzerConfig) {
    im.configs[analyzerID] = config
}

// InitializeAll initializes all registered analyzers
func (im *InitializationManager) InitializeAll(ctx context.Context) error {
    analyzerIDs := im.registry.List()

    for _, id := range analyzerIDs {
        analyzer, err := im.registry.Get(id)
        if err != nil {
            return fmt.Errorf("failed to get analyzer %s: %w", id, err)
        }

        config := im.configs[id]
        if config == nil {
            config = &AnalyzerConfig{} // Default config
        }

        if err := analyzer.Initialize(config); err != nil {
            return fmt.Errorf("failed to initialize analyzer %s: %w", id, err)
        }
    }

    return nil
}

// CleanupAll cleans up all registered analyzers
func (im *InitializationManager) CleanupAll() error {
    analyzerIDs := im.registry.List()

    var errors []error
    for _, id := range analyzerIDs {
        analyzer, err := im.registry.Get(id)
        if err != nil {
            errors = append(errors, err)
            continue
        }

        if err := analyzer.Cleanup(); err != nil {
            errors = append(errors, fmt.Errorf("cleanup %s: %w", id, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("cleanup errors: %v", errors)
    }

    return nil
}
```

---

## Language-Specific Analyzer Implementations

### 1. Go Analyzer Implementation

```go
package goanalyzer

import (
    "context"
    "go/ast"
    "go/parser"
    "go/token"
    "io/ioutil"

    "github.com/shipshape/analyzer"
)

// GoTestAnalyzer implements Analyzer for Go test files
type GoTestAnalyzer struct {
    config   *analyzer.AnalyzerConfig
    fset     *token.FileSet
    astCache map[string]*ast.File
}

// NewGoTestAnalyzer creates a new Go analyzer
func NewGoTestAnalyzer() *GoTestAnalyzer {
    return &GoTestAnalyzer{
        fset:     token.NewFileSet(),
        astCache: make(map[string]*ast.File),
    }
}

// Analyze implements Analyzer.Analyze
func (ga *GoTestAnalyzer) Analyze(
    ctx context.Context,
    req *analyzer.AnalysisRequest,
) (*analyzer.AnalysisResult, error) {
    // Parse the file
    file, err := ga.parseFile(req.File)
    if err != nil {
        return nil, err
    }

    result := &analyzer.AnalysisResult{
        FilePath:  req.File.Path,
        Language:  "go",
        Framework: "testing",
        Metadata: &analyzer.ResultMetadata{
            AnalyzerID:      "go-test-analyzer",
            AnalyzerVersion: "1.0.0",
        },
    }

    // Extract test functions
    testFuncs := ga.extractTestFunctions(file)
    result.TestFunctions = testFuncs

    // Detect table-driven tests
    for _, tf := range testFuncs {
        if ga.isTableDriven(file, tf) {
            tf.IsTable = true
        }
    }

    // Detect parallel tests
    for _, tf := range testFuncs {
        if ga.hasParallelCall(file, tf) {
            tf.IsParallel = true
        }
    }

    // Extract subtests
    for _, tf := range testFuncs {
        subtests := ga.extractSubtests(file, tf)
        tf.Subtests = subtests
    }

    // Calculate metrics
    result.Metrics = ga.calculateMetrics(testFuncs)

    return result, nil
}

// GetMetadata implements Analyzer.GetMetadata
func (ga *GoTestAnalyzer) GetMetadata() *analyzer.AnalyzerMetadata {
    return &analyzer.AnalyzerMetadata{
        ID:         "go-test-analyzer",
        Name:       "Go Test Analyzer",
        Languages:  []string{"go"},
        Extensions: []string{".go"},
        Frameworks: []string{"testing", "testify"},
        Version:    "1.0.0",
        Capabilities: []analyzer.Capability{
            analyzer.CapabilityTestDetection,
            analyzer.CapabilityComplexityAnalysis,
            analyzer.CapabilitySmellDetection,
        },
        Performance: &analyzer.PerformanceInfo{
            FilesPerSecond:   50,
            MemoryPerFile:    1.5,
            SupportsParallel: true,
        },
    }
}

// IsApplicable implements Analyzer.IsApplicable
func (ga *GoTestAnalyzer) IsApplicable(file *analyzer.FileInfo) bool {
    return file.Extension == ".go" &&
           (file.Name == "_test.go" ||
            len(file.Name) > 8 && file.Name[len(file.Name)-8:] == "_test.go")
}

// Initialize implements Analyzer.Initialize
func (ga *GoTestAnalyzer) Initialize(config *analyzer.AnalyzerConfig) error {
    ga.config = config
    return nil
}

// Cleanup implements Analyzer.Cleanup
func (ga *GoTestAnalyzer) Cleanup() error {
    // Clear AST cache
    ga.astCache = make(map[string]*ast.File)
    return nil
}

// parseFile parses a Go file into an AST
func (ga *GoTestAnalyzer) parseFile(file *analyzer.FileInfo) (*ast.File, error) {
    // Check cache
    if cached, exists := ga.astCache[file.Path]; exists {
        return cached, nil
    }

    // Parse file
    f, err := parser.ParseFile(ga.fset, file.Path, file.Content, parser.ParseComments)
    if err != nil {
        return nil, err
    }

    // Cache if enabled
    if ga.config != nil && ga.config.CacheAST {
        ga.astCache[file.Path] = f
    }

    return f, nil
}

// extractTestFunctions finds all test functions
func (ga *GoTestAnalyzer) extractTestFunctions(file *ast.File) []*analyzer.TestFunction {
    var tests []*analyzer.TestFunction

    ast.Inspect(file, func(n ast.Node) bool {
        if fn, ok := n.(*ast.FuncDecl); ok {
            if ga.isTestFunction(fn) {
                tf := &analyzer.TestFunction{
                    Name: fn.Name.Name,
                    Location: analyzer.SourceLocation{
                        Line:   ga.fset.Position(fn.Pos()).Line,
                        Column: ga.fset.Position(fn.Pos()).Column,
                    },
                }

                // Calculate complexity
                tf.Complexity = ga.calculateComplexity(fn)

                tests = append(tests, tf)
            }
        }
        return true
    })

    return tests
}

// isTestFunction checks if a function is a test
func (ga *GoTestAnalyzer) isTestFunction(fn *ast.FuncDecl) bool {
    if len(fn.Name.Name) < 4 {
        return false
    }

    // Check for Test*, Benchmark*, Example* prefixes
    name := fn.Name.Name
    return (name[:4] == "Test" ||
            name[:9] == "Benchmark" ||
            name[:7] == "Example") &&
           ga.hasTestingParameter(fn)
}

// hasTestingParameter checks for *testing.T parameter
func (ga *GoTestAnalyzer) hasTestingParameter(fn *ast.FuncDecl) bool {
    if fn.Type.Params.NumFields() != 1 {
        return false
    }

    param := fn.Type.Params.List[0]
    starExpr, ok := param.Type.(*ast.StarExpr)
    if !ok {
        return false
    }

    selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
    if !ok {
        return false
    }

    ident, ok := selectorExpr.X.(*ast.Ident)
    if !ok {
        return false
    }

    return ident.Name == "testing" && selectorExpr.Sel.Name == "T"
}

// isTableDriven detects table-driven test pattern
func (ga *GoTestAnalyzer) isTableDriven(file *ast.File, tf *analyzer.TestFunction) bool {
    // Look for test cases slice
    hasTestCases := false

    ast.Inspect(file, func(n ast.Node) bool {
        if assign, ok := n.(*ast.AssignStmt); ok {
            for _, lhs := range assign.Lhs {
                if ident, ok := lhs.(*ast.Ident); ok {
                    if ident.Name == "tests" || ident.Name == "testCases" {
                        hasTestCases = true
                        return false
                    }
                }
            }
        }
        return true
    })

    return hasTestCases
}

// hasParallelCall detects t.Parallel() calls
func (ga *GoTestAnalyzer) hasParallelCall(file *ast.File, tf *analyzer.TestFunction) bool {
    hasParallel := false

    ast.Inspect(file, func(n ast.Node) bool {
        if call, ok := n.(*ast.CallExpr); ok {
            if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
                if sel.Sel.Name == "Parallel" {
                    hasParallel = true
                    return false
                }
            }
        }
        return true
    })

    return hasParallel
}

// extractSubtests finds t.Run calls
func (ga *GoTestAnalyzer) extractSubtests(file *ast.File, tf *analyzer.TestFunction) []*analyzer.TestFunction {
    var subtests []*analyzer.TestFunction

    ast.Inspect(file, func(n ast.Node) bool {
        if call, ok := n.(*ast.CallExpr); ok {
            if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
                if sel.Sel.Name == "Run" && len(call.Args) >= 2 {
                    // Extract subtest name
                    if lit, ok := call.Args[0].(*ast.BasicLit); ok {
                        subtest := &analyzer.TestFunction{
                            Name: lit.Value,
                            Location: analyzer.SourceLocation{
                                Line: ga.fset.Position(call.Pos()).Line,
                            },
                        }
                        subtests = append(subtests, subtest)
                    }
                }
            }
        }
        return true
    })

    return subtests
}

// calculateComplexity calculates cyclomatic complexity
func (ga *GoTestAnalyzer) calculateComplexity(fn *ast.FuncDecl) int {
    complexity := 1 // Base complexity

    ast.Inspect(fn, func(n ast.Node) bool {
        switch n.(type) {
        case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
             *ast.CaseClause, *ast.CommClause:
            complexity++
        }
        return true
    })

    return complexity
}

// calculateMetrics computes test metrics
func (ga *GoTestAnalyzer) calculateMetrics(tests []*analyzer.TestFunction) *analyzer.TestMetrics {
    metrics := &analyzer.TestMetrics{
        TotalTests: len(tests),
    }

    totalComplexity := 0
    for _, test := range tests {
        totalComplexity += test.Complexity
    }

    if len(tests) > 0 {
        metrics.AverageComplexity = float64(totalComplexity) / float64(len(tests))
    }

    return metrics
}
```

### 2. Python Analyzer Implementation (tree-sitter)

```go
package pythonanalyzer

import (
    "context"

    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/python"

    "github.com/shipshape/analyzer"
)

// PythonTestAnalyzer implements Analyzer for Python test files
type PythonTestAnalyzer struct {
    config *analyzer.AnalyzerConfig
    parser *sitter.Parser
}

// NewPythonTestAnalyzer creates a new Python analyzer
func NewPythonTestAnalyzer() *PythonTestAnalyzer {
    parser := sitter.NewParser()
    parser.SetLanguage(python.GetLanguage())

    return &PythonTestAnalyzer{
        parser: parser,
    }
}

// Analyze implements Analyzer.Analyze
func (pa *PythonTestAnalyzer) Analyze(
    ctx context.Context,
    req *analyzer.AnalysisRequest,
) (*analyzer.AnalysisResult, error) {
    // Parse file with tree-sitter
    tree := pa.parser.Parse(nil, req.File.Content)
    defer tree.Close()

    root := tree.RootNode()

    result := &analyzer.AnalysisResult{
        FilePath:  req.File.Path,
        Language:  "python",
        Framework: pa.detectFramework(root, req.File.Content),
        Metadata: &analyzer.ResultMetadata{
            AnalyzerID:      "python-test-analyzer",
            AnalyzerVersion: "1.0.0",
        },
    }

    // Extract test functions
    testFuncs := pa.extractTestFunctions(root, req.File.Content)
    result.TestFunctions = testFuncs

    // Extract fixtures (pytest)
    fixtures := pa.extractFixtures(root, req.File.Content)
    result.Fixtures = fixtures

    // Detect parametrized tests
    pa.markParametrizedTests(testFuncs, root, req.File.Content)

    // Calculate metrics
    result.Metrics = pa.calculateMetrics(testFuncs)

    return result, nil
}

// GetMetadata implements Analyzer.GetMetadata
func (pa *PythonTestAnalyzer) GetMetadata() *analyzer.AnalyzerMetadata {
    return &analyzer.AnalyzerMetadata{
        ID:         "python-test-analyzer",
        Name:       "Python Test Analyzer",
        Languages:  []string{"python"},
        Extensions: []string{".py"},
        Frameworks: []string{"pytest", "unittest"},
        Version:    "1.0.0",
        Capabilities: []analyzer.Capability{
            analyzer.CapabilityTestDetection,
            analyzer.CapabilityFixtureDetection,
            analyzer.CapabilitySmellDetection,
        },
        Performance: &analyzer.PerformanceInfo{
            FilesPerSecond:   20,
            MemoryPerFile:    2.5,
            SupportsParallel: true,
        },
    }
}

// IsApplicable implements Analyzer.IsApplicable
func (pa *PythonTestAnalyzer) IsApplicable(file *analyzer.FileInfo) bool {
    return file.Extension == ".py" &&
           (file.Name[:5] == "test_" || file.Name[len(file.Name)-8:] == "_test.py")
}

// Initialize implements Analyzer.Initialize
func (pa *PythonTestAnalyzer) Initialize(config *analyzer.AnalyzerConfig) error {
    pa.config = config
    return nil
}

// Cleanup implements Analyzer.Cleanup
func (pa *PythonTestAnalyzer) Cleanup() error {
    return nil
}

// detectFramework determines if pytest or unittest
func (pa *PythonTestAnalyzer) detectFramework(root *sitter.Node, source []byte) string {
    // Query for pytest imports
    query := `
        (import_statement
            name: (dotted_name) @import
            (#match? @import "pytest"))
    `

    // Simple heuristic: check for pytest imports
    content := string(source)
    if contains(content, "import pytest") || contains(content, "from pytest") {
        return "pytest"
    }

    if contains(content, "import unittest") || contains(content, "from unittest") {
        return "unittest"
    }

    return "unknown"
}

// extractTestFunctions finds test functions
func (pa *PythonTestAnalyzer) extractTestFunctions(root *sitter.Node, source []byte) []*analyzer.TestFunction {
    var tests []*analyzer.TestFunction

    // Traverse AST to find function definitions
    var traverse func(*sitter.Node)
    traverse = func(node *sitter.Node) {
        if node.Type() == "function_definition" {
            nameNode := node.ChildByFieldName("name")
            if nameNode != nil {
                name := source[nameNode.StartByte():nameNode.EndByte()]
                if pa.isTestFunction(string(name)) {
                    tf := &analyzer.TestFunction{
                        Name: string(name),
                        Location: analyzer.SourceLocation{
                            Line:   int(node.StartPoint().Row) + 1,
                            Column: int(node.StartPoint().Column) + 1,
                        },
                    }
                    tests = append(tests, tf)
                }
            }
        }

        for i := 0; i < int(node.ChildCount()); i++ {
            traverse(node.Child(i))
        }
    }

    traverse(root)

    return tests
}

// isTestFunction checks if function name indicates a test
func (pa *PythonTestAnalyzer) isTestFunction(name string) bool {
    return len(name) > 5 && name[:5] == "test_"
}

// extractFixtures finds pytest fixtures
func (pa *PythonTestAnalyzer) extractFixtures(root *sitter.Node, source []byte) []*analyzer.Fixture {
    var fixtures []*analyzer.Fixture

    // Look for @pytest.fixture decorators
    var traverse func(*sitter.Node)
    traverse = func(node *sitter.Node) {
        if node.Type() == "decorated_definition" {
            // Check for pytest.fixture decorator
            decoratorList := node.ChildByFieldName("decorator")
            if decoratorList != nil && pa.hasPytestFixtureDecorator(decoratorList, source) {
                // Get function name
                definition := node.ChildByFieldName("definition")
                if definition != nil && definition.Type() == "function_definition" {
                    nameNode := definition.ChildByFieldName("name")
                    if nameNode != nil {
                        name := source[nameNode.StartByte():nameNode.EndByte()]
                        fixture := &analyzer.Fixture{
                            Name: string(name),
                            Location: analyzer.SourceLocation{
                                Line: int(node.StartPoint().Row) + 1,
                            },
                            Type:  analyzer.FixtureTypeSetup,
                            Scope: analyzer.FixtureScopeFunction,
                        }
                        fixtures = append(fixtures, fixture)
                    }
                }
            }
        }

        for i := 0; i < int(node.ChildCount()); i++ {
            traverse(node.Child(i))
        }
    }

    traverse(root)

    return fixtures
}

// hasPytestFixtureDecorator checks for pytest.fixture decorator
func (pa *PythonTestAnalyzer) hasPytestFixtureDecorator(node *sitter.Node, source []byte) bool {
    content := source[node.StartByte():node.EndByte()]
    return contains(string(content), "pytest.fixture")
}

// markParametrizedTests detects @pytest.mark.parametrize
func (pa *PythonTestAnalyzer) markParametrizedTests(
    tests []*analyzer.TestFunction,
    root *sitter.Node,
    source []byte,
) {
    // Implementation similar to extractFixtures
    // Mark tests with parametrization tags
}

// calculateMetrics computes test metrics
func (pa *PythonTestAnalyzer) calculateMetrics(tests []*analyzer.TestFunction) *analyzer.TestMetrics {
    return &analyzer.TestMetrics{
        TotalTests: len(tests),
    }
}

// Helper function
func contains(s, substr string) bool {
    return len(s) >= len(substr) &&
           findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return i
        }
    }
    return -1
}
```

### 3. JavaScript/TypeScript Analyzer Implementation

```go
package jsanalyzer

import (
    "context"

    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/javascript"
    "github.com/smacker/go-tree-sitter/typescript/typescript"

    "github.com/shipshape/analyzer"
)

// JSTestAnalyzer implements Analyzer for JavaScript/TypeScript test files
type JSTestAnalyzer struct {
    config    *analyzer.AnalyzerConfig
    jsParser  *sitter.Parser
    tsParser  *sitter.Parser
}

// NewJSTestAnalyzer creates a new JS/TS analyzer
func NewJSTestAnalyzer() *JSTestAnalyzer {
    jsParser := sitter.NewParser()
    jsParser.SetLanguage(javascript.GetLanguage())

    tsParser := sitter.NewParser()
    tsParser.SetLanguage(typescript.GetLanguage())

    return &JSTestAnalyzer{
        jsParser: jsParser,
        tsParser: tsParser,
    }
}

// Analyze implements Analyzer.Analyze
func (ja *JSTestAnalyzer) Analyze(
    ctx context.Context,
    req *analyzer.AnalysisRequest,
) (*analyzer.AnalysisResult, error) {
    // Choose parser based on file extension
    parser := ja.jsParser
    language := "javascript"

    if req.File.Extension == ".ts" || req.File.Extension == ".tsx" {
        parser = ja.tsParser
        language = "typescript"
    }

    // Parse file
    tree := parser.Parse(nil, req.File.Content)
    defer tree.Close()

    root := tree.RootNode()

    result := &analyzer.AnalysisResult{
        FilePath:  req.File.Path,
        Language:  language,
        Framework: ja.detectFramework(root, req.File.Content),
        Metadata: &analyzer.ResultMetadata{
            AnalyzerID:      "js-test-analyzer",
            AnalyzerVersion: "1.0.0",
        },
    }

    // Extract test blocks
    testBlocks := ja.extractTestBlocks(root, req.File.Content)
    result.TestFunctions = testBlocks

    // Extract hooks
    hooks := ja.extractHooks(root, req.File.Content)
    result.Hooks = hooks

    // Detect mocks
    // (Implementation omitted for brevity)

    // Calculate metrics
    result.Metrics = ja.calculateMetrics(testBlocks)

    return result, nil
}

// GetMetadata implements Analyzer.GetMetadata
func (ja *JSTestAnalyzer) GetMetadata() *analyzer.AnalyzerMetadata {
    return &analyzer.AnalyzerMetadata{
        ID:         "js-test-analyzer",
        Name:       "JavaScript/TypeScript Test Analyzer",
        Languages:  []string{"javascript", "typescript"},
        Extensions: []string{".js", ".jsx", ".ts", ".tsx"},
        Frameworks: []string{"jest", "vitest", "mocha"},
        Version:    "1.0.0",
        Capabilities: []analyzer.Capability{
            analyzer.CapabilityTestDetection,
            analyzer.CapabilityMockDetection,
            analyzer.CapabilitySmellDetection,
        },
        Performance: &analyzer.PerformanceInfo{
            FilesPerSecond:   15,
            MemoryPerFile:    3.0,
            SupportsParallel: true,
        },
    }
}

// IsApplicable implements Analyzer.IsApplicable
func (ja *JSTestAnalyzer) IsApplicable(file *analyzer.FileInfo) bool {
    ext := file.Extension
    name := file.Name

    if ext != ".js" && ext != ".jsx" && ext != ".ts" && ext != ".tsx" {
        return false
    }

    // Check naming patterns
    return contains(name, ".test.") ||
           contains(name, ".spec.") ||
           contains(name, "__tests__")
}

// Initialize implements Analyzer.Initialize
func (ja *JSTestAnalyzer) Initialize(config *analyzer.AnalyzerConfig) error {
    ja.config = config
    return nil
}

// Cleanup implements Analyzer.Cleanup
func (ja *JSTestAnalyzer) Cleanup() error {
    return nil
}

// detectFramework determines Jest vs Vitest
func (ja *JSTestAnalyzer) detectFramework(root *sitter.Node, source []byte) string {
    content := string(source)

    if contains(content, "vitest") || contains(content, "import { vi") {
        return "vitest"
    }

    if contains(content, "jest") || contains(content, "expect(") {
        return "jest"
    }

    if contains(content, "describe(") {
        return "mocha"
    }

    return "unknown"
}

// extractTestBlocks finds describe/it/test blocks
func (ja *JSTestAnalyzer) extractTestBlocks(root *sitter.Node, source []byte) []*analyzer.TestFunction {
    var tests []*analyzer.TestFunction

    var traverse func(*sitter.Node, int)
    traverse = func(node *sitter.Node, depth int) {
        if node.Type() == "call_expression" {
            // Check if it's describe, it, or test
            function := node.ChildByFieldName("function")
            if function != nil {
                funcName := source[function.StartByte():function.EndByte()]

                if ja.isTestBlock(string(funcName)) {
                    // Extract test name
                    arguments := node.ChildByFieldName("arguments")
                    if arguments != nil && arguments.ChildCount() > 0 {
                        firstArg := arguments.Child(1) // Skip '('
                        if firstArg != nil {
                            testName := source[firstArg.StartByte():firstArg.EndByte()]
                            tf := &analyzer.TestFunction{
                                Name: string(testName),
                                Location: analyzer.SourceLocation{
                                    Line:   int(node.StartPoint().Row) + 1,
                                    Column: int(node.StartPoint().Column) + 1,
                                },
                            }
                            tests = append(tests, tf)
                        }
                    }
                }
            }
        }

        for i := 0; i < int(node.ChildCount()); i++ {
            traverse(node.Child(i), depth+1)
        }
    }

    traverse(root, 0)

    return tests
}

// isTestBlock checks if function name is a test block
func (ja *JSTestAnalyzer) isTestBlock(name string) bool {
    return name == "describe" ||
           name == "it" ||
           name == "test" ||
           name == "context"
}

// extractHooks finds beforeEach, afterEach, etc.
func (ja *JSTestAnalyzer) extractHooks(root *sitter.Node, source []byte) []*analyzer.TestHook {
    var hooks []*analyzer.TestHook

    // Similar implementation to extractTestBlocks
    // Look for beforeEach, afterEach, beforeAll, afterAll calls

    return hooks
}

// calculateMetrics computes test metrics
func (ja *JSTestAnalyzer) calculateMetrics(tests []*analyzer.TestFunction) *analyzer.TestMetrics {
    return &analyzer.TestMetrics{
        TotalTests: len(tests),
    }
}
```

---

## Plugin Discovery and Lifecycle

### 1. Auto-Registration Pattern

```go
package main

import (
    "github.com/shipshape/analyzer"
    "github.com/shipshape/analyzer/goanalyzer"
    "github.com/shipshape/analyzer/pythonanalyzer"
    "github.com/shipshape/analyzer/jsanalyzer"
)

// RegisterBuiltinAnalyzers registers all built-in analyzers
func RegisterBuiltinAnalyzers(registry *analyzer.Registry) error {
    // Go analyzer
    goAnalyzer := goanalyzer.NewGoTestAnalyzer()
    if err := registry.Register("go-test", goAnalyzer); err != nil {
        return err
    }

    // Python analyzer
    pyAnalyzer := pythonanalyzer.NewPythonTestAnalyzer()
    if err := registry.Register("python-test", pyAnalyzer); err != nil {
        return err
    }

    // JavaScript/TypeScript analyzer
    jsAnalyzer := jsanalyzer.NewJSTestAnalyzer()
    if err := registry.Register("js-test", jsAnalyzer); err != nil {
        return err
    }

    return nil
}
```

### 2. Configuration-Based Enabling

```yaml
# .shipshape.yml
analyzers:
  go-test:
    enabled: true
    parallel: true
    cache_ast: true

  python-test:
    enabled: true
    approach: tree-sitter  # or 'subprocess'

  js-test:
    enabled: true
    handle_jsx: true
    frameworks: [jest, vitest]

  # Custom analyzer
  custom-java-test:
    enabled: false
    plugin_path: ./plugins/java-analyzer.so
```

### 3. Lifecycle Management

```go
package analyzer

import (
    "context"
    "fmt"
)

// Lifecycle manages analyzer lifecycle
type Lifecycle struct {
    registry *Registry
    configs  map[string]*AnalyzerConfig
}

// Start initializes all enabled analyzers
func (lc *Lifecycle) Start(ctx context.Context) error {
    for id, config := range lc.configs {
        if !config.Enabled {
            continue
        }

        analyzer, err := lc.registry.Get(id)
        if err != nil {
            return fmt.Errorf("analyzer %s not found: %w", id, err)
        }

        if err := analyzer.Initialize(config); err != nil {
            return fmt.Errorf("failed to initialize %s: %w", id, err)
        }
    }

    return nil
}

// Stop cleans up all analyzers
func (lc *Lifecycle) Stop() error {
    for id := range lc.configs {
        analyzer, err := lc.registry.Get(id)
        if err != nil {
            continue
        }

        if err := analyzer.Cleanup(); err != nil {
            fmt.Printf("Warning: cleanup failed for %s: %v\n", id, err)
        }
    }

    return nil
}
```

---

## Error Handling and Graceful Degradation

### 1. Error Types

```go
package analyzer

import "fmt"

// AnalysisError represents an analysis error
type AnalysisError struct {
    AnalyzerID  string
    FilePath    string
    Err         error
    Recoverable bool
}

func (e *AnalysisError) Error() string {
    return fmt.Sprintf("analyzer %s failed on %s: %v",
        e.AnalyzerID, e.FilePath, e.Err)
}

// Common error types
var (
    ErrNotApplicable     = fmt.Errorf("analyzer not applicable")
    ErrMissingDependency = fmt.Errorf("missing dependency")
    ErrParsingFailed     = fmt.Errorf("parsing failed")
    ErrTimeout           = fmt.Errorf("analysis timeout")
)
```

### 2. Graceful Degradation Strategy

```go
package analyzer

import (
    "context"
    "fmt"
)

// ExecuteWithFallback tries analyzers in order until one succeeds
func ExecuteWithFallback(
    ctx context.Context,
    analyzers []Analyzer,
    req *AnalysisRequest,
) (*AnalysisResult, error) {
    var errors []error

    for _, analyzer := range analyzers {
        // Check applicability
        if !analyzer.IsApplicable(req.File) {
            continue
        }

        // Try analysis
        result, err := analyzer.Analyze(ctx, req)
        if err == nil {
            return result, nil
        }

        errors = append(errors, &AnalysisError{
            AnalyzerID:  analyzer.GetMetadata().ID,
            FilePath:    req.File.Path,
            Err:         err,
            Recoverable: isRecoverable(err),
        })
    }

    // All analyzers failed
    return nil, fmt.Errorf("all analyzers failed: %v", errors)
}

// isRecoverable checks if error is recoverable
func isRecoverable(err error) bool {
    return err == ErrMissingDependency || err == ErrNotApplicable
}
```

### 3. Partial Results

```go
package analyzer

// PartialResult allows returning partial analysis results
type PartialResult struct {
    Result   *AnalysisResult
    Errors   []*AnalysisError
    Warnings []string
}

// IsComplete checks if analysis completed fully
func (pr *PartialResult) IsComplete() bool {
    return len(pr.Errors) == 0
}

// HasCriticalErrors checks for non-recoverable errors
func (pr *PartialResult) HasCriticalErrors() bool {
    for _, err := range pr.Errors {
        if !err.Recoverable {
            return true
        }
    }
    return false
}
```

---

## Performance and Concurrency

### 1. Parallel Execution

```go
package analyzer

import (
    "context"
    "sync"
)

// ParallelExecutor executes analyzers in parallel
type ParallelExecutor struct {
    registry   *Registry
    maxWorkers int
}

// NewParallelExecutor creates a parallel executor
func NewParallelExecutor(registry *Registry, maxWorkers int) *ParallelExecutor {
    return &ParallelExecutor{
        registry:   registry,
        maxWorkers: maxWorkers,
    }
}

// Execute analyzes multiple files in parallel
func (pe *ParallelExecutor) Execute(
    ctx context.Context,
    files []*FileInfo,
    opts *AnalysisOptions,
) ([]*AnalysisResult, error) {
    // Create semaphore for concurrency control
    sem := make(chan struct{}, pe.maxWorkers)
    results := make([]*AnalysisResult, len(files))
    errors := make([]error, len(files))

    var wg sync.WaitGroup

    for i, file := range files {
        wg.Add(1)

        go func(idx int, f *FileInfo) {
            defer wg.Done()

            // Acquire semaphore
            sem <- struct{}{}
            defer func() { <-sem }()

            // Get applicable analyzer
            analyzers := pe.registry.GetByExtension(f.Extension)
            if len(analyzers) == 0 {
                errors[idx] = fmt.Errorf("no analyzer for %s", f.Path)
                return
            }

            // Execute analysis
            req := &AnalysisRequest{
                File:    f,
                Options: opts,
            }

            result, err := analyzers[0].Analyze(ctx, req)
            results[idx] = result
            errors[idx] = err
        }(i, file)
    }

    wg.Wait()

    // Check for errors
    var firstError error
    for _, err := range errors {
        if err != nil && firstError == nil {
            firstError = err
        }
    }

    return results, firstError
}
```

### 2. Caching Strategy

```go
package analyzer

import (
    "crypto/sha256"
    "encoding/hex"
    "sync"
)

// ResultCache caches analysis results
type ResultCache struct {
    mu    sync.RWMutex
    cache map[string]*AnalysisResult
}

// NewResultCache creates a result cache
func NewResultCache() *ResultCache {
    return &ResultCache{
        cache: make(map[string]*AnalysisResult),
    }
}

// Get retrieves cached result
func (rc *ResultCache) Get(file *FileInfo) (*AnalysisResult, bool) {
    rc.mu.RLock()
    defer rc.mu.RUnlock()

    key := rc.cacheKey(file)
    result, exists := rc.cache[key]
    return result, exists
}

// Put stores result in cache
func (rc *ResultCache) Put(file *FileInfo, result *AnalysisResult) {
    rc.mu.Lock()
    defer rc.mu.Unlock()

    key := rc.cacheKey(file)
    rc.cache[key] = result
}

// cacheKey generates cache key from file content
func (rc *ResultCache) cacheKey(file *FileInfo) string {
    hash := sha256.Sum256(file.Content)
    return file.Path + ":" + hex.EncodeToString(hash[:])
}

// Clear clears the cache
func (rc *ResultCache) Clear() {
    rc.mu.Lock()
    defer rc.mu.Unlock()
    rc.cache = make(map[string]*AnalysisResult)
}
```

---

## Configuration and Extensibility

### 1. Configuration Structure

```go
package analyzer

import "time"

// AnalyzerConfig configures an analyzer
type AnalyzerConfig struct {
    // Enable/disable analyzer
    Enabled bool

    // Parallel execution
    Parallel bool

    // Cache ASTs
    CacheAST bool

    // Analysis timeout
    Timeout time.Duration

    // Language-specific options
    LanguageOptions map[string]interface{}

    // Custom patterns
    CustomPatterns []string

    // Exclusions
    ExcludePatterns []string

    // Plugin path (for external plugins)
    PluginPath string
}

// GlobalConfig configures the analyzer system
type GlobalConfig struct {
    // Analyzer configurations
    Analyzers map[string]*AnalyzerConfig

    // Parallel execution settings
    MaxConcurrency int

    // Result cache settings
    EnableCache bool
    CacheSize   int

    // Plugin directories
    PluginDirs []string

    // Performance settings
    Performance *PerformanceConfig
}

// PerformanceConfig configures performance parameters
type PerformanceConfig struct {
    // Maximum memory per analyzer (MB)
    MaxMemoryPerAnalyzer int

    // Enable profiling
    EnableProfiling bool

    // Profiling output directory
    ProfilingDir string
}
```

### 2. Configuration Loading

```go
package analyzer

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from file
func LoadConfig(path string) (*GlobalConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }

    var config GlobalConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    // Apply defaults
    applyDefaults(&config)

    return &config, nil
}

// applyDefaults sets default values
func applyDefaults(config *GlobalConfig) {
    if config.MaxConcurrency == 0 {
        config.MaxConcurrency = 4
    }

    if config.Performance == nil {
        config.Performance = &PerformanceConfig{
            MaxMemoryPerAnalyzer: 512,
            EnableProfiling:      false,
        }
    }
}
```

---

## Testing Strategy

### 1. Unit Tests

```go
package goanalyzer_test

import (
    "context"
    "testing"

    "github.com/shipshape/analyzer"
    "github.com/shipshape/analyzer/goanalyzer"
)

func TestGoAnalyzer_DetectTableDrivenTests(t *testing.T) {
    analyzer := goanalyzer.NewGoTestAnalyzer()

    testCode := `
package example

import "testing"

func TestExample(t *testing.T) {
    tests := []struct {
        name string
        input int
        want int
    }{
        {"case1", 1, 2},
        {"case2", 2, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test code
        })
    }
}
`

    file := &analyzer.FileInfo{
        Path:      "example_test.go",
        Name:      "example_test.go",
        Extension: ".go",
        Content:   []byte(testCode),
    }

    req := &analyzer.AnalysisRequest{
        File: file,
    }

    result, err := analyzer.Analyze(context.Background(), req)
    if err != nil {
        t.Fatalf("analysis failed: %v", err)
    }

    if len(result.TestFunctions) != 1 {
        t.Errorf("expected 1 test function, got %d", len(result.TestFunctions))
    }

    if !result.TestFunctions[0].IsTable {
        t.Error("expected table-driven test to be detected")
    }
}
```

### 2. Integration Tests

```go
package analyzer_test

import (
    "context"
    "testing"

    "github.com/shipshape/analyzer"
    "github.com/shipshape/analyzer/goanalyzer"
    "github.com/shipshape/analyzer/pythonanalyzer"
)

func TestRegistry_MultiLanguage(t *testing.T) {
    registry := analyzer.NewRegistry()

    // Register analyzers
    registry.Register("go-test", goanalyzer.NewGoTestAnalyzer())
    registry.Register("python-test", pythonanalyzer.NewPythonTestAnalyzer())

    tests := []struct {
        file     *analyzer.FileInfo
        expected string
    }{
        {
            file: &analyzer.FileInfo{
                Path:      "test.go",
                Extension: ".go",
            },
            expected: "go-test",
        },
        {
            file: &analyzer.FileInfo{
                Path:      "test.py",
                Extension: ".py",
            },
            expected: "python-test",
        },
    }

    for _, tt := range tests {
        analyzers := registry.GetByExtension(tt.file.Extension)
        if len(analyzers) == 0 {
            t.Errorf("no analyzer found for %s", tt.file.Extension)
            continue
        }

        metadata := analyzers[0].GetMetadata()
        if metadata.ID != tt.expected {
            t.Errorf("expected %s, got %s", tt.expected, metadata.ID)
        }
    }
}
```

### 3. Benchmark Tests

```go
package analyzer_test

import (
    "context"
    "testing"

    "github.com/shipshape/analyzer"
    "github.com/shipshape/analyzer/goanalyzer"
)

func BenchmarkGoAnalyzer_SingleFile(b *testing.B) {
    analyzer := goanalyzer.NewGoTestAnalyzer()

    file := &analyzer.FileInfo{
        Path:      "example_test.go",
        Extension: ".go",
        Content:   loadTestFile("testdata/example_test.go"),
    }

    req := &analyzer.AnalysisRequest{
        File: file,
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := analyzer.Analyze(context.Background(), req)
        if err != nil {
            b.Fatalf("analysis failed: %v", err)
        }
    }
}

func BenchmarkParallelExecution_100Files(b *testing.B) {
    registry := analyzer.NewRegistry()
    registry.Register("go-test", goanalyzer.NewGoTestAnalyzer())

    executor := analyzer.NewParallelExecutor(registry, 4)

    files := make([]*analyzer.FileInfo, 100)
    for i := 0; i < 100; i++ {
        files[i] = &analyzer.FileInfo{
            Path:      fmt.Sprintf("test%d.go", i),
            Extension: ".go",
            Content:   generateTestCode(),
        }
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := executor.Execute(context.Background(), files, nil)
        if err != nil {
            b.Fatalf("execution failed: %v", err)
        }
    }
}
```

---

## Implementation Roadmap

### Phase 1: Core Infrastructure (Weeks 1-2)
- [ ] Define all core interfaces
- [ ] Implement Registry with thread safety
- [ ] Implement Lifecycle management
- [ ] Create base test utilities
- [ ] Document interface contracts

### Phase 2: Go Analyzer (Week 3)
- [ ] Implement Go AST analyzer
- [ ] Add table-driven test detection
- [ ] Add parallel test detection
- [ ] Add subtest extraction
- [ ] Create comprehensive test suite
- [ ] Benchmark performance (target: 100 files in <2s)

### Phase 3: Python Analyzer (Week 4)
- [ ] Implement tree-sitter Python analyzer
- [ ] Add pytest fixture detection
- [ ] Add parametrized test detection
- [ ] Add unittest support
- [ ] Create comprehensive test suite
- [ ] Benchmark performance (target: 50 files in <5s)

### Phase 4: JavaScript/TypeScript Analyzer (Week 5)
- [ ] Implement tree-sitter JS/TS analyzer
- [ ] Add test block detection
- [ ] Add hook detection
- [ ] Add mock detection
- [ ] Support Jest and Vitest
- [ ] Create comprehensive test suite
- [ ] Benchmark performance (target: 100 files in <10s)

### Phase 5: Error Handling & Graceful Degradation (Week 6)
- [ ] Implement error type hierarchy
- [ ] Add fallback mechanisms
- [ ] Add partial result handling
- [ ] Test error scenarios
- [ ] Document error handling patterns

### Phase 6: Performance & Concurrency (Week 7)
- [ ] Implement parallel executor
- [ ] Add result caching
- [ ] Add AST caching
- [ ] Optimize memory usage
- [ ] Performance benchmarks across all analyzers

### Phase 7: Configuration & Extensibility (Week 8)
- [ ] Implement configuration system
- [ ] Add plugin discovery
- [ ] Create plugin development guide
- [ ] Add custom analyzer examples
- [ ] Integration tests

### Phase 8: Documentation & Examples (Week 9-10)
- [ ] API documentation
- [ ] Usage examples
- [ ] Plugin development guide
- [ ] Performance tuning guide
- [ ] Migration guide

---

## Examples and Usage Patterns

### Example 1: Basic Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/shipshape/analyzer"
    "github.com/shipshape/analyzer/goanalyzer"
)

func main() {
    // Create registry
    registry := analyzer.NewRegistry()

    // Register Go analyzer
    goAnalyzer := goanalyzer.NewGoTestAnalyzer()
    registry.Register("go-test", goAnalyzer)

    // Initialize analyzer
    config := &analyzer.AnalyzerConfig{
        Enabled:  true,
        CacheAST: true,
    }
    goAnalyzer.Initialize(config)

    // Analyze file
    file := &analyzer.FileInfo{
        Path:      "example_test.go",
        Extension: ".go",
        Content:   loadFile("example_test.go"),
    }

    req := &analyzer.AnalysisRequest{
        File: file,
    }

    result, err := goAnalyzer.Analyze(context.Background(), req)
    if err != nil {
        panic(err)
    }

    // Print results
    fmt.Printf("Found %d test functions\n", len(result.TestFunctions))
    for _, test := range result.TestFunctions {
        fmt.Printf("  - %s (complexity: %d)\n", test.Name, test.Complexity)
    }
}
```

### Example 2: Multi-Language Analysis

```go
package main

import (
    "context"
    "fmt"

    "github.com/shipshape/analyzer"
)

func analyzeRepository(repoPath string) error {
    // Create registry and register all analyzers
    registry := analyzer.NewRegistry()
    RegisterBuiltinAnalyzers(registry)

    // Discover files
    files, err := discoverTestFiles(repoPath)
    if err != nil {
        return err
    }

    // Create parallel executor
    executor := analyzer.NewParallelExecutor(registry, 4)

    // Execute analysis
    results, err := executor.Execute(context.Background(), files, nil)
    if err != nil {
        return err
    }

    // Aggregate results by language
    byLanguage := make(map[string][]*analyzer.AnalysisResult)
    for _, result := range results {
        byLanguage[result.Language] = append(byLanguage[result.Language], result)
    }

    // Print summary
    for lang, langResults := range byLanguage {
        totalTests := 0
        for _, r := range langResults {
            totalTests += r.Metrics.TotalTests
        }
        fmt.Printf("%s: %d tests in %d files\n", lang, totalTests, len(langResults))
    }

    return nil
}
```

### Example 3: Custom Analyzer Plugin

```go
// File: plugins/custom-analyzer.go
package main

import (
    "context"

    "github.com/shipshape/analyzer"
)

// CustomAnalyzer implements a custom analyzer
type CustomAnalyzer struct {
    config *analyzer.AnalyzerConfig
}

// NewAnalyzer is the plugin entry point
func NewAnalyzer(config *analyzer.AnalyzerConfig) (analyzer.Analyzer, error) {
    return &CustomAnalyzer{config: config}, nil
}

// Implement Analyzer interface...
func (ca *CustomAnalyzer) Analyze(
    ctx context.Context,
    req *analyzer.AnalysisRequest,
) (*analyzer.AnalysisResult, error) {
    // Custom analysis logic
    return &analyzer.AnalysisResult{
        FilePath: req.File.Path,
        Language: "custom",
    }, nil
}

// ... implement other Analyzer methods ...
```

Build as plugin:
```bash
go build -buildmode=plugin -o custom-analyzer.so plugins/custom-analyzer.go
```

Load plugin:
```go
discovery := analyzer.NewPluginDiscovery(registry, []string{"./plugins"})
discovery.DiscoverPlugins()
```

---

## Conclusion

This Analyzer Plugin System design provides a **robust, extensible, and performant architecture** for multi-language test analysis in Ship Shape. Key achievements:

1. **Clean Interfaces**: Well-defined contracts for analyzers
2. **Extensibility**: Easy to add new language analyzers
3. **Performance**: Parallel execution with caching
4. **Reliability**: Graceful degradation and error handling
5. **Flexibility**: Configuration-driven behavior

The design draws from proven patterns in AgentReady while addressing Ship Shape's specific requirements for multi-language AST parsing, monorepo support, and performance at scale.

**Next Steps**:
1. Review and approve this design
2. Begin Phase 1 implementation (core infrastructure)
3. Execute SPIKE-001 to validate technical approaches
4. Iterate based on spike findings

---

**Document Status**: Ready for Review
**Approval Required**: Architecture Team, Technical Leads
**Related Documents**:
- [Requirements](.project/requirements.md)
- [Architecture](.project/architecture.md)
- [SPIKE-001: AST Parsing Framework](.project/v1.0.0/spikes/SPIKE-001-ast-parsing-framework.md)
- [User Stories](.project/user-stories.md)
