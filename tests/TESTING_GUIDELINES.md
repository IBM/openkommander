# Testing Guidelines for OpenKommander

## Overview

This document outlines the testing standards and guidelines for the OpenKommander project. Following these guidelines ensures consistent, reliable, and maintainable tests across the codebase.

## Project Structure

```
openkommander/
├── tests/                          # Testing documentation and shared utilities
│   ├── TESTING_GUIDELINES.md      # This file
│   ├── fixtures/                   # Test data and fixtures
│   └── utils/                      # Shared testing utilities
├── pkg/                            # Main packages
│   ├── api/
│   │   ├── router.go
│   │   └── router_test.go          # Unit tests for router
│   ├── cli/
│   │   ├── cli.go
│   │   └── cli_test.go             # Unit tests for CLI
│   └── ...
└── internal/                       # Internal packages
    └── core/
        ├── commands/
        │   ├── topic.go
        │   └── topic_test.go       # Unit tests for commands
        └── ...
```

## Test File Naming Conventions

- **Tests**: `filename_test.go` (same package as source)

## Test Function Naming

### Tests

```go
func TestFunctionName(t *testing.T)
func TestFunctionName_SpecificScenario(t *testing.T)
func TestStructName_MethodName(t *testing.T)
```

### What Constitutes a Qualifying Test

A qualifying test must:

1. **Test a single unit of functionality** (function, method, or small component)
2. **Be isolated** - not depend on external systems when possible
3. **Be deterministic** - same input always produces same output
4. **Have clear assertions** - verify expected behavior explicitly
5. **Cover edge cases** - test boundary conditions and error scenarios

#### Example of a Good Test:

```go
func TestParseTopicName_ValidInput(t *testing.T) {
    input := "my-topic-name"
    expected := "my-topic-name"
  
    result, err := ParseTopicName(input)
  
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}

func TestParseTopicName_InvalidInput(t *testing.T) {
    testCases := []struct {
        name        string
        input       string
        expectError bool
    }{
        {"empty string", "", true},
        {"invalid characters", "topic@name", true},
        {"too long", strings.Repeat("a", 300), true},
    }
  
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := ParseTopicName(tc.input)
            if tc.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Coverage Requirements

### Minimum Coverage Target

- **Overall Project**: 80%

### Coverage Commands

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Check coverage percentage
go tool cover -func=coverage.out | grep total
```

## Testing Best Practices

### 1. Arrange, Act, Assert (AAA) Pattern

```go
func TestCalculateTotal(t *testing.T) {
    // Arrange
    items := []Item{{Price: 10}, {Price: 20}}
  
    // Act
    total := CalculateTotal(items)
  
    // Assert
    assert.Equal(t, 30, total)
}
```

### 2. Use Subtests for Logical Grouping

```go
func TestUserService(t *testing.T) {
    t.Run("CreateUser", func(t *testing.T) {
        t.Run("ValidInput", func(t *testing.T) { /* test */ })
        t.Run("InvalidInput", func(t *testing.T) { /* test */ })
    })
  
    t.Run("GetUser", func(t *testing.T) {
        t.Run("ExistingUser", func(t *testing.T) { /* test */ })
        t.Run("NonExistentUser", func(t *testing.T) { /* test */ })
    })
}
```

### 3. Mock External Dependencies

```go
type MockKafkaClient struct {
    mock.Mock
}

func (m *MockKafkaClient) Send(topic, message string) error {
    args := m.Called(topic, message)
    return args.Error(0)
}

func TestProduceMessage(t *testing.T) {
    mockClient := new(MockKafkaClient)
    mockClient.On("Send", "test-topic", "test-message").Return(nil)
  
    service := &MessageService{client: mockClient}
    err := service.ProduceMessage("test-topic", "test-message")
  
    assert.NoError(t, err)
    mockClient.AssertExpectations(t)
}
```

### 4. Test Error Conditions

```go
func TestDivide_ByZero(t *testing.T) {
    _, err := Divide(10, 0)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "division by zero")
}
```

## Required Testing Libraries

### Core Libraries

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)
```

### Additional Libraries for Specific Use Cases

- **HTTP Testing**: `net/http/httptest`
- **Context Testing**: `context`
- **Time Testing**: Use fixed time or time mocking
- **Kafka Testing**: `github.com/segmentio/kafka-go` test utilities

## Continuous Integration Requirements

### Pre-commit Checks

Tests must:

1. Pass all unit tests
2. Pass linting (`golangci-lint`)
3. Pass formatting (`gofmt`)
4. Maintain coverage thresholds

### Pipeline Stages

1. **Unit Tests**: Fast feedback (< 2 minutes)
2. **Integration Tests**: Medium feedback (< 10 minutes)
3. **E2E Tests**: Comprehensive feedback (< 30 minutes)
4. **Coverage Report**: Generate and publish coverage

## Documentation Requirements

### Test Documentation

Each test file should include:

1. Package-level comment explaining what's being tested
2. Complex test functions should have comments
3. Test data should be self-explanatory

### Example:

```go
// Package api_test contains unit tests for the API router functionality.
// These tests verify HTTP routing, middleware, and response handling.
package api_test

// TestRouter_HandleTopics verifies that the topics endpoint correctly
// routes requests and returns appropriate responses for various scenarios.
func TestRouter_HandleTopics(t *testing.T) {
    // Test implementation
}
```

## Getting Started Checklist

- [ ] Read and understand these guidelines
- [ ] Set up your development environment with testing tools
- [ ] Write your first test using the AAA pattern
- [ ] Ensure your test follows naming conventions
- [ ] Verify test coverage meets requirements
- [ ] Run tests locally before committing
