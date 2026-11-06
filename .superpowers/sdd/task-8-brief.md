# Task 8 Brief: Auth and Rate Limit Middleware

## Context

You are implementing Task 8 of 16. Tasks 1-7 are complete.

**Your task builds on:** 
- model.Cache (Task 6)
- common.Logger (Task 4)
- types.AppError (Task 3)

## Your Task

Implement authentication and rate limiting middleware.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/middleware/auth.go`** — Authentication middleware
2. **`server/middleware/rate-limit.go`** — Rate limiting middleware
3. **`server/middleware/auth_test.go`** — Auth tests
4. **`server/middleware/rate-limit_test.go`** — Rate limit tests (optional if combined)

## Requirements

### auth.go

```go
type AuthConfig struct {
    SessionSecret string
    TokenExpiry   int
}

func AuthMiddleware(cfg *AuthConfig, cache model.Cache, allowedRoles ...int) gin.HandlerFunc
func AdminAuth(cfg *AuthConfig, cache model.Cache) gin.HandlerFunc
func UserAuth(cfg *AuthConfig, cache model.Cache) gin.HandlerFunc
```

**AuthMiddleware behavior:**
1. Try session auth (return false for now - placeholder)
2. Try Bearer token from Authorization header
3. Parse JWT using `pkg/crypto.ParseToken(cfg.SessionSecret, token)`
4. Check role if allowedRoles specified
5. Set user context: `c.Set("user_id", claims.UserID)`, `c.Set("username", claims.Username)`, `c.Set("role", claims.Role)`
6. Return 401 for missing/invalid token
7. Return 403 for insufficient permissions

### rate-limit.go

```go
type RateLimitConfig struct {
    Enabled         bool
    DefaultMax      int
    DefaultDuration int // seconds
}

func RateLimit(cfg *RateLimitConfig, cache model.Cache) gin.HandlerFunc
```

**RateLimit behavior:**
1. If disabled, call c.Next()
2. Generate key: `"rate_limit:" + c.ClientIP()`
3. Use cache.Get/Set/Incr for sliding window
4. Return 429 if limit exceeded
5. On cache error, fail open (allow request)

### Test Requirements

```go
func TestAuthMiddleware_NoToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    cache := model.NewMemoryCache()
    cfg := &AuthConfig{SessionSecret: "test-secret"}
    
    w := &testingResponseWriter{}
    req, _ := http.NewRequest("GET", "/test", nil)
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    c.Set("cache", cache)
    
    handler := AuthMiddleware(cfg, cache, 0)
    handler(c)
    
    if c.Writer.Status() != http.StatusUnauthorized {
        t.Errorf("Expected 401, got %d", c.Writer.Status())
    }
}

func TestRateLimitMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    cache := model.NewMemoryCache()
    cfg := &RateLimitConfig{
        Enabled:         true,
        DefaultMax:      2,
        DefaultDuration: 60,
    }
    
    // First 2 requests should pass, 3rd should be rate limited
    // ... (similar to plan Task 8)
}
```

## Global Constraints

- Go 1.25+
- Use pkg/crypto for JWT parsing
- Use model.Cache for rate limit storage

## Test Verification

```bash
cd server && go test ./middleware -v
```

Expected: At least auth tests pass

## Commit Message

```
feat(middleware): add auth and rate limiting middleware

- Implement AuthMiddleware with JWT token support
- Add AdminAuth and UserAuth convenience functions
- Implement RateLimit middleware with sliding window
- Support both memory and Redis cache backends
- Add unit tests for authentication and rate limiting

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-8-report.md`
