# Testing Guidelines

## Test Implementation Guidelines

### JSON Handling in Tests

✅ **DO:**
- Use `json.NewDecoder()` to parse responses into structs
- Use `json.Marshal()` to create request bodies
- Define proper Go structs with JSON tags for request/response formats
- Compare struct fields directly (e.g., `if response.Text != expectedText`)

❌ **DON'T:**
- Use `strings.Contains()` or other string operations to check JSON content
- Manually parse JSON strings with `strings.Index()` or regex
- Compare raw JSON strings

### Example of Correct Pattern:

```go
// Define response struct
type NoteResponse struct {
    ID   string `json:"id"`
    Text string `json:"text"`
}

// Parse response properly
var note NoteResponse
if err := json.NewDecoder(response.Body).Decode(&note); err != nil {
    t.Fatalf("Failed to parse response: %v", err)
}

// Compare fields directly
if note.Text != expectedText {
    t.Errorf("Expected text %q, got %q", expectedText, note.Text)
}
```

## Test Structure

### Test File Organization
- Tests are organized by user story: `tests/us-XXXX/`
- Each test file corresponds to one acceptance criterion
- Test names should be descriptive and match the acceptance criterion

### Test Components
1. **Setup**: Initialize test environment using `setupTest()`
2. **Authentication**: Login to get valid session (for API tests)
3. **Test Logic**: Execute the functionality being tested
4. **Verification**: Check responses using proper JSON parsing
5. **Cleanup**: Remove test data using `os.RemoveAll("./test_data")`

### Example Test Structure

```go
func TestFeatureName(t *testing.T) {
    // Setup
    setupTest()
    
    // Cleanup
    defer func() {
        os.RemoveAll("./test_data")
    }()
    
    // Authentication (if needed)
    loginReq, err := http.NewRequest("POST", "/login", bytes.NewBufferString(`{"username": "testuser", "password": "testpass"}`))
    // ... login logic ...
    
    // Test Logic
    req, err := http.NewRequest("POST", "/endpoint", bytes.NewBuffer(requestBody))
    // ... request execution ...
    
    // Verification
    var response ResponseType
    if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    if response.Field != expectedValue {
        t.Errorf("Expected %q, got %q", expectedValue, response.Field)
    }
}
```

## TDD Process

### Red Phase (Test Writing)
1. Identify the first unimplemented acceptance criterion
2. Write a test that defines the desired behavior
3. **Verify the test fails** - This confirms it's testing new functionality
4. The test should fail with a clear error indicating missing implementation

### Green Phase (Implementation)
1. Write the minimal code to make the test pass
2. Ensure all existing tests still pass
3. Verify the new test now passes

### Refactor Phase
1. Improve code quality without breaking tests
2. All tests must continue to pass
3. Focus on readability, maintainability, and performance

## Test Data Management

### Test Data Location
- Test data is stored in `./test_data/` directory
- Each test should clean up its own data
- Use unique names/IDs to avoid conflicts between tests

### Cleanup Requirements
- Every test MUST include cleanup logic
- Use `defer` to ensure cleanup runs even if test fails
- Remove all created files and directories

```go
// Correct cleanup pattern
defer func() {
    os.RemoveAll("./test_data")
}()
```

## Common Test Patterns

### API Testing Pattern
```go
// 1. Setup test environment
setupTest()
defer os.RemoveAll("./test_data")

// 2. Authenticate (if needed)
loginReq := http.NewRequest("POST", "/login", bytes.NewBufferString(credentials))
// ... get session ID ...

// 3. Create request with proper JSON marshaling
requestBody, _ := json.Marshal(RequestStruct{...})
req := http.NewRequest("POST", "/endpoint", bytes.NewBuffer(requestBody))
req.Header.Set("X-Session-ID", sessionID)
req.Header.Set("Content-Type", "application/json")

// 4. Execute request
rr := httptest.NewRecorder()
http.DefaultServeMux.ServeHTTP(rr, req)

// 5. Verify response
if rr.Code != http.StatusExpected {
    t.Errorf("Expected status %d, got %d", http.StatusExpected, rr.Code)
}

var response ResponseStruct
json.NewDecoder(rr.Body).Decode(&response)

// 6. Verify response fields
if response.Field != expectedValue {
    t.Errorf("Expected %q, got %q", expectedValue, response.Field)
}
```

### CLI Testing Pattern
```go
// 1. Setup test environment
setupTest()
defer os.RemoveAll("./test_data")

// 2. Execute CLI command
cmd := exec.Command("go", "run", "./cmd/cli", "command", "--flags")
var stdout, stderr bytes.Buffer
cmd.Stdout = &stdout
cmd.Stderr = &stderr

// 3. Verify execution
err := cmd.Run()
if err != nil {
    t.Fatalf("Command failed: %v, stderr: %s", err, stderr.String())
}

// 4. Verify output
if !strings.Contains(stdout.String(), expectedOutput) {
    t.Errorf("Expected output to contain %q, got %q", expectedOutput, stdout.String())
}
```

## Test Execution

### Running Tests
```bash
# Run all tests
cd tests && go test ./...

# Run specific test suite
cd tests && go test ./us-XXXX -v

# Run specific test
cd tests/us-XXXX && go test -v -run TestName
```

### Test Coverage
```bash
# Generate coverage report
cd tests && go test ./... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

## Best Practices

### Test Quality
- ✅ Test one specific behavior per test
- ✅ Use descriptive test names
- ✅ Keep tests independent (no shared state)
- ✅ Test both happy path and error cases
- ✅ Use table-driven tests for similar cases

### Test Maintenance
- ✅ Update tests when requirements change
- ✅ Fix broken tests immediately
- ✅ Review test changes in code reviews
- ✅ Keep test execution fast

### Anti-Patterns
- ❌ Tests that depend on external services
- ❌ Tests that modify production data
- ❌ Tests with sleep/timeout dependencies
- ❌ Overly complex test setup
- ❌ Tests that test implementation details