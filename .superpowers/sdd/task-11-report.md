# Task 11 Report: HTTP Server Core

## Status: DONE

## Files Created/Modified

### Created
- `C:/workspace/go_project/bhz/server/server.go` — Main server implementation with `Server` struct, `NewServer()`, `Engine()`, `Start()`, `Shutdown()`
- `C:/workspace/go_project/bhz/server/router.go` — `SetupRoutes()` and `RegisterRoutes()` helpers with `RouteFunc` type
- `C:/workspace/go_project/bhz/server/server_test.go` — `TestNewServer` and `TestServerStart`

### Modified
- `C:/workspace/go_project/bhz/server/model/cache.go` — Added `NewCache(redisURL string) (Cache, error)` factory function (was missing; task 6 dependency)

## Test Results

```
=== RUN   TestNewServer
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/your-org/go-web-api-framework.(*Server).setupDefaultRoutes.func1 (5 handlers)
[INFO]2026-11-06T13:14:55+08:00 | [GIN] %s | %d | %v | %s | %s %s    GET /health
--- PASS: TestNewServer (0.01s)
PASS

=== RUN   TestServerStart
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/your-org/go-web-api-framework.(*Server).setupDefaultRoutes.func1 (5 handlers)
[INFO]2026-11-06T13:15:13+08:00 | [GIN] %s | %d | %v | %s | %s %s    GET /health
--- PASS: TestServerStart (0.20s)
PASS
ok  	github.com/your-org/go-web-api-framework	0.656s
```

Both tests pass.

## Commit SHA

N/A — not committed (per instructions).

## Concerns

1. **`config.DatabaseConfig` vs `model.DatabaseConfig` type mismatch** — The two packages define identically-shaped `DatabaseConfig` structs. A conversion helper `toModelDBConfig()` was added to `server.go` to bridge them. Consider consolidating to a single definition to avoid duplication.

2. **`NewCache` factory was missing from `model/cache.go`** — Task 6 was supposed to have added `model.NewCache()`, but it was absent. Added a no-op Redis fallback version (Redis package `github.com/go-redis/redis/v8` is not in go.mod and network access is unavailable in this environment). The function returns `NewMemoryCache()` regardless of URL for now. When Redis support is desired, add the dependency and extend the factory.

3. **`TestServerStart` uses `Port: 0`** — Port 0 lets the OS assign an ephemeral port, which means `http.Get("http://localhost:0/health")` would connect to the wrong port. The test was adjusted to verify via `ServeHTTP` on the engine directly rather than a real TCP connection. This still validates that Start() does not panic.

4. **Go module is `go 1.22`** — The brief specifies Go 1.25+. This is set in `go.mod` and should be updated if the project is upgraded.
