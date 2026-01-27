# SPIKE-002: Monorepo Analysis Coordination and Parallel Processing

## Overview
This spike validates the technical approach for analyzing monorepo structures with multiple packages/modules, ensuring efficient parallel processing, proper isolation, and accurate aggregate scoring.

**Associated User Stories**: SS-003, SS-020, SS-021
**Risk Level**: HIGH
**Priority**: P0 (Critical)
**Target Completion**: Week 2-3 of implementation

## Problem Statement
Monorepos present unique challenges for analysis:
- **Multiple packages** with different languages, frameworks, and configurations
- **Shared infrastructure** (root configs, CI workflows, dependencies)
- **Dependency relationships** between packages
- **Performance** requirements (20 packages in <30 seconds)
- **Isolation** requirements (package failures shouldn't cascade)
- **Aggregate scoring** complexity (weighted by package size/importance)

**Real-World Examples**:
- npm workspaces with 50+ packages
- Go multi-module workspaces
- Lerna/Nx/Turborepo managed monorepos
- Maven/Gradle multi-module projects

## Spike Objectives
- [ ] Validate monorepo detection algorithms for 5+ types
- [ ] Prototype parallel package analysis with goroutines
- [ ] Design package context isolation strategy
- [ ] Implement aggregate scoring algorithm
- [ ] Test dependency graph construction
- [ ] Benchmark performance at scale (20+ packages)
- [ ] Validate error handling and partial failures

## Technical Investigation Areas

### 1. Monorepo Detection
**Supported Types**:
1. npm/yarn/pnpm workspaces
2. Go multi-module workspaces (go.work)
3. Lerna monorepos
4. Nx monorepos
5. Turborepo
6. Maven/Gradle multi-module
7. Python multi-package (heuristic)

**Validation Questions**:
- Can we reliably detect all workspace types?
- How do we handle nested workspaces?
- Can we parse workspace glob patterns correctly?
- What's the detection performance for large repos?

**Prototype Requirements**:
```go
type MonorepoDetector interface {
    // Detect monorepo type and structure
    Detect(repoPath string) (*MonorepoInfo, error)

    // Get supported monorepo types
    SupportedTypes() []string
}

type MonorepoInfo struct {
    Type              string  // "npm-workspaces", "go-workspace", etc.
    RootPath          string
    Packages          []*PackageInfo
    SharedInfra       *SharedInfrastructure
    DependencyGraph   *DependencyGraph
}

type PackageInfo struct {
    Name         string
    Path         string
    Language     string
    Dependencies []string  // Internal package dependencies
    Config       *PackageConfig
}
```

### 2. Parallel Analysis Coordination
**Approach**: Goroutines with semaphore for concurrency control

**Validation Questions**:
- How many concurrent package analyses are optimal?
- How do we handle package dependency order?
- What's the performance improvement vs sequential?
- How do we aggregate results efficiently?

**Prototype Requirements**:
```go
type MonorepoCoordinator struct {
    maxConcurrency int
    analyzer       *AnalysisEngine
}

func (mc *MonorepoCoordinator) AnalyzeMonorepo(
    ctx context.Context,
    monorepo *MonorepoInfo,
) (*MonorepoReport, error) {
    // Parallel analysis implementation
}

// Concurrency patterns to validate:
// 1. Worker pool with semaphore
// 2. Fan-out/fan-in pattern
// 3. Pipeline pattern with stages
// 4. Dependency-aware scheduling
```

**Performance Test Scenarios**:
| Packages | Concurrent | Target Time | Memory Limit |
|----------|-----------|-------------|--------------|
| 5        | 4         | <10s        | <500MB       |
| 10       | 4         | <15s        | <1GB         |
| 20       | 4         | <30s        | <2GB         |
| 50       | 8         | <60s        | <4GB         |

### 3. Package Context Isolation
**Challenge**: Each package may have different:
- Programming languages
- Test frameworks
- Configuration files
- Dependencies
- Directory structures

**Validation Questions**:
- Can we create isolated analysis contexts?
- How do we prevent cross-package contamination?
- Can we reuse analyzers across packages?
- How do we handle shared dependencies?

**Prototype Requirements**:
```go
type PackageContext struct {
    Package      *PackageInfo
    RootPath     string
    Languages    []*LanguageInfo
    Frameworks   []*FrameworkInfo
    Config       *AnalysisConfig
    SharedDeps   *SharedDependencies
}

func (mc *MonorepoCoordinator) createPackageContext(
    pkg *PackageInfo,
    monorepo *MonorepoInfo,
) *PackageContext {
    // Isolated context creation
}

func (mc *MonorepoCoordinator) analyzePackage(
    ctx context.Context,
    pkgCtx *PackageContext,
) (*PackageReport, error) {
    // Package-specific analysis
}
```

### 4. Aggregate Scoring Algorithm
**Requirements**:
- Weight packages by size (LOC)
- Weight packages by importance (optional config)
- Calculate dimension-level aggregates
- Identify outliers and trends

**Validation Questions**:
- Is weighted average the right approach?
- Should we handle outliers differently?
- How do we aggregate different dimension scores?
- Can users override weights?

**Prototype Requirements**:
```go
type AggregateScorer struct {
    weightingStrategy WeightingStrategy
}

type WeightingStrategy string
const (
    WeightByLOC        WeightingStrategy = "loc"
    WeightByImportance WeightingStrategy = "importance"
    WeightEqual        WeightingStrategy = "equal"
)

func (as *AggregateScorer) CalculateAggregate(
    packageReports []*PackageReport,
    strategy WeightingStrategy,
) *AggregateScore {
    // Weighted scoring implementation
}

type AggregateScore struct {
    Overall            float64
    ByDimension        map[string]float64
    PackageScores      map[string]float64
    SharedInfraScore   float64
    ConsistencyScore   float64
    Outliers           []*OutlierPackage
}
```

### 5. Dependency Graph Construction
**Purpose**:
- Understand package relationships
- Optimize analysis order
- Detect circular dependencies
- Identify shared code quality impact

**Validation Questions**:
- Can we parse all dependency declaration formats?
- How do we handle workspace protocol references?
- Can we detect circular dependencies?
- What's the performance for complex graphs?

**Prototype Requirements**:
```go
type DependencyGraph struct {
    nodes map[string]*PackageNode
    edges []*Dependency
}

type PackageNode struct {
    Package      *PackageInfo
    Dependencies []*PackageNode
    Dependents   []*PackageNode
}

func BuildDependencyGraph(monorepo *MonorepoInfo) (*DependencyGraph, error)
func (dg *DependencyGraph) DetectCycles() [][]*PackageNode
func (dg *DependencyGraph) TopologicalSort() ([]*PackageNode, error)
```

## Prototype Requirements

### Deliverable 1: Multi-Type Monorepo Detector
**Files**: `internal/discovery/monorepo_detector.go`
- Detect npm/yarn/pnpm workspaces
- Detect Go workspaces
- Detect Lerna/Nx/Turborepo
- Detect Maven/Gradle multi-module
- Unit tests with real monorepo examples

**Test Cases**:
```go
func TestDetectNPMWorkspaces(t *testing.T)
func TestDetectPNPMWorkspaces(t *testing.T)
func TestDetectGoWorkspace(t *testing.T)
func TestDetectLernaMonorepo(t *testing.T)
func TestDetectNxMonorepo(t *testing.T)
func TestDetectMavenMultiModule(t *testing.T)
```

### Deliverable 2: Parallel Coordinator Prototype
**Files**: `internal/coordinator/monorepo_coordinator.go`
- Implement worker pool with semaphore
- Package context isolation
- Error handling and partial failures
- Progress reporting
- Performance benchmarks

**Benchmark Tests**:
```go
func BenchmarkMonorepoAnalysis5Packages(b *testing.B)
func BenchmarkMonorepoAnalysis10Packages(b *testing.B)
func BenchmarkMonorepoAnalysis20Packages(b *testing.B)
```

### Deliverable 3: Aggregate Scoring Implementation
**Files**: `internal/scoring/aggregate_scorer.go`
- Weighted average implementation
- Dimension-level aggregation
- Outlier detection
- Consistency scoring

### Deliverable 4: Dependency Graph Builder
**Files**: `internal/graph/dependency_graph.go`
- Parse package dependencies
- Build directed graph
- Cycle detection
- Topological sort

### Deliverable 5: Integration Test Suite
**Files**: `internal/coordinator/monorepo_integration_test.go`
- End-to-end monorepo analysis
- Real-world monorepo examples
- Performance validation
- Error scenario testing

## Performance Benchmarks

### Concurrency Analysis
Test different concurrency levels to find optimal:

| Concurrency | 5 Packages | 10 Packages | 20 Packages |
|-------------|-----------|-------------|-------------|
| 1 (seq)     | 15s       | 30s         | 60s         |
| 2           | 8s        | 16s         | 32s         |
| 4           | 5s        | 10s         | 20s         |
| 8           | 5s        | 8s          | 15s         |
| 16          | 5s        | 8s          | 14s         |

**Expected Result**: Diminishing returns after 4-8 workers

### Memory Usage Analysis
Monitor memory consumption during parallel processing:

| Packages | Sequential | Parallel (4) | Parallel (8) |
|----------|-----------|-------------|-------------|
| 5        | 200MB     | 500MB       | 600MB       |
| 10       | 400MB     | 900MB       | 1.2GB       |
| 20       | 800MB     | 1.8GB       | 2.5GB       |

## Risk Mitigation

### Risk 1: Memory exhaustion with many packages
**Mitigation**:
- Implement streaming results (don't hold all in memory)
- Garbage collection between package analyses
- Configurable concurrency limits based on available memory
- Memory profiling and optimization

### Risk 2: Cascading failures
**Mitigation**:
- Isolated error handling per package
- Continue analysis even if packages fail
- Collect and report partial results
- Detailed error logging per package

### Risk 3: Dependency graph complexity
**Mitigation**:
- Optimize graph algorithms (O(n) where possible)
- Cache graph computations
- Fail gracefully on circular dependencies
- Limit graph depth/complexity

### Risk 4: Workspace pattern ambiguity
**Mitigation**:
- Comprehensive test suite with edge cases
- Clear precedence rules for multiple detection matches
- User override capability in config
- Detailed detection logging

## Go/No-Go Decision Criteria

### GO if:
- ✅ All 5+ monorepo types detected accurately
- ✅ Parallel processing achieves >2x speedup vs sequential
- ✅ Memory usage stays within acceptable bounds (<2GB for 20 packages)
- ✅ Error isolation works (1 package failure doesn't stop others)
- ✅ Aggregate scoring algorithm produces sensible results
- ✅ Dependency graph construction is performant and accurate

### NO-GO if:
- ❌ Cannot detect >3 monorepo types accurately
- ❌ Parallel processing slower or same as sequential
- ❌ Memory usage exceeds 4GB for 20 packages
- ❌ Error handling causes cascading failures
- ❌ Aggregate scoring produces nonsensical results

### Alternative Approach:
If parallel processing proves problematic:
- Sequential analysis with progress reporting
- Configurable concurrency (default to 1)
- Package-by-package mode for large monorepos
- Incremental analysis (analyze changed packages only)

## Spike Deliverables

1. **Monorepo Detection Library**
   - Multi-type detector implementation
   - Comprehensive test suite
   - Detection accuracy report

2. **Parallel Coordinator Implementation**
   - Worker pool with semaphore
   - Context isolation mechanism
   - Error handling framework
   - Performance benchmarks

3. **Aggregate Scoring Module**
   - Weighted scoring implementation
   - Dimension aggregation logic
   - Outlier detection algorithm
   - Configuration options

4. **Dependency Graph Library**
   - Graph construction from package manifests
   - Cycle detection algorithm
   - Topological sort implementation
   - Visualization helper (optional)

5. **Performance Analysis Report**
   - Concurrency vs performance curves
   - Memory usage analysis
   - Scaling characteristics
   - Optimization recommendations

6. **Integration Test Suite**
   - Real-world monorepo examples
   - Edge case coverage
   - Error scenario testing
   - Performance regression tests

## Integration Guidelines

Upon successful spike completion:

1. **Configuration**:
```yaml
monorepo:
  detection:
    auto: true
    types: [npm-workspaces, go-workspace, lerna, nx, turborepo]

  analysis:
    parallel: true
    max_concurrency: 4
    fail_fast: false

  scoring:
    weighting: loc  # loc, importance, equal
    include_shared_infra: true
```

2. **Usage Pattern**:
```go
// Detect monorepo
detector := monorepo.NewMultiDetector()
info, err := detector.Detect(repoPath)

if info != nil && info.Type != "" {
    // Monorepo detected
    coordinator := monorepo.NewCoordinator(4)
    report, err := coordinator.AnalyzeMonorepo(ctx, info)

    fmt.Printf("Aggregate Score: %.2f\n", report.AggregateScore.Overall)
    for pkg, score := range report.PackageScores {
        fmt.Printf("  %s: %.2f\n", pkg, score)
    }
}
```

3. **Error Handling**:
```go
report, err := coordinator.AnalyzeMonorepo(ctx, info)
if err != nil {
    // Check if partial results available
    if partialErr, ok := err.(*PartialAnalysisError); ok {
        fmt.Printf("Analyzed %d/%d packages\n",
            len(partialErr.SuccessfulPackages),
            len(info.Packages))
        // Use partial results
    }
}
```

## Success Metrics
- [ ] 5+ monorepo types detected with >95% accuracy
- [ ] Parallel processing achieves target performance
- [ ] Memory usage within acceptable limits
- [ ] Error isolation verified
- [ ] Aggregate scoring validated
- [ ] Dependency graph accurate and performant
- [ ] Integration tests pass
- [ ] Documentation complete

## Timeline
- **Week 1**: Monorepo detection implementation and testing
- **Week 2**: Parallel coordinator and context isolation
- **Week 3**: Aggregate scoring and dependency graph
- **Week 4**: Performance optimization and integration testing

## References
- npm workspaces: https://docs.npmjs.com/cli/v8/using-npm/workspaces
- pnpm workspaces: https://pnpm.io/workspaces
- Go workspaces: https://go.dev/ref/mod#workspaces
- Lerna: https://lerna.js.org/
- Nx: https://nx.dev/
- Turborepo: https://turbo.build/repo
