# Task 15 Brief: Router Registration Helper

## Context

You are implementing Task 15 of 16.

## Your Task

Add versioning support to route registration.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/router.go`** (enhance existing)
2. **`server/router_test.go`**

## Requirements

### router.go

```go
type RouteFunc func(r *gin.RouterGroup)

func (s *Server) SetupRoutes(version string, routes ...RouteFunc) {
    api := s.engine.Group("/api")
    if version != "" {
        api = api.Group("/" + version)
    }
    for _, route := range routes {
        route(api)
    }
}

func (s *Server) RegisterRoutes(routes ...RouteFunc) {
    s.SetupRoutes("", routes...)
}
```

### Test

```go
func TestRegisterRoutes(t *testing.T) {
    gin.SetMode(gin.TestMode)
    cfg := &config.Config{...}
    srv, _ := NewServer(cfg)
    
    called := false
    srv.RegisterRoutes(func(r *gin.RouterGroup) {
        called = true
    })
    
    if !called { t.Error("Route function was not called") }
}
```

## Test Verification

```bash
cd server && go test -v
```

## Commit Message

```
feat(router): add versioned route registration

- Add version parameter to SetupRoutes
- Maintain backward compatibility with RegisterRoutes
- Add unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-15-report.md`
