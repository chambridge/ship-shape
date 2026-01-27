# Spike to User Story Mapping

This document maps each spike to its associated user stories and provides guidance on how spike outcomes should feed back into user story details.

## SPIKE-001: Multi-Language AST Parsing Framework

### Associated User Stories
- **SS-010**: Go Test Analysis with AST Parsing (13 points)
- **SS-011**: Python Test Analysis (pytest/unittest) (13 points)
- **SS-012**: JavaScript/TypeScript Test Analysis (Jest/Vitest) (13 points)

### Spike Outcomes to Feed into Stories

#### Upon Completion:
1. **Update Technical Requirements sections with**:
   - Chosen AST parsing library (go/ast, tree-sitter, or hybrid)
   - Performance benchmarks achieved
   - Pattern detection accuracy metrics
   - Known limitations

2. **Update Acceptance Criteria with**:
   - Specific detection accuracy percentages
   - Performance targets validated by spike
   - Edge cases discovered during spike

3. **Update Test Requirements with**:
   - Test files from spike prototypes
   - Benchmark suites
   - Known good/bad pattern examples

4. **Add Implementation Notes**:
   ```markdown
   ## Spike Outcomes
   - **Parsing Approach**: [tree-sitter/go-ast/hybrid]
   - **Accuracy**: [%] on benchmark suite
   - **Performance**: [X] files in [Y] seconds
   - **Known Limitations**: [list]
   - **Recommended Patterns**: [reference to spike code]
   ```

#### Example Update for SS-010:
```markdown
## Spike Validation Results (SPIKE-001)

**Parsing Library**: go/ast (Go standard library)
**Performance**: 100 Go test files parsed in 1.8 seconds
**Accuracy**:
- Table-driven test detection: 98.5%
- t.Parallel() detection: 100%
- Subtest detection: 97.2%

**Implementation Recommendations**:
- Use ast.Inspect() for recursive node traversal
- Cache parsed ASTs for unchanged files
- See `internal/analyzer/go_analyzer.go` from spike prototype

**Known Limitations**:
- Cannot detect table-driven tests using map structures (only slices)
- Requires valid Go syntax (fails on broken code)
```

---

## SPIKE-002: Monorepo Analysis Coordination

### Associated User Stories
- **SS-003**: Monorepo Structure Detection (8 points)
- **SS-020**: Monorepo Package-Level Analysis (13 points)
- **SS-021**: Monorepo Aggregate Scoring (5 points)

### Spike Outcomes to Feed into Stories

#### Upon Completion:
1. **Update SS-003 Technical Requirements**:
   - Confirmed monorepo types detectable
   - Workspace pattern resolution approach
   - Performance metrics for detection

2. **Update SS-020 Technical Requirements**:
   - Optimal concurrency level (from benchmarks)
   - Memory usage patterns
   - Error isolation strategy validated

3. **Update SS-021 with Scoring Algorithm**:
   - Weighting formula validated
   - Outlier handling approach
   - Configuration options

#### Example Update for SS-020:
```markdown
## Spike Validation Results (SPIKE-002)

**Concurrency Configuration**: 4 concurrent package analyses (optimal)
**Performance**: 20 packages analyzed in 28 seconds
**Memory Usage**: 1.6GB peak for 20 packages

**Monorepo Types Validated**:
- npm workspaces: ✅ 100% detection accuracy
- pnpm workspaces: ✅ 100% detection accuracy
- Go workspaces: ✅ 100% detection accuracy
- Lerna: ✅ 98% detection accuracy
- Nx: ✅ 95% detection accuracy

**Implementation Recommendations**:
- Use semaphore with max 4 workers
- Implement package context isolation per `internal/coordinator/monorepo_coordinator.go`
- See spike for dependency graph algorithm

**Known Limitations**:
- Nested monorepos not fully supported
- Custom workspace configs require explicit configuration
```

---

## SPIKE-003: Test Smell Detection Framework

### Associated User Stories
- **SS-030**: Test Smell Detection (13 points)

### Spike Outcomes to Feed into Stories

#### Upon Completion:
1. **Update with Detection Accuracy Metrics**:
   - Precision/recall per smell type
   - False positive rates
   - Confidence thresholds

2. **Update with Pattern Libraries**:
   - Language-specific patterns validated
   - Reference to pattern definitions
   - Examples of detected smells

3. **Add Remediation Database Reference**:
   - Link to remediation guidance YAML
   - Code examples per language

#### Example Update for SS-030:
```markdown
## Spike Validation Results (SPIKE-003)

**Detectors Implemented**: 11 test smell types
**Overall Accuracy**: 92.3% precision, 89.7% recall
**False Positive Rate**: 7.8%

**Per-Smell Accuracy**:
| Smell Type | Precision | Recall |
|-----------|-----------|--------|
| Mystery Guest | 95% | 92% |
| Assertion Roulette | 98% | 97% |
| Eager Test | 88% | 85% |
| Resource Optimism | 91% | 88% |
| Conditional Logic | 100% | 95% |

**Implementation Recommendations**:
- Use confidence scoring (threshold 0.75 for reporting)
- Pattern library in `internal/smell/patterns/`
- Remediation database in `data/remediations.yaml`

**Validated on Projects**:
- kubernetes (Go): 234 tests analyzed, 12 smells detected
- django (Python): 3,421 tests analyzed, 87 smells detected
- react (JavaScript): 1,892 tests analyzed, 43 smells detected
```

---

## SPIKE-004: Multi-Format Coverage Report Parsing

### Associated User Stories
- **SS-040**: Coverage Report Parsing (13 points)
- **SS-041**: Coverage Quality Assessment (5 points)

### Spike Outcomes to Feed into Stories

#### Upon Completion:
1. **Update SS-040 with Supported Formats**:
   - Parser accuracy per format
   - Performance benchmarks
   - Format auto-detection rules

2. **Update SS-041 with Assessment Algorithm**:
   - Critical path detection approach
   - Quality grading thresholds validated
   - GitHub API integration details

#### Example Update for SS-040:
```markdown
## Spike Validation Results (SPIKE-004)

**Formats Supported**: 4 (Cobertura XML, LCOV, Coverage.py JSON, Go profile)
**Parser Accuracy**: >99.5% for all formats

**Performance Benchmarks**:
| Format | File Size | Parse Time | Accuracy |
|--------|-----------|------------|----------|
| Cobertura XML | 10MB | 0.8s | 99.7% |
| LCOV | 10MB | 1.2s | 99.8% |
| Coverage.py JSON | 10MB | 0.6s | 99.9% |
| Go profile | 5MB | 0.3s | 100% |

**Implementation Recommendations**:
- Use streaming parsers for files >50MB
- Format auto-detection via `detectFormat()` function
- Parsers in `internal/coverage/parsers/`

**GitHub API Integration**:
- Artifact download working (tested on real repos)
- Rate limit handling implemented
- Historical data caching recommended
```

---

## SPIKE-005: GitHub Actions Integration and CI/CD Analysis

### Associated User Stories
- **SS-080**: GitHub Actions Workflow Detection (5 points)
- **SS-081**: GitHub Actions Integration (Ship Shape Action) (13 points)
- **SS-082**: Pre-commit Hook Integration (5 points)

### Spike Outcomes to Feed into Stories

#### Upon Completion:
1. **Update SS-081 with Action Configuration**:
   - action.yml final structure
   - Docker image size and build time
   - GitHub API permission requirements

2. **Update SS-080 with Workflow Analysis**:
   - Optimization patterns detected
   - YAML parsing approach
   - Performance recommendations accuracy

3. **Add Example Workflows**:
   - Basic usage
   - Advanced configuration
   - Monorepo usage

#### Example Update for SS-081:
```markdown
## Spike Validation Results (SPIKE-005)

**Action Performance**: 8 seconds overhead (target: <10s)
**Docker Image Size**: 145MB (compressed)
**Build Time**: 3 minutes

**GitHub API Integration**:
- Checks API: ✅ Working with annotations
- PR Comments: ✅ Single comment update strategy validated
- Artifacts API: ✅ Report upload working

**Workflow Analysis Capabilities**:
- Caching detection: 95% accuracy
- Parallelization opportunities: 88% accuracy
- Estimated speedup calculations: ±20% accuracy

**Implementation Recommendations**:
- Use Docker container action (not composite)
- PR comment marker: `<!-- ship-shape-analysis -->`
- See `action.yml` and `Dockerfile` in spike prototype

**Tested Repositories**:
- Public repo: ✅ Full functionality
- Private repo: ✅ With proper token permissions
- Monorepo: ✅ 15 packages analyzed
```

---

## SPIKE-006: Interactive HTML Report Generation

### Associated User Stories
- **SS-090**: HTML Report Generation (13 points)

### Spike Outcomes to Feed into Stories

#### Upon Completion:
1. **Update with Chart Library Details**:
   - go-echarts integration approach
   - Chart types implemented
   - Interactive features validated

2. **Update with Performance Metrics**:
   - Report generation time
   - HTML file sizes
   - Browser load times

3. **Add Accessibility Results**:
   - WCAG compliance level achieved
   - Screen reader testing results
   - Keyboard navigation support

#### Example Update for SS-090:
```markdown
## Spike Validation Results (SPIKE-006)

**Chart Library**: go-echarts v2 (Apache ECharts wrapper)
**Template System**: html/template (Go stdlib)

**Performance Benchmarks**:
| Report Size | Findings | Gen Time | HTML Size | Load Time |
|------------|----------|----------|-----------|-----------|
| Small | 100 | 0.7s | 450KB | 0.8s |
| Medium | 1,000 | 2.1s | 1.8MB | 1.4s |
| Large | 10,000 | 8.3s | 9.2MB | 4.1s |

**Chart Types Implemented**:
- Radar Chart (score dimensions): ✅
- Pie Chart (language distribution): ✅
- Line Chart (coverage trends): ✅
- Heatmap (file coverage): ✅
- Bar Chart (test pyramid): ✅

**Accessibility**:
- WCAG 2.1 AA Compliance: ✅ 94/100
- Lighthouse Score: 96/100
- Keyboard Navigation: ✅ Fully supported
- Screen Reader: ✅ Tested with NVDA and VoiceOver

**Theme Support**:
- Light theme: ✅
- Dark theme: ✅
- Auto theme detection: ✅
- Theme persistence: ✅ (localStorage)

**Implementation Recommendations**:
- Use go-echarts for all charts
- Embed chart JavaScript inline (standalone HTML)
- See `internal/report/html/` for template structure
```

---

## Feeding Spike Outcomes Back - Process

### 1. Immediate Updates (During Spike)
As key decisions are made during the spike:
- Update the "Open Source Tools" section with chosen libraries
- Update "Technical Requirements" with validated approaches
- Add performance data to relevant sections

### 2. Spike Completion Updates
Upon spike completion and Go/No-Go decision:
- Add "Spike Validation Results" section to each user story
- Update story point estimates if spike reveals different complexity
- Add "Known Limitations" discovered during spike
- Link to spike prototype code/examples

### 3. No-Go Scenario Updates
If spike results in No-Go decision:
- Document why original approach was rejected
- Update user story with alternative approach
- Adjust acceptance criteria based on new approach
- Update story points if complexity changed

### 4. Test Plan Updates
After each spike:
- Add spike test cases to story test requirements
- Include benchmark data as acceptance criteria
- Reference spike test fixtures

## Template for User Story Updates

```markdown
## 🔬 Spike Validation (SPIKE-XXX)

**Spike Status**: ✅ GO / ❌ NO-GO
**Completion Date**: YYYY-MM-DD
**Spike Document**: [Link to spike MD]

### Key Outcomes
- **Approach Validated**: [chosen approach]
- **Performance**: [benchmark results]
- **Accuracy**: [accuracy metrics]
- **Libraries/Tools**: [chosen tools with versions]

### Implementation Guidance
- [Key implementation recommendation 1]
- [Key implementation recommendation 2]
- [Reference to spike prototype code]

### Known Limitations
- [Limitation 1]
- [Limitation 2]

### Updated Acceptance Criteria
Based on spike findings, the following acceptance criteria have been refined:
- [Updated AC 1]
- [Updated AC 2]

### Changes from Original Story
- **Story Points**: [Original] → [Updated] (if changed)
- **Dependencies**: [Any new dependencies discovered]
- **Risks**: [Any new risks identified]
```

## Tracking Spike Outcomes

Create a tracking table in each user story:

```markdown
## Spike History

| Date | Spike | Outcome | Changes Made |
|------|-------|---------|--------------|
| 2026-02-15 | SPIKE-001 | GO | Updated parser choice, performance targets |
| - | - | - | - |
```

---

**Last Updated**: 2026-01-27
**Status**: Planning - Ready for spike execution
