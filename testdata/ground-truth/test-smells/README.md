# Test Smells Ground Truth Examples

This directory contains curated examples of all 11 test smell types that Ship Shape detects.

## Purpose

These examples serve multiple purposes:

1. **Validation**: Verify Ship Shape's detection accuracy (≥90% precision and recall targets)
2. **Documentation**: Demonstrate what each test smell looks like in real code
3. **Testing**: Integration tests for smell detection algorithms
4. **Regression Prevention**: Ensure new changes don't break existing detection

## Organization

Each test smell has its own subdirectory with language-specific examples:

```
test-smells/
├── mystery-guest/
│   ├── go/
│   ├── python/
│   ├── javascript/
│   └── java/
├── eager-test/
│   └── ...
└── [other smells]/
```

Within each language directory:

```
go/
├── example1/              # Simple, clear example
│   ├── code_test.go       # Test code demonstrating the smell
│   ├── code.go           # Supporting production code
│   └── metadata.yml      # Expected detection results
├── example2/              # More complex scenario
│   └── ...
└── negative/             # Counter-examples (should NOT be detected)
    └── ...
```

## Test Smell Catalog

### 1. Mystery Guest (External Dependencies)

Tests that depend on external resources without explicit setup, making them hard to understand and brittle.

**Characteristics**:
- File I/O without explicit file creation in test
- Database calls in unit tests
- HTTP calls to external services without mocks
- Environment variable dependencies not documented

**Example**:
```go
func TestUserData(t *testing.T) {
    // Mystery: Where does users.json come from?
    data := loadUserDataFromFile("users.json")
    assert.NotNil(t, data)
}
```

**Fix**: Use explicit test fixtures or mock external dependencies.

### 2. Eager Test (Multiple Concerns)

Tests that verify multiple unrelated concerns in a single test function.

**Characteristics**:
- >5 assertions testing different functionality
- Testing multiple code paths in one test
- Single test covering create, update, delete operations

**Example**:
```go
func TestUserRegistration(t *testing.T) {
    // Creating user
    user := CreateUser("alice")
    assert.NotNil(t, user)

    // Validating email
    assert.True(t, user.EmailValid())

    // Sending welcome email
    err := SendWelcomeEmail(user)
    assert.NoError(t, err)

    // Cleaning up database
    err = DeleteUser(user.ID)
    assert.NoError(t, err)
}
```

**Fix**: Split into separate focused tests.

### 3. Lazy Test (Multiple Scenarios)

Multiple test scenarios combined in a single test function instead of separate tests or table-driven tests.

**Characteristics**:
- Manual loops iterating over test cases
- Sequential testing of multiple scenarios
- Not using language's test framework features (table-driven, parametrized)

**Example**:
```go
func TestMath(t *testing.T) {
    // Should use table-driven test instead
    assert.Equal(t, 4, Add(2, 2))
    assert.Equal(t, 0, Add(-1, 1))
    assert.Equal(t, 100, Add(50, 50))
}
```

**Fix**: Use table-driven tests (Go), parametrized tests (Python/Jest), or data providers (JUnit).

### 4. Obscure Test (Unclear Intent)

Tests with unclear purpose or assertions that are hard to understand.

**Characteristics**:
- Missing or unclear assertion messages
- Complex setup without comments
- Unclear test naming (Test1, Test2)
- Magic numbers and strings

**Example**:
```go
func TestUserCreation(t *testing.T) {
    u := NewUser("test")
    assert.NotNil(t, u)  // What exactly are we testing?
}
```

**Fix**: Clear naming, descriptive assertion messages, documented setup.

### 5. Conditional Test Logic

Tests containing if/else statements or loops that affect test behavior.

**Characteristics**:
- if/else branches in test code
- Loops modifying test behavior
- Different assertions based on runtime conditions

**Example**:
```go
func TestDataProcessing(t *testing.T) {
    data := LoadData()
    if len(data) > 0 {
        assert.True(t, Validate(data))
    } else {
        t.Skip("No data available")
    }
}
```

**Fix**: Split into separate tests or use proper skip conditions.

### 6. General Fixture (Overly Broad Setup)

Test fixtures that provide more setup than needed for individual tests.

**Characteristics**:
- Heavyweight setup used by only some tests
- Shared fixtures with unused components
- Setup time much longer than test execution

**Example**:
```python
class TestUserService(unittest.TestCase):
    def setUp(self):
        # Creates entire database, email service, cache, etc.
        # even though most tests only need user creation
        self.db = FullDatabaseSetup()
        self.email = EmailServiceSetup()
        self.cache = CacheSetup()
```

**Fix**: Use granular fixtures or setup only what each test needs.

### 7. Code Duplication

Repeated test code that should be refactored into helpers.

**Characteristics**:
- Copy-pasted setup code across tests
- Identical assertion patterns
- Repeated test structure

**Example**:
```go
func TestUserA(t *testing.T) {
    db := setupDB()
    defer db.Close()
    user := createTestUser()
    // ... test logic
}

func TestUserB(t *testing.T) {
    db := setupDB()  // Duplicated
    defer db.Close()
    user := createTestUser()  // Duplicated
    // ... different test logic
}
```

**Fix**: Extract common setup to helper functions or use test fixtures.

### 8. Assertion Roulette

Multiple assertions without descriptive messages, making failures hard to diagnose.

**Characteristics**:
- No assertion messages
- Generic error messages ("expected true")
- Ambiguous failure output

**Example**:
```go
func TestUserValidation(t *testing.T) {
    assert.True(t, len(user.Name) > 0)
    assert.True(t, user.Age >= 0)
    assert.True(t, user.Email != "")
    // If one fails, which validation failed?
}
```

**Fix**: Add descriptive messages to all assertions.

### 9. Sensitive Equality

Fragile equality assertions that break on irrelevant changes.

**Characteristics**:
- Deep object equality without field selection
- String matching on formatted output (e.g., timestamps)
- Exact time equality without tolerance

**Example**:
```go
func TestUserJSON(t *testing.T) {
    json := user.ToJSON()
    // Breaks if whitespace or field order changes
    expected := `{"name":"Alice","age":30,"created":"2024-01-01T00:00:00Z"}`
    assert.Equal(t, expected, json)
}
```

**Fix**: Use structured comparison, field-specific assertions, or approximate equality.

### 10. Resource Optimism

Tests that assume resources are available without verification or cleanup.

**Characteristics**:
- No cleanup of created files/databases
- Hardcoded file paths that may not exist
- Port number conflicts in parallel tests
- Assumes unlimited memory/disk

**Example**:
```go
func TestFileProcessing(t *testing.T) {
    // Assumes /tmp/test.dat exists and is writable
    ProcessFile("/tmp/test.dat")
    // No cleanup - file left behind
}
```

**Fix**: Create resources in test, verify availability, clean up afterward.

### 11. Flakiness

Tests that pass or fail non-deterministically.

**Characteristics**:
- Race conditions (missing synchronization)
- Time-dependent assertions (sleeps, exact time checks)
- Dependency on external services
- Improper cleanup causing state leakage

**Example**:
```go
func TestAsyncOperation(t *testing.T) {
    go ProcessData()
    time.Sleep(100 * time.Millisecond)  // Flaky: assumes 100ms is enough
    assert.True(t, IsComplete())
}
```

**Fix**: Proper synchronization, retries with timeout, or avoid time-dependent logic.

## Adding New Examples

When contributing a new example:

1. Choose the appropriate smell subdirectory
2. Create language-specific directory if it doesn't exist
3. Add numbered example directory (example1, example2, etc.)
4. Include:
   - Test code demonstrating the smell
   - Supporting production code if needed
   - `metadata.yml` with expected detections
5. Add negative counter-examples in `negative/` subdirectory
6. Get verification from 2+ engineers
7. Update this README if adding a new smell type

## Validation

Run validation suite to verify examples:

```bash
make validate-smells

# Expected output per smell:
# ✓ Mystery Guest: 8/8 examples detected (100%)
# ✓ Eager Test: 12/12 examples detected (100%)
# ...
```

## References

- [Metadata Schema](../metadata.schema.yml)
- [User Story SS-120](.project/user-stories.md#ss-120)
- [Test Smell Detection Requirements](.project/requirements.md)

---

**Last Updated**: 2026-01-28
**Examples**: 0 (structure only, examples to be added)
**Languages**: Go, Python, JavaScript, Java (planned)
