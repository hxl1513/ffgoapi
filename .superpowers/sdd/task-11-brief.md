# Task 11 Brief: HTTP Server Core

## Context

You are implementing Task 11 of 16. Tasks 1-10 are complete.

**Your task builds on:** All previous components

## Your Task

Implement HTTP server initialization, middleware setup, and routing.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/server.go`** — Main server implementation
2. **`server/router.go`** — Route registration helpers
3. **`server/server_test.go`** — Server tests

## Requirements

### server.go

```go
type Server struct {
    engine *gin.Engine
    config *config.Config
    logger *common.Logger
    db     *model.DB
    cache  model.Cache
}

func NewServer(cfg *config.Config) (*Server, error)
func (s *Server) Engine() *gin.Engine
func (s *Server) Start() error
func (s *Server) Shutdown(ctx context.Context) error
```

### NewServer Implementation

1. Set Gin mode (ReleaseMode or DebugMode based on cfg.Server.Mode)
2. Create `gin.New()` engine
3. Initialize logger: `common.NewLogger(gin.DefaultWriter, cfg.Logger.Level)`
4. Initialize database: `model.NewDB(&cfg.Database)`
5. Initialize cache: `model.NewCache(cfg.Redis.URL)`
6. Call `setupMiddleware()` — register:
   - `middleware.Recover()`
   - `middleware.RequestID()`
   - `middleware.Logger(s.logger)`
   - `middleware.CORS()`
7. Call `setupDefaultRoutes()` — register health check:
   ```go
   s.engine.GET("/health", func(c *gin.Context) {
       c.JSON(200, gin.H{"status": "ok"})
   })
   ```

### router.go

```go
type RouteFunc func(r *gin.RouterGroup)

func (s *Server) SetupRoutes(version string, routes ...RouteFunc)
func (s *Server) RegisterRoutes(routes ...RouteFunc)
```

- SetupRoutes creates `/api` group, optionally with version prefix
- RegisterRoutes calls SetupRoutes with empty version

### Test Requirements

```go
func TestNewServer(t *testing.T) {
    cfg := &config.Config{
        Server: config.ServerConfig{Port: 0},
        Database: config.DatabaseConfig{
            Type:   "sqlite",
            DBName: ":memory:",
        },
        Auth: config.AuthConfig{
            SessionSecret: "test-secret",
        },
    }
    
    srv, err := NewServer(cfg)
    if err != nil { t.Fatalf("NewServer() error = %v", err) }
    defer srv.Shutdown(context.Background())
    
    if srv.Engine() == nil { t.Error("Engine() returned nil") }
}

func TestServerStart(t *testing.T) {
    // Start server in background, verify it starts
    cfg := &config.Config{...}
    srv, _ := NewServer(cfg)
    
    go func() {
        srv.Start() // Will fail on port 0 but proves it started
    }()
    
    time.Sleep(100 * time.Millisecond)
    srv.Shutdown(context.Background())
}
```

## Global Constraints

- Go 1.25+
- Use `fmt.Sprintf` for address formatting

## Test Verification

```bash
cd server && go test -v -run TestNewServer
```

Expected: Tests pass

## Commit Message

```
feat(server): add HTTP server core

- Implement Server with Gin engine
- Add automatic middleware registration
- Add default health check route
- Implement Start() and Shutdown()
- Add graceful shutdown support

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-11-report.md`
