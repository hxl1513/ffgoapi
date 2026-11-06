# Task 7 Brief: Middleware Core (RequestID, CORS, Logger, Recover)

## Context

You are implementing Task 7 of 16. Tasks 1-6 are complete.

**Your task builds on:** common.Logger from Task 4

## Your Task

Implement core middleware: RequestID, CORS, Logger, and Recover.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/middleware/request-id.go`** — Request ID generation
2. **`server/middleware/cors.go`** — CORS support
3. **`server/middleware/logger.go`** — Request logging
4. **`server/middleware/recover.go`** — Panic recovery
5. **`server/middleware/request-id_test.go`** — RequestID tests

## Requirements

### request-id.go

```go
const RequestIDKey = "X-Request-ID"

func RequestID() gin.HandlerFunc
func GetRequestID(ctx context.Context) string

// RequestID generates a unique 32-char hex string via crypto/rand
// Sets c.Set(RequestIDKey, id) and c.Header(RequestIDKey, id)
```

### cors.go

```go
func CORS() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowCredentials = true
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"*"}
    return cors.New(config)
}
```

Import: `github.com/gin-contrib/cors`

### logger.go

```go
func Logger(logger *common.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery
        
        c.Next()
        
        latency := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        
        logger.Info("[GIN] %s | %d | %v | %s | %s %s",
            clientIP, statusCode, latency, method, path, query)
    }
}
```

### recover.go

```go
func Recover() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                debug.PrintStack()
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "success": false,
                    "message": fmt.Sprintf("panic: %v", err),
                })
            }
        }()
        c.Next()
    }
}
```

## Test Requirements (request-id_test.go only)

```go
func TestRequestID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := &testingResponseWriter{}
    req, _ := http.NewRequest("GET", "/test", nil)
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    
    handler := RequestID()
    handler(c)
    
    requestID := c.GetString("X-Request-ID")
    if requestID == "" { t.Error("RequestID not set") }
    
    if w.Header().Get("X-Request-ID") == "" { t.Error("RequestID not in header") }
}
```

Minimal testingResponseWriter:

```go
type testingResponseWriter struct{}
func (w *testingResponseWriter) Header() http.Header { return http.Header{} }
func (w *testingResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (w *testingResponseWriter) WriteHeader(int) {}
```

## Global Constraints

- Go 1.25+
- Import cors: `github.com/gin-contrib/cors`

## Test Verification

```bash
cd server && go test ./middleware -v
```

Expected: At least RequestID test passes

## Commit Message

```
feat(middleware): add core middleware

- Implement RequestID middleware with UUID generation
- Add CORS middleware with configurable origins
- Implement Logger middleware with request tracking
- Add Recover middleware for panic handling
- Add unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-7-report.md`
