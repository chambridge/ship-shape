# Validation Directory - Real-World OSS Project Testing

This directory contains references to real-world open-source projects used for validating Ship Shape against production code.

## Purpose

Real-world validation ensures:

1. **Production Accuracy**: Ship Shape works on actual codebases, not just curated examples
2. **Performance Validation**: Handles large repositories efficiently
3. **Edge Case Discovery**: Find issues not covered by unit tests
4. **Benchmark Baselines**: Establish performance and accuracy baselines
5. **Regression Detection**: Ensure updates don't break real-world analysis

## Approach

Rather than cloning entire repositories (which would bloat the test suite), we maintain:

- **Project References**: Links to specific commits/versions
- **Expected Results**: Cached analysis results for comparison
- **Test Scripts**: Automated validation against live repos

## Validated Projects

### Tier 1: Core Validation (Always Run)

Small to medium projects that run quickly in CI:

```yaml
projects:
  - name: go-sqlmock
    language: go
    repository: https://github.com/DATA-DOG/go-sqlmock
    commit: stable-tag-or-commit
    size: small
    test_framework: testing
    reason: "Well-tested Go library with standard patterns"

  - name: requests
    language: python
    repository: https://github.com/psf/requests
    commit: v2.31.0
    size: medium
    test_framework: pytest
    reason: "Popular Python library with comprehensive tests"

  - name: lodash
    language: javascript
    repository: https://github.com/lodash/lodash
    commit: 4.17.21
    size: medium
    test_framework: jest
    reason: "Widely-used JS utility library with extensive tests"
```

### Tier 2: Extended Validation (Nightly/Weekly)

Larger projects for comprehensive validation:

```yaml
projects:
  - name: kubernetes
    language: go
    repository: https://github.com/kubernetes/kubernetes
    commit: v1.28.0
    size: large
    test_framework: testing
    reason: "Complex monorepo with extensive test suite"

  - name: django
    language: python
    repository: https://github.com/django/django
    commit: stable/4.2.x
    size: large
    test_framework: unittest
    reason: "Major framework with mature testing practices"

  - name: vscode
    language: typescript
    repository: https://github.com/microsoft/vscode
    commit: latest-stable
    size: xlarge
    test_framework: mocha
    reason: "Large TypeScript project with complex testing"
```

### Tier 3: Specific Feature Validation

Projects selected for specific features:

```yaml
projects:
  - name: gin
    language: go
    repository: https://github.com/gin-gonic/gin
    commit: v1.9.1
    focus: "benchmark tests, table-driven patterns"

  - name: flask
    language: python
    repository: https://github.com/pallets/flask
    focus: "pytest fixtures, doctest usage"

  - name: react
    language: javascript
    repository: https://github.com/facebook/react
    focus: "monorepo structure, jest patterns"
```

## Directory Structure

```
validation/
├── README.md                    # This file
├── projects.yml                 # Project definitions
├── expected-results/            # Cached expected results
│   ├── go-sqlmock-v1.5.0.yml
│   ├── requests-v2.31.0.yml
│   └── lodash-4.17.21.yml
├── scripts/
│   ├── validate-projects.sh     # Validation runner
│   ├── update-baselines.sh      # Update expected results
│   └── clone-and-analyze.sh     # Clone repo and run analysis
└── reports/
    └── validation-report.html   # Latest validation results
```

## Project Definition Format

Each project in `projects.yml`:

```yaml
- id: go-sqlmock
  name: go-sqlmock
  description: "Mock library for database/sql"
  language: go
  repository: https://github.com/DATA-DOG/go-sqlmock
  commit: 7ae8f5ca866e70c86e17f9d30f8a5dc55ca73045
  tag: v1.5.0
  tier: 1
  size: small
  estimated_analysis_time_seconds: 5
  test_framework: testing
  validation_focus:
    - table-driven-tests
    - mock-patterns
    - test-coverage
  known_issues:
    - "Uses some build tags that may affect coverage"
  notes: |
    Well-maintained library with good test practices.
    Useful for validating mock detection.
```

## Expected Results Format

Each project has an `expected-results` file:

```yaml
# expected-results/go-sqlmock-v1.5.0.yml
project: go-sqlmock
commit: 7ae8f5ca866e70c86e17f9d30f8a5dc55ca73045
analysis_date: "2026-01-28"
ship_shape_version: "1.0.0"

results:
  languages:
    - language: go
      files: 12
      test_files: 8
      percentage: 100.0

  test_frameworks:
    - framework: testing
      files: 8

  coverage:
    overall: 87.5
    packages:
      - path: "."
        coverage: 87.5

  test_smells:
    total: 2
    by_severity:
      medium: 2
    findings:
      - file: "rows_test.go"
        line: 156
        smell: eager-test
        severity: medium

  quality_score:
    overall: 85.0
    breakdown:
      test_coverage: 87.5
      test_quality: 82.0
      tool_adoption: 85.0

  performance:
    analysis_time_seconds: 3.2
    files_analyzed: 12
    lines_of_code: 2847

validation_criteria:
  coverage_tolerance: 2.0
  score_tolerance: 5.0
  smell_count_tolerance: 1
  max_analysis_time_seconds: 10.0
```

## Validation Script

Automated validation against real projects:

```bash
#!/bin/bash
# scripts/validate-projects.sh

set -euo pipefail

TIER=${1:-1}  # Default to tier 1 projects

echo "=== Ship Shape Real-World Validation ==="
echo "Tier: $TIER"

# Read projects from projects.yml
projects=$(yq eval ".[] | select(.tier == $TIER) | .id" projects.yml)

for project in $projects; do
    echo "Validating: $project"

    # Clone project
    ./scripts/clone-and-analyze.sh "$project"

    # Compare with expected results
    actual="reports/${project}-actual.yml"
    expected="expected-results/${project}.yml"

    # Run comparison
    ./scripts/compare-results.sh "$expected" "$actual"
done

echo "=== Validation Complete ==="
```

## Validation Criteria

Projects are validated against:

1. **Accuracy Metrics**:
   - Coverage parsing accuracy: ±2%
   - Test smell detection: ±1 finding
   - Quality score: ±5 points

2. **Performance Metrics**:
   - Analysis time within expected range
   - Memory usage within limits
   - No crashes or hangs

3. **Completeness Metrics**:
   - All expected languages detected
   - All test frameworks identified
   - Complete analysis (no skipped files)

## Updating Baselines

When Ship Shape improves or fixes bugs, update expected results:

```bash
# Update baseline for specific project
./scripts/update-baselines.sh go-sqlmock

# Update all tier 1 baselines
./scripts/update-baselines.sh --tier 1

# Update after reviewing changes
git add expected-results/
git commit -m "chore: update validation baselines for improved detection"
```

## Validation Frequency

- **Tier 1**: Every PR (fast, critical validation)
- **Tier 2**: Nightly builds
- **Tier 3**: Weekly or before releases

## Adding New Projects

When adding a project for validation:

1. **Choose Wisely**: Select projects with:
   - Good test coverage
   - Active maintenance
   - Representative patterns
   - Stable APIs/structure

2. **Document Purpose**: Explain why this project is valuable for validation

3. **Establish Baseline**: Run Ship Shape and capture expected results

4. **Set Tolerance**: Define acceptable variance (coverage, scores, etc.)

5. **Review Regularly**: Ensure projects remain representative

## Known Challenges

### Large Repositories

Large repos (Kubernetes, Chrome) require:
- Shallow clones to save space
- Longer analysis times
- More memory
- May only run in nightly builds

### Dependency Installation

Some projects require:
- Go modules download
- npm install
- pip install
- May fail if dependencies unavailable

### Build Systems

Projects using:
- Bazel
- Custom build tools
- May require additional setup

## Validation Report

After each validation run, generate report:

```
=== Ship Shape Validation Report ===
Date: 2026-01-28
Ship Shape Version: 1.0.0
Projects Validated: 10/10

Results:
✓ go-sqlmock    : PASS (score: 85/85, coverage: 87.5%/87.5%)
✓ requests      : PASS (score: 82/80, coverage: 89%/90%)
✓ lodash        : PASS (score: 88/90, coverage: 95%/95%)
...

Performance:
Average analysis time: 4.2s (expected: 5.0s)
Peak memory usage: 182MB (limit: 256MB)

Accuracy:
Coverage parsing: 100% within tolerance
Smell detection: 95% exact match
Quality scores: 100% within tolerance

Issues Found: 0
Regressions: 0

Overall: PASS ✓
```

## CI Integration

Validation runs in CI:

```yaml
# .github/workflows/validation.yml
name: Real-World Validation

on:
  pull_request:
  schedule:
    - cron: '0 2 * * *'  # Nightly at 2am

jobs:
  validate-tier1:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Tier 1 Validation
        run: ./testdata/validation/scripts/validate-projects.sh 1
```

---

**Last Updated**: 2026-01-28
**Tier 1 Projects**: 0 (to be added)
**Tier 2 Projects**: 0 (to be added)
**Status**: Structure ready, project validation to be implemented
