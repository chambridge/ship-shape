# Test Patterns Ground Truth Examples

This directory contains examples of idiomatic test patterns for each supported language.

## Purpose

These examples demonstrate:

1. **Best Practices**: Idiomatic test patterns for each language
2. **Framework Usage**: Proper use of testing frameworks
3. **Pattern Recognition**: Help Ship Shape identify good test practices
4. **Documentation**: Reference examples for developers

## Language-Specific Patterns

### Go Test Patterns

**Location**: `go/`

#### Table-Driven Tests
The idiomatic Go testing pattern for testing multiple scenarios.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -1, -2},
        {"zero", 0, 5, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

#### Subtests
Using `t.Run()` for hierarchical test organization.

#### Test Helpers
Using `t.Helper()` to mark helper functions.

#### TestMain
Setup and teardown for entire test suite.

#### Benchmarks
Performance benchmarking with `testing.B`.

### Python Test Patterns

**Location**: `python/`

#### pytest Fixtures
Dependency injection for test setup.

```python
@pytest.fixture
def user_service():
    service = UserService()
    yield service
    service.cleanup()

def test_create_user(user_service):
    user = user_service.create("alice")
    assert user.name == "alice"
```

#### Parametrized Tests
Testing multiple scenarios with `@pytest.mark.parametrize`.

```python
@pytest.mark.parametrize("input,expected", [
    (2, 4),
    (3, 9),
    (4, 16),
])
def test_square(input, expected):
    assert square(input) == expected
```

#### unittest.TestCase
Class-based testing with setUp/tearDown.

#### doctest
Embedded tests in docstrings.

### JavaScript Test Patterns

**Location**: `javascript/`

#### Jest describe/it Blocks
Hierarchical test organization.

```javascript
describe('UserService', () => {
    let service;

    beforeEach(() => {
        service = new UserService();
    });

    describe('createUser', () => {
        it('should create user with valid name', () => {
            const user = service.createUser('alice');
            expect(user.name).toBe('alice');
        });

        it('should throw error for empty name', () => {
            expect(() => service.createUser('')).toThrow();
        });
    });
});
```

#### Async/Await Testing
Testing asynchronous code.

```javascript
test('fetches user data', async () => {
    const data = await fetchUser(123);
    expect(data).toHaveProperty('id', 123);
});
```

#### Mock and Spy Patterns
Using Jest mocks and spies.

```javascript
const mockFn = jest.fn();
mockFn('arg');
expect(mockFn).toHaveBeenCalledWith('arg');
```

### Java Test Patterns

**Location**: `java/`

#### JUnit 5 Annotations
Modern JUnit testing with `@Test`, `@BeforeEach`, etc.

```java
class UserServiceTest {
    private UserService service;

    @BeforeEach
    void setUp() {
        service = new UserService();
    }

    @Test
    void createUser_ValidName_ReturnsUser() {
        User user = service.createUser("alice");
        assertEquals("alice", user.getName());
    }
}
```

#### Parameterized Tests
Data-driven testing with `@ParameterizedTest`.

```java
@ParameterizedTest
@ValueSource(strings = {"alice", "bob", "charlie"})
void createUser_ValidNames_ReturnsUser(String name) {
    User user = service.createUser(name);
    assertNotNull(user);
}
```

#### AssertJ Assertions
Fluent assertion library.

```java
assertThat(user)
    .isNotNull()
    .hasFieldOrPropertyWithValue("name", "alice")
    .extracting("age")
    .isEqualTo(30);
```

## Directory Structure

```
test-patterns/
├── go/
│   ├── table-driven/
│   │   ├── example_test.go
│   │   └── metadata.yml
│   ├── subtests/
│   ├── helpers/
│   ├── test-main/
│   └── benchmarks/
├── python/
│   ├── pytest-fixtures/
│   ├── parametrize/
│   ├── unittest/
│   └── doctest/
├── javascript/
│   ├── jest-describe/
│   ├── async-await/
│   ├── mocking/
│   └── snapshot/
└── java/
    ├── junit5/
    ├── parameterized/
    └── assertj/
```

## Example Metadata

```yaml
version: "1.0.0"
language: go
category: test-pattern
pattern_type: table-driven
description: "Demonstrates idiomatic Go table-driven test pattern"
verified_by:
  - "engineer1@example.com"
  - "engineer2@example.com"
verified_date: "2026-01-28"

expected_detections:
  - type: pattern
    pattern: table-driven-test
    file: "example_test.go"
    function: "TestAdd"
    reason: "Uses struct slice with test cases and t.Run for subtests"
    confidence: high

  - type: best-practice
    pattern: test-helper-usage
    file: "example_test.go"
    function: "assertNoError"
    reason: "Helper function properly calls t.Helper()"
    confidence: high

tags: ["idiomatic", "best-practice", "table-driven"]
notes: |
  This is the recommended Go testing pattern for testing multiple scenarios.
  Benefits:
  - Clear separation of test data and test logic
  - Easy to add new test cases
  - Subtests run independently
  - Clear failure messages
```

## Pattern Recognition

Ship Shape should recognize these patterns and:

1. **Identify Best Practices**: Reward use of idiomatic patterns
2. **Suggest Improvements**: Recommend patterns when appropriate
3. **Framework Detection**: Correctly identify testing frameworks
4. **Quality Scoring**: Higher scores for idiomatic patterns

## Anti-Patterns

Examples may also demonstrate anti-patterns to avoid:

- Not using language's testing features
- Reinventing framework functionality
- Ignoring framework conventions
- Over-complicated test structure

## Validation

Pattern recognition tests verify Ship Shape can:

```bash
# Detect table-driven tests
shipshape analyze testdata/ground-truth/test-patterns/go/table-driven
# Expected: ✓ Detected 1 table-driven test pattern

# Detect pytest fixtures
shipshape analyze testdata/ground-truth/test-patterns/python/pytest-fixtures
# Expected: ✓ Detected 3 fixture usage patterns
```

## Adding New Patterns

When adding test pattern examples:

1. **Follow Language Idioms**: Use canonical examples from official docs
2. **Document Benefits**: Explain why this pattern is recommended
3. **Show Complete Example**: Include all necessary code
4. **Add Metadata**: Specify expected detections
5. **Include Counter-examples**: Show what NOT to do

## References

### Go
- [Go Testing Documentation](https://go.dev/doc/testing)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)

### Python
- [pytest Documentation](https://docs.pytest.org/)
- [pytest Fixtures](https://docs.pytest.org/en/stable/fixture.html)

### JavaScript
- [Jest Documentation](https://jestjs.io/)
- [Testing Best Practices](https://github.com/goldbergyoni/javascript-testing-best-practices)

### Java
- [JUnit 5 User Guide](https://junit.org/junit5/docs/current/user-guide/)
- [AssertJ Documentation](https://assertj.github.io/doc/)

---

**Last Updated**: 2026-01-28
**Patterns**: 15+ planned across 4 languages
**Status**: Structure only, pattern examples to be added
