# Ship Shape - Technical Architecture Document

**Version**: 1.0.0
**Date**: 2026-01-27
**Status**: Draft
**Author**: Senior Software Engineer

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Core Architecture](#core-architecture)
3. [Component Architecture](#component-architecture)
4. [Type System and Interfaces](#type-system-and-interfaces)
5. [Repository Discovery Engine](#repository-discovery-engine)
6. [Analyzer Framework](#analyzer-framework)
7. [Assessor Framework](#assessor-framework)
8. [Execution Engine](#execution-engine)
9. [Report Generation System](#report-generation-system)
10. [Quality Gate System](#quality-gate-system)
11. [Data Models](#data-models)
12. [Plugin Architecture](#plugin-architecture)
13. [Multi-Language Support](#multi-language-support)
14. [Monorepo Handling](#monorepo-handling)
15. [Storage and Persistence](#storage-and-persistence)
16. [Integration Architecture](#integration-architecture)
17. [Security Architecture](#security-architecture)
18. [Deployment Architecture](#deployment-architecture)
19. [Future API Design](#future-api-design)

---

## System Overview

### Architecture Philosophy

Ship Shape is designed as a **modular, extensible, repository-aware testing quality analysis platform** built on the following principles:

1. **Context-First Analysis**: Discover repository structure before analysis
2. **Language Isolation**: Each language analyzer is independent and composable
3. **Type-Safe Core**: Go-based core for performance and reliability
4. **Plugin Architecture**: Extensible through well-defined interfaces
5. **Parallel Execution**: Concurrent analysis where dependencies allow
6. **Graceful Degradation**: Partial results on failures, not complete failure

### High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Ship Shape CLI                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Command    │  │    Config    │  │   Output     │              │
│  │   Handler    │  │   Manager    │  │   Writer     │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Core Engine (coordinator)                       │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │              Repository Context (shared state)                 │  │
│  │  - Languages, Frameworks, Structure, Monorepo Config          │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐             │
│  │  Discovery   │─▶│  Execution   │─▶│   Report     │             │
│  │   Engine     │  │   Engine     │  │  Generator   │             │
│  └──────────────┘  └──────────────┘  └──────────────┘             │
└─────────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌────────────────┐    ┌──────────────┐
│   Analyzer    │    │   Assessor     │    │  Tool        │
│   Registry    │    │   Registry     │    │  Detector    │
└───────────────┘    └────────────────┘    └──────────────┘
        │                     │                     │
        │                     │                     │
   ┌────┴────┬────────┬──────┴──────┬──────────────┴─────┐
   │         │        │             │                     │
   ▼         ▼        ▼             ▼                     ▼
┌─────┐  ┌─────┐  ┌─────┐       ┌─────┐            ┌─────────┐
│ Go  │  │Python│ │JS/TS│  ...  │Score│    ...     │ Tool DB │
│Anlzr│  │Anlzr │ │Anlzr│       │Assor│            │         │
└─────┘  └─────┘  └─────┘       └─────┘            └─────────┘
```

### Data Flow

```
Repository Path
      │
      ▼
[1. Discovery Phase]
      │
      ├─▶ Language Detection ─────┐
      ├─▶ Framework Detection ─────┤
      ├─▶ Monorepo Detection ──────┤──▶ Repository Context
      ├─▶ Structure Mapping ───────┤
      └─▶ Dependency Analysis ─────┘
      │
      ▼
[2. Analyzer Selection]
      │
      ├─▶ Match Analyzers to Languages
      ├─▶ Build Dependency Graph
      └─▶ Create Execution Plan
      │
      ▼
[3. Parallel Execution]
      │
      ├─▶ Layer 1: Test Discovery
      ├─▶ Layer 2: Structure Analysis
      ├─▶ Layer 3: Quality Analysis
      ├─▶ Layer 4: Coverage Analysis
      ├─▶ Layer 5: Execution Analysis (optional)
      ├─▶ Layer 6: Mutation Testing (optional)
      └─▶ Layer 7: CI/CD Analysis
      │
      ▼
[4. Assessment]
      │
      ├─▶ Calculate Dimension Scores
      ├─▶ Apply Language-Specific Weights
      ├─▶ Generate Findings
      └─▶ Compute Overall Score
      │
      ▼
[5. Report Generation]
      │
      ├─▶ HTML Report
      ├─▶ JSON Output
      ├─▶ CLI Summary
      └─▶ Quality Gate Evaluation
      │
      ▼
Exit Code + Reports
```

---

## Core Architecture

### Package Structure

```
shipshape/
├── cmd/
│   └── shipshape/           # CLI entry point
│       ├── main.go
│       ├── analyze.go       # analyze command
│       ├── gate.go          # gate command
│       ├── tools.go         # tools command
│       └── report.go        # report command
│
├── pkg/
│   ├── core/                # Core types and interfaces
│   │   ├── context.go       # Repository context
│   │   ├── analyzer.go      # Analyzer interface
│   │   ├── assessor.go      # Assessor interface
│   │   ├── finding.go       # Finding types
│   │   ├── result.go        # Result types
│   │   └── config.go        # Configuration types
│   │
│   ├── discovery/           # Repository discovery
│   │   ├── discovery.go     # Main discovery engine
│   │   ├── language.go      # Language detection
│   │   ├── framework.go     # Framework detection
│   │   ├── monorepo.go      # Monorepo detection
│   │   └── structure.go     # Directory structure mapping
│   │
│   ├── engine/              # Execution engine
│   │   ├── engine.go        # Main execution orchestrator
│   │   ├── scheduler.go     # Parallel execution scheduler
│   │   ├── registry.go      # Analyzer/Assessor registry
│   │   └── pipeline.go      # Analysis pipeline
│   │
│   ├── analyzers/           # Language-specific analyzers
│   │   ├── base/            # Base analyzer utilities
│   │   ├── go/              # Go analyzers
│   │   │   ├── testing.go
│   │   │   ├── coverage.go
│   │   │   └── quality.go
│   │   ├── python/          # Python analyzers
│   │   │   ├── pytest.go
│   │   │   ├── unittest.go
│   │   │   └── coverage.go
│   │   ├── javascript/      # JavaScript/TypeScript analyzers
│   │   │   ├── jest.go
│   │   │   ├── vitest.go
│   │   │   └── cypress.go
│   │   └── ...              # Other languages
│   │
│   ├── assessors/           # Scoring assessors
│   │   ├── coverage.go      # Coverage scoring
│   │   ├── quality.go       # Quality scoring
│   │   ├── performance.go   # Performance scoring
│   │   ├── tools.go         # Tool adoption scoring
│   │   └── pyramid.go       # Test pyramid scoring
│   │
│   ├── tools/               # Tool detection and analysis
│   │   ├── detector.go      # Tool detector
│   │   ├── registry.go      # Tool registry/database
│   │   └── recommender.go   # Tool recommender
│   │
│   ├── coverage/            # Coverage integration
│   │   ├── parser.go        # Coverage report parser
│   │   ├── cobertura.go     # Cobertura XML parser
│   │   ├── lcov.go          # LCOV parser
│   │   └── go.go            # Go coverage parser
│   │
│   ├── quality/             # Test quality analysis
│   │   ├── smells.go        # Test smell detection
│   │   ├── assertions.go    # Assertion analysis
│   │   ├── naming.go        # Naming conventions
│   │   └── independence.go  # Test independence checker
│   │
│   ├── report/              # Report generation
│   │   ├── html.go          # HTML report generator
│   │   ├── json.go          # JSON report generator
│   │   ├── cli.go           # CLI output formatter
│   │   └── templates/       # HTML templates
│   │
│   ├── gates/               # Quality gates
│   │   ├── evaluator.go     # Gate evaluator
│   │   ├── config.go        # Gate configuration
│   │   └── trend.go         # Trend-based gates
│   │
│   ├── storage/             # Data persistence
│   │   ├── sqlite.go        # SQLite storage (historical data)
│   │   ├── file.go          # File-based storage
│   │   └── models.go        # Storage models
│   │
│   └── util/                # Utilities
│       ├── ast.go           # AST parsing helpers
│       ├── exec.go          # Command execution helpers
│       ├── fs.go            # Filesystem helpers
│       └── parallel.go      # Parallel execution helpers
│
├── internal/
│   └── research/            # Research data and citations
│       ├── citations.go     # Citation database
│       └── thresholds.go    # Evidence-based thresholds
│
├── configs/
│   ├── default.yml          # Default configuration
│   └── examples/            # Example configurations
│
├── data/
│   └── tools/               # Tool database (YAML/JSON)
│       ├── python.yml
│       ├── go.yml
│       ├── javascript.yml
│       └── ...
│
└── scripts/
    ├── build.sh
    └── install.sh
```

---

## Component Architecture

### 1. CLI Layer

**Responsibility**: Command-line interface and user interaction

```go
// cmd/shipshape/main.go
package main

type RootCommand struct {
    ConfigPath string
    Verbose    bool
    NoColor    bool
}

type AnalyzeCommand struct {
    Path         string
    Config       string
    Output       string
    Format       string // html, json, text, markdown
    FailOn       []string // severity levels to fail on
    SkipLayers   []string // layers to skip (e.g., "mutation", "execution")
    Parallel     int // max parallel workers
}

type GateCommand struct {
    Config       string
    Baseline     string // path to baseline results
    TrendWindow  int    // number of runs for trend analysis
}

type ToolsCommand struct {
    Action       string // list, setup, check
    Tool         string // specific tool name
    Language     string // filter by language
}
```

### 2. Core Engine

**Responsibility**: Orchestrates the entire analysis workflow

```go
// pkg/core/engine.go
package core

type Engine struct {
    config         *Config
    discoveryEng   *discovery.Engine
    executionEng   *engine.Executor
    reportGen      *report.Generator
    gateEvaluator  *gates.Evaluator
    storage        storage.Storage
}

func NewEngine(config *Config) (*Engine, error) {
    return &Engine{
        config:        config,
        discoveryEng:  discovery.New(),
        executionEng:  engine.New(config.Parallel),
        reportGen:     report.New(),
        gateEvaluator: gates.New(config.Gates),
        storage:       storage.NewSQLite(config.StoragePath),
    }, nil
}

func (e *Engine) Analyze(ctx context.Context, repoPath string) (*Report, error) {
    // 1. Discover repository context
    repoCtx, err := e.discoveryEng.Discover(ctx, repoPath)
    if err != nil {
        return nil, fmt.Errorf("discovery failed: %w", err)
    }

    // 2. Select and execute analyzers
    results, err := e.executionEng.Execute(ctx, repoCtx)
    if err != nil {
        return nil, fmt.Errorf("execution failed: %w", err)
    }

    // 3. Generate report
    report, err := e.reportGen.Generate(repoCtx, results)
    if err != nil {
        return nil, fmt.Errorf("report generation failed: %w", err)
    }

    // 4. Evaluate quality gates (if configured)
    if e.config.Gates != nil {
        gateResults, err := e.gateEvaluator.Evaluate(report)
        if err != nil {
            return nil, fmt.Errorf("gate evaluation failed: %w", err)
        }
        report.GateResults = gateResults
    }

    // 5. Store results for trend analysis
    if e.storage != nil {
        if err := e.storage.StoreReport(report); err != nil {
            // Log warning but don't fail
            log.Warnf("failed to store report: %v", err)
        }
    }

    return report, nil
}
```

### 3. Discovery Engine

**Responsibility**: Analyzes repository structure and builds context

```go
// pkg/discovery/discovery.go
package discovery

type Engine struct {
    langDetector      *LanguageDetector
    frameworkDetector *FrameworkDetector
    monorepoDetector  *MonorepoDetector
    structureMapper   *StructureMapper
}

func (e *Engine) Discover(ctx context.Context, repoPath string) (*RepositoryContext, error) {
    ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()

    repoCtx := &RepositoryContext{
        RootPath:  repoPath,
        StartTime: time.Now(),
    }

    // Parallel discovery (all independent)
    g, gctx := errgroup.WithContext(ctx)

    // Language detection
    g.Go(func() error {
        langs, err := e.langDetector.Detect(gctx, repoPath)
        if err != nil {
            return err
        }
        repoCtx.Languages = langs
        return nil
    })

    // Structure mapping
    g.Go(func() error {
        structure, err := e.structureMapper.Map(gctx, repoPath)
        if err != nil {
            return err
        }
        repoCtx.Structure = structure
        return nil
    })

    // Monorepo detection
    g.Go(func() error {
        monorepo, err := e.monorepoDetector.Detect(gctx, repoPath)
        if err != nil {
            return err
        }
        repoCtx.Monorepo = monorepo
        return nil
    })

    if err := g.Wait(); err != nil {
        return nil, err
    }

    // Framework detection (depends on language detection)
    frameworks, err := e.frameworkDetector.Detect(ctx, repoCtx)
    if err != nil {
        return nil, err
    }
    repoCtx.Frameworks = frameworks

    repoCtx.DiscoveryDuration = time.Since(repoCtx.StartTime)
    return repoCtx, nil
}
```

### 4. Execution Engine

**Responsibility**: Schedules and executes analyzers in correct order

```go
// pkg/engine/engine.go
package engine

type Executor struct {
    registry      *Registry
    scheduler     *Scheduler
    maxWorkers    int
}

type ExecutionPlan struct {
    Layers        []*LayerPlan
    TotalAnalyzers int
}

type LayerPlan struct {
    Name          string
    Analyzers     []Analyzer
    Dependencies  []string // layer names this depends on
}

func (e *Executor) Execute(ctx context.Context, repoCtx *RepositoryContext) (*AnalysisResults, error) {
    // 1. Select applicable analyzers based on repository context
    analyzers := e.registry.SelectAnalyzers(repoCtx)

    // 2. Build execution plan (topological sort of dependencies)
    plan, err := e.buildExecutionPlan(analyzers)
    if err != nil {
        return nil, err
    }

    // 3. Execute layers in order, analyzers within layer in parallel
    results := &AnalysisResults{
        RepositoryContext: repoCtx,
        LayerResults:      make(map[string]*LayerResult),
    }

    for _, layer := range plan.Layers {
        layerResult, err := e.executeLayer(ctx, layer, repoCtx, results)
        if err != nil {
            // Check if error is critical or allows graceful degradation
            if layer.Required {
                return nil, fmt.Errorf("critical layer %s failed: %w", layer.Name, err)
            }
            // Log and continue for non-critical layers
            log.Warnf("layer %s failed but continuing: %v", layer.Name, err)
        }
        results.LayerResults[layer.Name] = layerResult
    }

    return results, nil
}

func (e *Executor) executeLayer(ctx context.Context, layer *LayerPlan,
                                repoCtx *RepositoryContext,
                                previousResults *AnalysisResults) (*LayerResult, error) {

    layerResult := &LayerResult{
        Name:      layer.Name,
        StartTime: time.Now(),
        Findings:  make([]*Finding, 0),
    }

    // Execute analyzers in parallel
    resultsChan := make(chan *AnalyzerResult, len(layer.Analyzers))
    errorsChan := make(chan error, len(layer.Analyzers))

    sem := make(chan struct{}, e.maxWorkers)

    for _, analyzer := range layer.Analyzers {
        sem <- struct{}{} // acquire
        go func(a Analyzer) {
            defer func() { <-sem }() // release

            result, err := a.Analyze(ctx, repoCtx, previousResults)
            if err != nil {
                errorsChan <- fmt.Errorf("%s: %w", a.Name(), err)
                return
            }
            resultsChan <- result
        }(analyzer)
    }

    // Wait for all analyzers
    var errs []error
    for i := 0; i < len(layer.Analyzers); i++ {
        select {
        case result := <-resultsChan:
            layerResult.Findings = append(layerResult.Findings, result.Findings...)
        case err := <-errorsChan:
            errs = append(errs, err)
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }

    layerResult.Duration = time.Since(layerResult.StartTime)

    if len(errs) > 0 {
        return layerResult, fmt.Errorf("layer had %d errors: %v", len(errs), errs)
    }

    return layerResult, nil
}
```

---

## Type System and Interfaces

### Core Interfaces

```go
// pkg/core/analyzer.go
package core

// Analyzer interface - all analyzers must implement this
type Analyzer interface {
    // Name returns the unique identifier for this analyzer
    Name() string

    // AppliesTo determines if this analyzer should run based on repository context
    AppliesTo(ctx *RepositoryContext) bool

    // Layer returns which analysis layer this analyzer belongs to
    Layer() AnalysisLayer

    // Dependencies returns analyzers that must run before this one
    Dependencies() []string

    // Analyze performs the analysis and returns findings
    Analyze(ctx context.Context, repoCtx *RepositoryContext,
            previousResults *AnalysisResults) (*AnalyzerResult, error)
}

// AnalysisLayer represents the execution phase
type AnalysisLayer int

const (
    LayerDiscovery AnalysisLayer = iota  // Built-in, not an analyzer
    LayerStructure                        // Test file discovery and categorization
    LayerQuality                          // Test quality analysis (smells, naming, etc.)
    LayerCoverage                         // Coverage analysis
    LayerExecution                        // Test execution and performance
    LayerMutation                         // Mutation testing
    LayerCICD                             // CI/CD integration analysis
)

// Assessor interface - scoring components
type Assessor interface {
    // Name returns the unique identifier for this assessor
    Name() string

    // Dimension returns which scoring dimension this assessor contributes to
    Dimension() ScoreDimension

    // AppliesTo determines if this assessor should run
    AppliesTo(ctx *RepositoryContext) bool

    // Assess calculates a score based on analysis results
    Assess(ctx *RepositoryContext, results *AnalysisResults) (*DimensionScore, error)
}

// ScoreDimension represents a scoring category
type ScoreDimension string

const (
    DimensionCoverage      ScoreDimension = "test_coverage"
    DimensionQuality       ScoreDimension = "test_quality"
    DimensionPerformance   ScoreDimension = "test_performance"
    DimensionTools         ScoreDimension = "tool_adoption"
    DimensionMaintain      ScoreDimension = "maintainability"
    DimensionCodeQuality   ScoreDimension = "code_quality"
)
```

### Core Data Types

```go
// pkg/core/context.go
package core

// RepositoryContext holds all discovered information about the repository
type RepositoryContext struct {
    RootPath           string
    StartTime          time.Time
    DiscoveryDuration  time.Duration

    // Repository type
    Type               RepositoryType // SingleLanguage, MultiLanguage, Monorepo

    // Language information
    Languages          []*LanguageInfo

    // Framework information
    Frameworks         []*FrameworkInfo

    // Monorepo information (nil if not a monorepo)
    Monorepo           *MonorepoInfo

    // Directory structure
    Structure          *StructureInfo

    // Build system
    BuildSystem        *BuildSystemInfo

    // Version control
    VCS                *VCSInfo
}

type RepositoryType string

const (
    TypeSingleLanguage RepositoryType = "single-language"
    TypeMultiLanguage  RepositoryType = "multi-language"
    TypeMonorepo       RepositoryType = "monorepo"
)

type LanguageInfo struct {
    Name            string
    Percentage      float64
    FileCount       int
    IsPrimary       bool // >10% presence
    Version         string // if detectable

    // Detected infrastructure
    TestFrameworks  []string
    CoverageTools   []string
    Linters         []string
    Formatters      []string
    BuildTools      []string
}

type FrameworkInfo struct {
    Name            string
    Language        string
    Version         string
    ConfigFiles     []string
    Type            FrameworkType // Testing, Coverage, Linting, etc.
}

type FrameworkType string

const (
    FrameworkTypeTesting    FrameworkType = "testing"
    FrameworkTypeCoverage   FrameworkType = "coverage"
    FrameworkTypeLinting    FrameworkType = "linting"
    FrameworkTypeFormatting FrameworkType = "formatting"
    FrameworkTypeBuild      FrameworkType = "build"
    FrameworkTypeE2E        FrameworkType = "e2e"
)

type MonorepoInfo struct {
    Type            string // npm-workspaces, lerna, nx, go-workspaces, etc.
    ConfigFile      string
    Packages        []*PackageInfo
    SharedInfra     *SharedInfrastructure
}

type PackageInfo struct {
    Name            string
    Path            string
    Language        string
    TestFramework   string
    Dependencies    []string // other package names in monorepo
    Private         bool
}

type SharedInfrastructure struct {
    RootConfigs     map[string]string // config type -> file path
    CommonCI        bool
    SharedDeps      []string
}

type StructureInfo struct {
    SourceDirs      []string
    TestDirs        []string
    ConfigDirs      []string
    TestFiles       []*TestFileInfo
    TotalLOC        int64
    TestLOC         int64
}

type TestFileInfo struct {
    Path            string
    Language        string
    Framework       string
    TestCount       int
    TestTypes       []TestType // unit, integration, e2e, etc.
}

type TestType string

const (
    TestTypeUnit        TestType = "unit"
    TestTypeIntegration TestType = "integration"
    TestTypeE2E         TestType = "e2e"
    TestTypePerformance TestType = "performance"
    TestTypeContract    TestType = "contract"
    TestTypeSmoke       TestType = "smoke"
    TestTypeRegression  TestType = "regression"
)
```

### Finding and Result Types

```go
// pkg/core/finding.go
package core

type Finding struct {
    ID              string           // Unique identifier
    CheckID         string           // Analyzer check identifier
    Type            FindingType
    Severity        Severity
    Title           string
    Description     string
    Location        *Location
    Evidence        []string         // Code snippets, log lines, etc.
    Remediation     *Remediation
    References      []*Reference
    Confidence      float64          // 0.0-1.0
    Metadata        map[string]interface{}
}

type FindingType string

const (
    FindingTypeCoverage       FindingType = "coverage"
    FindingTypeQuality        FindingType = "quality"
    FindingTypePerformance    FindingType = "performance"
    FindingTypeMaintainability FindingType = "maintainability"
    FindingTypeCodeQuality    FindingType = "code_quality"
    FindingTypeDocumentation  FindingType = "documentation"
    FindingTypeBestPractice   FindingType = "best_practice"
    FindingTypeSecurity       FindingType = "security"
)

type Severity string

const (
    SeverityInfo     Severity = "info"
    SeverityLow      Severity = "low"
    SeverityMedium   Severity = "medium"
    SeverityHigh     Severity = "high"
    SeverityCritical Severity = "critical"
)

type Location struct {
    FilePath        string
    StartLine       int
    EndLine         int
    StartColumn     int
    EndColumn       int
}

type Remediation struct {
    Description     string
    Effort          EffortLevel
    AutoFix         *AutoFix
    Examples        []string
}

type EffortLevel string

const (
    EffortMinimal   EffortLevel = "minimal"   // <15 min
    EffortLow       EffortLevel = "low"       // 15-60 min
    EffortMedium    EffortLevel = "medium"    // 1-4 hours
    EffortHigh      EffortLevel = "high"      // 4-8 hours
    EffortVeryHigh  EffortLevel = "very_high" // >8 hours
)

type AutoFix struct {
    Available       bool
    Description     string
    Command         string
    Safe            bool // can be applied automatically without review
}

type Reference struct {
    Title           string
    URL             string
    DOI             string
    Summary         string
}

// pkg/core/result.go
package core

type AnalysisResults struct {
    RepositoryContext *RepositoryContext
    LayerResults      map[string]*LayerResult
    StartTime         time.Time
    EndTime           time.Time
}

type LayerResult struct {
    Name              string
    StartTime         time.Time
    Duration          time.Duration
    Findings          []*Finding
    Metadata          map[string]interface{}
}

type AnalyzerResult struct {
    AnalyzerName      string
    Success           bool
    Error             error
    Findings          []*Finding
    Metrics           map[string]interface{} // raw metrics for assessors
    Duration          time.Duration
}

type DimensionScore struct {
    Dimension         ScoreDimension
    Score             float64 // 0-100
    Weight            float64 // contribution to overall score
    Components        map[string]float64 // sub-scores
    Rationale         string
}

type Report struct {
    RepositoryContext *RepositoryContext
    Results           *AnalysisResults
    Scores            *ScoreCard
    Findings          []*Finding
    ToolAnalysis      *ToolAnalysis
    GateResults       *GateResults
    GeneratedAt       time.Time
}

type ScoreCard struct {
    Overall           float64
    Grade             string // A+, A, B, C, D, F
    Dimensions        map[ScoreDimension]*DimensionScore
}

type ToolAnalysis struct {
    OverallScore      float64
    Categories        map[string]*ToolCategoryScore
    Detected          []*DetectedTool
    Missing           []*MissingTool
}

type ToolCategoryScore struct {
    Category          string // Testing, Quality, Security, Automation
    Score             float64
    DetectedCount     int
    AvailableCount    int
    EssentialMissing  int
}

type DetectedTool struct {
    Name              string
    Version           string
    Category          string
    Configured        bool
    ConfigQuality     float64 // 0-1
}

type MissingTool struct {
    Name              string
    Category          string
    Priority          string // Essential, Recommended, Advanced
    Benefit           string
    Effort            EffortLevel
    SetupGuide        string
}
```

---

## Repository Discovery Engine

### Language Detection

```go
// pkg/discovery/language.go
package discovery

type LanguageDetector struct {
    extensionMap      map[string]string // .go -> Go, .py -> Python
    configDetectors   map[string]*ConfigDetector
}

func (ld *LanguageDetector) Detect(ctx context.Context, repoPath string) ([]*LanguageInfo, error) {
    // 1. Count files by extension
    extensionCounts := make(map[string]int)
    err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }

        // Skip common ignored directories
        if shouldSkip(path) {
            return nil
        }

        ext := filepath.Ext(path)
        if lang, ok := ld.extensionMap[ext]; ok {
            extensionCounts[lang]++
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    // 2. Convert counts to percentages
    totalFiles := 0
    for _, count := range extensionCounts {
        totalFiles += count
    }

    languages := make([]*LanguageInfo, 0)
    for lang, count := range extensionCounts {
        percentage := (float64(count) / float64(totalFiles)) * 100
        isPrimary := percentage >= 10.0

        langInfo := &LanguageInfo{
            Name:       lang,
            Percentage: percentage,
            FileCount:  count,
            IsPrimary:  isPrimary,
        }

        // 3. Detect language-specific infrastructure
        ld.detectInfrastructure(repoPath, langInfo)

        languages = append(languages, langInfo)
    }

    // Sort by percentage (descending)
    sort.Slice(languages, func(i, j int) bool {
        return languages[i].Percentage > languages[j].Percentage
    })

    return languages, nil
}

func (ld *LanguageDetector) detectInfrastructure(repoPath string, lang *LanguageInfo) {
    switch lang.Name {
    case "Python":
        ld.detectPythonInfra(repoPath, lang)
    case "Go":
        ld.detectGoInfra(repoPath, lang)
    case "JavaScript", "TypeScript":
        ld.detectJSInfra(repoPath, lang)
    // ... other languages
    }
}

func (ld *LanguageDetector) detectPythonInfra(repoPath string, lang *LanguageInfo) {
    // Check for pytest
    if fileExists(filepath.Join(repoPath, "pytest.ini")) ||
       fileExists(filepath.Join(repoPath, "pyproject.toml")) {
        lang.TestFrameworks = append(lang.TestFrameworks, "pytest")
    }

    // Check for coverage.py
    if fileExists(filepath.Join(repoPath, ".coveragerc")) {
        lang.CoverageTools = append(lang.CoverageTools, "coverage.py")
    }

    // Check requirements.txt or pyproject.toml for tools
    deps := ld.parsePythonDependencies(repoPath)
    for _, dep := range deps {
        if strings.HasPrefix(dep, "pytest") {
            lang.TestFrameworks = appendUnique(lang.TestFrameworks, "pytest")
        }
        if dep == "coverage" {
            lang.CoverageTools = appendUnique(lang.CoverageTools, "coverage.py")
        }
        // ... more tool detection
    }
}

func shouldSkip(path string) bool {
    skipDirs := []string{
        "node_modules", ".git", ".venv", "venv", "__pycache__",
        "vendor", "target", "build", "dist", ".next",
    }
    for _, skip := range skipDirs {
        if strings.Contains(path, skip) {
            return true
        }
    }
    return false
}
```

### Framework Detection

```go
// pkg/discovery/framework.go
package discovery

type FrameworkDetector struct {
    detectors map[string]FrameworkDetectorFunc
}

type FrameworkDetectorFunc func(repoPath string, lang *LanguageInfo) ([]*FrameworkInfo, error)

func (fd *FrameworkDetector) Detect(ctx context.Context, repoCtx *RepositoryContext) ([]*FrameworkInfo, error) {
    frameworks := make([]*FrameworkInfo, 0)

    for _, lang := range repoCtx.Languages {
        if detector, ok := fd.detectors[lang.Name]; ok {
            detected, err := detector(repoCtx.RootPath, lang)
            if err != nil {
                log.Warnf("framework detection failed for %s: %v", lang.Name, err)
                continue
            }
            frameworks = append(frameworks, detected...)
        }
    }

    return frameworks, nil
}

func detectPytestFramework(repoPath string, lang *LanguageInfo) ([]*FrameworkInfo, error) {
    frameworks := make([]*FrameworkInfo, 0)

    // Check for pytest.ini
    pytestIni := filepath.Join(repoPath, "pytest.ini")
    if fileExists(pytestIni) {
        frameworks = append(frameworks, &FrameworkInfo{
            Name:        "pytest",
            Language:    "Python",
            ConfigFiles: []string{pytestIni},
            Type:        FrameworkTypeTesting,
        })
    }

    // Check pyproject.toml for [tool.pytest.ini_options]
    pyproject := filepath.Join(repoPath, "pyproject.toml")
    if fileExists(pyproject) {
        content, _ := os.ReadFile(pyproject)
        if strings.Contains(string(content), "[tool.pytest.ini_options]") {
            frameworks = append(frameworks, &FrameworkInfo{
                Name:        "pytest",
                Language:    "Python",
                ConfigFiles: []string{pyproject},
                Type:        FrameworkTypeTesting,
            })
        }
    }

    return frameworks, nil
}
```

### Monorepo Detection

```go
// pkg/discovery/monorepo.go
package discovery

type MonorepoDetector struct {
    detectors []MonorepoDetectorFunc
}

type MonorepoDetectorFunc func(repoPath string) (*MonorepoInfo, error)

func (md *MonorepoDetector) Detect(ctx context.Context, repoPath string) (*MonorepoInfo, error) {
    for _, detector := range md.detectors {
        info, err := detector(repoPath)
        if err != nil {
            continue
        }
        if info != nil {
            return info, nil
        }
    }
    return nil, nil // not a monorepo
}

func detectNPMWorkspaces(repoPath string) (*MonorepoInfo, error) {
    packageJSON := filepath.Join(repoPath, "package.json")
    if !fileExists(packageJSON) {
        return nil, nil
    }

    var pkg struct {
        Workspaces interface{} `json:"workspaces"`
    }

    data, err := os.ReadFile(packageJSON)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(data, &pkg); err != nil {
        return nil, err
    }

    if pkg.Workspaces == nil {
        return nil, nil
    }

    // Parse workspaces (can be array or object)
    var workspacePatterns []string
    switch v := pkg.Workspaces.(type) {
    case []interface{}:
        for _, pattern := range v {
            if str, ok := pattern.(string); ok {
                workspacePatterns = append(workspacePatterns, str)
            }
        }
    case map[string]interface{}:
        if packages, ok := v["packages"].([]interface{}); ok {
            for _, pattern := range packages {
                if str, ok := pattern.(string); ok {
                    workspacePatterns = append(workspacePatterns, str)
                }
            }
        }
    }

    if len(workspacePatterns) == 0 {
        return nil, nil
    }

    // Find all package.json files matching workspace patterns
    packages := make([]*PackageInfo, 0)
    for _, pattern := range workspacePatterns {
        matches, err := filepath.Glob(filepath.Join(repoPath, pattern, "package.json"))
        if err != nil {
            continue
        }

        for _, match := range matches {
            pkgInfo, err := parseNPMPackage(match)
            if err != nil {
                log.Warnf("failed to parse package %s: %v", match, err)
                continue
            }
            packages = append(packages, pkgInfo)
        }
    }

    return &MonorepoInfo{
        Type:       "npm-workspaces",
        ConfigFile: packageJSON,
        Packages:   packages,
        SharedInfra: detectSharedInfra(repoPath, packages),
    }, nil
}

func detectGoWorkspaces(repoPath string) (*MonorepoInfo, error) {
    goWork := filepath.Join(repoPath, "go.work")
    if !fileExists(goWork) {
        return nil, nil
    }

    // Parse go.work file
    content, err := os.ReadFile(goWork)
    if err != nil {
        return nil, err
    }

    packages := make([]*PackageInfo, 0)
    lines := strings.Split(string(content), "\n")
    inUse := false

    for _, line := range lines {
        line = strings.TrimSpace(line)

        if strings.HasPrefix(line, "use") {
            inUse = true
            continue
        }

        if inUse && line != "" && !strings.HasPrefix(line, ")") {
            modulePath := strings.Trim(line, "./")
            pkgPath := filepath.Join(repoPath, modulePath)

            pkgInfo := &PackageInfo{
                Name:     modulePath,
                Path:     pkgPath,
                Language: "Go",
            }

            // Detect test framework by checking for *_test.go files
            if hasGoTestFiles(pkgPath) {
                pkgInfo.TestFramework = "testing"
            }

            packages = append(packages, pkgInfo)
        }

        if strings.HasPrefix(line, ")") {
            inUse = false
        }
    }

    return &MonorepoInfo{
        Type:       "go-workspaces",
        ConfigFile: goWork,
        Packages:   packages,
        SharedInfra: detectSharedInfra(repoPath, packages),
    }, nil
}
```

---

## Analyzer Framework

### Base Analyzer

```go
// pkg/analyzers/base/analyzer.go
package base

// BaseAnalyzer provides common functionality for all analyzers
type BaseAnalyzer struct {
    name         string
    layer        core.AnalysisLayer
    dependencies []string
}

func (ba *BaseAnalyzer) Name() string {
    return ba.name
}

func (ba *BaseAnalyzer) Layer() core.AnalysisLayer {
    return ba.layer
}

func (ba *BaseAnalyzer) Dependencies() []string {
    return ba.dependencies
}

// Helper methods for creating findings
func (ba *BaseAnalyzer) NewFinding(checkID, title, description string,
                                   severity core.Severity,
                                   findingType core.FindingType) *core.Finding {
    return &core.Finding{
        ID:          generateID(),
        CheckID:     checkID,
        Type:        findingType,
        Severity:    severity,
        Title:       title,
        Description: description,
        Confidence:  1.0,
        Metadata:    make(map[string]interface{}),
    }
}
```

### Language-Specific Analyzer Examples

```go
// pkg/analyzers/go/testing.go
package golang

type GoTestingAnalyzer struct {
    base.BaseAnalyzer
}

func NewGoTestingAnalyzer() *GoTestingAnalyzer {
    return &GoTestingAnalyzer{
        BaseAnalyzer: base.BaseAnalyzer{
            name:  "go-testing",
            layer: core.LayerStructure,
        },
    }
}

func (gta *GoTestingAnalyzer) AppliesTo(ctx *core.RepositoryContext) bool {
    for _, lang := range ctx.Languages {
        if lang.Name == "Go" && lang.IsPrimary {
            return true
        }
    }
    return false
}

func (gta *GoTestingAnalyzer) Analyze(ctx context.Context,
                                       repoCtx *core.RepositoryContext,
                                       previousResults *core.AnalysisResults) (*core.AnalyzerResult, error) {

    result := &core.AnalyzerResult{
        AnalyzerName: gta.Name(),
        Success:      true,
        Findings:     make([]*core.Finding, 0),
        Metrics:      make(map[string]interface{}),
    }

    startTime := time.Now()

    // Find all *_test.go files
    testFiles, err := gta.findTestFiles(repoCtx.RootPath)
    if err != nil {
        result.Success = false
        result.Error = err
        return result, err
    }

    // Analyze each test file
    tableTestCount := 0
    parallelTestCount := 0
    totalTests := 0

    for _, testFile := range testFiles {
        analysis, err := gta.analyzeTestFile(testFile)
        if err != nil {
            log.Warnf("failed to analyze %s: %v", testFile, err)
            continue
        }

        totalTests += analysis.TestCount
        tableTestCount += analysis.TableTests
        parallelTestCount += analysis.ParallelTests

        // Generate findings for issues
        result.Findings = append(result.Findings, analysis.Findings...)
    }

    // Calculate metrics
    tableTestPercentage := 0.0
    if totalTests > 0 {
        tableTestPercentage = (float64(tableTestCount) / float64(totalTests)) * 100
    }

    result.Metrics["total_tests"] = totalTests
    result.Metrics["table_tests"] = tableTestCount
    result.Metrics["table_test_percentage"] = tableTestPercentage
    result.Metrics["parallel_tests"] = parallelTestCount

    // Add finding if table-driven test usage is low
    if tableTestPercentage < 50.0 && totalTests > 10 {
        finding := gta.NewFinding(
            "go-table-tests",
            "Low table-driven test usage",
            fmt.Sprintf("Only %.1f%% of tests use table-driven pattern. "+
                       "Go best practice recommends table-driven tests for better coverage and maintainability.",
                       tableTestPercentage),
            core.SeverityMedium,
            core.FindingTypeBestPractice,
        )
        finding.Remediation = &core.Remediation{
            Description: "Convert tests to table-driven pattern",
            Effort:      core.EffortMedium,
            Examples: []string{
                "tests := []struct{name string; input int; want int}{{\"positive\", 1, 1}, {\"negative\", -1, 1}}",
                "for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { ... }) }",
            },
        }
        finding.References = []*core.Reference{
            {
                Title:   "Table Driven Tests - Go Wiki",
                URL:     "https://go.dev/wiki/TableDrivenTests",
                Summary: "Official Go documentation on table-driven testing pattern",
            },
        }
        result.Findings = append(result.Findings, finding)
    }

    result.Duration = time.Since(startTime)
    return result, nil
}

func (gta *GoTestingAnalyzer) analyzeTestFile(filePath string) (*testFileAnalysis, error) {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }

    analysis := &testFileAnalysis{
        FilePath: filePath,
        Findings: make([]*core.Finding, 0),
    }

    // Walk the AST
    ast.Inspect(node, func(n ast.Node) bool {
        switch fn := n.(type) {
        case *ast.FuncDecl:
            if gta.isTestFunction(fn) {
                analysis.TestCount++

                // Check for table-driven pattern
                if gta.isTableDrivenTest(fn) {
                    analysis.TableTests++
                }

                // Check for t.Parallel()
                if gta.callsParallel(fn) {
                    analysis.ParallelTests++
                }

                // Check for test smells
                findings := gta.detectTestSmells(fn, fset)
                analysis.Findings = append(analysis.Findings, findings...)
            }
        }
        return true
    })

    return analysis, nil
}

type testFileAnalysis struct {
    FilePath      string
    TestCount     int
    TableTests    int
    ParallelTests int
    Findings      []*core.Finding
}
```

```go
// pkg/analyzers/python/pytest.go
package python

type PytestAnalyzer struct {
    base.BaseAnalyzer
}

func NewPytestAnalyzer() *PytestAnalyzer {
    return &PytestAnalyzer{
        BaseAnalyzer: base.BaseAnalyzer{
            name:  "pytest",
            layer: core.LayerQuality,
        },
    }
}

func (pa *PytestAnalyzer) AppliesTo(ctx *core.RepositoryContext) bool {
    for _, lang := range ctx.Languages {
        if lang.Name == "Python" {
            for _, fw := range lang.TestFrameworks {
                if fw == "pytest" {
                    return true
                }
            }
        }
    }
    return false
}

func (pa *PytestAnalyzer) Analyze(ctx context.Context,
                                   repoCtx *core.RepositoryContext,
                                   previousResults *core.AnalysisResults) (*core.AnalyzerResult, error) {

    result := &core.AnalyzerResult{
        AnalyzerName: pa.Name(),
        Success:      true,
        Findings:     make([]*core.Finding, 0),
        Metrics:      make(map[string]interface{}),
    }

    // 1. Analyze pytest configuration
    configFindings := pa.analyzeConfiguration(repoCtx.RootPath)
    result.Findings = append(result.Findings, configFindings...)

    // 2. Analyze test files for pytest-specific patterns
    testFiles := pa.findPytestFiles(repoCtx)

    fixtureLeaks := 0
    parametrizedTests := 0
    totalTests := 0

    for _, testFile := range testFiles {
        analysis := pa.analyzeTestFile(testFile)
        fixtureLeaks += analysis.FixtureLeaks
        parametrizedTests += analysis.ParametrizedTests
        totalTests += analysis.TotalTests
        result.Findings = append(result.Findings, analysis.Findings...)
    }

    result.Metrics["fixture_leaks"] = fixtureLeaks
    result.Metrics["parametrized_tests"] = parametrizedTests
    result.Metrics["total_tests"] = totalTests
    result.Metrics["parametrization_percentage"] = (float64(parametrizedTests) / float64(totalTests)) * 100

    return result, nil
}

func (pa *PytestAnalyzer) analyzeConfiguration(repoPath string) []*core.Finding {
    findings := make([]*core.Finding, 0)

    // Check for pytest.ini or pyproject.toml
    pytestIni := filepath.Join(repoPath, "pytest.ini")
    pyprojectToml := filepath.Join(repoPath, "pyproject.toml")

    hasConfig := false
    configFile := ""

    if fileExists(pytestIni) {
        hasConfig = true
        configFile = "pytest.ini"
    } else if fileExists(pyprojectToml) {
        content, _ := os.ReadFile(pyprojectToml)
        if strings.Contains(string(content), "[tool.pytest.ini_options]") {
            hasConfig = true
            configFile = "pyproject.toml"
        }
    }

    if !hasConfig {
        finding := &core.Finding{
            ID:       generateID(),
            CheckID:  "pytest-no-config",
            Type:     core.FindingTypeBestPractice,
            Severity: core.SeverityMedium,
            Title:    "No pytest configuration file found",
            Description: "pytest configuration file (pytest.ini or pyproject.toml) is recommended " +
                        "for consistent test execution and coverage settings.",
            Remediation: &core.Remediation{
                Description: "Create pytest.ini or add [tool.pytest.ini_options] to pyproject.toml",
                Effort:      core.EffortMinimal,
                Examples: []string{
                    "[pytest]\naddopts = --strict-markers --cov=src --cov-report=html",
                    "testpaths = tests",
                },
            },
        }
        findings = append(findings, finding)
    } else {
        // Analyze configuration quality
        configFindings := pa.analyzeConfigQuality(filepath.Join(repoPath, configFile))
        findings = append(findings, configFindings...)
    }

    return findings
}
```

---

## Assessor Framework

### Coverage Assessor

```go
// pkg/assessors/coverage.go
package assessors

type CoverageAssessor struct {
    name      string
    dimension core.ScoreDimension
}

func NewCoverageAssessor() *CoverageAssessor {
    return &CoverageAssessor{
        name:      "coverage-assessor",
        dimension: core.DimensionCoverage,
    }
}

func (ca *CoverageAssessor) Name() string {
    return ca.name
}

func (ca *CoverageAssessor) Dimension() core.ScoreDimension {
    return ca.dimension
}

func (ca *CoverageAssessor) AppliesTo(ctx *core.RepositoryContext) bool {
    // Applies to all repositories with tests
    return len(ctx.Structure.TestFiles) > 0
}

func (ca *CoverageAssessor) Assess(ctx *core.RepositoryContext,
                                    results *core.AnalysisResults) (*core.DimensionScore, error) {

    score := &core.DimensionScore{
        Dimension:  ca.dimension,
        Weight:     0.30, // 30% of overall score
        Components: make(map[string]float64),
    }

    // Extract coverage metrics from results
    coverageMetrics := ca.extractCoverageMetrics(results)
    if coverageMetrics == nil {
        score.Score = 0
        score.Rationale = "No coverage data available"
        return score, nil
    }

    // Calculate component scores
    lineCovScore := ca.scoreCoverage(coverageMetrics.LinePercentage, []threshold{
        {min: 90, score: 100},
        {min: 80, score: 85},
        {min: 70, score: 70},
        {min: 60, score: 50},
        {min: 0, score: 25},
    })

    branchCovScore := ca.scoreCoverage(coverageMetrics.BranchPercentage, []threshold{
        {min: 85, score: 100},
        {min: 70, score: 85},
        {min: 60, score: 70},
        {min: 0, score: 40},
    })

    functionCovScore := ca.scoreCoverage(coverageMetrics.FunctionPercentage, []threshold{
        {min: 95, score: 100},
        {min: 85, score: 85},
        {min: 75, score: 70},
        {min: 0, score: 50},
    })

    // Critical path coverage (binary: 100 or 0)
    criticalPathScore := 100.0
    if coverageMetrics.UncoveredCriticalPaths > 0 {
        criticalPathScore = 0.0
    }

    // Weighted component scores
    score.Components["line_coverage"] = lineCovScore
    score.Components["branch_coverage"] = branchCovScore
    score.Components["function_coverage"] = functionCovScore
    score.Components["critical_path_coverage"] = criticalPathScore

    // Calculate overall coverage score
    score.Score = (lineCovScore * 0.40) +
                  (branchCovScore * 0.30) +
                  (functionCovScore * 0.15) +
                  (criticalPathScore * 0.15)

    score.Rationale = fmt.Sprintf(
        "Line: %.1f%%, Branch: %.1f%%, Function: %.1f%%, Critical paths uncovered: %d",
        coverageMetrics.LinePercentage,
        coverageMetrics.BranchPercentage,
        coverageMetrics.FunctionPercentage,
        coverageMetrics.UncoveredCriticalPaths,
    )

    return score, nil
}

type threshold struct {
    min   float64
    score float64
}

func (ca *CoverageAssessor) scoreCoverage(percentage float64, thresholds []threshold) float64 {
    for _, t := range thresholds {
        if percentage >= t.min {
            return t.score
        }
    }
    return 0
}

type coverageMetrics struct {
    LinePercentage          float64
    BranchPercentage        float64
    FunctionPercentage      float64
    UncoveredCriticalPaths  int
}

func (ca *CoverageAssessor) extractCoverageMetrics(results *core.AnalysisResults) *coverageMetrics {
    // Extract from coverage analyzer results
    coverageLayer, ok := results.LayerResults["coverage"]
    if !ok {
        return nil
    }

    // Find coverage analyzer result
    for _, result := range coverageLayer.AnalyzerResults {
        if strings.Contains(result.AnalyzerName, "coverage") {
            return &coverageMetrics{
                LinePercentage:         result.Metrics["line_coverage"].(float64),
                BranchPercentage:       result.Metrics["branch_coverage"].(float64),
                FunctionPercentage:     result.Metrics["function_coverage"].(float64),
                UncoveredCriticalPaths: result.Metrics["uncovered_critical"].(int),
            }
        }
    }

    return nil
}
```

### Tool Adoption Assessor

```go
// pkg/assessors/tools.go
package assessors

type ToolAdoptionAssessor struct {
    name      string
    dimension core.ScoreDimension
    toolDB    *tools.Registry
}

func NewToolAdoptionAssessor(toolDB *tools.Registry) *ToolAdoptionAssessor {
    return &ToolAdoptionAssessor{
        name:      "tool-adoption-assessor",
        dimension: core.DimensionTools,
        toolDB:    toolDB,
    }
}

func (taa *ToolAdoptionAssessor) Assess(ctx *core.RepositoryContext,
                                         results *core.AnalysisResults) (*core.DimensionScore, error) {

    score := &core.DimensionScore{
        Dimension:  taa.dimension,
        Weight:     0.10, // 10% of overall score
        Components: make(map[string]float64),
    }

    categoryScores := make(map[string]float64)

    // For each primary language, assess tool adoption
    for _, lang := range ctx.Languages {
        if !lang.IsPrimary {
            continue
        }

        // Get recommended tools for this language
        recommendedTools := taa.toolDB.GetRecommendedTools(lang.Name)

        // Categorize tools
        categories := map[string][]string{
            "testing":    {},
            "quality":    {},
            "security":   {},
            "automation": {},
        }

        for _, tool := range recommendedTools {
            categories[tool.Category] = append(categories[tool.Category], tool.Name)
        }

        // Score each category
        for category, toolList := range categories {
            detected := 0
            essential := 0
            configured := 0

            for _, toolName := range toolList {
                tool := taa.toolDB.GetTool(toolName)
                if tool.Priority == "Essential" {
                    essential++
                }

                // Check if tool is detected
                if taa.isToolDetected(ctx, lang, toolName) {
                    detected++

                    // Check if properly configured
                    if taa.isToolConfigured(ctx, lang, toolName) {
                        configured++
                    }
                }
            }

            // Calculate category score
            baseScore := (float64(detected) / float64(len(toolList))) * 100

            // Configuration bonus
            if detected > 0 {
                configBonus := (float64(configured) / float64(detected)) * 20
                baseScore += configBonus
            }

            // Essential tools weight
            if essential > 0 {
                essentialDetected := taa.countEssentialDetected(ctx, lang, toolList)
                essentialWeight := (float64(essentialDetected) / float64(essential)) * 30
                baseScore = (baseScore * 0.7) + essentialWeight
            }

            categoryScores[category] = math.Min(baseScore, 100)
        }
    }

    // Aggregate category scores
    score.Components = categoryScores
    score.Score = 0
    for _, catScore := range categoryScores {
        score.Score += catScore
    }
    score.Score /= float64(len(categoryScores))

    return score, nil
}
```

---

## Execution Engine

### Scheduler

```go
// pkg/engine/scheduler.go
package engine

type Scheduler struct {
    maxWorkers int
}

func (s *Scheduler) BuildExecutionPlan(analyzers []core.Analyzer) (*ExecutionPlan, error) {
    // Group analyzers by layer
    layerMap := make(map[core.AnalysisLayer][]core.Analyzer)
    for _, analyzer := range analyzers {
        layer := analyzer.Layer()
        layerMap[layer] = append(layerMap[layer], analyzer)
    }

    // Create layer plans in order
    layers := []core.AnalysisLayer{
        core.LayerStructure,
        core.LayerQuality,
        core.LayerCoverage,
        core.LayerExecution,
        core.LayerMutation,
        core.LayerCICD,
    }

    plan := &ExecutionPlan{
        Layers: make([]*LayerPlan, 0),
    }

    for _, layer := range layers {
        if analyzers, ok := layerMap[layer]; ok && len(analyzers) > 0 {
            layerPlan := &LayerPlan{
                Name:      layer.String(),
                Analyzers: analyzers,
                Dependencies: s.getLayerDependencies(layer),
            }
            plan.Layers = append(plan.Layers, layerPlan)
            plan.TotalAnalyzers += len(analyzers)
        }
    }

    return plan, nil
}

func (s *Scheduler) getLayerDependencies(layer core.AnalysisLayer) []string {
    deps := make([]string, 0)

    switch layer {
    case core.LayerStructure:
        // No dependencies
    case core.LayerQuality:
        deps = append(deps, core.LayerStructure.String())
    case core.LayerCoverage:
        deps = append(deps, core.LayerStructure.String())
    case core.LayerExecution:
        deps = append(deps, core.LayerStructure.String())
    case core.LayerMutation:
        deps = append(deps, core.LayerCoverage.String())
    case core.LayerCICD:
        // Can run independently
    }

    return deps
}
```

### Registry

```go
// pkg/engine/registry.go
package engine

type Registry struct {
    analyzers map[string]core.Analyzer
    assessors map[core.ScoreDimension][]core.Assessor
    mu        sync.RWMutex
}

func NewRegistry() *Registry {
    return &Registry{
        analyzers: make(map[string]core.Analyzer),
        assessors: make(map[core.ScoreDimension][]core.Assessor),
    }
}

func (r *Registry) RegisterAnalyzer(analyzer core.Analyzer) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.analyzers[analyzer.Name()] = analyzer
}

func (r *Registry) RegisterAssessor(assessor core.Assessor) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.assessors[assessor.Dimension()] = append(r.assessors[assessor.Dimension()], assessor)
}

func (r *Registry) SelectAnalyzers(ctx *core.RepositoryContext) []core.Analyzer {
    r.mu.RLock()
    defer r.mu.RUnlock()

    selected := make([]core.Analyzer, 0)
    for _, analyzer := range r.analyzers {
        if analyzer.AppliesTo(ctx) {
            selected = append(selected, analyzer)
        }
    }
    return selected
}

func (r *Registry) SelectAssessors(ctx *core.RepositoryContext) []core.Assessor {
    r.mu.RLock()
    defer r.mu.RUnlock()

    selected := make([]core.Assessor, 0)
    for _, assessorList := range r.assessors {
        for _, assessor := range assessorList {
            if assessor.AppliesTo(ctx) {
                selected = append(selected, assessor)
            }
        }
    }
    return selected
}

// Initialize registers all built-in analyzers and assessors
func (r *Registry) Initialize() {
    // Register Go analyzers
    r.RegisterAnalyzer(golang.NewGoTestingAnalyzer())
    r.RegisterAnalyzer(golang.NewGoCoverageAnalyzer())
    r.RegisterAnalyzer(golang.NewGoQualityAnalyzer())

    // Register Python analyzers
    r.RegisterAnalyzer(python.NewPytestAnalyzer())
    r.RegisterAnalyzer(python.NewPythonCoverageAnalyzer())
    r.RegisterAnalyzer(python.NewPythonQualityAnalyzer())

    // Register JavaScript analyzers
    r.RegisterAnalyzer(javascript.NewJestAnalyzer())
    r.RegisterAnalyzer(javascript.NewVitestAnalyzer())
    r.RegisterAnalyzer(javascript.NewJSCoverageAnalyzer())

    // ... register other language analyzers

    // Register assessors
    r.RegisterAssessor(assessors.NewCoverageAssessor())
    r.RegisterAssessor(assessors.NewQualityAssessor())
    r.RegisterAssessor(assessors.NewPerformanceAssessor())
    r.RegisterAssessor(assessors.NewToolAdoptionAssessor(toolDB))
    r.RegisterAssessor(assessors.NewMaintainabilityAssessor())
    r.RegisterAssessor(assessors.NewCodeQualityAssessor())
}
```

---

## Report Generation System

### Report Generator

```go
// pkg/report/generator.go
package report

type Generator struct {
    htmlGen     *HTMLGenerator
    jsonGen     *JSONGenerator
    cliGen      *CLIGenerator
    mdGen       *MarkdownGenerator
}

func New() *Generator {
    return &Generator{
        htmlGen: NewHTMLGenerator(),
        jsonGen: NewJSONGenerator(),
        cliGen:  NewCLIGenerator(),
        mdGen:   NewMarkdownGenerator(),
    }
}

func (g *Generator) Generate(ctx *core.RepositoryContext,
                             results *core.AnalysisResults) (*core.Report, error) {

    report := &core.Report{
        RepositoryContext: ctx,
        Results:           results,
        GeneratedAt:       time.Now(),
    }

    // Calculate scores
    scores, err := g.calculateScores(ctx, results)
    if err != nil {
        return nil, err
    }
    report.Scores = scores

    // Aggregate findings
    report.Findings = g.aggregateFindings(results)

    // Tool analysis
    report.ToolAnalysis = g.analyzeTools(ctx, results)

    return report, nil
}

func (g *Generator) calculateScores(ctx *core.RepositoryContext,
                                     results *core.AnalysisResults) (*core.ScoreCard, error) {

    scoreCard := &core.ScoreCard{
        Dimensions: make(map[core.ScoreDimension]*core.DimensionScore),
    }

    // Run all applicable assessors
    registry := engine.GetRegistry() // singleton
    assessors := registry.SelectAssessors(ctx)

    overall := 0.0
    for _, assessor := range assessors {
        dimScore, err := assessor.Assess(ctx, results)
        if err != nil {
            log.Warnf("assessor %s failed: %v", assessor.Name(), err)
            continue
        }
        scoreCard.Dimensions[assessor.Dimension()] = dimScore
        overall += dimScore.Score * dimScore.Weight
    }

    scoreCard.Overall = overall
    scoreCard.Grade = g.calculateGrade(overall)

    return scoreCard, nil
}

func (g *Generator) calculateGrade(score float64) string {
    switch {
    case score >= 95:
        return "A+"
    case score >= 90:
        return "A"
    case score >= 80:
        return "B"
    case score >= 70:
        return "C"
    case score >= 60:
        return "D"
    default:
        return "F"
    }
}

func (g *Generator) ExportHTML(report *core.Report, outputPath string) error {
    return g.htmlGen.Generate(report, outputPath)
}

func (g *Generator) ExportJSON(report *core.Report, outputPath string) error {
    return g.jsonGen.Generate(report, outputPath)
}

func (g *Generator) PrintCLI(report *core.Report) error {
    return g.cliGen.Print(report)
}
```

### HTML Generator

```go
// pkg/report/html.go
package report

type HTMLGenerator struct {
    templates *template.Template
}

func NewHTMLGenerator() *HTMLGenerator {
    tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
    return &HTMLGenerator{
        templates: tmpl,
    }
}

func (hg *HTMLGenerator) Generate(report *core.Report, outputPath string) error {
    data := hg.prepareTemplateData(report)

    f, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer f.Close()

    return hg.templates.ExecuteTemplate(f, "report.html", data)
}

func (hg *HTMLGenerator) prepareTemplateData(report *core.Report) map[string]interface{} {
    return map[string]interface{}{
        "Title":       "Ship Shape Analysis Report",
        "GeneratedAt": report.GeneratedAt.Format(time.RFC3339),
        "Repository":  report.RepositoryContext,
        "Scores":      report.Scores,
        "Findings":    hg.groupFindingsBySeverity(report.Findings),
        "Tools":       report.ToolAnalysis,
        "Charts":      hg.generateChartData(report),
    }
}

func (hg *HTMLGenerator) groupFindingsBySeverity(findings []*core.Finding) map[core.Severity][]*core.Finding {
    grouped := make(map[core.Severity][]*core.Finding)
    for _, f := range findings {
        grouped[f.Severity] = append(grouped[f.Severity], f)
    }
    return grouped
}

func (hg *HTMLGenerator) generateChartData(report *core.Report) map[string]interface{} {
    charts := make(map[string]interface{})

    // Radar chart for dimensions
    charts["dimensions"] = map[string]interface{}{
        "labels": []string{
            "Coverage", "Quality", "Performance",
            "Tools", "Maintainability", "Code Quality",
        },
        "data": []float64{
            report.Scores.Dimensions[core.DimensionCoverage].Score,
            report.Scores.Dimensions[core.DimensionQuality].Score,
            report.Scores.Dimensions[core.DimensionPerformance].Score,
            report.Scores.Dimensions[core.DimensionTools].Score,
            report.Scores.Dimensions[core.DimensionMaintain].Score,
            report.Scores.Dimensions[core.DimensionCodeQuality].Score,
        },
    }

    // Language distribution pie chart
    langLabels := make([]string, 0)
    langData := make([]float64, 0)
    for _, lang := range report.RepositoryContext.Languages {
        langLabels = append(langLabels, lang.Name)
        langData = append(langData, lang.Percentage)
    }
    charts["languages"] = map[string]interface{}{
        "labels": langLabels,
        "data":   langData,
    }

    return charts
}
```

---

## Quality Gate System

```go
// pkg/gates/evaluator.go
package gates

type Evaluator struct {
    config *GateConfig
}

type GateConfig struct {
    Blocking []Gate `yaml:"blocking"`
    Warning  []Gate `yaml:"warning"`
    Trend    []TrendGate `yaml:"trend"`
}

type Gate struct {
    Metric    string  `yaml:"metric"`
    Threshold float64 `yaml:"threshold"`
    Operator  string  `yaml:"operator"` // >=, <=, ==, !=, >, <
    Message   string  `yaml:"message"`
}

type TrendGate struct {
    Metric   string  `yaml:"metric"`
    Delta    float64 `yaml:"delta"`
    Operator string  `yaml:"operator"`
    Message  string  `yaml:"message"`
}

func (e *Evaluator) Evaluate(report *core.Report) (*core.GateResults, error) {
    results := &core.GateResults{
        Passed:   true,
        Blocking: make([]*GateResult, 0),
        Warning:  make([]*GateResult, 0),
        Trend:    make([]*GateResult, 0),
    }

    // Evaluate blocking gates
    for _, gate := range e.config.Blocking {
        result := e.evaluateGate(gate, report)
        results.Blocking = append(results.Blocking, result)
        if !result.Passed {
            results.Passed = false
        }
    }

    // Evaluate warning gates
    for _, gate := range e.config.Warning {
        result := e.evaluateGate(gate, report)
        results.Warning = append(results.Warning, result)
    }

    // Evaluate trend gates (if historical data available)
    if e.hasHistoricalData(report) {
        for _, gate := range e.config.Trend {
            result := e.evaluateTrendGate(gate, report)
            results.Trend = append(results.Trend, result)
            if !result.Passed {
                results.Passed = false
            }
        }
    }

    return results, nil
}

func (e *Evaluator) evaluateGate(gate Gate, report *core.Report) *GateResult {
    value := e.extractMetric(gate.Metric, report)
    passed := e.compareValues(value, gate.Threshold, gate.Operator)

    return &GateResult{
        Gate:    gate.Metric,
        Passed:  passed,
        Value:   value,
        Threshold: gate.Threshold,
        Message: gate.Message,
    }
}

func (e *Evaluator) extractMetric(metricPath string, report *core.Report) float64 {
    // metricPath example: "test_coverage.line_coverage"
    parts := strings.Split(metricPath, ".")

    switch parts[0] {
    case "test_coverage":
        dim := report.Scores.Dimensions[core.DimensionCoverage]
        if len(parts) == 2 {
            return dim.Components[parts[1]]
        }
        return dim.Score
    case "test_quality":
        dim := report.Scores.Dimensions[core.DimensionQuality]
        if len(parts) == 2 {
            return dim.Components[parts[1]]
        }
        return dim.Score
    // ... other dimensions
    }

    return 0
}

func (e *Evaluator) compareValues(value, threshold float64, operator string) bool {
    switch operator {
    case ">=":
        return value >= threshold
    case "<=":
        return value <= threshold
    case "==":
        return value == threshold
    case "!=":
        return value != threshold
    case ">":
        return value > threshold
    case "<":
        return value < threshold
    default:
        return false
    }
}
```

---

(Continuing in next message due to length...)

## Data Models

### Storage Schema (SQLite)

```sql
-- Historical analysis results for trend tracking
CREATE TABLE analysis_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_path TEXT NOT NULL,
    repository_hash TEXT, -- git commit hash
    timestamp DATETIME NOT NULL,
    overall_score REAL,
    grade TEXT,
    duration_seconds INTEGER,
    analyzer_version TEXT
);

CREATE TABLE dimension_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    dimension TEXT NOT NULL,
    score REAL NOT NULL,
    weight REAL NOT NULL,
    rationale TEXT,
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE TABLE findings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    finding_id TEXT NOT NULL,
    check_id TEXT NOT NULL,
    type TEXT NOT NULL,
    severity TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    file_path TEXT,
    line_number INTEGER,
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE TABLE metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    metric_name TEXT NOT NULL,
    metric_value REAL NOT NULL,
    metric_category TEXT, -- coverage, quality, performance, etc.
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE TABLE monorepo_packages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    package_name TEXT NOT NULL,
    package_path TEXT NOT NULL,
    language TEXT,
    score REAL,
    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE INDEX idx_runs_repo ON analysis_runs(repository_path, timestamp);
CREATE INDEX idx_findings_run ON findings(run_id);
CREATE INDEX idx_metrics_run ON metrics(run_id);
CREATE INDEX idx_packages_run ON monorepo_packages(run_id);
```

### Tool Database Schema (YAML)

```yaml
# data/tools/python.yml
language: Python
tools:
  - name: pytest
    category: testing
    priority: Essential
    description: Python testing framework with fixtures and plugins
    installation:
      methods:
        - pip: pip install pytest
        - pipx: pipx install pytest
    configuration:
      files:
        - pytest.ini
        - pyproject.toml
      example: |
        [pytest]
        addopts = --strict-markers --cov=src
        testpaths = tests
    detection:
      dependencies:
        - pytest
        - pytest-cov
      config_files:
        - pytest.ini
        - pyproject.toml
      imports:
        - pytest
    benefits:
      - Powerful fixture system
      - Extensive plugin ecosystem
      - Parametrized testing
    effort: Low
    references:
      - url: https://docs.pytest.org/
        title: pytest documentation

  - name: coverage.py
    category: testing
    priority: Essential
    description: Code coverage measurement for Python
    installation:
      methods:
        - pip: pip install coverage[toml]
    configuration:
      files:
        - .coveragerc
        - pyproject.toml
      example: |
        [tool.coverage.run]
        source = ["src"]
        branch = true
    detection:
      dependencies:
        - coverage
      config_files:
        - .coveragerc
      commands:
        - coverage
    benefits:
      - Branch coverage tracking
      - HTML report generation
      - Integration with pytest
    effort: Minimal
    references:
      - url: https://coverage.readthedocs.io/
        title: Coverage.py documentation

  - name: black
    category: quality
    priority: Recommended
    description: Uncompromising Python code formatter
    # ... similar structure
```

---

## Plugin Architecture

### Plugin Interface

```go
// pkg/core/plugin.go
package core

// Plugin represents an external analyzer or assessor
type Plugin interface {
    // Metadata
    Name() string
    Version() string
    Author() string
    Description() string

    // Lifecycle
    Initialize(config map[string]interface{}) error
    Shutdown() error
}

// AnalyzerPlugin extends Plugin for analyzer plugins
type AnalyzerPlugin interface {
    Plugin
    Analyzer
}

// AssessorPlugin extends Plugin for assessor plugins
type AssessorPlugin interface {
    Plugin
    Assessor
}

// PluginLoader loads and manages plugins
type PluginLoader struct {
    pluginDir string
    plugins   map[string]Plugin
    mu        sync.RWMutex
}

func NewPluginLoader(pluginDir string) *PluginLoader {
    return &PluginLoader{
        pluginDir: pluginDir,
        plugins:   make(map[string]Plugin),
    }
}

func (pl *PluginLoader) LoadPlugins() error {
    // Scan plugin directory for .so files (compiled Go plugins)
    files, err := filepath.Glob(filepath.Join(pl.pluginDir, "*.so"))
    if err != nil {
        return err
    }

    for _, file := range files {
        if err := pl.loadPlugin(file); err != nil {
            log.Warnf("failed to load plugin %s: %v", file, err)
            continue
        }
    }

    return nil
}

func (pl *PluginLoader) loadPlugin(path string) error {
    // Open plugin
    p, err := plugin.Open(path)
    if err != nil {
        return err
    }

    // Look for NewPlugin symbol
    symPlugin, err := p.Lookup("NewPlugin")
    if err != nil {
        return err
    }

    // Type assert to plugin constructor
    newPlugin, ok := symPlugin.(func() Plugin)
    if !ok {
        return fmt.Errorf("NewPlugin has incorrect signature")
    }

    // Create plugin instance
    pluginInstance := newPlugin()

    // Initialize plugin
    if err := pluginInstance.Initialize(nil); err != nil {
        return err
    }

    pl.mu.Lock()
    defer pl.mu.Unlock()
    pl.plugins[pluginInstance.Name()] = pluginInstance

    // Register with appropriate registry
    switch p := pluginInstance.(type) {
    case AnalyzerPlugin:
        engine.GetRegistry().RegisterAnalyzer(p)
    case AssessorPlugin:
        engine.GetRegistry().RegisterAssessor(p)
    }

    return nil
}
```

### Example Plugin Implementation

```go
// example-plugin/main.go
package main

import "github.com/shipshape/pkg/core"

type CustomAnalyzer struct {
    core.BaseAnalyzer
}

func (ca *CustomAnalyzer) Name() string {
    return "custom-analyzer"
}

func (ca *CustomAnalyzer) Version() string {
    return "1.0.0"
}

func (ca *CustomAnalyzer) Author() string {
    return "Custom Team"
}

func (ca *CustomAnalyzer) Description() string {
    return "Custom analyzer for specific use case"
}

func (ca *CustomAnalyzer) Initialize(config map[string]interface{}) error {
    // Custom initialization
    return nil
}

func (ca *CustomAnalyzer) Shutdown() error {
    // Cleanup
    return nil
}

func (ca *CustomAnalyzer) AppliesTo(ctx *core.RepositoryContext) bool {
    // Custom logic
    return true
}

func (ca *CustomAnalyzer) Analyze(ctx context.Context, 
                                   repoCtx *core.RepositoryContext,
                                   previousResults *core.AnalysisResults) (*core.AnalyzerResult, error) {
    // Custom analysis logic
    return &core.AnalyzerResult{
        AnalyzerName: ca.Name(),
        Success:      true,
        Findings:     make([]*core.Finding, 0),
        Metrics:      make(map[string]interface{}),
    }, nil
}

// Plugin entry point
func NewPlugin() core.Plugin {
    return &CustomAnalyzer{}
}
```

---

## Multi-Language Support

### Language Abstraction Layer

```go
// pkg/languages/interface.go
package languages

// LanguageSupport defines language-specific operations
type LanguageSupport interface {
    // Detection
    Name() string
    FileExtensions() []string
    DetectVersion(repoPath string) (string, error)

    // Test framework detection
    DetectTestFrameworks(repoPath string) ([]string, error)
    DetectCoverageTools(repoPath string) ([]string, error)

    // AST parsing (optional, language-specific)
    ParseTestFile(filePath string) (*TestFileAST, error)

    // Metrics
    GetRecommendedCoverageThreshold() float64
    GetTestFilePatterns() []string
}

type TestFileAST struct {
    Functions []*TestFunction
    Classes   []*TestClass
    Imports   []string
}

type TestFunction struct {
    Name       string
    LineNumber int
    Type       TestType
    Assertions int
}

type TestClass struct {
    Name      string
    Methods   []*TestFunction
    SetUp     bool
    TearDown  bool
}

// Registry of language support
type LanguageRegistry struct {
    languages map[string]LanguageSupport
}

func (lr *LanguageRegistry) Register(lang LanguageSupport) {
    lr.languages[lang.Name()] = lang
}

func (lr *LanguageRegistry) Get(name string) (LanguageSupport, bool) {
    lang, ok := lr.languages[name]
    return lang, ok
}
```

### Example Language Support Implementation

```go
// pkg/languages/python/support.go
package python

type PythonSupport struct{}

func (ps *PythonSupport) Name() string {
    return "Python"
}

func (ps *PythonSupport) FileExtensions() []string {
    return []string{".py", ".pyw"}
}

func (ps *PythonSupport) DetectVersion(repoPath string) (string, error) {
    // Check for .python-version, runtime.txt, pyproject.toml
    pythonVersion := filepath.Join(repoPath, ".python-version")
    if fileExists(pythonVersion) {
        content, err := os.ReadFile(pythonVersion)
        if err == nil {
            return strings.TrimSpace(string(content)), nil
        }
    }

    // Try to detect from pyproject.toml
    pyproject := filepath.Join(repoPath, "pyproject.toml")
    if fileExists(pyproject) {
        // Parse TOML and extract python version requirement
        version := extractPythonVersionFromPyproject(pyproject)
        if version != "" {
            return version, nil
        }
    }

    return "", nil
}

func (ps *PythonSupport) DetectTestFrameworks(repoPath string) ([]string, error) {
    frameworks := make([]string, 0)

    // Check for pytest
    if fileExists(filepath.Join(repoPath, "pytest.ini")) {
        frameworks = append(frameworks, "pytest")
    }

    // Check dependencies
    deps := parsePythonDependencies(repoPath)
    if contains(deps, "pytest") {
        frameworks = appendUnique(frameworks, "pytest")
    }
    if contains(deps, "unittest") || containsPattern(deps, "test") {
        frameworks = appendUnique(frameworks, "unittest")
    }

    return frameworks, nil
}

func (ps *PythonSupport) GetRecommendedCoverageThreshold() float64 {
    return 80.0 // Python projects typically target 80% coverage
}

func (ps *PythonSupport) GetTestFilePatterns() []string {
    return []string{
        "test_*.py",
        "*_test.py",
        "tests/**/*.py",
    }
}

func (ps *PythonSupport) ParseTestFile(filePath string) (*languages.TestFileAST, error) {
    // Use Python AST parser (call Python script or use Go-based parser)
    return parsePythonTestFile(filePath)
}
```

---

## Monorepo Handling

### Monorepo Coordinator

```go
// pkg/monorepo/coordinator.go
package monorepo

type Coordinator struct {
    registry *engine.Registry
}

func NewCoordinator(registry *engine.Registry) *Coordinator {
    return &Coordinator{
        registry: registry,
    }
}

// AnalyzeMonorepo orchestrates analysis across all packages
func (c *Coordinator) AnalyzeMonorepo(ctx context.Context,
                                       repoCtx *core.RepositoryContext) (*MonorepoReport, error) {

    if repoCtx.Monorepo == nil {
        return nil, fmt.Errorf("not a monorepo")
    }

    report := &MonorepoReport{
        MonorepoType:   repoCtx.Monorepo.Type,
        PackageReports: make(map[string]*core.Report),
        SharedInfra:    c.analyzeSharedInfrastructure(repoCtx),
    }

    // Analyze each package in parallel
    results := make(chan *packageResult, len(repoCtx.Monorepo.Packages))
    sem := make(chan struct{}, 4) // max 4 concurrent package analyses

    for _, pkg := range repoCtx.Monorepo.Packages {
        sem <- struct{}{}
        go func(p *core.PackageInfo) {
            defer func() { <-sem }()

            pkgReport, err := c.analyzePackage(ctx, repoCtx, p)
            results <- &packageResult{
                packageName: p.Name,
                report:      pkgReport,
                err:         err,
            }
        }(pkg)
    }

    // Collect results
    for i := 0; i < len(repoCtx.Monorepo.Packages); i++ {
        result := <-results
        if result.err != nil {
            log.Warnf("package %s analysis failed: %v", result.packageName, result.err)
            continue
        }
        report.PackageReports[result.packageName] = result.report
    }

    // Calculate aggregate scores
    report.AggregateScore = c.calculateAggregateScore(report.PackageReports)

    // Analyze cross-package consistency
    report.Consistency = c.analyzeConsistency(report.PackageReports)

    return report, nil
}

func (c *Coordinator) analyzePackage(ctx context.Context,
                                      repoCtx *core.RepositoryContext,
                                      pkg *core.PackageInfo) (*core.Report, error) {

    // Create package-specific context
    pkgCtx := &core.RepositoryContext{
        RootPath:  pkg.Path,
        Type:      core.TypeSingleLanguage,
        Languages: []*core.LanguageInfo{
            {
                Name:           pkg.Language,
                Percentage:     100,
                IsPrimary:      true,
                TestFrameworks: []string{pkg.TestFramework},
            },
        },
    }

    // Run discovery for package
    discoveryEng := discovery.New()
    pkgStructure, err := discoveryEng.DiscoverStructure(ctx, pkg.Path)
    if err != nil {
        return nil, err
    }
    pkgCtx.Structure = pkgStructure

    // Run analysis
    executor := engine.New(4)
    results, err := executor.Execute(ctx, pkgCtx)
    if err != nil {
        return nil, err
    }

    // Generate report
    reportGen := report.New()
    pkgReport, err := reportGen.Generate(pkgCtx, results)
    if err != nil {
        return nil, err
    }

    return pkgReport, nil
}

func (c *Coordinator) calculateAggregateScore(packageReports map[string]*core.Report) *core.ScoreCard {
    if len(packageReports) == 0 {
        return &core.ScoreCard{Overall: 0, Grade: "F"}
    }

    // Weight by package size (LOC)
    totalLOC := 0
    weightedScore := 0.0

    for _, report := range packageReports {
        pkgLOC := report.RepositoryContext.Structure.TotalLOC
        totalLOC += int(pkgLOC)
        weightedScore += report.Scores.Overall * float64(pkgLOC)
    }

    overallScore := weightedScore / float64(totalLOC)

    return &core.ScoreCard{
        Overall: overallScore,
        Grade:   calculateGrade(overallScore),
    }
}

func (c *Coordinator) analyzeConsistency(packageReports map[string]*core.Report) *ConsistencyAnalysis {
    analysis := &ConsistencyAnalysis{
        TestFrameworks:   make(map[string]int),
        CoverageTools:    make(map[string]int),
        QualityTools:     make(map[string]int),
        ScoreVariance:    0,
    }

    scores := make([]float64, 0)

    for _, report := range packageReports {
        scores = append(scores, report.Scores.Overall)

        // Track tool usage across packages
        for _, lang := range report.RepositoryContext.Languages {
            for _, fw := range lang.TestFrameworks {
                analysis.TestFrameworks[fw]++
            }
            for _, tool := range lang.CoverageTools {
                analysis.CoverageTools[tool]++
            }
        }
    }

    // Calculate score variance
    analysis.ScoreVariance = calculateVariance(scores)

    return analysis
}

type packageResult struct {
    packageName string
    report      *core.Report
    err         error
}

type MonorepoReport struct {
    MonorepoType    string
    PackageReports  map[string]*core.Report
    AggregateScore  *core.ScoreCard
    SharedInfra     *SharedInfraScore
    Consistency     *ConsistencyAnalysis
}

type SharedInfraScore struct {
    ConfigQuality float64
    CIIntegration float64
    SharedDeps    float64
    Overall       float64
}

type ConsistencyAnalysis struct {
    TestFrameworks  map[string]int // framework -> count of packages using it
    CoverageTools   map[string]int
    QualityTools    map[string]int
    ScoreVariance   float64
}
```

---

## Storage and Persistence

### Storage Interface

```go
// pkg/storage/interface.go
package storage

type Storage interface {
    // Store analysis result
    StoreReport(report *core.Report) error

    // Retrieve historical results
    GetRecentReports(repoPath string, limit int) ([]*core.Report, error)
    GetReportByID(id int64) (*core.Report, error)
    GetReportsByDateRange(repoPath string, start, end time.Time) ([]*core.Report, error)

    // Trend analysis
    GetMetricTrend(repoPath, metricName string, days int) ([]MetricPoint, error)

    // Cleanup
    DeleteOldReports(olderThan time.Time) error
    Close() error
}

type MetricPoint struct {
    Timestamp time.Time
    Value     float64
}

// SQLiteStorage implements Storage using SQLite
type SQLiteStorage struct {
    db *sql.DB
}

func NewSQLite(dbPath string) (*SQLiteStorage, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    storage := &SQLiteStorage{db: db}
    if err := storage.initialize(); err != nil {
        return nil, err
    }

    return storage, nil
}

func (s *SQLiteStorage) initialize() error {
    // Create tables if they don't exist
    _, err := s.db.Exec(schemaSQL) // from data models section
    return err
}

func (s *SQLiteStorage) StoreReport(report *core.Report) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert analysis run
    result, err := tx.Exec(`
        INSERT INTO analysis_runs (repository_path, repository_hash, timestamp, 
                                   overall_score, grade, duration_seconds, analyzer_version)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, report.RepositoryContext.RootPath,
        getGitCommitHash(report.RepositoryContext.RootPath),
        report.GeneratedAt,
        report.Scores.Overall,
        report.Scores.Grade,
        report.Results.EndTime.Sub(report.Results.StartTime).Seconds(),
        version.Version,
    )
    if err != nil {
        return err
    }

    runID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Insert dimension scores
    for _, dimScore := range report.Scores.Dimensions {
        _, err = tx.Exec(`
            INSERT INTO dimension_scores (run_id, dimension, score, weight, rationale)
            VALUES (?, ?, ?, ?, ?)
        `, runID, dimScore.Dimension, dimScore.Score, dimScore.Weight, dimScore.Rationale)
        if err != nil {
            return err
        }
    }

    // Insert findings (top severity only to save space)
    for _, finding := range report.Findings {
        if finding.Severity == core.SeverityCritical || finding.Severity == core.SeverityHigh {
            _, err = tx.Exec(`
                INSERT INTO findings (run_id, finding_id, check_id, type, severity, 
                                     title, description, file_path, line_number)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            `, runID, finding.ID, finding.CheckID, finding.Type, finding.Severity,
                finding.Title, finding.Description,
                finding.Location.FilePath, finding.Location.StartLine)
            if err != nil {
                return err
            }
        }
    }

    return tx.Commit()
}

func (s *SQLiteStorage) GetMetricTrend(repoPath, metricName string, days int) ([]MetricPoint, error) {
    cutoff := time.Now().AddDate(0, 0, -days)

    rows, err := s.db.Query(`
        SELECT timestamp, metric_value
        FROM metrics m
        JOIN analysis_runs r ON m.run_id = r.id
        WHERE r.repository_path = ?
          AND m.metric_name = ?
          AND r.timestamp >= ?
        ORDER BY r.timestamp ASC
    `, repoPath, metricName, cutoff)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    points := make([]MetricPoint, 0)
    for rows.Next() {
        var point MetricPoint
        if err := rows.Scan(&point.Timestamp, &point.Value); err != nil {
            return nil, err
        }
        points = append(points, point)
    }

    return points, nil
}
```

---

## Integration Architecture

### CI/CD Integration Patterns

```go
// pkg/ci/integration.go
package ci

// CIDetector identifies CI/CD platform
type CIDetector struct{}

func (cd *CIDetector) DetectPlatform(repoPath string) (string, error) {
    detectors := map[string]func(string) bool{
        "github-actions": cd.isGitHubActions,
        "gitlab-ci":      cd.isGitLabCI,
        "jenkins":        cd.isJenkins,
        "circleci":       cd.isCircleCI,
        "travis-ci":      cd.isTravisCI,
    }

    for platform, detector := range detectors {
        if detector(repoPath) {
            return platform, nil
        }
    }

    return "", nil
}

func (cd *CIDetector) isGitHubActions(repoPath string) bool {
    workflowsDir := filepath.Join(repoPath, ".github", "workflows")
    if !dirExists(workflowsDir) {
        return false
    }

    files, _ := filepath.Glob(filepath.Join(workflowsDir, "*.yml"))
    return len(files) > 0
}

// CIAnalyzer analyzes CI/CD configuration
type CIAnalyzer struct {
    platform string
}

func (ca *CIAnalyzer) Analyze(repoPath string) (*CIAnalysis, error) {
    switch ca.platform {
    case "github-actions":
        return ca.analyzeGitHubActions(repoPath)
    case "gitlab-ci":
        return ca.analyzeGitLabCI(repoPath)
    default:
        return &CIAnalysis{}, nil
    }
}

func (ca *CIAnalyzer) analyzeGitHubActions(repoPath string) (*CIAnalysis, error) {
    analysis := &CIAnalysis{
        Platform:    "github-actions",
        ConfigFiles: make([]string, 0),
        Findings:    make([]*core.Finding, 0),
    }

    workflowsDir := filepath.Join(repoPath, ".github", "workflows")
    files, err := filepath.Glob(filepath.Join(workflowsDir, "*.yml"))
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        analysis.ConfigFiles = append(analysis.ConfigFiles, file)

        // Parse workflow file
        workflow := parseGitHubWorkflow(file)

        // Check for test job
        hasTestJob := false
        for _, job := range workflow.Jobs {
            if strings.Contains(strings.ToLower(job.Name), "test") {
                hasTestJob = true
                analysis.TestJobFound = true

                // Analyze test job configuration
                jobFindings := ca.analyzeTestJob(job)
                analysis.Findings = append(analysis.Findings, jobFindings...)
            }
        }

        if !hasTestJob {
            analysis.Findings = append(analysis.Findings, &core.Finding{
                CheckID:     "ci-no-test-job",
                Type:        core.FindingTypeBestPractice,
                Severity:    core.SeverityCritical,
                Title:       "No test job found in CI workflow",
                Description: "GitHub Actions workflow should include a job that runs tests",
            })
        }
    }

    return analysis, nil
}

type CIAnalysis struct {
    Platform      string
    ConfigFiles   []string
    TestJobFound  bool
    Parallelized  bool
    CacheEnabled  bool
    Findings      []*core.Finding
}
```

---

## Security Architecture

### Security Considerations

```go
// pkg/security/scanner.go
package security

// SecureExecutor wraps command execution with safety checks
type SecureExecutor struct {
    allowedCommands map[string]bool
    sandboxed       bool
}

func NewSecureExecutor(sandboxed bool) *SecureExecutor {
    return &SecureExecutor{
        allowedCommands: map[string]bool{
            "pytest":   true,
            "go":       true,
            "npm":      true,
            "coverage": true,
            "jest":     true,
        },
        sandboxed: sandboxed,
    }
}

func (se *SecureExecutor) Execute(cmd string, args []string, dir string) (*ExecResult, error) {
    // Validate command is allowed
    if !se.allowedCommands[cmd] {
        return nil, fmt.Errorf("command not allowed: %s", cmd)
    }

    // Sanitize arguments
    for _, arg := range args {
        if !se.isArgSafe(arg) {
            return nil, fmt.Errorf("unsafe argument: %s", arg)
        }
    }

    // Execute in sandbox if enabled
    if se.sandboxed {
        return se.executeInSandbox(cmd, args, dir)
    }

    return se.executeDirect(cmd, args, dir)
}

func (se *SecureExecutor) isArgSafe(arg string) bool {
    // Check for command injection attempts
    dangerous := []string{";", "|", "&", "$", "`", "(", ")", ">", "<"}
    for _, d := range dangerous {
        if strings.Contains(arg, d) {
            return false
        }
    }
    return true
}

func (se *SecureExecutor) executeInSandbox(cmd string, args []string, dir string) (*ExecResult, error) {
    // Use Docker or similar containerization
    dockerArgs := []string{
        "run", "--rm",
        "-v", fmt.Sprintf("%s:/workspace", dir),
        "-w", "/workspace",
        "--network", "none", // no network access
        "--memory", "2g",    // memory limit
        "--cpus", "2",       // CPU limit
        "shipshape-sandbox",
        cmd,
    }
    dockerArgs = append(dockerArgs, args...)

    return exec.Command("docker", dockerArgs...), nil
}
```

---

## Deployment Architecture

### Build and Distribution

```makefile
# Makefile
.PHONY: build test install clean

# Build for current platform
build:
	go build -o bin/shipshape ./cmd/shipshape

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/shipshape-darwin-amd64 ./cmd/shipshape
	GOOS=darwin GOARCH=arm64 go build -o bin/shipshape-darwin-arm64 ./cmd/shipshape
	GOOS=linux GOARCH=amd64 go build -o bin/shipshape-linux-amd64 ./cmd/shipshape
	GOOS=linux GOARCH=arm64 go build -o bin/shipshape-linux-arm64 ./cmd/shipshape
	GOOS=windows GOARCH=amd64 go build -o bin/shipshape-windows-amd64.exe ./cmd/shipshape

# Run tests
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Install locally
install:
	go install ./cmd/shipshape

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.txt
```

### Docker Support

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shipshape ./cmd/shipshape

# Final image with language runtimes
FROM ubuntu:22.04

# Install language runtimes
RUN apt-get update && apt-get install -y \
    python3 python3-pip \
    nodejs npm \
    default-jdk \
    && rm -rf /var/lib/apt/lists/*

# Copy shipshape binary
COPY --from=builder /app/shipshape /usr/local/bin/

ENTRYPOINT ["shipshape"]
```

---

## Future API Design

### REST API (Planned Q4 2026)

```go
// pkg/api/server.go
package api

type Server struct {
    engine *core.Engine
    router *chi.Mux
    port   int
}

func NewServer(engine *core.Engine, port int) *Server {
    s := &Server{
        engine: engine,
        router: chi.NewRouter(),
        port:   port,
    }
    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    s.router.Use(middleware.Logger)
    s.router.Use(middleware.Recoverer)

    s.router.Post("/api/v1/analyze", s.handleAnalyze)
    s.router.Get("/api/v1/results/{id}", s.handleGetResult)
    s.router.Get("/api/v1/trends/{repo}", s.handleGetTrends)
    s.router.Post("/api/v1/gates/evaluate", s.handleEvaluateGates)
}

func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
    var req AnalyzeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Trigger analysis
    report, err := s.engine.Analyze(r.Context(), req.RepositoryPath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(AnalyzeResponse{
        ID:     report.ID,
        Status: "completed",
        Report: report,
    })
}

type AnalyzeRequest struct {
    RepositoryPath string                 `json:"repository_path"`
    Config         map[string]interface{} `json:"config"`
}

type AnalyzeResponse struct {
    ID     int64        `json:"id"`
    Status string       `json:"status"`
    Report *core.Report `json:"report"`
}
```

---

## Summary

This technical architecture document provides a comprehensive blueprint for Ship Shape implementation:

**Key Architectural Decisions:**
1. **Go-based core** for performance, type safety, and cross-platform support
2. **Layered analysis** with clear dependencies and parallel execution
3. **Plugin architecture** for extensibility
4. **Repository-aware design** with monorepo first-class support
5. **Multi-language abstraction** for consistent cross-language analysis
6. **Evidence-based scoring** with research citations
7. **Flexible storage** for historical trend analysis
8. **Security-first** execution with sandboxing support

**Implementation Roadmap:**
- **Phase 1 (Q1)**: Core engine, discovery, basic analyzers (Go, Python, JS)
- **Phase 2 (Q2)**: All language support, tool adoption, quality gates
- **Phase 3 (Q3)**: Mutation testing, CI/CD optimization, plugins
- **Phase 4 (Q4)**: REST API, web dashboard, organization features

**Next Steps:**
1. Detailed interface specifications
2. Test strategy and validation approach
3. Performance benchmarking plan
4. User stories and acceptance criteria

---

**Document Version**: 1.0.0  
**Last Updated**: 2026-01-27  
**Status**: Ready for Technical Review
