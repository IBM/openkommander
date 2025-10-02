# Test Utilities

This directory contains shared testing utilities and helper functions used across the OpenKommander test suite.

## Best Practices

1. **Always use defer for cleanup**:

   ```go
   helper := testutils.NewTestHelper(t)
   defer helper.Cleanup()
   ```
2. **Use temp directories for file operations**:

   ```go
   tempDir := helper.CreateTempDir("mytest-*")
   configFile := filepath.Join(tempDir, "config.json")
   helper.WriteFile(configFile, `{"test": true}`)
   ```
3. **Mock external dependencies**:

   ```go
   mockServer := testutils.NewMockHTTPServer()
   defer mockServer.Close()
   // Use mockServer.URL() in your service configuration
   ```
4. **Wait for services to be ready**:

   ```go
   testutils.WaitForPort(t, "localhost", port, 30*time.Second)
   ```

## Import in Tests

```go
import (
    "testing"
    "github.com/IBM/openkommander/tests/utils"
)

func TestExample(t *testing.T) {
    helper := testutils.NewTestHelper(t)
    defer helper.Cleanup()
  
    // Use utilities here
}
```

## Guidelines

1. Keep utilities generic and reusable
2. Always provide cleanup mechanisms
3. Use meaningful names and good documentation
4. Follow the same quality standards as production code
5. Add tests for complex utility functions
