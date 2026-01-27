# SPIKE-004: Multi-Format Coverage Report Parsing

## Overview
This spike validates the technical approach for parsing code coverage reports from multiple tools and formats, normalizing data into a unified structure, and leveraging GitHub API for historical coverage tracking.

**Associated User Stories**: SS-040, SS-041
**Risk Level**: HIGH
**Priority**: P0 (Critical)
**Target Completion**: Week 4-5 of implementation

## Problem Statement
Ship Shape needs to parse coverage reports from various language-specific tools (Coverage.py, Istanbul/nyc, gocov, JaCoCo, etc.) in multiple formats (XML, JSON, LCOV, HTML) and provide unified coverage analysis with historical trend tracking.

**Key Challenges**:
1. Multiple coverage formats (Cobertura XML, LCOV, JSON, etc.)
2. Language-specific coverage tools with different schemas
3. Branch vs line vs function coverage normalization
4. Missing coverage files (graceful degradation)
5. Large coverage files (performance)
6. Historical data tracking and trends

## Existing Tools and Standards

### Coverage Report Formats (Industry Standard)
1. **Cobertura XML** - Multi-language standard
   - Supported by: GitLab CI, Jenkins, Azure DevOps
   - Languages: Java (JaCoCo), Python (coverage.py), JavaScript (istanbul), C# (coverlet)
   - Schema well-defined and consistent

2. **LCOV** - Text-based format
   - Originated from gcov (C/C++)
   - Supported by: JavaScript (c8, nyc), Rust (cargo-tarpaulin)
   - Line-based coverage information

3. **JSON formats** - Tool-specific
   - Coverage.py JSON (Python)
   - Istanbul JSON (JavaScript)
   - gocov JSON (Go)
   - Varies by tool

### Coverage Tools by Language

**Python**:
- **Coverage.py 7.13.2** (latest 2026) - Industry standard
  - Formats: XML (Cobertura), JSON, HTML, LCOV
  - Python 3.10-3.15 support
  - [Documentation](https://coverage.readthedocs.io/)

**JavaScript/TypeScript**:
- **nyc (Istanbul)** - Most popular
  - Formats: JSON, LCOV, HTML, Cobertura, text
  - [npm package](https://www.npmjs.com/package/nyc)
- **c8** - Modern V8 coverage
  - Native V8 coverage API
  - Formats: LCOV, JSON, HTML

**Go**:
- **Built-in coverage** (`go test -cover`)
  - Format: Native Go coverage profile
  - Tools: gocov, go tool cover

**Java**:
- **JaCoCo** - Standard for Java
  - Formats: XML, HTML, CSV
  - Maven/Gradle integration

**Rust**:
- **cargo-tarpaulin** - Coverage for Rust
  - Formats: LCOV, Cobertura XML, JSON

### GitHub API for Coverage Data
- **GitHub Actions Artifacts API**: Download coverage reports from CI runs
- **GitHub Checks API**: Coverage check runs and annotations
- **GitHub Commits API**: Historical commit data for trend analysis
- **GitHub Repository API**: Repository metadata and statistics

## Spike Objectives
- [ ] Parse Cobertura XML format (priority - most common)
- [ ] Parse LCOV format
- [ ] Parse Coverage.py JSON format
- [ ] Parse Go coverage profiles
- [ ] Design unified coverage data model
- [ ] Implement coverage quality assessment
- [ ] Integrate GitHub API for historical data
- [ ] Benchmark parsing performance (100MB files)
- [ ] Validate accuracy against native tools

## Technical Investigation Areas

### 1. Unified Coverage Data Model

```go
// Universal coverage representation
type CoverageReport struct {
    Source        string              // Tool that generated report
    Format        string              // cobertura, lcov, json, etc.
    GeneratedAt   time.Time
    Repository    *RepositoryInfo
    OverallMetrics *CoverageMetrics
    Packages      []*PackageCoverage
    Files         []*FileCoverage
}

type CoverageMetrics struct {
    LineRate      float64  // 0.0 - 1.0
    BranchRate    float64  // 0.0 - 1.0
    FunctionRate  float64  // 0.0 - 1.0 (if available)
    Complexity    float64  // Average complexity (if available)

    LinesCovered  int
    LinesValid    int
    BranchesCovered int
    BranchesValid   int
    FunctionsCovered int
    FunctionsValid   int
}

type FileCoverage struct {
    Path          string
    Package       string
    Metrics       *CoverageMetrics
    Lines         []*LineCoverage
    Branches      []*BranchCoverage
    Functions     []*FunctionCoverage
}

type LineCoverage struct {
    Number        int
    Hits          int
    IsCovered     bool
    IsBranch      bool
    BranchInfo    *BranchInfo
}

type BranchCoverage struct {
    LineNumber    int
    BranchNumber  int
    Hits          int
    IsCovered     bool
}

type FunctionCoverage struct {
    Name          string
    LineNumber    int
    Hits          int
    IsCovered     bool
}
```

### 2. Format-Specific Parsers

#### Cobertura XML Parser
**Libraries**: Use Go's `encoding/xml` standard library

```go
type CoberturaParser struct{}

func (cp *CoberturaParser) Parse(reader io.Reader) (*CoverageReport, error) {
    var cobertura struct {
        XMLName xml.Name `xml:"coverage"`
        LineRate float64 `xml:"line-rate,attr"`
        BranchRate float64 `xml:"branch-rate,attr"`
        Packages struct {
            Package []struct {
                Name string `xml:"name,attr"`
                Classes struct {
                    Class []struct {
                        Filename string `xml:"filename,attr"`
                        Lines struct {
                            Line []struct {
                                Number int `xml:"number,attr"`
                                Hits int `xml:"hits,attr"`
                                Branch bool `xml:"branch,attr"`
                                ConditionCoverage string `xml:"condition-coverage,attr"`
                            } `xml:"line"`
                        } `xml:"lines"`
                    } `xml:"class"`
                } `xml:"classes"`
            } `xml:"package"`
        } `xml:"packages"`
    }

    decoder := xml.NewDecoder(reader)
    err := decoder.Decode(&cobertura)
    if err != nil {
        return nil, fmt.Errorf("failed to parse Cobertura XML: %w", err)
    }

    return convertCoberturaToUnified(&cobertura), nil
}
```

**Existing Library Option**:
- [cobertura-go](https://github.com/t-yuki/gocover-cobertura) - Cobertura XML generator/parser

#### LCOV Parser
**Format**: Text-based, line-by-line parsing

```go
type LCOVParser struct{}

func (lp *LCOVParser) Parse(reader io.Reader) (*CoverageReport, error) {
    scanner := bufio.NewScanner(reader)
    report := &CoverageReport{Format: "lcov"}
    var currentFile *FileCoverage

    for scanner.Scan() {
        line := scanner.Text()

        switch {
        case strings.HasPrefix(line, "SF:"):
            // Source file
            filename := strings.TrimPrefix(line, "SF:")
            currentFile = &FileCoverage{Path: filename}
            report.Files = append(report.Files, currentFile)

        case strings.HasPrefix(line, "DA:"):
            // Line coverage: DA:line_number,hit_count
            parts := strings.Split(strings.TrimPrefix(line, "DA:"), ",")
            lineNum, _ := strconv.Atoi(parts[0])
            hits, _ := strconv.Atoi(parts[1])
            currentFile.Lines = append(currentFile.Lines, &LineCoverage{
                Number: lineNum,
                Hits: hits,
                IsCovered: hits > 0,
            })

        case strings.HasPrefix(line, "BRDA:"):
            // Branch coverage: BRDA:line,block,branch,taken
            parts := strings.Split(strings.TrimPrefix(line, "BRDA:"), ",")
            lineNum, _ := strconv.Atoi(parts[0])
            branchNum, _ := strconv.Atoi(parts[2])
            taken := parts[3] != "-" && parts[3] != "0"
            // Add branch coverage...

        case line == "end_of_record":
            // Calculate file metrics
            calculateFileMetrics(currentFile)
        }
    }

    calculateOverallMetrics(report)
    return report, nil
}
```

**Existing Library Option**:
- [go-lcov](https://github.com/wadey/gocovmerge) - LCOV utilities

#### Coverage.py JSON Parser

```go
type CoveragePyJSONParser struct{}

func (cpj *CoveragePyJSONParser) Parse(reader io.Reader) (*CoverageReport, error) {
    var coveragePyData struct {
        Meta struct {
            Version string `json:"version"`
        } `json:"meta"`
        Files map[string]struct {
            ExecutedLines []int `json:"executed_lines"`
            MissingLines  []int `json:"missing_lines"`
            Summary struct {
                CoveredLines    int     `json:"covered_lines"`
                NumStatements   int     `json:"num_statements"`
                PercentCovered  float64 `json:"percent_covered"`
                MissingLines    int     `json:"missing_lines"`
            } `json:"summary"`
        } `json:"files"`
        Totals struct {
            CoveredLines   int     `json:"covered_lines"`
            NumStatements  int     `json:"num_statements"`
            PercentCovered float64 `json:"percent_covered"`
        } `json:"totals"`
    }

    decoder := json.NewDecoder(reader)
    err := decoder.Decode(&coveragePyData)
    if err != nil {
        return nil, err
    }

    return convertCoveragePyToUnified(&coveragePyData), nil
}
```

#### Go Coverage Profile Parser

```go
type GoCoverageParser struct{}

func (gcp *GoCoverageParser) Parse(reader io.Reader) (*CoverageReport, error) {
    profiles, err := cover.ParseProfiles(reader)
    if err != nil {
        return nil, err
    }

    report := &CoverageReport{
        Format: "go-coverage-profile",
        Files:  make([]*FileCoverage, 0),
    }

    for _, profile := range profiles {
        fileCov := &FileCoverage{
            Path:  profile.FileName,
            Lines: make([]*LineCoverage, 0),
        }

        for _, block := range profile.Blocks {
            // Each block represents covered lines
            for lineNum := block.StartLine; lineNum <= block.EndLine; lineNum++ {
                fileCov.Lines = append(fileCov.Lines, &LineCoverage{
                    Number:    lineNum,
                    Hits:      block.Count,
                    IsCovered: block.Count > 0,
                })
            }
        }

        report.Files = append(report.Files, fileCov)
    }

    calculateOverallMetrics(report)
    return report, nil
}
```

**Existing Library**: Go standard library `golang.org/x/tools/cover`

### 3. Parser Registry and Auto-Detection

```go
type CoverageParser interface {
    Parse(reader io.Reader) (*CoverageReport, error)
    CanParse(filename string, content []byte) bool
    SupportedFormats() []string
}

type ParserRegistry struct {
    parsers []CoverageParser
}

func (pr *ParserRegistry) ParseFile(filePath string) (*CoverageReport, error) {
    // Read file
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    // Try each parser
    for _, parser := range pr.parsers {
        if parser.CanParse(filePath, content) {
            reader := bytes.NewReader(content)
            return parser.Parse(reader)
        }
    }

    return nil, fmt.Errorf("no parser found for %s", filePath)
}

// Auto-detect format
func detectFormat(filename string, content []byte) string {
    switch {
    case strings.Contains(string(content[:100]), "<coverage"):
        return "cobertura"
    case strings.HasPrefix(string(content), "TN:"):
        return "lcov"
    case strings.Contains(string(content[:50]), `"meta"`):
        return "coverage-py-json"
    case strings.HasPrefix(string(content), "mode:"):
        return "go-coverage-profile"
    default:
        return "unknown"
    }
}
```

### 4. GitHub API Integration for Historical Coverage

```go
type GitHubCoverageTracker struct {
    client *github.Client
    owner  string
    repo   string
}

// Download coverage artifacts from GitHub Actions
func (gct *GitHubCoverageTracker) GetHistoricalCoverage(
    ctx context.Context,
    days int,
) ([]*HistoricalCoverage, error) {
    // Get recent workflow runs
    runs, _, err := gct.client.Actions.ListWorkflowRunsByFileName(
        ctx, gct.owner, gct.repo, "test.yml",
        &github.ListWorkflowRunsOptions{
            Status: "success",
            ListOptions: github.ListOptions{
                PerPage: 100,
            },
        },
    )
    if err != nil {
        return nil, err
    }

    historical := make([]*HistoricalCoverage, 0)

    for _, run := range runs.WorkflowRuns {
        // Get artifacts from this run
        artifacts, _, err := gct.client.Actions.ListWorkflowRunArtifacts(
            ctx, gct.owner, gct.repo, run.GetID(), nil,
        )
        if err != nil {
            continue
        }

        // Find coverage artifact
        for _, artifact := range artifacts.Artifacts {
            if strings.Contains(artifact.GetName(), "coverage") {
                // Download and parse
                coverage, err := gct.downloadAndParseCoverage(ctx, artifact)
                if err == nil {
                    historical = append(historical, &HistoricalCoverage{
                        Date:     run.GetCreatedAt().Time,
                        Commit:   run.GetHeadSHA(),
                        Coverage: coverage,
                    })
                }
            }
        }
    }

    return historical, nil
}

type HistoricalCoverage struct {
    Date     time.Time
    Commit   string
    Coverage *CoverageReport
}
```

**GitHub API Libraries**:
- [google/go-github](https://github.com/google/go-github) - Official Go client for GitHub API v3

### 5. Coverage Quality Assessment

```go
type CoverageAssessor struct {
    thresholds *CoverageThresholds
}

type CoverageThresholds struct {
    LineExcellent   float64  // 90%
    LineGood        float64  // 80%
    LineAcceptable  float64  // 70%

    BranchExcellent float64  // 85%
    BranchGood      float64  // 70%
    BranchAcceptable float64 // 60%
}

func (ca *CoverageAssessor) Assess(report *CoverageReport) *CoverageAssessment {
    assessment := &CoverageAssessment{
        OverallGrade: calculateGrade(report.OverallMetrics),
        UncoveredCriticalPaths: findUncoveredCriticalPaths(report),
        Gaps: findCoverageGaps(report),
        Recommendations: generateRecommendations(report),
    }

    return assessment
}

type CoverageAssessment struct {
    OverallGrade           string
    LineGrade              string
    BranchGrade            string
    UncoveredCriticalPaths []*UncoveredPath
    Gaps                   []*CoverageGap
    Recommendations        []*Recommendation
}

func findUncoveredCriticalPaths(report *CoverageReport) []*UncoveredPath {
    paths := []*UncoveredPath{}

    for _, file := range report.Files {
        // Identify critical files (main logic, security, data integrity)
        if isCriticalFile(file.Path) {
            for _, line := range file.Lines {
                if !line.IsCovered {
                    paths = append(paths, &UncoveredPath{
                        File: file.Path,
                        Line: line.Number,
                        Reason: determineCriticality(file.Path),
                    })
                }
            }
        }
    }

    return paths
}
```

## Prototype Requirements

### Deliverable 1: Core Parser Framework
**Files**: `internal/coverage/parser.go`, `internal/coverage/model.go`
- Unified coverage data model
- Parser interface
- Parser registry
- Format auto-detection

### Deliverable 2: Format Parsers
**Files**: `internal/coverage/parsers/*.go`
- Cobertura XML parser
- LCOV parser
- Coverage.py JSON parser
- Go coverage profile parser

### Deliverable 3: GitHub Integration
**Files**: `internal/coverage/github.go`
- GitHub API client setup
- Artifact download
- Historical coverage tracking
- Trend analysis

### Deliverable 4: Coverage Quality Assessor
**Files**: `internal/coverage/assessor.go`
- Quality grading algorithm
- Critical path detection
- Gap analysis
- Recommendation engine

### Deliverable 5: Test Suite
**Files**: `internal/coverage/parsers/*_test.go`
- Sample coverage files for each format
- Parser accuracy tests
- Performance benchmarks
- Real-world file testing

## Performance Benchmarks

| Format     | File Size | Target Parse Time | Memory Usage |
|-----------|-----------|-------------------|--------------|
| Cobertura | 10MB      | <1s               | <50MB        |
| LCOV      | 10MB      | <2s               | <50MB        |
| JSON      | 10MB      | <1s               | <100MB       |
| Go Profile| 5MB       | <500ms            | <30MB        |

## Validation Strategy

### 1. Parser Accuracy
Compare parsed metrics with native tool output:
```bash
# Python
coverage.py report  # Get official metrics
# Compare with Ship Shape parsed metrics

# JavaScript
nyc report --reporter=text-summary
# Compare with Ship Shape

# Go
go tool cover -func=coverage.out
# Compare with Ship Shape
```

Target: <1% difference in overall metrics

### 2. Format Compatibility
Test with real coverage files from popular projects:
- Python: django, flask, requests
- JavaScript: react, vue, express
- Go: kubernetes, docker, prometheus

## Risk Mitigation

### Risk 1: Large coverage files (>100MB)
**Mitigation**:
- Streaming parsers (don't load entire file)
- Incremental parsing
- Memory limits and monitoring
- File size warnings

### Risk 2: Malformed coverage files
**Mitigation**:
- Robust error handling
- Partial parsing (extract what's parsable)
- Clear error messages
- Validation against schemas

### Risk 3: GitHub API rate limiting
**Mitigation**:
- Respect rate limits (5000 requests/hour)
- Caching of historical data
- Incremental updates (not full refresh)
- Token-based authentication

## Go/No-Go Decision Criteria

### GO if:
- ✅ All 4 formats parsed with >99% accuracy
- ✅ Performance benchmarks met
- ✅ Unified model handles all format nuances
- ✅ GitHub integration functional
- ✅ Quality assessment produces useful insights

### NO-GO if:
- ❌ Parser accuracy <95%
- ❌ Performance >5x slower than targets
- ❌ Unified model too complex
- ❌ Cannot handle real-world edge cases

## Spike Deliverables

1. **Parser Framework**
   - Unified data model
   - Parser registry
   - Auto-detection

2. **Format Parsers** (4 total)
   - Cobertura, LCOV, Coverage.py JSON, Go profile
   - Accuracy validation reports
   - Performance benchmarks

3. **GitHub Integration**
   - Historical coverage tracking
   - Artifact download
   - Trend analysis

4. **Quality Assessor**
   - Grading algorithm
   - Critical path detection
   - Recommendations

## Success Metrics
- [ ] 4 coverage formats supported
- [ ] >99% parsing accuracy
- [ ] Performance targets met
- [ ] GitHub integration working
- [ ] Real-world validation passed
- [ ] Documentation complete

## Timeline
- **Week 1**: Core framework and Cobertura parser
- **Week 2**: LCOV and JSON parsers
- **Week 3**: Go parser and GitHub integration
- **Week 4**: Quality assessor and validation

## Sources and References
- [Coverage.py Documentation](https://coverage.readthedocs.io/)
- [Cobertura XML Format](https://cobertura.github.io/cobertura/)
- [LCOV Format Specification](http://ltp.sourceforge.net/coverage/lcov.php)
- [Istanbul/nyc](https://istanbul.js.org/)
- [GitHub Actions Artifacts API](https://docs.github.com/en/rest/actions/artifacts)
- [go-github Library](https://github.com/google/go-github)
- [Best Code Coverage Tools 2026](https://www.codeant.ai/blogs/best-code-test-coverage-tools-2025)
