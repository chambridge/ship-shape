# SPIKE-005: GitHub Actions Integration and CI/CD Analysis

## Overview
This spike validates the technical approach for integrating Ship Shape into GitHub Actions workflows, analyzing existing CI configurations, and providing actionable optimization recommendations.

**Associated User Stories**: SS-080, SS-081, SS-082
**Risk Level**: HIGH
**Priority**: P0 (Critical)
**Target Completion**: Week 5-6 of implementation

## Problem Statement
Ship Shape needs to:
1. **Run as a GitHub Action** in CI/CD pipelines with quality gates
2. **Analyze existing workflows** for test execution patterns
3. **Comment on PRs** with analysis results
4. **Provide optimization recommendations** for CI performance
5. **Integrate with GitHub Checks API** for status reporting

**Key Challenges**:
- GitHub Actions marketplace distribution
- Authentication and permissions
- PR commenting without spam
- Workflow YAML parsing and analysis
- Performance optimization detection
- Check runs integration

## Existing Tools and Standards

### GitHub Actions Ecosystem
- **actions/checkout** - Standard for checking out code
- **actions/upload-artifact** - Artifact upload
- **actions/download-artifact** - Artifact download
- **github-script** - Run JavaScript in workflows
- **peter-evans/create-or-update-comment** - PR commenting

### GitHub APIs
- **Checks API** - Create check runs with annotations
- **Pull Requests API** - Comment on PRs
- **Actions API** - Access workflow runs and artifacts
- **Repository API** - Access repo metadata

### CI/CD Analysis Tools
- **actionlint** - GitHub Actions workflow linter
- **super-linter** - Multi-language linter action
- **dependency-review-action** - Dependency scanning

## Spike Objectives
- [ ] Create GitHub Action (action.yml)
- [ ] Implement workflow YAML parsing
- [ ] Design CI optimization detector
- [ ] Prototype PR commenting
- [ ] Integrate GitHub Checks API
- [ ] Test in real GitHub repos
- [ ] Benchmark action performance
- [ ] Create marketplace listing

## Technical Investigation Areas

### 1. GitHub Action Implementation

**Action Types**:
1. **Docker Container Action** - Full control, larger image
2. **JavaScript Action** - Fast, limited to Node.js
3. **Composite Action** - Reusable workflows

**Recommendation**: Docker Container Action for full Go binary execution

**action.yml Structure**:
```yaml
name: 'Ship Shape Analysis'
description: 'Analyze test quality and coverage for your codebase'
author: 'Ship Shape Team'

branding:
  icon: 'check-circle'
  color: 'blue'

inputs:
  config:
    description: 'Path to Ship Shape configuration file'
    required: false
    default: '.shipshape.yml'

  fail-on:
    description: 'Fail build on severity level (critical, high, medium, low)'
    required: false
    default: 'critical'

  github-token:
    description: 'GitHub token for PR comments and checks'
    required: true

  comment-pr:
    description: 'Comment on PR with results'
    required: false
    default: 'true'

  upload-report:
    description: 'Upload HTML report as artifact'
    required: false
    default: 'true'

outputs:
  overall-score:
    description: 'Overall quality score (0-100)'

  grade:
    description: 'Letter grade (A+, A, B, C, D, F)'

  passed-gates:
    description: 'Whether all quality gates passed'

runs:
  using: 'docker'
  image: 'Dockerfile'
  args:
    - 'analyze'
    - '--config'
    - ${{ inputs.config }}
    - '--format'
    - 'github-action'
    - '--github-token'
    - ${{ inputs.github-token }}
```

**Dockerfile**:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /shipshape ./cmd/shipshape

FROM alpine:latest
RUN apk add --no-cache git

COPY --from=builder /shipshape /usr/local/bin/shipshape
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
```

**entrypoint.sh**:
```bash
#!/bin/sh
set -e

echo "🚢 Running Ship Shape Analysis..."

# Run Ship Shape
shipshape "$@"

# Exit with appropriate code
exit $?
```

### 2. Workflow YAML Parser and Analyzer

```go
type WorkflowAnalyzer struct {
    parser *WorkflowParser
}

type WorkflowAnalysis struct {
    Workflows          []*Workflow
    TestJobs           []*TestJob
    Optimizations      []*Optimization
    BestPractices      []*BestPracticeViolation
    EstimatedRuntime   time.Duration
}

type Workflow struct {
    Name     string
    FilePath string
    Triggers []string
    Jobs     []*Job
}

type Job struct {
    Name          string
    RunsOn        string
    Steps         []*Step
    Strategy      *MatrixStrategy
    Uses          string  // If it's a reusable workflow
    CachingUsed   bool
    Parallelized  bool
}

type Step struct {
    Name    string
    Uses    string  // Action reference
    Run     string  // Shell command
    With    map[string]string
    Env     map[string]string
}

// Analyze workflow for optimizations
func (wa *WorkflowAnalyzer) AnalyzeWorkflow(
    workflowPath string,
) (*WorkflowAnalysis, error) {
    content, err := os.ReadFile(workflowPath)
    if err != nil {
        return nil, err
    }

    var workflow struct {
        Name string
        On   interface{}
        Jobs map[string]struct {
            RunsOn   string `yaml:"runs-on"`
            Strategy struct {
                Matrix map[string][]interface{} `yaml:"matrix"`
            } `yaml:"strategy"`
            Steps []struct {
                Name string
                Uses string
                Run  string
                With map[string]string
            } `yaml:"steps"`
        } `yaml:"jobs"`
    }

    err = yaml.Unmarshal(content, &workflow)
    if err != nil {
        return nil, err
    }

    analysis := &WorkflowAnalysis{
        Optimizations: make([]*Optimization, 0),
    }

    // Detect missing caching
    for jobName, job := range workflow.Jobs {
        hasCache := false
        for _, step := range job.Steps {
            if step.Uses == "actions/cache@v3" || step.Uses == "actions/setup-node@v3" {
                hasCache = true
                break
            }
        }

        if !hasCache && containsTestSteps(job.Steps) {
            analysis.Optimizations = append(analysis.Optimizations, &Optimization{
                Type:        OptimizationTypeCaching,
                Job:         jobName,
                Description: "No dependency caching detected",
                Recommendation: "Add actions/cache for faster builds",
                Impact:      "Potential 30-60s speedup per run",
            })
        }
    }

    // Detect missing parallelization
    for jobName, job := range workflow.Jobs {
        if containsTestSteps(job.Steps) && job.Strategy.Matrix == nil {
            analysis.Optimizations = append(analysis.Optimizations, &Optimization{
                Type:        OptimizationTypeParallelization,
                Job:         jobName,
                Description: "Tests run sequentially",
                Recommendation: "Use matrix strategy to parallelize tests",
                Impact:      "Potential 2-5x speedup",
            })
        }
    }

    return analysis, nil
}

type OptimizationType string
const (
    OptimizationTypeCaching         OptimizationType = "caching"
    OptimizationTypeParallelization OptimizationType = "parallelization"
    OptimizationTypeArtifacts       OptimizationType = "artifacts"
    OptimizationTypeConditions      OptimizationType = "conditions"
)
```

### 3. GitHub Checks API Integration

```go
type GitHubCheckReporter struct {
    client *github.Client
    owner  string
    repo   string
    sha    string
}

func (gcr *GitHubCheckReporter) CreateCheckRun(
    ctx context.Context,
    report *Report,
) error {
    checkRun := &github.CreateCheckRunOptions{
        Name:       "Ship Shape Analysis",
        HeadSHA:    gcr.sha,
        Status:     github.String("completed"),
        Conclusion: github.String(determineConclusion(report)),
        Output: &github.CheckRunOutput{
            Title:   github.String(fmt.Sprintf("Score: %d/100 (Grade: %s)", report.Score, report.Grade)),
            Summary: github.String(generateSummary(report)),
            Text:    github.String(generateDetailedText(report)),
            Annotations: convertFindingsToAnnotations(report.Findings),
        },
    }

    _, _, err := gcr.client.Checks.CreateCheckRun(ctx, gcr.owner, gcr.repo, *checkRun)
    return err
}

func convertFindingsToAnnotations(findings []*Finding) []*github.CheckRunAnnotation {
    annotations := make([]*github.CheckRunAnnotation, 0)

    for _, finding := range findings {
        if finding.Location != nil {
            annotations = append(annotations, &github.CheckRunAnnotation{
                Path:            github.String(finding.Location.File),
                StartLine:       github.Int(finding.Location.Line),
                EndLine:         github.Int(finding.Location.Line),
                AnnotationLevel: github.String(severityToLevel(finding.Severity)),
                Message:         github.String(finding.Description),
                Title:           github.String(finding.Title),
            })
        }
    }

    return annotations
}

func severityToLevel(severity Severity) string {
    switch severity {
    case SeverityCritical, SeverityHigh:
        return "failure"
    case SeverityMedium:
        return "warning"
    default:
        return "notice"
    }
}
```

### 4. PR Comment Integration

**Strategy**: Single comment, updated on subsequent runs (avoid spam)

```go
type PRCommentManager struct {
    client *github.Client
    owner  string
    repo   string
    prNum  int
}

func (pcm *PRCommentManager) CreateOrUpdateComment(
    ctx context.Context,
    report *Report,
) error {
    commentBody := generatePRComment(report)

    // Find existing Ship Shape comment
    comments, _, err := pcm.client.Issues.ListComments(
        ctx, pcm.owner, pcm.repo, pcm.prNum, nil,
    )
    if err != nil {
        return err
    }

    commentMarker := "<!-- ship-shape-analysis -->"
    var existingComment *github.IssueComment

    for _, comment := range comments {
        if strings.Contains(comment.GetBody(), commentMarker) {
            existingComment = comment
            break
        }
    }

    fullBody := commentMarker + "\n" + commentBody

    if existingComment != nil {
        // Update existing comment
        _, _, err = pcm.client.Issues.EditComment(
            ctx, pcm.owner, pcm.repo, existingComment.GetID(),
            &github.IssueComment{Body: github.String(fullBody)},
        )
    } else {
        // Create new comment
        _, _, err = pcm.client.Issues.CreateComment(
            ctx, pcm.owner, pcm.repo, pcm.prNum,
            &github.IssueComment{Body: github.String(fullBody)},
        )
    }

    return err
}

func generatePRComment(report *Report) string {
    var sb strings.Builder

    sb.WriteString("## 🚢 Ship Shape Analysis Results\n\n")
    sb.WriteString(fmt.Sprintf("**Overall Score**: %d/100 (%s)\n\n", report.Score, report.Grade))

    // Score breakdown
    sb.WriteString("### Score Breakdown\n\n")
    sb.WriteString("| Dimension | Score | Status |\n")
    sb.WriteString("|-----------|-------|--------|\n")
    for dim, score := range report.DimensionScores {
        status := scoreToEmoji(score)
        sb.WriteString(fmt.Sprintf("| %s | %.0f/100 | %s |\n", dim, score, status))
    }
    sb.WriteString("\n")

    // Top issues
    if len(report.TopIssues) > 0 {
        sb.WriteString("### 🔍 Top Issues\n\n")
        for i, issue := range report.TopIssues[:min(5, len(report.TopIssues))] {
            sb.WriteString(fmt.Sprintf("%d. **%s**: %s\n", i+1, issue.Severity, issue.Title))
        }
        sb.WriteString("\n")
    }

    // Quality gates
    if len(report.GateResults) > 0 {
        sb.WriteString("### ⚡ Quality Gates\n\n")
        for _, gate := range report.GateResults {
            status := "✅"
            if !gate.Passed {
                status = "❌"
            }
            sb.WriteString(fmt.Sprintf("%s %s\n", status, gate.Name))
        }
    }

    sb.WriteString("\n📊 [View Full Report](link-to-artifact)\n")

    return sb.String()
}

func scoreToEmoji(score float64) string {
    switch {
    case score >= 90:
        return "🟢 Excellent"
    case score >= 80:
        return "🔵 Good"
    case score >= 70:
        return "🟡 Acceptable"
    default:
        return "🔴 Needs Improvement"
    }
}
```

### 5. Quality Gate Enforcement

```go
type QualityGateEvaluator struct {
    config *GateConfig
}

func (qge *QualityGateEvaluator) Evaluate(report *Report) (*GateResult, error) {
    result := &GateResult{
        Passed:  true,
        Failed:  make([]*GateCheck, 0),
        Warnings: make([]*GateCheck, 0),
    }

    for _, gate := range qge.config.BlockingGates {
        if !gate.Check(report) {
            result.Passed = false
            result.Failed = append(result.Failed, gate)
        }
    }

    for _, gate := range qge.config.WarningGates {
        if !gate.Check(report) {
            result.Warnings = append(result.Warnings, gate)
        }
    }

    return result, nil
}

// Exit code determination
func (qge *QualityGateEvaluator) ExitCode(result *GateResult) int {
    if !result.Passed {
        return 1  // Fail build
    }
    if len(result.Warnings) > 0 {
        return 2  // Warning exit code
    }
    return 0  // Success
}
```

## Prototype Requirements

### Deliverable 1: GitHub Action Package
**Files**: `action.yml`, `Dockerfile`, `entrypoint.sh`
- Action metadata
- Docker configuration
- Entry point script
- Input/output definitions

### Deliverable 2: Workflow Analyzer
**Files**: `internal/github/workflow_analyzer.go`
- YAML parser for workflows
- Optimization detector
- Best practices checker
- Performance estimator

### Deliverable 3: GitHub Integrations
**Files**: `internal/github/*.go`
- Checks API client
- PR comment manager
- Artifact uploader
- Authentication handler

### Deliverable 4: Quality Gate Enforcer
**Files**: `internal/gates/evaluator.go`
- Gate evaluation logic
- Exit code determination
- Gate result reporting

### Deliverable 5: Example Workflows
**Files**: `.github/workflows/examples/*.yml`
- Basic usage example
- Advanced configuration
- Monorepo usage
- Pre-commit integration

## Performance Benchmarks

| Repo Size | Analysis Time | Action Overhead | Total Time |
|-----------|---------------|-----------------|------------|
| Small     | 10s           | 5s              | 15s        |
| Medium    | 30s           | 5s              | 35s        |
| Large     | 60s           | 10s             | 70s        |

**Target**: Action overhead <10s for all sizes

## Validation Strategy

### 1. Real GitHub Repository Testing
Test in actual repos:
- Public open-source project
- Private repository (permissions testing)
- Monorepo with multiple languages
- Repository with existing checks

### 2. GitHub Actions Simulation
Use [act](https://github.com/nektos/act) for local testing:
```bash
act pull_request -j ship-shape-analysis
```

### 3. Marketplace Review Checklist
- Clear documentation
- Proper versioning
- Security scanning passed
- Example workflows included
- Branding assets

## Risk Mitigation

### Risk 1: GitHub API rate limiting
**Mitigation**:
- Use authenticated requests (higher limits)
- Cache API responses
- Batch operations where possible
- Monitor rate limit headers

### Risk 2: Action execution timeout (6 hours max)
**Mitigation**:
- Optimize analysis performance
- Provide configuration for subset analysis
- Incremental analysis option
- Clear timeout warnings

### Risk 3: Large PR comments
**Mitigation**:
- Truncate long comments
- Link to full report artifact
- Configurable detail level
- Summary-first approach

### Risk 4: Permission issues
**Mitigation**:
- Clear documentation of required permissions
- Graceful degradation if permissions missing
- Helpful error messages
- Permission check at start

## Go/No-Go Decision Criteria

### GO if:
- ✅ Action runs successfully in test repos
- ✅ PR comments display correctly
- ✅ Checks API integration works
- ✅ Workflow analysis detects optimizations
- ✅ Quality gates enforce correctly
- ✅ Performance acceptable (<10s overhead)
- ✅ Ready for marketplace publishing

### NO-GO if:
- ❌ Action fails in test repos
- ❌ Permissions issues unresolvable
- ❌ Performance unacceptable (>30s overhead)
- ❌ Cannot integrate with Checks API
- ❌ Security concerns

## Spike Deliverables

1. **GitHub Action Package**
   - action.yml definition
   - Docker image
   - Documentation

2. **Workflow Analyzer**
   - YAML parser
   - Optimization detector
   - Performance recommendations

3. **GitHub Integrations**
   - Checks API client
   - PR comment manager
   - Working examples

4. **Quality Gate System**
   - Gate evaluator
   - Exit code logic
   - Configuration examples

5. **Marketplace Assets**
   - README.md
   - Example workflows
   - Branding assets
   - Security documentation

## Integration Guidelines

**Example Workflow**:
```yaml
name: Ship Shape Analysis

on:
  pull_request:
    branches: [main]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Ship Shape
        uses: shipshape/action@v1
        with:
          config: .shipshape.yml
          fail-on: high
          github-token: ${{ secrets.GITHUB_TOKEN }}
          comment-pr: true
          upload-report: true

      - name: Upload Report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: shipshape-report
          path: shipshape-report.html
```

## Success Metrics
- [ ] Action runs in test repos
- [ ] PR comments work correctly
- [ ] Checks API integration functional
- [ ] Workflow analysis accurate
- [ ] Quality gates enforce correctly
- [ ] Performance targets met
- [ ] Marketplace ready

## Timeline
- **Week 1**: Action package and basic integration
- **Week 2**: Workflow analyzer and optimization detection
- **Week 3**: PR comments and Checks API
- **Week 4**: Quality gates and marketplace prep

## Sources and References
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Creating a Docker Container Action](https://docs.github.com/en/actions/creating-actions/creating-a-docker-container-action)
- [GitHub Checks API](https://docs.github.com/en/rest/checks)
- [go-github Library](https://github.com/google/go-github)
- [actionlint](https://github.com/rhysd/actionlint)
- [act - Local GitHub Actions](https://github.com/nektos/act)
