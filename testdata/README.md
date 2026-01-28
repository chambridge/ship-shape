# Ship Shape Test Data

This directory contains curated test data for validating Ship Shape's analysis capabilities.

## Directory Structure

```
testdata/
├── ground-truth/           # Curated examples with known, verified results
│   ├── test-smells/        # Test smell detection examples
│   ├── coverage-reports/   # Sample coverage reports for parsing validation
│   └── test-patterns/      # Language-specific test pattern examples
├── integration/            # Integration test repositories
└── validation/             # Real-world OSS project samples for validation
```

## Ground Truth Dataset

The `ground-truth/` directory contains carefully curated examples used to validate Ship Shape's detection accuracy. Each example includes:

- Source code demonstrating a specific pattern or smell
- Metadata file (`metadata.yml`) with expected results
- Documentation explaining the example

### Validation Requirements

All ground truth examples must meet these criteria:

1. **Verified by 2+ engineers**: Each example reviewed and agreed upon
2. **Documented expectations**: Clear metadata describing what should be detected
3. **Positive and negative examples**: Both instances of the pattern and counterexamples
4. **Real-world relevance**: Examples drawn from actual code patterns
5. **Version controlled**: All changes tracked with justification

### Metadata Format

Each example directory should contain a `metadata.yml` file:

```yaml
# Example: testdata/ground-truth/test-smells/eager-test/go/example1/metadata.yml
version: "1.0"
language: go
category: test-smell
smell_type: eager-test
description: "Test function testing multiple unrelated concerns"
verified_by:
  - "engineer1@example.com"
  - "engineer2@example.com"
verified_date: "2026-01-28"
expected_detections:
  - type: test-smell
    smell: eager-test
    file: "example_test.go"
    function: "TestUserRegistration"
    severity: medium
    reason: "Tests user creation, validation, email sending, and database cleanup in one test"
false_positives_expected: []
notes: |
  This is a classic eager test that could be split into:
  - TestUserCreation
  - TestUserValidation
  - TestEmailNotification
  - TestDatabaseCleanup
```

## Test Smells Directory

The `test-smells/` directory contains examples for all 11 detected test smell types:

### 1. Mystery Guest
Tests that depend on external resources without explicit setup.
- File I/O without explicit file creation
- Database calls in unit tests
- HTTP calls without mocks

### 2. Eager Test
Tests that verify multiple unrelated concerns in a single test function.
- Multiple assertions testing different functionality
- Testing unrelated code paths
- >5 assertions in one test

### 3. Lazy Test
Multiple test scenarios combined in a single test function instead of separate tests.
- Loop-based testing instead of table-driven
- Manual scenario switching instead of subtests

### 4. Obscure Test
Tests with unclear purpose or assertions that are hard to understand.
- Missing assertion messages
- Complex setup without comments
- Unclear test naming

### 5. Conditional Test Logic
Tests containing if/else statements or loops that affect test behavior.
- Conditional assertions
- Loop-based verification
- Test behavior varies by runtime conditions

### 6. General Fixture
Overly broad test fixtures that provide more setup than needed.
- Fixture used by only a subset of tests
- Heavy initialization for simple tests
- Shared state between unrelated tests

### 7. Code Duplication
Repeated test code that should be refactored into helpers.
- Copy-pasted setup code
- Duplicate assertions
- Repeated test patterns

### 8. Assertion Roulette
Multiple assertions without descriptive messages making failures hard to diagnose.
- No assertion messages
- Generic error messages
- Ambiguous failure output

### 9. Sensitive Equality
Fragile equality assertions that break on irrelevant changes.
- Deep object equality without field selection
- String matching on formatted output
- Time-based equality without tolerance

### 10. Resource Optimism
Tests that assume resources are available without verification.
- No cleanup of created resources
- Hardcoded file paths
- Port number conflicts

### 11. Flakiness
Tests that pass or fail non-deterministically.
- Race conditions
- Time-dependent assertions
- Dependency on external services
- Improper cleanup

## Coverage Reports Directory

The `coverage-reports/` directory contains sample coverage report files for testing Ship Shape's parsing capabilities:

- **Go**: `go test -coverprofile` format
- **Python**: coverage.py XML and JSON formats
- **JavaScript**: Istanbul/nyc LCOV and JSON formats
- **Java**: JaCoCo XML format

Each format includes:
- Minimal valid report
- Complex multi-file report
- Edge cases (0% coverage, 100% coverage)

## Test Patterns Directory

The `test-patterns/` directory demonstrates idiomatic test patterns for each language:

### Go Patterns
- Table-driven tests
- Subtests with t.Run
- TestMain usage
- Benchmark tests
- Example tests

### Python Patterns
- pytest fixtures
- parametrize decorators
- unittest.TestCase
- doctest examples

### JavaScript Patterns
- Jest describe/it blocks
- Mocha/Chai patterns
- async/await testing
- Mock and spy patterns

### Java Patterns
- JUnit 5 annotations
- Parameterized tests
- Test lifecycle methods
- AssertJ assertions

## Integration Test Repositories

The `integration/` directory contains complete mini-repositories for end-to-end testing:

- Simple single-language projects
- Multi-language monorepos
- Projects with various test frameworks
- Edge cases (no tests, mixed patterns)

## Validation Directory

The `validation/` directory contains references to real-world open-source projects used for validation:

- Project references (not full clones)
- Expected analysis results
- Performance benchmarks
- Regression test cases

## Usage in Tests

### Example: Testing Test Smell Detection

```go
func TestEagerTestDetection(t *testing.T) {
    // Load ground truth example
    examplePath := "testdata/ground-truth/test-smells/eager-test/go/example1"
    metadata := loadMetadata(t, filepath.Join(examplePath, "metadata.yml"))

    // Run Ship Shape analyzer
    results := analyzer.Analyze(examplePath)

    // Verify expected detections
    for _, expected := range metadata.ExpectedDetections {
        assert.Contains(t, results.Findings, expected)
    }

    // Verify no false positives
    assert.Len(t, results.Findings, len(metadata.ExpectedDetections))
}
```

## Contributing Ground Truth Examples

When adding new examples:

1. Create example in appropriate subdirectory
2. Write clear, self-contained code
3. Add comprehensive `metadata.yml`
4. Get verification from 2+ engineers
5. Document in PR why this example is valuable
6. Include both positive and negative cases

## Quality Standards

Ground truth dataset must maintain:

- **Precision target**: ≥90% (low false positive rate)
- **Recall target**: ≥90% (low false negative rate)
- **Coverage**: All supported smells, languages, and frameworks
- **Maintenance**: Regular review and updates with Ship Shape changes

## Validation Metrics

Ship Shape validates against ground truth in CI:

```bash
# Run validation suite
make validate-ground-truth

# Expected output:
# Ground Truth Validation Results:
# ✓ Test Smells: 127/130 detected (97.7% recall)
# ✓ False Positives: 2/129 (1.5% FP rate)
# ✓ Coverage Parsing: 48/48 formats (100%)
# ✓ Pattern Recognition: 89/92 patterns (96.7%)
```

## References

- [User Story SS-120: Ground Truth Dataset Management](.project/user-stories.md#ss-120)
- [Architecture: Validation Strategy](.project/architecture.md#validation)
- [Test Strategy Documentation](docs/v1.0.0/testing.md)

---

**Last Updated**: 2026-01-28
**Version**: 1.0.0
**Maintainer**: Ship Shape Development Team
