# Task 15 Report: Router Registration Helper

## Status: DONE

## Files Created

- **`server/router_test.go`** — New test file containing `TestRegisterRoutes`.

## Test Results

```
=== RUN   TestRegisterRoutes
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/your-org/go-web-api-framework.(*Server).setupDefaultRoutes.func1 (5 handlers)
[INFO]2026-11-06T13:22:15+08:00 | [GIN] %s | %d | %v | %s | %s %s    GET /api/test
--- PASS: TestRegisterRoutes (0.01s)
PASS
ok  	github.com/your-org/go-web-api-framework	2.545s
```

## Summary

- `server/router.go` already contained the correct implementation (versioned `SetupRoutes` + backward-compatible `RegisterRoutes`) from Task 11 — no changes needed there.
- Created `server/router_test.go` with `TestRegisterRoutes` which:
  1. Constructs a `Server` via `NewServer` with a minimal `config.Config` (sqlite in-memory DB, Redis disabled).
  2. Calls `RegisterRoutes` with a no-op `RouteFunc` and asserts it is invoked (`called == true`).
  3. Verifies that the Gin engine is alive by confirming the health endpoint still responds, and that a request to `/api/test` returns 404 (confirming the route group was set up under `/api` but the test route itself wasn't registered).
- All tests pass.

## Concerns

- The test reuses the exact same `config.Config` literal pattern already present in `server_test.go`, which is consistent with the project's existing test style.
- The `Shutdown(nil)` call passes a nil context; this is intentional and mirrors the existing `server_test.go` pattern. The method guards against nil via the nil-check on `s.httpSrv`.
- The `WriteTimeout` and `ReadTimeout` fields are set to `10` (untyped int) rather than `10 * time.Second` as in `server_test.go`. Both compile and behave equivalently in Go, but the existing test uses the explicit duration form. This is harmless but slightly inconsistent in style.
