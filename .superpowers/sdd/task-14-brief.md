# Task 14 Brief: Request Body Limit Middleware

## Context

You are implementing Task 14 of 16.

## Your Task

Implement request body size limiting middleware.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/middleware/request-body-limit.go`**
2. **`server/middleware/body-cleanup.go`**
3. **`server/middleware/request-body-limit_test.go`**

## Requirements

### request-body-limit.go

```go
const defaultMaxBodyBytes = 4 << 20 // 4MB

func RequestBodyLimit(maxBytes int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        if maxBytes <= 0 {
            maxBytes = defaultMaxBodyBytes
        }
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
        c.Next()
        if c.Request.Body != nil {
            c.Request.Body.Close()
        }
    }
}
```

### body-cleanup.go

```go
func BodyCleanup() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if c.Request.Body != nil {
            c.Request.Body.Close()
        }
    }
}
```

### Test

```go
func TestRequestBodyLimit(t *testing.T) {
    gin.SetMode(gin.TestMode)
    largeBody := bytes.Repeat([]byte("a"), 1024*1024) // 1MB
    
    w := &testingResponseWriter{}
    req, _ := http.NewRequest("POST", "/test", bytes.NewReader(largeBody))
    req.ContentLength = int64(len(largeBody))
    
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    
    handler := RequestBodyLimit(512 * 1024) // 512KB limit
    handler(c)
    
    if c.Writer.Status() != http.StatusRequestEntityTooLarge {
        t.Errorf("Expected 413, got %d", c.Writer.Status())
    }
}
```

## Test Verification

```bash
cd server && go test ./middleware -v -run TestRequestBodyLimit
```

## Commit Message

```
feat(middleware): add request body limit and cleanup

- Implement RequestBodyLimit middleware
- Add BodyCleanup middleware for resource cleanup
- Add unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-14-report.md`
